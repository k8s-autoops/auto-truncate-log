package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	envDry, _  = strconv.ParseBool(os.Getenv("CFG_DRY"))
	envOnce, _ = strconv.ParseBool(os.Getenv("CFG_ONCE"))

	now = time.Now()

	todayMarks = []string{
		fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day()),
		fmt.Sprintf("%04d.%02d.%02d", now.Year(), now.Month(), now.Day()),
		fmt.Sprintf("%04d_%02d_%02d", now.Year(), now.Month(), now.Day()),
	}

	regexpLogFile           = regexp.MustCompile(`(?i)\.log$`)
	regexpHistoricalLogFile = []*regexp.Regexp{
		regexp.MustCompile(`(?i)ROT.+\.log.*$`),
		regexp.MustCompile(`(?i)\.log.*\.gz.*$`),
		regexp.MustCompile(`(?i)\.log[_.-]\d+$`),
		regexp.MustCompile(`(?i)\.log[_.-]\d{4}[_.-]\d{2}[_.-]\d{2}.*$`),
		regexp.MustCompile(`(?i)[_.-]\d+\.log$`),
		regexp.MustCompile(`(?i)\d{4}[_.-]\d{2}[_.-]\d{2}.*\.log$`),
	}
)

type FileType int

const (
	FileTypeNone FileType = iota
	FileTypeActiveLog
	FileTypeHistoryLog
)

func determineFileType(name string) FileType {
	for _, p := range regexpHistoricalLogFile {
		if p.MatchString(name) {
			for _, todayMark := range todayMarks {
				if strings.Contains(name, todayMark) {
					return FileTypeActiveLog
				}
			}
			return FileTypeHistoryLog
		}
	}
	if regexpLogFile.MatchString(name) {
		return FileTypeActiveLog
	}
	return FileTypeNone
}

func main() {
	if envOnce {
		execute()
	}

	c := cron.New()
	if _, err := c.AddFunc(os.Getenv("CFG_CRON"), execute); err != nil {
		log.Println("failed to initialize cron:", err.Error())
		os.Exit(1)
	}
	c.Start()
	defer c.Stop()

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)
	sig := <-chSig
	log.Println("signal caught:", sig.String())
}

func execute() {
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
		// determine file type
		fileType := determineFileType(info.Name())
		switch fileType {
		case FileTypeActiveLog:
			if envDry {
				log.Printf("%s: will truncate", path)
				return nil
			}
			if err := os.Truncate(path, 0); err != nil {
				log.Printf("%s: %s", path, err.Error())
			} else {
				log.Printf("%s: truncated", path)
			}
		case FileTypeHistoryLog:
			if envDry {
				log.Printf("%s: will delete", path)
				return nil
			}
			if err := os.Remove(path); err != nil {
				log.Printf("%s: %s", path, err.Error())
			} else {
				log.Printf("%s: deleted", path)
			}
		default:
			log.Printf("%s: ignored", path)
		}
		return nil
	})
	if err != nil {
		log.Println("failed to iterate files:", err.Error())
	}
}
