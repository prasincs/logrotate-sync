package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/prasincs/logrotate-sync/processing"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		dirs           = []string{}
		compressGzip   = false
		uploadPath     = ""
		statePath      = ""
		uploadPattern  = "2006/01/02"
		filterLogTypes = []string{}
		hostname       = ""
		awsRegion      = ""
		err            error
	)
	flag.StringSliceVar(&dirs, "dir", []string{"."}, "Dir(s) to read the log files from")
	flag.BoolVar(&compressGzip, "compress-gzip", false, "Compress log files if not already gzipped")
	flag.StringVar(&uploadPath, "upload-path", "", "s3://<bucket>/<prefix>")
	flag.StringVar(&statePath, "state-path", "/var/run/logrotatesync.json", "Where to store the json of state of files that are available")
	flag.StringVar(&uploadPattern, "upload-time-pattern", "2006/01/02/15", "Directory pattern in Go time format")
	flag.StringSliceVar(&filterLogTypes, "logtype", []string{}, "Types of logs to send to S3, by default sends everything")
	flag.StringVar(&hostname, "hostname", "", "Hostname to use, takes the default from os.Hostname() by default")
	flag.StringVar(&awsRegion, "aws-region", "", "Pick the AWS Region for the s3 buckets")
	flag.Parse()

	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			log.Fatalf("Unable to read the hostname: %s", err)
		}
	}

	logPattern := regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)
	if !strings.HasPrefix("s3://") {
		panic(errors.New("Only S3 upload paths are supported right now"))
	}

	for _, d := range dirs {
		dir := processing.Dir(d)
		files, err := dir.LogFiles(logPattern)
		if err != nil {
			log.Fatalf("Error getting files: %v", files)
		}
		for _, f := range files {
			tm, err := f.ParseTime()
			if err != nil {
				log.Fatalf("Error parsing time: %s", err)
			}

		}
	}
}
