package processing

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

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

// mapDirOpenError maps the provided non-nil error from opening name
// to a possibly better non-nil error. In particular, it turns OS-specific errors
// about opening files in non-directories into os.ErrNotExist. See Issue 18984.
func mapDirOpenError(originalErr error, name string) error {
	if os.IsNotExist(originalErr) || os.IsPermission(originalErr) {
		return originalErr
	}

	parts := strings.Split(name, string(filepath.Separator))
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		fi, err := os.Stat(strings.Join(parts[:i+1], string(filepath.Separator)))
		if err != nil {
			return originalErr
		}
		if !fi.IsDir() {
			return os.ErrNotExist
		}
	}
	return originalErr
}

func (d Dir) Open(name string) (File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("logsync: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, mapDirOpenError(err, fullName)
	}
	return f, nil
}

// Returns the LogFiles in the directory that match the pattern
func (d Dir) LogFiles(matchPattern *regexp.Regexp) ([]fileMatch, error) {
	fileInfos, err := ioutil.ReadDir(string(d))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to enumerate files from %s", d)
	}

	files := []fileMatch{}
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if matchPattern.MatchString(name) {
			files = append(files, fileMatch{
				path:    path.Join(string(d), name),
				matches: matchFileNames(name, matchPattern),
			})
		}
	}
	return files, nil
}
