package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	PatternLogFile           = regexp.MustCompile(`(?i)\.log$`)
	PatternCompressedLogFile = regexp.MustCompile(`(?i)[_.-]?\d*\.log[_.-]?\d*\.gz[_.-]?\d*$`)
)

func main() {
	err := filepath.Walk("/mnt", func(path string, info os.FileInfo, err error) error {
		// ignore error
		if err != nil {
			log.Printf("%s: %s", path, err.Error())
			return nil
		}
		// skip non-regular file
		if info.Mode()&os.ModeType != 0 {
			log.Printf("%s: ignored for file type '%s'", path, info.Mode().String())
			return nil
		}
		// truncate log file
		if PatternLogFile.MatchString(info.Name()) {
			if err := os.Truncate(path, 0); err != nil {
				log.Printf("%s: %s", path, err.Error())
			} else {
				log.Printf("%s: truncated", path)
			}
			return nil
		}
		// delete gzipped log file
		if PatternCompressedLogFile.MatchString(info.Name()) {
			if err := os.Remove(path); err != nil {
				log.Printf("%s: %s", path, err.Error())
			} else {
				log.Printf("%s: deleted", path)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}
}
