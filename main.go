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

	// Fill data structures
	// Local Current Dir Map
	theSyncData, err = su.FillSyncData("theConf.json")
	if err != nil {
		log.Fatal.Printf("Error Filling SyncData %v", err)
	}
	su.StoreSync(&theSyncData.CurrentLocalSync, "WrittenCurrentLocalSync")
	log.Trace.Println("Filled sync data")
	log.Trace.Println("Finished SyncFileUtility")
}

