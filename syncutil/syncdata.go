package syncutil

import "github.com/GavinGrayCT/SyncFileUtility/Sync02/log"

type LastMaster struct {
	LastSyncMaster string
}
type SyncData struct {
	TheConf *Conf
	LocalCurrentDirMap DirMap
	TheLastMaster LastMaster
}


var TheSyncData SyncData

func FillData(theConf *Conf) (*SyncData, error) {
	TheSyncData.TheConf = theConf
	TheSyncData.LocalCurrentDirMap = make(map[string]*Dir)
	var err error
	for r := range theConf.Roots {

		err = FillDirMap(&(theConf.Roots[r].Include), theConf.Roots[r].Excludes, TheSyncData.LocalCurrentDirMap)
		if err != nil {
			log.Fatal.Println("Fatal error filling Local Current Dirmap. Error is: ", err)
		}
	}

	return &TheSyncData, err
}
