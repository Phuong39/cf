package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/tj/go-update"
	githubUpdateStore "github.com/tj/go-update/stores/github"
	"runtime"
)

var (
	CFVersion  = "0.3.4"
	version    = "v" + CFVersion
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
		_, _, _ = CheckVersion()
	} else if IsFlushCache(oldTimeStamp) {
		Check, latest, _ := CheckVersion()
		if Check {
			log.Warnln("发现 %s 新版本，可以使用 upgrade 命令进行更新 (Found a new version of %s, use the upgrade command to update)", latest.Version, latest.Version)
		} else {
			log.Infoln("未发现新版本 (No new versions found)")
		}
	} else {
		TimeDifference(oldTimeStamp)
	}
}

func CheckVersion() (bool, *update.Release, *update.Manager) {
	WriteTimeStamp(ReturnVersionTimeStampFile())
	var command string
	switch runtime.GOOS {
	case "windows":
		command = "cf.exe"
	default:
		command = "cf"
	}
	m := &update.Manager{
		Command: command,
		Store: &githubUpdateStore.Store{
			Owner:   "teamssix",
			Repo:    "cf",
			Version: CFVersion,
		},
	}
	releases, err := m.LatestReleases()
	if err != nil {
		HandleErr(err)
	}
	if len(releases) == 0 {
		return false, nil, nil
	} else {
		latest := releases[0]
		return true, latest, m
	}
}
