package pubutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	UserId           string
	UserName         string
	Password         string
	LoginUrl         string
	CreateTime       string
}

func GetUserDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func GetConfigFilePath() string {
	home, _ := GetCFHomeDir()
	CreateFolder(home)
	configFilePath := filepath.Join(home, "cache.db")
	return configFilePath
}

func GetCFHomeDir() (string, error) {
	return filepath.Join(GetUserDir(), global.AppDirName), nil
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

func ReadFile(filePath string) (bool, string) {
	if FileExists(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		return true, string(content)
	} else {
		return false, ""
	}
}

func StringClean(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "\n", "", -1)
	return str
}

func MaskAK(ak string) string {
	if len(ak) > 7 {
		prefix := ak[:2]
		suffix := ak[len(ak)-6:]
		return prefix + strings.Repeat("*", 18) + suffix
	} else {
		return ak
	}
}

func IN(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}
