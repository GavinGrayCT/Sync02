package main

import (
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	su "github.com/GavinGrayCT/SyncFileUtility/Sync02/syncutil"
	"os"
)

var theConf *su.Conf
var theSyncData *su.SyncData

func main() {
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Trace.Println("Started SyncFileUtility")

	// Read config
	var err error
	theConf, err = su.GetTheConf("gghome.yaml")
	if err != nil {
		log.Fatal.Println("Fatal Error getting config: ", err)
	}

	// Fill data structures
	// Local Current Dir Map
	theSyncData, err = su.FillData(theConf)
	log.Trace.Println("Filled sync data")
}

