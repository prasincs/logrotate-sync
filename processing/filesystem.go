package processing

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
)

// A FileSystem implements access to a collection of named files.
// The elements in a file path are separated by slash ('/', U+002F)
// characters, regardless of host operating system convention.
type FileSystem interface {
	Open(name string) (File, error)
}

// A File is returned by a FileSystem's Open method and can be
// served by the FileServer implementation.
//
// The methods should behave the same as those on an *os.File.
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

// A Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// While the FileSystem.Open method takes '/'-separated paths, a Dir's string
// value is a filename on the native file system, not a URL, so it is separated
// by filepath.Separator, which isn't necessarily '/'.
//
// Note that Dir will allow access to files and directories starting with a
// period, which could expose sensitive directories like a .git directory or
// sensitive files like .htpasswd. To exclude files with a leading period,
// remove the files/directories from the server or create a custom FileSystem
// implementation.
//
// An empty Dir is treated as ".".
type Dir string

// LogFiles returns files in the directory ordered in descending order by the rotated time
func (d Dir) LogFiles(matchPattern *regexp.Regexp) ([]FileMatch, error) {
	fileInfos, err := ioutil.ReadDir(string(d))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to enumerate files from %s", d)
	}

	files := []FileMatch{}
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if matchPattern.MatchString(name) {
			fm, err := NewFileMatch(path.Join(string(d), name), matchFileNames(name, matchPattern))
			if err != nil {
				// Usually we don't care
				//log.Printf("Failed to match file %s", name)
				continue
			}
			files = append(files, *fm)
		}
	}
	return files, nil
}
