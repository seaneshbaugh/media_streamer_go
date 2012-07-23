package fileutils

import "os"

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func IsFile(path string) (bool, error) {
	fileExists, err := FileExists(path)

	if fileExists == false {
		return false, err
	}

	fileInfo, err := os.Stat(path)

	if fileInfo.Mode()&os.ModeDir == 0 {
		return true, err
	}

	return false, err
}

func IsDirectory(path string) (bool, error) {
	fileExists, err := FileExists(path)

	if fileExists == false {
		return false, err
	}

	fileInfo, err := os.Stat(path)

	if fileInfo.Mode()&os.ModeDir == 1 {
		return true, err
	}

	return false, err
}

func IsSocket(path string) (bool, error) {
	fileExists, err := FileExists(path)

	if fileExists == false {
		return false, err
	}

	fileInfo, err := os.Stat(path)

	if fileInfo.Mode()&os.ModeSocket == 1 {
		return true, err
	}

	return false, err
}

func IsDevice(path string) (bool, error) {
	fileExists, err := FileExists(path)

	if fileExists == false {
		return false, err
	}

	fileInfo, err := os.Stat(path)

	if fileInfo.Mode()&os.ModeDevice == 1 {
		return true, err
	}

	return false, err
}