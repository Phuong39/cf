package cmdutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
	"os"
	"runtime"
)

func Upgrade() {
	Check, latest, m := util.CheckVersion()
	if Check {
		currentOS := runtime.GOOS
		var final *update.Asset
		switch runtime.GOOS {
		case "windows":
			final = latest.FindZip(currentOS, runtime.GOARCH)
		default:
			final = latest.FindTarball(currentOS, runtime.GOARCH)
		}
		if final == nil {
			log.Errorln("未找到该系统的二进制文件，更新失败 (The binary file for this system was not found and the update failed)")
			os.Exit(0)
		}
		log.Infoln("正在下载中…… (Downloading...)")
		tarball, err := final.DownloadProxy(progress.Reader)
		if err != nil {
			util.HandleErr(err)
		}
		currentPath, _ := os.Getwd()
		if err = m.InstallTo(tarball, currentPath); err != nil {
			util.HandleErr(err)
		}
		log.Infoln("升级完成 (Upgrade completed)")
	} else {
		log.Infof("当前版本为最新版本，无需升级 (The current version is the latest version, no need to upgrade)")
	}
}
