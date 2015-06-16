package merge

import (
	. "github.com/smartystreets/goconvey/convey"
	c "github.com/vrecan/MergeForward/c"
	"testing"
)

var conf c.Conf

func TestMerge(t *testing.T) {
	Convey("override an exact match", t, func() {
		conf := c.GetConf("")
		conf.ConfigOverrides["override"] = ` "value:override:value:value"`
		d := &Value{Key: "override", Value: `:  "value:value:value:value"`}
		override(d, conf)
		So(d.Value, ShouldEqual, `: "value:override:value:value"`)
	})
	Convey("override contains match", t, func() {
		conf := c.GetConf("")
		conf.ConfigOverrides["override"] = ` "value:override:value:value"`
		d := &Value{Key: "stuff override stuff", Value: `: "value:value:value:value"`}
		override(d, conf)
		So(d.Value, ShouldEqual, `: "value:override:value:value"`)
	})
	Convey("Merge simple string", t, func() {
		src := "yamlString: value"
		dst := "yamlString: newValue"
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, src)
	})

	Convey("Merge with = delimiter", t, func() {
		src := "yamlString= value"
		dst := "yamlString= newValue"
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, "=", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, src)
	})

	Convey("Merge with multiline conf", t, func() {
		src := `yamlString: value
		value: "value:change:value:value"
		other: "other:other:other:other"
		`
		dst := `yamlString: value
		value: "value:value:value:value"
		other: "other:other:other:other"
		new: "new:field:woo"
		`

		shouldResult := `yamlString: value
		value: "value:change:value:value"
		other: "other:other:other:other"
		new: "new:field:woo"
		`
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, shouldResult)
	})
	Convey("Merge with multiline conf and overrides", t, func() {
		src := `yamlString : value
		value : "value:change:value:value"
		override : "value:value:value:value"
		`
		dst := `yamlString : value
		value : "value:value:value:value"
		override : "value:value:value:value"
		new : "new:field:woo"
		`

		shouldResult := `yamlString : value
		value : "value:change:value:value"
		override : "value:override:value:value"
		new : "new:field:woo"
		`
		conf := c.GetConf("")
		conf.ConfigOverrides["override"] = ` "value:override:value:value"`
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, shouldResult)
	})
	Convey("Ignore strings with no delimiter", t, func() {
		src := "!!YAMLthingy.com.company.objectname"
		dst := `!!YAMLthingy.com.company.objectname
		somethingElse: "value"
		`

		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, "=", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, dst)
	})

	Convey("Test yaml lists with same object should replace in order", t, func() {
		src := `!!com.company.configelements.DispatchConfigMessage
indexerConfList:
   - !!com.company.configelements.DispatchConfigMessage$IndexerConfigMessage
      inputQueue: tcp://0.0.0.0:13109
      inputQueueTimeoutSec: 5000
      passthroughQName: tcp://127.0.0.1:13104
   - !!com.company.configelements.DispatchConfigMessage$IndexerConfigMessage
      inputQueue: tcp://127.0.0.1:13102
      inputQueueTimeoutSec: 5000
      passthroughQName: "customePass"`
		dst := `!!com.company.configelements.DispatchConfigMessage
indexerConfList:
   - !!com.company.configelements.DispatchConfigMessage$IndexerConfigMessage
      inputQueue: tcp://*:13100
      inputQueueTimeoutSec: 5000
      passthroughQName: tcp://127.0.0.1:13104
   - !!com.company.configelements.DispatchConfigMessage$IndexerConfigMessage
      inputQueue: tcp://127.0.0.1:13102
      inputQueueTimeoutSec: 5000
      passthroughQName: ""`

		shouldResult := `!!com.company.configelements.DispatchConfigMessage
indexerConfList:
   - !!com.company.configelements.DispatchConfigMessage$IndexerConfigMessage
      inputQueue: tcp://0.0.0.0:13109
      inputQueueTimeoutSec: 5000
      passthroughQName: tcp://127.0.0.1:13104
   - !!com.company.configelements.DispatchConfigMessage$IndexerConfigMessage
      inputQueue: tcp://127.0.0.1:13102
      inputQueueTimeoutSec: 5000
      passthroughQName: "customePass"`
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, shouldResult)
	})

	Convey("empty value test", t, func() {
		src := `something:`
		dst := `something:`
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, dst)
	})

	Convey("comment block test update", t, func() {
		src := `//comment block`
		dst := `//comment`
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, dst)
	})
	Convey("comment block test update with :", t, func() {
		src := `something: woo
		//comment explaining something else
		something else: "new"`
		dst := `something: woo
		something else: "woo"`
		shouldResult := `something: woo
		something else: "new"`
		conf := c.GetConf("")
		r, err := SimpleMerge(src, dst, ":", conf)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, shouldResult)
	})

}
