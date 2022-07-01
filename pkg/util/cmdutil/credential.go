package cmdutil

import (
	"errors"
	"os"
	"path/filepath"

	"githubu.com/teamssix/cf/pkg/util"
)

func GetCFHomeDir() (string, error) {
	home := os.Getenv(util.CFHomeEnvVar)
	if home != "" {
		return home, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get user home dir")
	}
	return filepath.Join(home, util.AppDirName), nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
