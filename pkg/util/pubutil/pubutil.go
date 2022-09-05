package pubutil

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/env"
	"os"
	"path/filepath"
)

type TimestampCache struct {
	TimestampType string
	Timestamp     int64
}

type OSSCache struct {
	AccessKeyId  string
	SN           string
	Name         string
	BucketACL    string
	ObjectNumber string
	ObjectSize   string
	Region       string
	BucketURL    string
}

type ECSCache struct {
	AccessKeyId      string
	SN               string
	InstanceId       string
	InstanceName     string
	OSName           string
	OSType           string
	Status           string
	PrivateIpAddress string
	PublicIpAddress  string
	RegionId         string
}

type RDSCache struct {
	AccessKeyId      string
	SN               string
	DBInstanceId     string
	Engine           string
	EngineVersion    string
	DBInstanceStatus string
	RegionId         string
}

type TakeoverConsoleCache struct {
	Provider         string
	AccessKeyId      string
	PrimaryAccountID string
	UserName         string
	Password         string
	LoginUrl         string
	CreateTime       string
}

func GetConfigFilePath() string {
	home, _ := GetCFHomeDir()
	CreateFolder(home)
	configFilePath := filepath.Join(home, "cache.db")
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
