package main

import (
	"github.com/GavinGrayCT/SyncFileUtility/Sync01/syncutil"
	"github.com/GavinGrayCT/SyncFileUtility/Sync01/log"
	"os"
	)

var testPaths = []string{"C:/Junk/junk", "longName.xlsx", "longName.xlsxq"}
var excludes = []string{".*junk$", `.*\.xls.$`}

func main() {
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Info.Println("In tryit main")
	var matched bool
	var excludeString string
	for i := range testPaths {
		log.Trace.Println("Considering ", testPaths[i])
		matched, excludeString, _ = syncutil.CheckExclude(testPaths[i], excludes)
		if matched {
			log.Trace.Println(testPaths[i], " is excluded by ", excludeString)
		}
	}
}
