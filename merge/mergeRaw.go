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

//Merge the old values(src) into the new values (dst)
func SimpleMerge(src string, dst string) (result string, err error) {
	reader := bytes.NewBufferString(src)
	scanner := bufio.NewScanner(reader)
	merge := &Merge{}
	for scanner.Scan() {
		srcParts := strings.Split(scanner.Text(), ":")
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
	cnt := 2
	size := len(parts)
	for _, p := range parts[1:] {

		buffer.WriteString(p)
		if size > cnt {
			buffer.WriteString(":")
		}
		cnt++
	}
	value.Value = buffer.String()
	s.Values = append(s.Values, &value)
}

//Merge values from the old values map into the new string.
func (s *Merge) MergeInto(dst string) (string, error) {

	reader := bytes.NewBufferString(dst)
	scanner := bufio.NewScanner(reader)
	merge := &Merge{}
	for scanner.Scan() {
		dstParts := strings.Split(scanner.Text(), ":")
		merge.AddValues(dstParts)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	merge = Combine(s, merge)
	var result bytes.Buffer
	for _, v := range merge.Values {
		result.WriteString(v.Key)
		if len(v.Value) > 0 {
			// fmt.Println("FOUND MATCH===>", v.Key, v.Value)
			result.WriteString(":")
			result.WriteString(v.Value)

		}
		result.WriteString("\n")
	}
	return result.String(), nil
}

//Combine our results and return the result
func Combine(src *Merge, dst *Merge) *Merge {
	for _, s := range src.Values {
		for _, d := range dst.Values {
			if s.Key == d.Key {
				d.Value = s.Value
			}
		}
	}
	return dst
}
