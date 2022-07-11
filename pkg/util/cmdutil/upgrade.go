package cmdutil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
)

var cfWorkDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func Upgrade(version string) {
	check, newVersion := util.CheckVersion(version)
	if check {
		log.Infof("当前版本为 %s ，发现 %s 新版本，正在下载新版本 (The current version is %s , Found %s new version, downloading new version now)", version, newVersion, version, newVersion)
		fileName := fmt.Sprintf("cf-%s-%s-%s.tar.gz", newVersion, runtime.GOOS, runtime.GOARCH)
		downloadURL := fmt.Sprintf("https://ghproxy.com/github.com/teamssix/cf/releases/download/%s/%s", newVersion, fileName)
		path, _ := os.Executable()
		_, oldFileName := filepath.Split(path)
		oldBakFileName := cfWorkDir + "/" + oldFileName + ".bak"
		downloadFile(downloadURL, fileName)
		err := os.Rename(cfWorkDir+"/"+oldFileName, oldBakFileName)
		util.HandleErr(err)
		unzipFile(fileName)
		err = os.Remove(fileName)
		util.HandleErr(err)
		log.Infof("更新完成，历史版本已被重命名为 %s (The update is complete and the previous version has been renamed to %s)", oldFileName+".bak", oldFileName+".bak")
	} else {
		log.Infof("当前 %s 版本为最新版本，无需升级 (The current %s version is the latest version, no need to upgrade)", version, version)
	}
}

func downloadFile(downloadURL string, fileName string) {
	log.Debugln("下载地址 (download url): " + downloadURL)
	req, err := http.NewRequest("GET", downloadURL, nil)
	util.HandleErr(err)
	resp, err := http.DefaultClient.Do(req)
	util.HandleErr(err)
	defer resp.Body.Close()
	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	bar := progressbar.NewOptions64(resp.ContentLength,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	log.Debugln("下载完成 (Download completed)")
}

func unzipFile(fileName string) {
	log.Debugf("解压 %s 文件 (Unzip the %s file)", fileName, fileName)
	gzipStream, err := os.Open(fileName)
	util.HandleErr(err)
	defer gzipStream.Close()
	uncompressedStream, err := gzip.NewReader(gzipStream)
	util.HandleErr(err)
	defer uncompressedStream.Close()
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		util.HandleErr(err)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		filename := "./cfcache/" + header.Name
		file, err := createFile(filename)
		util.HandleErr(err)
		io.Copy(file, tarReader)
	}
	newCfPath := cfWorkDir + "/cf"
	log.Debugln("解压完成 (Unzip completed)")
	log.Tracef("将 ./cfcache/cf 文件移动到 %s (Move the ./cfcache/cf file to %s)", newCfPath, newCfPath)
	os.Rename("./cfcache/cf", newCfPath)
	log.Traceln("删除 ./cfcache/ 文件夹 (Delete ./cfcache/ folder)")
	err = os.RemoveAll("./cfcache/")
	util.HandleErr(err)
	log.Traceln("为 ./cf 文件赋予可执行权限 (Grant execute permission to ./cf file)")
	f, err := os.Open(newCfPath)
	util.HandleErr(err)
	defer f.Close()
	err = f.Chmod(0755)
	util.HandleErr(err)
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
