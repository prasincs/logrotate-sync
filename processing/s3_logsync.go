package processing

import (
	"regexp"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var defaultLogPattern = regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)

// S3Logsync is used for reading from the directory if it matches the log pattern
// if no matchPattern isn't provided, defaultLogPattern is using
// defaultLogPattern = regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)
// it will save using the upload pattern which uses Go time Format for time/date
type S3Logsync struct {
	// Directory to scan for
	Directory Dir
	// Compress the files if not compressed already
	CompressGzip bool
	// S3Bucket to upload to
	S3Bucket string
	// S3Prefix to save the files in
	S3Prefix string
	// S3 Client to use -- using the S3API to be able to mock it
	S3Client s3iface.S3API
	// MatchPattern to use
	MatchPattern *regexp.Regexp
	// Upload Pattern to use
	UploadPattern string
	// Log Types to filter on, if you don't want to upload everything
	FilterLogTypes []string
	// Hostname to prefix by
	Hostname string
}

func (s *S3Logsync) Process(deleteLocalAfterUpload bool) error {
	// matchPattern := defaultLogPattern
	// if s.MatchPattern != nil {
	// 	matchPattern = s.MatchPattern
	// }
	// filesMatch, err := s.Directory.LogFiles(matchPattern)
	// if err != nil {
	// 	return errors.Wrapf(err, "Failed to list log files")
	// }
	return nil
}
