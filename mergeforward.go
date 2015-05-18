package main

// Merges files together and prints result to stdout

import (
	"flag"
	"fmt"
	"github.com/vrecan/MergeForward/merge"
	"io/ioutil"
	"os"
)

var src = flag.String("src", "", "Source configuration file. Src values are prefered over dest.")
var dst = flag.String("dst", "", "Destination configuration file.")
var split = flag.String("split", ":", "What character we should use to split key value pairs")

type conf interface{}

func main() {
	flag.Parse()
	if len(*src) == 0 {
		fmt.Println("No source file.")
		os.Exit(1)
	}
	if len(*dst) == 0 {
		fmt.Println("No destination file.")
		os.Exit(1)
	}

	srcBytes, err := ioutil.ReadFile(*src)
	if nil != err {
		fmt.Println("Unable to read src file: ", err)
		os.Exit(1)
	}

	dstBytes, err := ioutil.ReadFile(*dst)
	if nil != err {
		fmt.Println("Unable to read destination file: ", err)
		os.Exit(1)
	}

	result, err := merge.SimpleMerge(string(srcBytes), string(dstBytes), *split)
	if nil != err {
		fmt.Println("Unable to merge files: ", err)
		os.Exit(1)
	} else {
		fmt.Print(result)
	}
}
