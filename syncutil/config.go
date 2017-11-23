package syncutil

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	"fmt"
)

const hostNameFileName = "Hostname.yaml"

type Root struct {
	Include string
	Excludes []string
}

type Conf struct {
	LocalWorkspace string
	RemovableWorkspace string
	Roots []Root
	Hostname string
}

var theConf Conf

var wConf = Conf{"abc", "bcd", []Root{{"inc1", []string{"ex1-1", "ex1-2"}}}, "hosty"}
var wConf1 = Conf{"abc", "bcd", []Root{{"inc1", []string{"ex1-1", "ex1-2"}}, {"inc2", []string{"ex2-1", "ex2-2"}}}, "hosty"}
var wConf2 = Conf{"abc", "bcd", []Root{{"inc1", []string{"ex1-1"}}}, "hosty"}

func GetTheConf(configPathFile string) (*Conf, error) {
	log.Trace.Println("In GetConf - configPathFile is ", configPathFile)

	d, err := yaml.Marshal(&wConf)
	if err != nil {
		log.Fatal.Printf("error: %v", err)
	}
	log.Trace.Printf("--- m dump:\n%s\n\n", string(d))

	d, err = yaml.Marshal(&wConf1)
	if err != nil {
		log.Fatal.Printf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))

	d, err = yaml.Marshal(&wConf2)
	if err != nil {
		log.Fatal.Printf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))


	err = ReadYaml(configPathFile, &theConf)
	if err == nil{
		log.Trace.Println("In GetConf - getting Host from ", theConf.LocalWorkspace + "\\" + hostNameFileName)
		err = ReadYaml(theConf.LocalWorkspace + "\\" + hostNameFileName, &theConf)
		log.Trace.Println("Host is ", theConf.Hostname, "  and err is ", err)
	}
	log.Trace.Println("Done GetConf - err is ", err)
	return &theConf, err
}

func ReadYaml(yamPathFile string, whereStruct interface{}) error {
	yamlFile, err := ioutil.ReadFile( yamPathFile)
	if err != nil {
		log.Error.Printf("yamlFile.Get err %v ", err)
	} else {
		err = yaml.Unmarshal(yamlFile, whereStruct)
		if err != nil {
			log.Error.Printf("Unmarshal: %v", err)
		}
	}
	return err
}

