package pubutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/global"
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
	home := os.Getenv(global.CFHomeEnvVar)
	if home != "" {
		return home, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get user home dir")
	}
	return filepath.Join(home, global.AppDirName), nil
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

func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2f B", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f GB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f TB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
