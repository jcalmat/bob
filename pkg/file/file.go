package file

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func Copy(from, to string) error {
	//TODO: Replace by actual golang code
	cmd := exec.Command("cp", "-R", from, to)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
