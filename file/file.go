package file

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func WriteOutputFile(fileName string, records [][]string) (err error) {

	exist, _ := Exist(fileName)
	if exist {
		Delete(fileName)
	}

	directory := filepath.Dir(fileName)
	if _, err = os.Stat(directory); os.IsNotExist(err) {
		if err = Mkdir(directory); err != nil {
			return
		}
	}
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	w := csv.NewWriter(file)

	for _, record := range records {
		if err = w.Write(record); err != nil {
			return

		}
	}
	w.Flush()
	return nil
}

// Exist checks if folder or file exist
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
	}
	return true, err
}

// Mkdir make directory if directory does not exist
func Mkdir(path string) error {
	if exist, _ := Exist(path); !exist {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// Delete delete file
func Delete(filename string) error {
	return os.Remove(filename)
}
