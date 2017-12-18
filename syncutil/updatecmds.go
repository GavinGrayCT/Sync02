package syncutil

import (
	"github.com/GavinGrayCT/SyncFileUtility/Sync02/log"
	"os"
	"io/ioutil"
)

func AddCommand(theDir *Dir, fs *FileStates) (err error) {
	log.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	log.Trace.Printf("In Dir %s, file %s --- ", theDir.RelPath, fs.clf.Name)
	if fs.clf.FModDate == fs.llf.FModDate &&
		fs.clf.FModDate == fs.llrf.FModDate &&
		fs.clf.FModDate == fs.crf.FModDate &&
		fs.clf.FModDate == fs.lorf.FModDate {
		log.Trace.Println("No Action cmd")
	} else {
	log.Trace.Println("Some Action cmd")
		log.Trace.Println("fs.clf.FModDate: ", fs.clf.FModDate)
		log.Trace.Println("fs.llf.FModDate: ", fs.llf.FModDate)
		log.Trace.Println("fs.crf.FModDate: ", fs.crf.FModDate)
		log.Trace.Println("fs.lorf.FModDate: ", fs.lorf.FModDate)
	}
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	return err
}
