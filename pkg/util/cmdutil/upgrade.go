package cmdutil

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/errutil"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var cfWorkDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func Upgrade(version string) {
	var (
		downloadURL string
		fileName    string
	)
	check, newVersion, err := util.CheckVersion(version)
	if check {
		log.Infof("当前版本为 %s ，发现 %s 新版本，正在下载新版本 (The current version is %s , Found %s new version, downloading new version now)", version, newVersion, version, newVersion)
		if runtime.GOOS == "windows" {
			fileName = fmt.Sprintf("cf_%s_%s_%s.zip", newVersion, runtime.GOOS, runtime.GOARCH)
		} else {
			fileName = fmt.Sprintf("cf_%s_%s_%s.tar.gz", newVersion, runtime.GOOS, runtime.GOARCH)
		}
		if isMainLand() {
			downloadURL = fmt.Sprintf("https://ghproxy.com/github.com/teamssix/cf/releases/download/%s/%s", newVersion, fileName)
		} else {
			downloadURL = fmt.Sprintf("https://github.com/teamssix/cf/releases/download/%s/%s", newVersion, fileName)
		}
		path, _ := os.Executable()
		_, oldFileName := filepath.Split(path)
		oldBakFileName := cfWorkDir + "/" + oldFileName + ".bak"
		downloadFile(downloadURL, fileName)
		err := os.Rename(cfWorkDir+"/"+oldFileName, oldBakFileName)
		errutil.HandleErr(err)
		unzipFile(fileName)
		err = os.Remove(fileName)
		errutil.HandleErr(err)
		log.Infof("更新完成，历史版本已被重命名为 %s (The update is complete and the previous version has been renamed to %s)", oldFileName+".bak", oldFileName+".bak")
	} else if err == nil {
		log.Infof("当前 %s 版本为最新版本，无需升级 (The current %s version is the latest version, no need to upgrade)", version, version)
	}
}

func downloadFile(downloadURL string, fileName string) {
	log.Debugln("下载地址 (download url): " + downloadURL)
	req, err := http.NewRequest("GET", downloadURL, nil)
	errutil.HandleErr(err)
	resp, err := http.DefaultClient.Do(req)
	errutil.HandleErr(err)
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
	log.Debugf("正在解压 %s 文件 (Unpacking %s file now)", fileName, fileName)
	cacheFolder := ReturnCacheDict()
	if runtime.GOOS == "windows" {
		archive, err := zip.OpenReader(fileName)
		if err != nil {
			errutil.HandleErr(err)
		}
		defer archive.Close()
		for _, f := range archive.File {
			filePath := filepath.Join(cacheFolder, f.Name)
			if !strings.HasPrefix(filePath, filepath.Clean(cacheFolder)+string(os.PathSeparator)) {
				log.Errorln("无效的文件路径 (invalid file path)")
			}
			if f.FileInfo().IsDir() {
				os.MkdirAll(filePath, os.ModePerm)
			}

			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				errutil.HandleErr(err)
			}
			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				errutil.HandleErr(err)
			}
			fileInArchive, err := f.Open()
			if err != nil {
				errutil.HandleErr(err)
			}
			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				errutil.HandleErr(err)
			}
			dstFile.Close()
			fileInArchive.Close()
		}
		oldCfPath := filepath.Join(cacheFolder, "cf.exe")
		newCfPath := filepath.Join(cfWorkDir, "cf.exe")
		log.Tracef("将 %s 文件移动到 %s (Move the %s file to %s)", oldCfPath, newCfPath, oldCfPath, newCfPath)
		moveFile(oldCfPath, newCfPath)
	} else {
		gzipStream, err := os.Open(fileName)
		errutil.HandleErr(err)
		defer gzipStream.Close()
		uncompressedStream, err := gzip.NewReader(gzipStream)
		errutil.HandleErr(err)
		defer uncompressedStream.Close()
		tarReader := tar.NewReader(uncompressedStream)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			errutil.HandleErr(err)
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			subFileName := filepath.Join(cacheFolder, header.Name)
			file, err := createFile(subFileName)
			errutil.HandleErr(err)
			io.Copy(file, tarReader)
		}
		log.Debugln("解压完成 (Unzip completed)")
		oldCfPath := filepath.Join(cacheFolder, "cf")
		newCfPath := filepath.Join(cfWorkDir, "cf")
		log.Tracef("将 %s 文件移动到 %s (Move the %s file to %s)", oldCfPath, newCfPath, oldCfPath, newCfPath)
		moveFile(oldCfPath, newCfPath)
		log.Traceln("为 ./cf 文件赋予可执行权限 (Grant execute permission to ./cf file)")
		f, err := os.Open(newCfPath)
		errutil.HandleErr(err)
		defer f.Close()
		err = f.Chmod(0755)
		errutil.HandleErr(err)
	}
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func isMainLand() bool {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://cip.cc", nil)
	req.Header.Add("User-Agent", "curl/7.64.1")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	address := strings.Split(string(body), "\n")[1]
	log.Tracef("当前地址为 (The current address is): %s", strings.Split(address, ":")[1])
	if isRegions(address) && strings.Contains(address, "中国") {
		return true
	} else {
		return false
	}
}

func isRegions(address string) bool {
	var (
		cityList = [3]string{"香港", "澳门", "台湾"}
		num      = 0
	)
	for _, v := range cityList {
		if strings.Contains(address, v) {
			num = num + 1
		}
	}
	if num == 0 {
		return true
	} else {
		return false
	}
}

func moveFile(oldCfPath string, newCfPath string) {
	oldByte, err := ioutil.ReadFile(oldCfPath)
	errutil.HandleErr(err)
	err = ioutil.WriteFile(newCfPath, oldByte, 0644)
	errutil.HandleErr(err)
	err = os.Remove(oldCfPath)
	errutil.HandleErr(err)
}
