package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetWorkingDirectory returns the current path
func GetWorkingDirectory() string {
	dir, _ := os.Getwd()
	return dir
}

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

// Move moves all the files from {from} path to {to} and skip the files and
// folders in {skip}
func Move(from, to string, skip []string) error {
	skippedFiles := make(map[string]struct{})

	for _, v := range skip {
		// populate the skippedFiles map
		skippedFiles[v] = struct{}{}
	}

	absPath, err := filepath.Abs(to)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(from)
	if err != nil {
		return err
	}

	for _, f := range files {
		if _, ok := skippedFiles[f.Name()]; ok {
			continue
		}
		err = os.Rename(filepath.Join(from, f.Name()), filepath.Join(absPath, f.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// Copy copies a file or folder from {from} to {to}
func Copy(from, to string) error {

	err := filepath.Walk(from, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		relPath := strings.Replace(path, from, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(to, relPath), 0755)
		}
		data, err := ioutil.ReadFile(filepath.Join(from, relPath))
		if err != nil {
			return err
		}
		return ioutil.WriteFile(filepath.Join(to, relPath), data, 0600)
	})
	return err
}
