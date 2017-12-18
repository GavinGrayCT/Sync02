package main

import (
	"testing"
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	su "github.com/GavinGrayCT/SyncFileUtility/Sync02/syncutil"
	"time"
	"io/ioutil"
	"os"
	"path"
)

type File struct {
	Present  bool
	Mtime    time.Time
	Contents string
}

type FileTestSet struct {
	RelPath string
	Name    string
	Lf      File
	Rf      File
	Of      File
}

var testConf su.Conf

var time1 = time.Date(2017, 7, 1, 10, 11, 12, 123456789, time.UTC)

var theFileTestSet = []FileTestSet{
	{"dir1", "fname1.txt", File{true, time1, "Donkey1"}, File{true, time1, "Donkey"}, File{true, time1, "Donkey"}},
	{"dir2", "fname2.txt", File{true, time1, "Donkey2"}, File{true, time1, "Donkey"}, File{true, time1, "Donkey"}},
	}

func init() {
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	testConf.GetTheConf("TestConf.json")
}

func TestWriteTestFiles(t *testing.T) {
	log.Trace.Println("In WriteTestFiles")
	su.WriteJson("TestFileData.json", theFileTestSet)
	var err error

	for _, afileTestSet := range theFileTestSet {
		err = writeFile(testConf.LocalRoot.Include, afileTestSet.RelPath, afileTestSet.Name, afileTestSet.Lf)
		if err != nil {
			log.Trace.Panicln("Error writing file", err)
			t.Error("Writing Test Files - false")
		}
	}

	log.Trace.Println("Done WriteTestFiles")
}
func writeFile(lr string, rp string, name string, f File) (err error) {
	var ap = path.Join(lr, rp)
	log.Trace.Printf("Writing file. Dir: %s  Filename: %s\n", ap, name)
	var pf = ap+"/"+name
	log.Trace.Printf("pf: %s\n", pf)
	err = os.MkdirAll(ap, 0777)
	if err == nil {
		err = ioutil.WriteFile(pf, []byte(f.Contents), 0777)
		if err == nil {
			err = os.Chtimes(pf, f.Mtime, f.Mtime)
		}
	}
	return err
}

