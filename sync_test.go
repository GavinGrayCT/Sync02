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

var theFileTestSet []FileTestSet

func init() {
	var err error
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	err = testConf.GetTheConf("TestConf.json")
	su.Check(err, "Reading Conf")
	su.ReadJson("TestFileData.json", &theFileTestSet)
	su.Check(err, "Reading Testdata")
}

func TestWriteTestFiles(t *testing.T) {
	log.Trace.Println("In WriteTestFiles")
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

