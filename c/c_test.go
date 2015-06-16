package merge

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMerge(t *testing.T) {
	Convey("reads config file literally for value", t, func() {
		var buffer bytes.Buffer
		buffer.WriteString(
			`foo=bar
baz="fubar"
boo='bar'
fu=bar=bar
space= a`)
		c := readFromBuffer(buffer.Bytes())
		So(c.ConfigOverrides, ShouldNotBeEmpty)
		So(c.ConfigOverrides["foo"], ShouldEqual, "bar")
		So(c.ConfigOverrides["baz"], ShouldEqual, `"fubar"`)
		So(c.ConfigOverrides["boo"], ShouldEqual, `'bar'`)
		So(c.ConfigOverrides["fu"], ShouldEqual, `bar=bar`)
		So(c.ConfigOverrides["space"], ShouldEqual, ` a`)
	})
	Convey("Read from invalid file", t, func() {
		var buffer bytes.Buffer
		buffer.WriteString(
			`foo bar`)
		c := readFromBuffer(buffer.Bytes())
		So(c.ConfigOverrides, ShouldBeEmpty)
	})
	Convey("Read from non existing file", t, func() {
		c := GetConf("/there/is/no/way/this/file/exists")
		So(c.ConfigOverrides, ShouldBeEmpty)
	})
	Convey("Read from empty path", t, func() {
		c := GetConf("")
		So(c.ConfigOverrides, ShouldBeEmpty)
	})
}
