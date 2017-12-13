package main

import (
	"testing"
	su "github.com/GavinGrayCT/SyncFileUtility/Sync02/syncutil"
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	"os"
	"io/ioutil"
)

const tRootPath = "c:\\Users\\Gavin\\6Saxon\\junk"

//var theConf *su.Conf

func init() {
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
}

func printOutDirMap(aDirMap *su.DirMapEx) {
	log.Trace.Printf("Printing out dirMap\n")
	for _, aDir := range aDirMap.Dirmap {
		log.Trace.Printf("Dir %s\n", aDir.Path)
		for _, aFile := range aDir.Files {
			log.Trace.Printf("File %s\n", aFile.Name)
		}
	}
}

func removeFileFromDirMap(aDirMap *su.DirMapEx) {
	log.Trace.Printf("removeFileFromDirMap\n")
	for _, aDir := range aDirMap.Dirmap {
		if _, ok := aDir.Files["WiredTigerLog.0000000001"]; ok {
			delete(aDir.Files, "WiredTigerLog.0000000001")
			log.Trace.Println("Removed WiredTigerLog.0000000001")
		}
	}
}

func TestConfig1(t *testing.T) {
	log.Trace.Println("In TestConfig1")
	var err error

	theConf, err = su.GetTheConf("gghome.yaml")
	if err != nil {
		log.Error.Println("Error message: ", err)
	} else {
		log.Trace.Printf("Config is %s", theConf)
		log.Trace.Printf("Local Workspace is %s\n", theConf.LocalWorkspace)
		log.Trace.Printf("Removable Workspace is %s\n", theConf.RemovableWorkspace)
		for r := range theConf.LocalRoot {
			log.Trace.Printf("Includes %s\n", theConf.LocalRoot[r].Include)
			for i := range theConf.LocalRoot[r].Excludes {
				log.Trace.Printf("Type of exclude is %T\n", i)
				log.Trace.Printf("Excludes %s\n", theConf.LocalRoot[r].Excludes[i])
			}
		}
	}
	if err != nil {
		t.Error("Testing getting config - false")
	}
	log.Trace.Println("Done TestConfig1")
}

