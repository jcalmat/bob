package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/daemon/graphdriver/copy"
)

// RenameFile replaces the regexp result with {replace} in {s}
func RenameFile(path, expression, replace string) (string, error) {
	re := regexp.MustCompile(expression)
	if re.MatchString(path) {
		to := re.ReplaceAllString(path, replace)
		err := os.Rename(path, to)
		return to, err
	}
	return path, nil
}

// Move moves all the files from {src} path to {dest} and skip the files and
// folders in {skip}
func Move(src, dest string, skip []string) error {
	skippedFiles := make(map[string]struct{})

	for _, v := range skip {
		// populate the skippedFiles map
		skippedFiles[v] = struct{}{}
	}

	absPath, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if _, ok := skippedFiles[f.Name()]; ok {
			continue
		}
		err = os.Rename(filepath.Join(src, f.Name()), filepath.Join(absPath, f.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// Copy copies a file or folder from {src} to {dest}
func Copy(src, dest string) error {
	return copy.DirCopy(src, dest, copy.Content, true)
}
