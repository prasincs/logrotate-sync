package processing

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Logsync Interfface to use for uploading to different services
type Logsync interface {
	Process(bool) error
}

// FileMatch is used to keep track of path, matched strings, and rotated Time
type FileMatch struct {
	path        string
	matches     map[string]string
	rotatedTime time.Time
}

// NewFileMatch takes the path and the matched map, and parses the timestamp from the matches
func NewFileMatch(path string, matches map[string]string) (*FileMatch, error) {
	f := &FileMatch{
		path:    path,
		matches: matches,
	}

	tm, err := f.ParseTime()
	if err != nil {
		return f, errors.Wrapf(err, "Failed to parse rotated time for %s", path)
	}
	f.rotatedTime = tm
	return f, nil
}

// Parse the Key from matches
func (f *FileMatch) parseInt(key string) (int, error) {
	extras, ok := f.matches[key]
	if !ok {
		return -1, fmt.Errorf("No field(%s) matched, unexpected filename %s", key, f.path)
	}
	extraInt, err := strconv.Atoi(extras)
	if err != nil {
		return -1, errors.Wrapf(err, "Cannot parse the field(%s) as int", key)
	}
	return extraInt, nil
}

// ParseTime relies on the matches map containing
// year, month, day and extras field
// if extras can be parsed as Unix timestamp, that's picked
// otherwise the last field is checked if it's an 0 <= hour <= 23
// the time is constructed if
func (f *FileMatch) ParseTime() (time.Time, error) {
	extraInt, err := f.parseInt("extra")
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "Failed to parse for extra")
	}
	day, err := f.parseInt("day")
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "Failed to parse for day")
	}
	month, err := f.parseInt("month")
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "Failed to parse for month")
	}
	year, err := f.parseInt("year")
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "Failed to parse for year")
	}
	if extraInt >= 0 && extraInt <= 23 {
		hour := extraInt
		return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC), nil
	}
	// assume it's unix time
	tm := time.Unix(int64(extraInt), 0)
	if tm.Year() != year || tm.Month() != time.Month(month) || tm.Day() != day {
		return time.Time{}, errors.Wrapf(err, "Unix time doesn't match the Y/M/D, probably invalid file pattern")
	}
	return tm, nil

}

func matchFileNames(name string, matchPattern *regexp.Regexp) map[string]string {
	matches := matchPattern.FindAllStringSubmatch(name, -1)
	if len(matches) == 0 {
		return nil
	}
	subExps := matchPattern.SubexpNames()
	elements := matches[0]
	md := map[string]string{}
	for i := 1; i < len(elements); i++ {
		md[subExps[i]] = elements[i]
	}
	return md
}
