package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/prasincs/logrotate-sync/processing"
)

type stringSlice []string

// Implement the Value interface
func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Value() []string {
	return *s
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var (
		dirs           = stringSlice{}
		compressGzip   = false
		uploadPath     = ""
		uploadPattern  = "2006/01/02"
		filterLogTypes = stringSlice{}
		//hostname       = ""
	)
	flag.Var(&dirs, "dir", "Dir(s) to read the log files from")
	flag.BoolVar(&compressGzip, "compress-gzip", false, "Compress log files if not already gzipped")
	flag.StringVar(&uploadPath, "uploadpath", "", "s3://<bucket>/<prefix>")
	flag.StringVar(&uploadPattern, "upload.pattern", "2006/01/02", "Directory pattern in Go time format")
	flag.Var(&filterLogTypes, "logtype", "Types of logs to send to S3, by default sends everything")
	flag.Parse()

	logPattern := regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)
	for _, d := range dirs {
		dir := processing.Dir(d)
		files, err := dir.LogFiles(logPattern)
		if err != nil {
			log.Fatalf("Error getting files: %v", files)
		}
		//log.Printf("%#v", files)
		for _, f := range files {
			tm, err := f.ParseTime()
			if err != nil {
				log.Fatalf("Error parsing time: %s", err)
			}
			fmt.Println(tm)
		}
	}
}
