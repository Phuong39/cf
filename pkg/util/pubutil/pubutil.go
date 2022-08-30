package pubutil

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/env"
	"os"
	"path/filepath"
)

func GetConfigFilePath() string {
	home, _ := GetCFHomeDir()
	CreateFolder(home)
	configFilePath := filepath.Join(home, "config.db")
	return configFilePath
}

func GetCFHomeDir() (string, error) {
	home := os.Getenv(env.CFHomeEnvVar)
	if home != "" {
		return home, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get user home dir")
	}
	return filepath.Join(home, env.AppDirName), nil
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

func CreateFolder(folder string) {
	if !FileExists(folder) {
		log.Tracef("创建 %s 目录 (Create %s directory): ", folder, folder)
		_ = os.MkdirAll(folder, 0700)
	}
}
