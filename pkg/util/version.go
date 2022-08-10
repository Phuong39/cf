package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type latestReleasesStruct struct {
	NewVersion string `json:"tag_name"`
}

var (
	version    = "v0.3.4"
	updateTime = "2022.8.10"
)

func GetCurrentVersion() string {
	return version
}

func GetUpdateTime() string {
	return updateTime
}

func AlertUpdateInfo() {
	oldTimeStamp := ReadTimeStamp(ReturnVersionTimeStampFile())
	if oldTimeStamp == 0 {
		CheckVersion(version)
	} else if IsFlushCache(oldTimeStamp) {
		check, newVersion := CheckVersion(version)
		if check {
			log.Warnf("发现 %s 新版本，可以使用 upgrade 命令进行更新 (Found a new version of %s, use the upgrade command to update)\n", newVersion, newVersion)
		} else {
			log.Debugln("未发现新版本 (No new versions found)")
		}
	} else {
		TimeDifference(oldTimeStamp)
	}
}

func CheckVersion(version string) (bool, string) {
	WriteTimeStamp(ReturnVersionTimeStampFile())
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
