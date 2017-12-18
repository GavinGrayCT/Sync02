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

type Afile struct {
	Name     string
	Size     int64
	FModDate time.Time
	Present  bool
}

type Dir struct {
	Path     string
	RelPath string
	Files    map[string]Afile
	DModDate time.Time
	Present  bool
}

var defaultFormat = "2006-01-02 15:04:05.000"

type DirMap map[string]*Dir

type RootedDirMap struct {
	rootDir  string
	excludes []string
	aDirMap  DirMap
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

func (aRootedDirMap RootedDirMap) visit(path string, f os.FileInfo, err error) error {
	path = filepath.ToSlash(path)
	log.Trace.Printf("Visiting: %s  name: %s  size:%d modified:%s isDir:%t\n", path, f.Name(), f.Size(), f.ModTime().Format(defaultFormat), f.IsDir())
	steps := strings.Split(path, string("/"))
	log.Trace.Println("Steps are: ", steps)
	relpath := path[len(aRootedDirMap.rootDir):]
	if len(relpath) == 0 {
		relpath = "/"
	}

	// Check if not excluded
	if excluded, _, _ := CheckExclude(path, aRootedDirMap.excludes); !excluded {
		if f.IsDir() {
			if aDir, ok := aRootedDirMap.aDirMap[path]; !ok {
				log.Trace.Printf("Adding dir %s + rel dir %s to dirMap\n", path, relpath)
				aDir = new(Dir)
				aDir.Path = path
				aDir.RelPath = relpath
				aDir.Files = make(map[string]Afile)
				aDir.DModDate = f.ModTime()
				aDir.Present = true
				aRootedDirMap.aDirMap[relpath] = aDir
			} else {
				log.Trace.Printf("Existing dir %s", relpath)
			}
		} else {
			theDirname := filepath.ToSlash(filepath.Dir(relpath))
			afile := Afile{f.Name(), f.Size(), f.ModTime(), true}
			log.Trace.Printf("Adding file %s to dir %s\n", afile, theDirname)
			if _, ok := aRootedDirMap.aDirMap[theDirname].Files[afile.Name]; !ok {
				log.Trace.Printf("Adding %s (size %d,  ModTime %s) to dir %s\n", f.Name(), f.Size(), f.ModTime(), theDirname)
				aRootedDirMap.aDirMap[theDirname].Files[afile.Name] = afile
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
	if checkContainsDirMap(aDirMap, bDirMap) {
		log.Trace.Printf("A contains B")
		if checkContainsDirMap(bDirMap, aDirMap) {
			log.Trace.Printf("and B contains A - so equal")
			return true
		}
	}
	return false
}

func checkContainsDirMap(aDirMap DirMap, bDirMap DirMap) bool {
	log.Trace.Println("Checking 'Contains Criteria'")
	for _, aDir := range aDirMap {
		if bDir, ok := bDirMap[aDir.Path]; ok {
			for _, aFile := range aDir.Files {
				bFile, ok := bDir.Files[aFile.Name]
				if !ok {
					log.Trace.Printf("Within %s, A file '%s' exists and B file '%s' does not exist\n", aDir.Path, aFile.Name, aFile.Name)
					return false
				}
				if aFile.Size != bFile.Size {
					log.Trace.Printf("Different sizes. ", aFile.Name, ":", aFile.Size, "  ", bFile.Name, ":", bFile.Size)
					return false
				}
				if !aFile.FModDate.Equal(bFile.FModDate) {
					log.Trace.Printf("Different ModDates. ", aFile.Name, ":", aFile.FModDate, "  ", bFile.Name, ":", bFile.FModDate)
					return false
				}
			}
		} else {
			log.Trace.Printf("A dir '%s' exists, B dir '%s' does not exist\n", aDir.Path, aDir.Path)
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

func FillDirMap(aRootedDirMap RootedDirMap) error {
	err := filepath.Walk(aRootedDirMap.rootDir, aRootedDirMap.visit)
	return err
}

func InitLocalWorkSpace(localWorkSpace string) {
	if _, err := os.Stat(localWorkSpace); os.IsNotExist(err) {
		fmt.Printf("Creating ", localWorkSpace)
	}
}

func WriteJson(fn string, i interface{}) (err error) {
	j, err := json.MarshalIndent(i, "", "   ")
	err = ioutil.WriteFile(fn, j, 0644)
	log.Trace.Println("Written json. Error is: ", err)
	return err
}

func ReadJson(pf string, s interface{}) error {
	j, err := ioutil.ReadFile(pf)
	if err == nil {
		err = json.Unmarshal(j, s)
		if err != nil {
			log.Error.Printf("Error json unmarshalling %s. Error: %v\n", pf, err)
		}
	} else {
		log.Error.Printf("Error reading file. Error: %v\n", err)

	}
	return err
}
