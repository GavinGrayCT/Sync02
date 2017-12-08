package syncutil

import (
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	"time"
	"io/ioutil"
	"encoding/json"
)

type Sync struct {
	Hostname  string
	SyncTime  time.Time
	TheDirMap DirMap
}

type SyncData struct {
	TheConf              *Conf
	CurrentLocalSync     Sync
	LastLocalSync        Sync
	CurrentRemovableSync Sync
	LastRemovableSync    Sync
}

func check(err error, m string) {
	if err != nil {
		log.Fatal.Printf("%s. Error is: %v", m, err)
	}
}

func FillSyncData(configPathFile string) (*SyncData, error) {
	var err error
	var TheSyncData SyncData
	TheSyncData.TheConf, err = GetTheConf(configPathFile)
	check(err, "Fatal Error getting config")

	err = fillCurrentLocalSync(TheSyncData.TheConf.Hostname, TheSyncData.TheConf.LocalRoots, &TheSyncData.CurrentLocalSync)
	check(err, "Fatal Error filling CurrentLocalSync")

	err = readSyncFromFile(TheSyncData.TheConf.LocalWorkspace+"/"+"LastLocalSync.json", &TheSyncData.LastLocalSync)
	check(err, "Fatal Error reading LastLocalSync.json")

	err = fillCurrentLocalSync(TheSyncData.TheConf.Hostname, TheSyncData.TheConf.RemovableRoots, &TheSyncData.CurrentRemovableSync)
	check(err, "Fatal Error filling CurrentRemovableSync")

	err = readSyncFromFile(TheSyncData.TheConf.RemovableWorkspace+"/"+"LastRemovableSync.json", &TheSyncData.LastRemovableSync)
	check(err, "Fatal Error reading LastRemovableSync.json")

	if compare2Sync(TheSyncData.CurrentLocalSync, TheSyncData.LastLocalSync) {
		log.Trace.Println("Current and Last Local Syncs are identical")
	} else {
		log.Trace.Println("Current and Last Local Syncs are different")

	}

	return &TheSyncData, err
}

func fillCurrentLocalSync(hostname string, Roots []Root, theSync *Sync) error {
	var err error
	theSync.Hostname = hostname
	theSync.TheDirMap = make(map[string]*Dir)
	for r := range Roots {

		err = FillDirMap(&(Roots[r].Include), Roots[r].Excludes, theSync.TheDirMap)
		if err != nil {
			log.Fatal.Println("Fatal error filling Local Current Dirmap. Error is:", err)
		}
	}
	theSync.SyncTime = time.Now()
	return err
}

func readSyncFromFile(pf string, theSync *Sync) (err error) {
	log.Trace.Printf("Reading ", pf, "into LastLocalSync")
	err = readJson(pf, &theSync)
	return err
}

func StoreSync(theSync *Sync, pathFile string) error {
	j, err := json.MarshalIndent(theSync, "", "   ")
	check(err, "Error in Json marshall")
	err = ioutil.WriteFile(pathFile+".json", j, 0644)
	check(err, "Error in writing to file")
	return err
}

func writeConfJson(theConf *Conf) {
	j, err := json.MarshalIndent(theConf, "", "   ")
	err = ioutil.WriteFile("theConf.json", j, 0644)
	log.Trace.Println("Written json. Error is: ", err)
}

func compare2Sync(syncA, syncB Sync) (eq bool) {
	eq = false
	log.Trace.Printf("syncA hostname %s,  SyncB hostname %s\n", syncA.Hostname, syncB.Hostname)
	if syncA.Hostname == syncB.Hostname {
		// SyncTime will always be different
		if CompDirMap(syncA.TheDirMap, syncB.TheDirMap) {
			eq = true
		}
	}
	return eq
}
