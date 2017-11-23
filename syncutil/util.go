package syncutil

import (
	"time"
	"os"
	"strings"
	"path/filepath"
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
)

type afile struct {
	Name     string
	Size     int64
	ModeDate time.Time
}

type Dir struct {
	Path     string
	Files    map[string]afile
	ModeDate time.Time
}

var defaultFormat = "2006-01-02 15:04:05.000"

type DirMap map[string]*Dir

type visitor struct {
	aDirMap DirMap
	excludes []string
}

func CheckExclude(pathfile string, excludes []string) (bool, string, error) {
	var matched bool
	var err error
	for i := range excludes {
		log.Trace.Println("Matching ", pathfile, "and", excludes[i])
		matched, err = regexp.Match(excludes[i], []byte(pathfile))
		if matched {
			return true, excludes[i], err
		}
	}
	return false, "", err
}

func (aVisitor visitor) visit(path string, f os.FileInfo, err error) error {
	log.Trace.Printf("Visiting: %s  name: %s  size:%d modified:%s isDir:%t\n", path, f.Name(), f.Size(), f.ModTime().Format(defaultFormat), f.IsDir())
	steps := strings.Split(path, string(filepath.Separator))
	log.Trace.Println("Steps are: ", steps)

	// Check if not excluded
	if excluded, _, _ := CheckExclude(path, aVisitor.excludes); !excluded {
		if f.IsDir() {
			if aDir, ok := aVisitor.aDirMap[path]; !ok {
				log.Trace.Printf("Adding %s to dirMap\n", path)
				aDir = new(Dir)
				aDir.Path = path
				aDir.Files = make(map[string]afile)
				aDir.ModeDate = f.ModTime()
				aVisitor.aDirMap[path] = aDir
			} else {
				log.Trace.Printf("Existing dir %s", path)
			}
		} else {
			theDirname := filepath.Dir(path)
			afile := afile{f.Name(), f.Size(), f.ModTime()}
			if _, ok := aVisitor.aDirMap[theDirname].Files[afile.Name]; !ok {
				log.Trace.Printf("Adding %s to dir %s\n", f.Name(), theDirname)
				aVisitor.aDirMap[theDirname].Files[afile.Name] = afile
			} else {
				log.Fatal.Printf("Existing file: %s in %s\n", f.Name(), theDirname)
			}
		}
	} else {
		log.Info.Printf("Excluding %s\n", path)
	}
	return err

}

func StoreDirMap(pathfile string, aDirMap DirMap) error {
	log.Trace.Println("Storing the dirMap in ", pathfile)
	j, err := json.MarshalIndent(aDirMap, "", "    ")
	log.Trace.Printf("Length of j:%d\n", len(j))
	err = ioutil.WriteFile(pathfile, j, 0644)
	if err != nil {
		log.Error.Println("Could not write to ", pathfile)
	}
	log.Trace.Println(string(j), "\n Marshal err:", err)
	return err
}

func CompDirMap(aDirMap DirMap, bDirMap DirMap) bool {
	log.Trace.Println("Comparing 2x dirMap")
	for _, aDir := range aDirMap {
		if bDir, ok := bDirMap[aDir.Path]; ok {
			for _, aFile := range aDir.Files {
				if _, ok := bDir.Files[aFile.Name]; !ok {
					return false
				}
			}
		} else {
			return false
		}

	}
	return true
}

func Retrieve(pathfile *string, aDirMap DirMap) {
	log.Trace.Println("Retrieving the dirMap from ", pathfile)
	j, err := ioutil.ReadFile(*pathfile)
	if err != nil {
		log.Error.Println("Could not read from ", pathfile)
	}
	aDirMap = make(map[string]*Dir)
	if err := json.Unmarshal(j, &aDirMap); err != nil {
		panic(err)
	}
}

func FillDirMap(rootPath *string, excludes []string, aDirMap DirMap) error {
	aVisitor := visitor{aDirMap, excludes}
	err := filepath.Walk(*rootPath, aVisitor.visit)
	return err
}

func InitLocalWorkSpace(localWorkSpace string) {
	if _, err := os.Stat(localWorkSpace); os.IsNotExist(err) {
		fmt.Printf("Creating ", localWorkSpace)
	}
}
