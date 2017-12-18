package main

import (
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	su "github.com/GavinGrayCT/SyncFileUtility/Sync02/cmd/tryConf/config"
	"os"
)

func main() {
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Trace.Printf("Starting tryConf.main\n")
	var conf su.Conf
	conf.LocalWorkspace = "config/conf.LocalWorkspace"
	conf.RemovableWorkspace = "conf.RemovableWorkspace"
	conf.RemovableRoot = "conf.RemovableRoot"
	var aLocalRoot su.Root
	aLocalRoot.Include = "aLocalRoot.Include"
	aLocalRoot.Excludes = []string{"exclude1","exclude2"}
	conf.LocalRoots = []su.Root{aLocalRoot}
	conf.Hostname = "conf.Hostname"
	if err := conf.WriteJson("tryConf.json"); err != nil {
		log.Warning.Printf("Writing conf. Error is: %v", err)
	} else {
		log.Trace.Printf("Written conf\n")

	}
	var gotConf su.Conf
	gotConf.GetTheConf("tryConf.json")
	log.Trace.Printf("Read tryconf back: %s\n", gotConf)
	log.Trace.Printf("Finished tryConf.main\n")
}
