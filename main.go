package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	envDry, _ = strconv.ParseBool(os.Getenv("DRY"))
)

var (
	patternLogFile        = regexp.MustCompile(`(?i)\.log$`)
	patternHistoryLogFile = []*regexp.Regexp{
		regexp.MustCompile(`(?i)ROT.+\.log.*$`),
		regexp.MustCompile(`(?i)\.log.*\.gz.*$`),
		regexp.MustCompile(`(?i)\.log[_.-]\d+$`),
		regexp.MustCompile(`(?i)\.log[_.-]\d{4}[_.-]\d{2}[_.-]\d{2}.*$`),
		regexp.MustCompile(`(?i)[_.-]\d+\.log$`),
		regexp.MustCompile(`(?i)\d{4}[_.-]\d{2}[_.-]\d{2}.*\.log$`),
	}
)

func isActiveLogFile(path string) bool {
	for _, p := range patternHistoryLogFile {
		if p.MatchString(path) {
			return false
		}
	}
	return patternLogFile.MatchString(path)
}

func isHistoryLogFile(path string) bool {
	for _, p := range patternHistoryLogFile {
		if p.MatchString(path) {
			return true
		}
	}
	return false
}

func main() {
	err := filepath.Walk("/mnt", func(path string, info os.FileInfo, err error) error {
		// ignore error
		if err != nil {
			log.Printf("%s: %s", path, err.Error())
			return nil
		}
		// skip non-regular file
		if info.Mode()&os.ModeType != 0 {
			return nil
		}
		if isActiveLogFile(info.Name()) {
			// truncate active log file
			if envDry {
				log.Printf("%s: will truncate", path)
				return nil
			}
			if err := os.Truncate(path, 0); err != nil {
				log.Printf("%s: %s", path, err.Error())
			} else {
				log.Printf("%s: truncated", path)
			}
			return nil
		} else if isHistoryLogFile(info.Name()) {
			// delete history log file
			if envDry {
				log.Printf("%s: will delete", path)
				return nil
			}
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
