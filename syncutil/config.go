package syncutil

import (
	"io/ioutil"
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	"encoding/json"
)

const hostNameFileName = "Hostname.json"

type Root struct {
	Include  string
	Excludes []string
}

type Conf struct {
	LocalWorkspace     string
	RemovableWorkspace string
	LocalRoots         []Root
	RemovableRoots     []Root
	Hostname           string
}

func GetTheConf(configPathFile string) (*Conf, error) {
	var theConf Conf
	log.Trace.Println("In GetConf - configPathFile is ", configPathFile)
	err := readJson("theConf.json", &theConf)
	if err != nil {
		log.Trace.Println("In GetConf - reading json conf. Error: ", err)
	} else {
		log.Trace.Println("In GetConf - getting Host from ", theConf.LocalWorkspace+"/"+hostNameFileName)
		err = readJson(theConf.LocalWorkspace+"/"+hostNameFileName, &theConf)
		log.Trace.Println("Host is ", theConf.Hostname, "  and err is ", err)
	}
	log.Trace.Println("Config is \n", theConf, "\n")
	log.Trace.Println("Done GetConf - err is ", err)

	return &theConf, err
}

func readJson(pf string, s interface{}) error {
	j, err := ioutil.ReadFile(pf)
	if err == nil {
		err = json.Unmarshal(j, s)
		if err != nil {
			log.Error.Printf("Error json unmarshalling %s. Error: %v\n", pf, err)
		}
	} else {
		log.Error.Printf("Error reading %s. Error: %v\n", pf, err)

	}
	return err
}
