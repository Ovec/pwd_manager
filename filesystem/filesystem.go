package filesystem

import (
	"fmt"
	"os"
)

func homeDir() (string, error) {
	return os.UserHomeDir()
}

func FileExists(f string) bool {
	fmt.Println(f)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func DataDir(d string) (string, error) {
	home, err := homeDir()

	if err != nil {
		return "", err
	}

	return home + d, nil
}

func CreateFiles(dataDir string, files []string) error {
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, f := range files {
		if FileExists(dataDir + "/" + f) {
			file, err := os.Create(dataDir + "/" + f)
			if err != nil {
				return err
			}
			defer file.Close()
		}
	}

	return nil
}
