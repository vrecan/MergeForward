package merge

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMerge(t *testing.T) {
	Convey("Merge simple string", t, func() {
		src := "yamlString: value"
		dst := "yamlString: newValue"
		r, err := SimpleMerge(src, dst, ":")
		So(err, ShouldBeNil)
		So(r, ShouldEqual, src)
	})

	Convey("Merge with = delimiter", t, func() {
		src := "yamlString= value"
		dst := "yamlString= newValue"
		r, err := SimpleMerge(src, dst, "=")
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
		r, err := SimpleMerge(src, dst, ":")
		So(err, ShouldBeNil)
		So(r, ShouldEqual, shouldResult)
	})

	Convey("Ignore strings with no delimiter", t, func() {
		src := "!!YAMLthingy.com.company.objectname"
		dst := `!!YAMLthingy.com.company.objectname
		somethingElse: "value"
		`

		r, err := SimpleMerge(src, dst, "=")
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

		r, err := SimpleMerge(src, dst, ":")
		So(err, ShouldBeNil)
		So(r, ShouldEqual, shouldResult)
	})

	Convey("empty value test", t, func() {
		src := `something:`
		dst := `something:`

		r, err := SimpleMerge(src, dst, ":")
		So(err, ShouldBeNil)
		So(r, ShouldEqual, dst)
	})

}
