package syncutil

import (
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
)

const hostNameFileName = "Hostname.json"

type Root struct {
	Include  string
	Excludes []string
}

type Conf struct {
	LocalWorkspace     string
	RemovableWorkspace string
	RemovableRoot      string
	LocalRoot          Root
	Hostname           string
	OtherHostname      string
}

func (theConf *Conf) GetTheConf(configPathFile string) (error) {
	log.Trace.Println("In GetConf - configPathFile is ", configPathFile)
	err := ReadJson(configPathFile, &theConf)
	if err != nil {
		log.Trace.Println("In GetConf - reading json conf. Error: ", err)
	} else {
		log.Trace.Println("In GetConf - getting Host from ", theConf.LocalWorkspace+"/"+hostNameFileName)
		err = ReadJson(theConf.LocalWorkspace+"/"+hostNameFileName, &theConf)
		log.Trace.Println("Host is ", theConf.Hostname, "  and err is ", err)
		log.Trace.Println("Other Host is ", theConf.OtherHostname, "  and err is ", err)
	}
	log.Trace.Println("Config is \n", theConf, "\n")
	log.Trace.Println("Done GetConf - err is ", err)

	return err
}

