package main

// Merges files together and prints result to stdout

import (
	"flag"
	"fmt"
	c "github.com/vrecan/MergeForward/c"
	"github.com/vrecan/MergeForward/merge"
	"io/ioutil"
	"os"
	"log"
	"io"

)

var src = flag.String("src", "", "Source configuration file. Src values are prefered over dest.")
var dst = flag.String("dst", "", "Destination configuration file.")
var split = flag.String("split", ":", "Splitter for key value pairs")
var config = flag.String("config", "./c.ini", "MergeForward configuration ini file")
var outputdir = flag.String("outputdir", "", "The output directory of the log file")

type conf interface{}

func main() {
	flag.Parse()

	_ = os.Remove(*outputdir + "MergeForward.log")
	logFile, err := os.OpenFile(*outputdir + "MergeForward.log", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error making log file: " + *outputdir + string(os.PathSeparator) + "MergeForward.log")
	}
	defer logFile.Close()
	
	dualLogger := log.New(io.MultiWriter(logFile,os.Stdout), "", 0)

	if len(*src) == 0 {
		dualLogger.Fatalln("No source file")
	}
	if len(*dst) == 0 {
		dualLogger.Fatalln("No destination file.")
	}
	conf := c.GetConf(*config)

	srcBytes, err := ioutil.ReadFile(*src)
	if nil != err {
		dualLogger.Fatalln("Unable to read src file: ", err)
	}

	dstBytes, err := ioutil.ReadFile(*dst)
	if nil != err {
		dualLogger.Fatalln("Unable to read destination file: ", err)
	}

	result, err := merge.SimpleMerge(string(srcBytes), string(dstBytes), *split, conf, logFile)
	if nil != err {
		dualLogger.Fatalln("Unable to merge files: ", err)
	} else {
		dualLogger.Println(result)
	}
}
