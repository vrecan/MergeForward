package merge

import (
	"bufio"
	"bytes"
	c "github.com/vrecan/MergeForward/c"
	"strings"
	"log"
	"os"
)

type Value struct {
	Key   string
	Value string
}

type Merge struct {
	Values []*Value
	conf   c.Conf
}

var SPLIT = ":"

//Merge the old values(src) into the new values (dst)
func SimpleMerge(src string, dst string, split string, conf c.Conf, logFile *os.File) (result string, err error) {
	log.SetOutput(logFile)
	SPLIT = split
	reader := bytes.NewBufferString(src)
	scanner := bufio.NewScanner(reader)
	merge := &Merge{conf: conf}
	
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
				if srcCnt[s.Key] == dstCnt[s.Key] && s.Value != d.Value {
					log.Println("Replacing new {", d.Key + d.Value, "}")
					log.Println("with the old  {", s.Key + s.Value, "}")
					d.Value = s.Value
				}
			}
			override(d, src.conf)
		}
	}
	return dst
}

// override the value if its key is in the conf.ConfigOverrides
func override(d *Value, conf c.Conf) {
	for k, v := range conf.ConfigOverrides {
		if strings.Contains(d.Key, k) && len(d.Value) != 0 {
			var buffer bytes.Buffer
			buffer.WriteString(string([]rune(d.Value)[0]))
			buffer.WriteString(v)
			d.Value = buffer.String()
			return
		}
	}
}
