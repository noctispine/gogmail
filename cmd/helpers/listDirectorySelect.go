package helper

import (
	"errors"
	"io/ioutil"
)

func MakeSliceFromDir(dirPath string) ([]string, error) {
	var fileNames []string
	files, err := ioutil.ReadDir(dirPath)

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	fileNames = append(fileNames, "../")
	err = stringSliceReplace(fileNames, len(fileNames)-1, 1)
	if err != nil {
		return nil, err
	}

	fileNames = append(fileNames, "Quit")
	err = stringSliceReplace(fileNames, len(fileNames)-1, 0)
	if err != nil {
		return nil, err
	}

	return fileNames, nil
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
