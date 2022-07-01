package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"

	log "github.com/sirupsen/logrus"
)

type latestReleasesStruct struct {
	NewVersion string `json:"tag_name"`
}

var (
	version    = "v0.0.2"
	updateTime = "2022.7.1"
)

func GetCurrentVersion() string {
	return version
}

func GetUpdateTime() string {
	return updateTime
}

func AlertUpdateInfo() {
	oldTimeStamp := ReadTimeStamp()
	nowTimeStamp := time.Now().Unix()
	if nowTimeStamp > oldTimeStamp+86400 {
		check, newVersion := CheckVersion(version)
		if check {
			color.LightYellow.Printf("\n发现 %s 新版本，可以使用 upgrade 命令进行更新 (Found a new version of %s, use the upgrade command to update)\n\n", newVersion, newVersion)
		}
	}
}

func CheckVersion(version string) (bool, string) {
	WriteTimeStamp()
	url := "https://api.github.com/repos/teamssix/cf/releases/latest"
	spaceClient := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return reqErr(err)
	} else {
		res, err := spaceClient.Do(req)
		if err != nil {
			return reqErr(err)
		} else {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return reqErr(err)
			} else {
				latestReleases := latestReleasesStruct{}
				err := json.Unmarshal(body, &latestReleases)
				if err != nil {
					return reqErr(err)
				} else {
					newVersion := latestReleases.NewVersion
					versionNumber := caclVersionNumber(version)
					newVersionNumber := caclVersionNumber(newVersion)
					if versionNumber >= newVersionNumber {
						return false, newVersion
					} else {
						return true, newVersion
					}
				}
			}

		}
	}
}

func caclVersionNumber(version string) int {
	version = version[1:]
	versionSplit := strings.Split(version, ".")
	versionNumber := Atoi(versionSplit[0])*10000 + Atoi(versionSplit[1])*100 + Atoi(versionSplit[2])
	return versionNumber
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	HandleErr(err)
	return i
}

func reqErr(err error) (bool, string) {
	log.Debugln("获取最新的 releases 版本失败 (Failed to get the latest releases version) : ", err)
	return false, ""
}

func ReturnCacheTimeStampFile() string {
	cacheDict := ReturnCacheDict() + "/timeStamp.txt"
	return cacheDict
}

func ReturnCacheDict() string {
	home, err := GetCFHomeDir()
	HandleErr(err)
	cacheDict := home + "/cache"
	return cacheDict
}

func GetCFHomeDir() (string, error) {
	home := os.Getenv(CFHomeEnvVar)
	if home != "" {
		return home, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get user home dir")
	}
	return filepath.Join(home, AppDirName), nil
}

func WriteTimeStamp() {
	log.Traceln("写入时间戳文件 (Writing to a timestamp file)")
	if !fileExists(ReturnCacheDict()) {
		err := os.MkdirAll(ReturnCacheDict(), 0700)
		HandleErr(err)
	}
	content := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	err := ioutil.WriteFile(ReturnCacheTimeStampFile(), content, 0644)
	HandleErr(err)
}

func ReadTimeStamp() int64 {
	log.Traceln("读取时间戳文件 (Reading to a timestamp file)")
	if !fileExists(ReturnCacheDict()) {
		CheckVersion(version)
	}
	file, err := os.Open(ReturnCacheTimeStampFile())
	HandleErr(err)
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	HandleErr(err)
	i, err := strconv.ParseInt(string(content), 10, 64)
	HandleErr(err)
	return i
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
