package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileFunc func(filename string) error
type FileNameFilter func(filename string) bool

// ApplyToPathWithFilter executes a 'FileFunc' function for each file in a given 'path'. If 'path' is a regular file itself 'FileFunc' will be applied to it directly. If 'path' is a folder, the function will be applied to each regular file inside the folder. This behaviour can be made recursive by setting 'recursive' to true.
func ApplyToPathWithFilter(path string, recursive bool, function FileFunc, filter FileNameFilter) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return ApplyToFolder(path, recursive, function, filter)
	}

	return ApplyToFile(path, function, filter)
}

func ApplyToFile(path string, function FileFunc, filter FileNameFilter) error {

	if filter == nil || filter(path) {
		return function(path)
	}
	return nil
}

func ApplyToFolder(folder string, recursive bool, function FileFunc, filter FileNameFilter) error {
	filenames, err := getFolderFilenames(folder, recursive)
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		// TODO aggregate errors?
		err = ApplyToFile(filename, function, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFolderFilenames(folder string, recursive bool) ([]string, error) {
	var filenames []string = make([]string, 0)
	if recursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filenames = append(filenames, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			filenames = append(filenames, filepath.Join(folder, file.Name()))
		}
	}
	return filenames, nil
}

func IsYamlFilter(filename string) bool {
	extensions := []string{".yaml", "yml"}
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
