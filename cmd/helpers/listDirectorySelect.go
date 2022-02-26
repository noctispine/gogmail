package helper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func MakeSliceFromDir(dirPath string) ([]string, error) {
	files, err := filePathWalkDir(dirPath)
	if err != nil {
		return nil, err
	}

	files = append(files, "../")
	err = stringSliceReplace(files, len(files)-1, 1)
	if err != nil {
		return nil, err
	}

	files = append(files, "Quit")
	err = stringSliceReplace(files, len(files)-1, 0)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fmt.Println(file)
	}

	return files, nil
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func stringSliceReplace(slice []string, pos1, pos2 int) error {
	if (pos1 > len(slice)) || (pos2 > len(slice)) {
		return errors.New("One of the given position exeeceds slice length")
	}
	temp := slice[pos1]
	slice[pos1] = slice[pos2]
	slice[pos2] = temp

	return nil
}
