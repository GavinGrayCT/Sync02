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
	TheConf                  Conf
	CurrentLocalSync         Sync
	LastLocalSync            Sync
	CurrentRemovableSync     Sync
	LocalLastSyncOnRemovable Sync
	OtherLastSyncOnRemovable Sync
}

type FileStates struct {
	cld  *Dir  // current local Dir
	clf  Afile // current local File
	llf  Afile // last local File
	llrf Afile // last local on removable File
	crf  Afile // current on removable File
	lorf Afile // last other on removable File
}

func check(err error, m string) {
	if err != nil {
		log.Fatal.Printf("%s. Error is: %v", m, err)
	}
}

func FillSyncData(configPathFile string) (*SyncData, error) {
	var err error
	var TheSyncData SyncData
	TheSyncData.TheConf.GetTheConf(configPathFile)
	check(err, "Fatal Error getting config")

	err = fillCurrentSync(TheSyncData.TheConf.Hostname, TheSyncData.TheConf.LocalRoot, &TheSyncData.CurrentLocalSync)
	check(err, "Fatal Error filling CurrentLocalSync")

	err = readSyncFromFile(TheSyncData.TheConf.LocalWorkspace+"/"+"LastLocalSync.json", &TheSyncData.LastLocalSync)
	check(err, "Fatal Error reading LastLocalSync.json")

	err = fillCurrentSync(TheSyncData.TheConf.Hostname, Root{TheSyncData.TheConf.RemovableRoot, []string{}}, &TheSyncData.CurrentRemovableSync)
	check(err, "Fatal Error filling CurrentRemovableSync")

	err = readSyncFromFile(TheSyncData.TheConf.RemovableWorkspace+"/"+TheSyncData.TheConf.Hostname+"LastSyncOnRemovable.json", &TheSyncData.LocalLastSyncOnRemovable)
	check(err, "Fatal Error reading LastSyncOnRemovable.json for host "+TheSyncData.TheConf.Hostname)

	err = readSyncFromFile(TheSyncData.TheConf.RemovableWorkspace+"/"+TheSyncData.TheConf.OtherHostname+"LastSyncOnRemovable.json", &TheSyncData.OtherLastSyncOnRemovable)
	check(err, "Fatal Error reading LastSyncOnRemovable.json for host "+TheSyncData.TheConf.OtherHostname)

	if compare2Sync(TheSyncData.CurrentLocalSync, TheSyncData.LastLocalSync) {
		log.Trace.Println("Current and Last Local Syncs are identical")
	} else {
		log.Trace.Println("Current and Last Local Syncs are different")
	}

	return &TheSyncData, err
}

func fillCurrentSync(hostname string, Root Root, theSync *Sync) error {
	var err error
	theSync.Hostname = hostname
	theSync.TheDirMap = make(map[string]*Dir)
	aRootedDirMap := RootedDirMap{Root.Include, Root.Excludes, theSync.TheDirMap}
	err = FillDirMap(aRootedDirMap)
	if err != nil {
		log.Fatal.Println("Fatal error filling Local Current Dirmap. Error is:", err)
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

func ProcessSyncData(theSyncData *SyncData) {
	log.Trace.Println("Starting ProcessSyncData")
	var fs FileStates
	for _, aCLD := range theSyncData.CurrentLocalSync.TheDirMap {
		log.Trace.Println("Dir: ", aCLD.Path)
		for _, aCLF := range aCLD.Files {
			log.Trace.Println("File: ", aCLF.Name, "  FModDate:", aCLF.FModDate, "  Size:", aCLF.Size, "  Present:", aCLF.Present)
			fs.populate(aCLD, aCLF, theSyncData)
			// now process fs
			AddCommand(aCLD, &fs)
		}
	}
}

func (fs *FileStates) populate(cld *Dir, clf Afile, theSyncData *SyncData) {
	var ok bool
	fs.cld = cld
	fs.clf = clf
	log.Trace.Println("fs.cld.RelPath:", fs.cld.RelPath, "fs.clf.Name:", fs.clf.Name)
	fs.llf, ok = theSyncData.LastLocalSync.TheDirMap[fs.cld.RelPath].Files[fs.clf.Name]
	if !ok {
		fs.llf = Afile{fs.clf.Name, 0, time.Now(), false}
	}

	var llrd *Dir
	llrd, ok = theSyncData.LocalLastSyncOnRemovable.TheDirMap[fs.cld.RelPath]
	if ok {
		fs.llrf, ok = llrd.Files[fs.clf.Name]
		if !ok {
			fs.llrf = Afile{fs.clf.Name, 0, time.Now(), false}
		} else {
			log.Trace.Println("llrf --- File: ", fs.llrf.Name, "  FModDate:", fs.llrf.FModDate, "  Size:", fs.llrf.Size, "  Present:", fs.llrf.Present)
		}
	} else {
		log.Trace.Printf("Dir %s does not exist in theSyncData.LocalLastSyncOnRemovable.TheDirMap\n", fs.cld.RelPath)
	}

	var crd *Dir
	crd, ok = theSyncData.CurrentRemovableSync.TheDirMap[fs.cld.RelPath]
	if ok {
		fs.crf, ok = crd.Files[fs.clf.Name]
		if !ok {
			fs.crf = Afile{fs.clf.Name, 0, time.Now(), false}
		} else {
			log.Trace.Printf("Dir %s does not exist in theSyncData.LocalLastSyncOnRemovable.TheDirMap\n", fs.cld.RelPath)
		}

		var lord *Dir
		lord, ok = theSyncData.OtherLastSyncOnRemovable.TheDirMap[fs.cld.RelPath]
		if ok {
			fs.lorf, ok = lord.Files[fs.clf.Name]
			if !ok {
				log.Trace.Printf("lorf %s does not exist in theSyncData.OtherLastSyncOnRemovable.TheDirMap dir %s\n", fs.clf.Name, fs.cld.RelPath)
				fs.lorf = Afile{fs.clf.Name, 0, time.Now(), false}
			} else {
				log.Trace.Printf("lorf %s does exist in theSyncData.OtherLastSyncOnRemovable.TheDirMap dir %s\n", fs.clf.Name, fs.cld.RelPath)
			}
		} else {
			log.Trace.Printf("Dir %s does not exist in theSyncData.OtherLastSyncOnRemovable.TheDirMap\n", fs.cld.RelPath)
		}

		log.Trace.Println("FileStates: ", fs)
	}
}