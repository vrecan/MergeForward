package merge

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

type Conf struct {
	ConfigOverrides map[string]string
}

// String version of the conf
func (c Conf) String() string {
	return fmt.Sprint("ConfigOverrides:", c.ConfigOverrides, "\n")
}

// GetConf from the file path sepcified
func GetConf(path string) (c Conf) {
	c = Conf{}
	if len(path) == 0 {
		c.ConfigOverrides = make(map[string]string)
		return c
	}
	data, err := ioutil.ReadFile(path)
	if nil != err {
		c.ConfigOverrides = make(map[string]string)
		return c
	}
	return readFromBuffer(data)
}

// readFromBuffer given and populate the conf
func readFromBuffer(data []byte) (c Conf) {
	c.ConfigOverrides = make(map[string]string)
	entries := strings.Split(string(data), "\n")
	for _, line := range entries {
		keyAndValue := strings.Split(line, "=")
		if len(keyAndValue) > 1 {
			var buffer bytes.Buffer
			buffer.WriteString(keyAndValue[1])
			for i := 2; i < len(keyAndValue); i++ {
				buffer.WriteString("=")
				buffer.WriteString(keyAndValue[i])
			}
			c.ConfigOverrides[keyAndValue[0]] = buffer.String()
		}
	}
	return c
}
