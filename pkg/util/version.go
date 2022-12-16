package util

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/global"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetCurrentVersion() string {
	return global.Version
}

func GetUpdateTime() string {
	return global.UpdateTime
}

type latestReleasesStruct struct {
	NewVersion string `json:"tag_name"`
}

func AlertUpdateInfo() {
	oldTimestamp := ReadTimestamp(ReturnVersionTimestampFile())
	if oldTimestamp == 0 {
		CheckVersion(global.Version)
	} else if IsFlushCache(oldTimestamp) {
		check, newVersion, err := CheckVersion(global.Version)
		if check {
			log.Warnf("发现 %s 新版本，可以使用 upgrade 命令进行更新 (Found a new version of %s, use the upgrade command to update)", newVersion, newVersion)
		} else if err == nil {
			log.Debugln("未发现新版本 (No new versions found)")
		}
	} else {
		TimeDifference(oldTimestamp)
	}
}

func CheckVersion(version string) (bool, string, error) {
	WriteTimestamp(ReturnVersionTimestampFile())
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
						return false, newVersion, err
					} else {
						return true, newVersion, err
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
	errutil.HandleErr(err)
	return i
}

func reqErr(err error) (bool, string, error) {
	log.Errorln("获取最新版本失败 (Failed to get the latest version) : ", err)
	return false, "", err
}
