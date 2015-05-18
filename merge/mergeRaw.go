package merge

import (
	"bufio"
	"bytes"
	// "fmt"
	"strings"
)

type Value struct {
	Key   string
	Value string
}

type Merge struct {
	Values []*Value
}

var SPLIT = ":"

//Merge the old values(src) into the new values (dst)
func SimpleMerge(src string, dst string, split string) (result string, err error) {
	SPLIT = split
	reader := bytes.NewBufferString(src)
	scanner := bufio.NewScanner(reader)
	merge := &Merge{}
	for scanner.Scan() {
		srcParts := strings.Split(scanner.Text(), SPLIT)
		merge.AddValues(srcParts)
	}
	if err = scanner.Err(); err != nil {
		return result, err
	}
	result, err = merge.MergeInto(dst)
	return
}

//Take slice of split with everything after the first split being the value
func (s *Merge) AddValues(parts []string) {
	var value Value
	if len(parts) <= 0 {
		return //no values to add
	}
	if len(parts) > 0 {
		value.Key = parts[0]
	}

	var buffer bytes.Buffer
	size := len(parts) - 1
	for i, p := range parts {
		if i != 0 { // don't write the key twice
			buffer.WriteString(p)
		}
		if i < size { //ignore the last one
			buffer.WriteString(SPLIT)
		}

	}
	r := buffer.String()
	if len(r) > 0 {
		value.Value = r
	}
	s.Values = append(s.Values, &value)

}

//Merge values from the old values map into the new string.
func (s *Merge) MergeInto(dst string) (string, error) {

	reader := bytes.NewBufferString(dst)
	scanner := bufio.NewScanner(reader)
	merge := &Merge{}
	for scanner.Scan() {
		dstParts := strings.Split(scanner.Text(), SPLIT)
		merge.AddValues(dstParts)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	merge = Combine(s, merge)
	var result bytes.Buffer
	cnt := 1
	size := len(merge.Values)
	for _, v := range merge.Values {
		result.WriteString(v.Key)
		if len(v.Value) > 0 {
			// fmt.Println("FOUND MATCH===>", v.Key, v.Value)
			result.WriteString(v.Value)

		}
		if cnt < size {
			result.WriteString("\n")
		}
		cnt++
	}
	return result.String(), nil
}

//Combine our results and return the result
func Combine(src *Merge, dst *Merge) *Merge {
	srcCnt := make(map[string]int, 0)

	for _, s := range src.Values {

		_, ok := srcCnt[s.Key]
		if !ok {
			srcCnt[s.Key] = 1
		} else {
			srcCnt[s.Key]++
		}
		dstCnt := make(map[string]int, 0)
		for _, d := range dst.Values {
			if s.Key == d.Key {
				_, ok := dstCnt[s.Key]
				if !ok {
					dstCnt[s.Key] = 1
				} else {
					dstCnt[s.Key]++
				}
				if srcCnt[s.Key] == dstCnt[s.Key] {
					d.Value = s.Value
				}
			}
		}
	}
	return dst
}
