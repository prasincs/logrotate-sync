package processing

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/pkg/errors"
)

// this log pattern handles the common logrotate datesuffix format and YYYY-MM-DD-HH
var defaultLogPattern = regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)

type Config struct {
	// Directory to scan for
	Directory Dir
	// Compress the files if not compressed already
	CompressGzip bool
	// MatchPattern to use
	MatchPattern *regexp.Regexp
	// Upload Pattern to use
	UploadDatePattern string
	// Log Types to filter on, if you don't want to upload everything
	FilterLogTypes []string
	// Hostname to prefix the files by
	Hostname string
	// AWS Region to use
	AWSRegion string
	//S3Endpoint to use if given
	S3Endpoint string
	// MaxFilesToKeep is the number of files to keep if the Process
	// is called with deleteLocal option
	// this will strictly be last N files, whether its done hourly or not
	MaxFilesToKeep int
}

func splitS3URL(rs string) (bucket, key string) {
	if strings.HasPrefix(rs, "s3://") {
		rs := (rs)[5:] // strip off the "s3://" prefix

		// split the ROUTE_SOURCE into bucket/key
		idx := strings.IndexRune(rs, '/')
		if idx <= 0 || idx+1 >= len(rs) {
			log.Printf("Error parsing iptun route source %q: no / separator for bucket and key\n", rs)
			return
		}
		bucket = rs[:idx]
		key = rs[idx+1:]
	} // else return "",""
	return bucket, key
}

// S3Logsync is used for reading from the directory if it matches the log pattern
// if no matchPattern isn't provided, defaultLogPattern is using
// defaultLogPattern = regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)
// it will save using the upload pattern which uses Go time Format for time/date
type S3Logsync struct {
	// config to put for uploading the files, etc
	config Config
	// S3Bucket to upload to
	bucket string
	// S3Prefix to save the files in
	prefix string
	// S3 Client to use -- using the S3API to be able to mock it
	s3Client s3iface.S3API
}

func NewS3Logsync(config Config, path string) (*S3Logsync, error) {
	bucket, prefix := splitS3URL(path)
	if bucket == "" || prefix == "" {
		return nil, fmt.Errorf("Invalid path: %s provided", path)
	}
	awsConfig := aws.NewConfig()
	if config.AWSRegion != "" {
		awsConfig.Region = aws.String(config.AWSRegion)
	}
	if config.S3Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.S3Endpoint)
	}
	sess := session.New(awsConfig)
	client := s3.New(sess, awsConfig)
	return &S3Logsync{
		config:   config,
		bucket:   bucket,
		prefix:   prefix,
		s3Client: client,
	}, nil
}

// this should take a file with
func (s *S3Logsync) remoteFileKey(fm *FileMatch) (string, error) {
	base := path.Base(fm.Path())
	tm, err := fm.ParseTime()
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse the file, cannot generate remote file key")
	}
	newpath := tm.Format(s.config.UploadDatePattern) + "/" + s.config.Hostname + "-" + base
	return s.prefix + "/" + newpath, nil
}

// Process finds the rotated files, and handles the uploads,
// if the the deleteLocalAfterUpload is true, it will delete the local files
func (s *S3Logsync) Process(deleteLocalAfterUpload bool) error {
	matchPattern := defaultLogPattern
	if s.config.MatchPattern != nil {
		matchPattern = s.config.MatchPattern
	}
	filesMatch, err := s.config.Directory.LogFiles(matchPattern)
	if err != nil {
		return errors.Wrapf(err, "Failed to list log files")
	}
	for _, fileMatch := range filesMatch {
		fmt.Printf("%s\n", fileMatch)
	}
	return nil
}
