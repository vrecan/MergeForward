package main

// Merges files together and prints result to stdout

import (
	"flag"
	"fmt"
	c "github.com/vrecan/MergeForward/c"
	"github.com/vrecan/MergeForward/merge"
	"io/ioutil"
	"os"
	"github.com/cihub/seelog"
	"path/filepath"
	"log"
)

var src = flag.String("src", "", "Source configuration file. Src values are prefered over dest.")
var dst = flag.String("dst", "", "Destination configuration file.")
var split = flag.String("split", ":", "Splitter for key value pairs")
var config = flag.String("config", "./c.ini", "MergeForward configuration ini file")

type conf interface{}

func main() {

	currentDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err, "- This error happened while automatically detecting the current directory of mergeforward")
	}

	fmt.Print("Printing current directory: ")
	fmt.Println(currentDirectory)

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

	flag.Parse()

	defer logger.Close()

	if len(*src) == 0 {
		logger.Critical("No source file")
		os.Exit(1)
	}
	if len(*dst) == 0 {
		logger.Critical("No destination file.")
		os.Exit(1)
	}

	logger.Info("Merging ", *src, " into ", *dst, "...")

	conf := c.GetConf(*config)

	srcBytes, err := ioutil.ReadFile(*src)
	if nil != err {
		logger.Critical("Unable to read src file: ", err)
		os.Exit(1)
	}

	dstBytes, err := ioutil.ReadFile(*dst)
	if nil != err {
		logger.Critical("Unable to read destination file: ", err)
		os.Exit(1)
	}

	result, err := merge.SimpleMerge(string(srcBytes), string(dstBytes), *split, conf, logger)
	if nil != err {
		logger.Critical("Unable to merge files: ", err)
		os.Exit(1)
	} else {
		fmt.Println(result)
		logger.Info("Final file output:\n" + result + "\n")
	}
}
