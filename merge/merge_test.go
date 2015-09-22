package merge

import (
	. "github.com/smartystreets/goconvey/convey"
	c "github.com/vrecan/MergeForward/c"
	"testing"
	"github.com/cihub/seelog"
	"log"
)

var conf c.Conf

func TestMerge(t *testing.T) {

	Convey("Build seelog logger", t, func() {

		logger, err := seelog.LoggerFromConfigAsString(
			`<seelog type="asynctimer" asyncinterval="1000000">` +
				`<outputs formatid="all">` +
					`<filter levels="info" formatid="fmtinfo">` +
						`<console/>` +
						`<rollingfile type="size" filename="/var/log/persistent/MergeForward.log" maxsize="20000000" maxrolls="5" />` +
					`</filter>` +
					`<filter levels="warn" formatid="fmtwarn">` +
						`<console/>` +
						`<rollingfile type="size" filename="/var/log/persistent/MergeForward.log" maxsize="20000000" maxrolls="5" />` +
					`</filter>` +
					`<filter levels="error,critical" formatid="fmterror">` +
						`<console/>` +
						`<rollingfile type="size" filename="/var/log/persistent/MergeForward.log" maxsize="20000000" maxrolls="5" />` +
					`</filter>` +
				`</outputs>` +
				`<formats>` +
					`<format id="fmtinfo" format="%EscM(32)[%Level]%EscM(0) [%Date %Time] [%File] %Msg%n"/>` +
					`<format id="fmterror" format="%EscM(31)[%LEVEL]%EscM(0) [%Date %Time] [%FuncShort @ %File.%Line] %Msg%n"/>` +
					`<format id="fmtwarn" format="%EscM(33)[%LEVEL]%EscM(0) [%Date %Time] [%FuncShort @ %File.%Line] %Msg%n"/>` +
					`<format id="all" format="%EscM(2)[%LEVEL]%EscM(0) [%Date %Time] [%FuncShort @ %File.%Line] %Msg%n"/>` +
				`</formats>` +
			`</seelog>`)

		if err != nil {
			log.Fatal(err, "- This error happened while automatically detecting the current directory of mergeforward")
		}

		defer logger.Close()

		Convey("override an exact match", func() {
			conf := c.GetConf("")
			conf.ConfigOverrides["override"] = ` "value:override:value:value"`
			d := &Value{Key: "override", Value: `:  "value:value:value:value"`}
			override(d, conf)
			So(d.Value, ShouldEqual, `: "value:override:value:value"`)
		})
		Convey("override contains match", func() {
			conf := c.GetConf("")
			conf.ConfigOverrides["override"] = ` "value:override:value:value"`
			d := &Value{Key: "stuff override stuff", Value: `: "value:value:value:value"`}
			override(d, conf)
			So(d.Value, ShouldEqual, `: "value:override:value:value"`)
		})
		Convey("override contains match but no separator on matched line", func() {
			conf := c.GetConf("")
			conf.ConfigOverrides["data"] = ` D:\data`
			d := &Value{Key: "   data", Value: ``}
			override(d, conf)
			So(d.Value, ShouldEqual, ``)
		})
		Convey("Merge simple string", func() {
			src := "yamlString: value"
			dst := "yamlString: newValue"
			conf := c.GetConf("")
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, src)
		})

		Convey("Merge with = delimiter", func() {
			src := "yamlString= value"
			dst := "yamlString= newValue"
			conf := c.GetConf("")
			r, err := SimpleMerge(src, dst, "=", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, src)
		})

		Convey("Merge with multiline conf", func() {
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
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, shouldResult)
		})
		Convey("Merge with multiline conf and overrides", func() {
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
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, shouldResult)
		})
		Convey("Ignore strings with no delimiter", func() {
			src := "!!YAMLthingy.com.company.objectname"
			dst := `!!YAMLthingy.com.company.objectname
		somethingElse: "value"
		`

			conf := c.GetConf("")
			r, err := SimpleMerge(src, dst, "=", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, dst)
		})

		Convey("Test yaml lists with same object should replace in order", func() {
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
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, shouldResult)
		})

		Convey("empty value test", func() {
			src := `something:`
			dst := `something:`
			conf := c.GetConf("")
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, dst)
		})

		Convey("comment block test update", func() {
			src := `//comment block`
			dst := `//comment`
			conf := c.GetConf("")
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, dst)
		})
		Convey("comment block test update with :", func() {
			src := `something: woo
		//comment explaining something else
		something else: "new"`
			dst := `something: woo
		something else: "woo"`
			shouldResult := `something: woo
		something else: "new"`
			conf := c.GetConf("")
			r, err := SimpleMerge(src, dst, ":", conf, logger)
			So(err, ShouldBeNil)
			So(r, ShouldEqual, shouldResult)
		})
	})
}
