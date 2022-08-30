package util

import (
	"errors"
	"github.com/teamssix/cf/pkg/util/env"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func WriteTimeStamp(file string) {
	log.Tracef("写入时间戳文件 %s (Writing to a timestamp file %s)", file, file)
	if !pubutil.FileExists(file) {
		log.Traceln("未找到时间戳文件，正在创建时间戳文件 (Timestamp file not found, being created the timestamp file)")
		err := os.MkdirAll(ReturnCacheDict(), 0700)
		HandleErr(err)
	}
	content := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	err := ioutil.WriteFile(file, content, 0644)
	HandleErr(err)
}

func ReadTimeStamp(timeStampFile string) int64 {
	log.Tracef("读取时间戳文件 %s (Reading to a timestamp file %s)", timeStampFile, timeStampFile)
	if !pubutil.FileExists(timeStampFile) {
		log.Traceln("未找到时间戳文件，正在获取最新数据 (Timestamp file not found, Getting the latest data)")
		return 0
	}
	file, err := os.Open(timeStampFile)
	HandleErr(err)
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	HandleErr(err)
	i, err := strconv.ParseInt(string(content), 10, 64)
	HandleErr(err)
	return i
}

func ReturnVersionTimeStampFile() string {
	cacheDict := ReturnCacheDict() + "/versionTimeStamp.txt"
	return cacheDict
}

func ReturnOSSTimeStampFile() string {
	cacheDict := ReturnCacheDict() + "/ossTimeStamp.txt"
	return cacheDict
}

func ReturnECSTimeStampFile() string {
	cacheDict := ReturnCacheDict() + "/ecsTimeStamp.txt"
	return cacheDict
}

func ReturnRDSTimeStampFile() string {
	cacheDict := ReturnCacheDict() + "/rdsTimeStamp.txt"
	return cacheDict
}

func ReturnCacheDict() string {
	home, err := GetCFHomeDir()
	HandleErr(err)
	cacheDict := home + "/cache"
	return cacheDict
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

func IsFlushCache(oldTimeStamp int64) bool {
	nowTimeStamp := time.Now().Unix()
	if nowTimeStamp > oldTimeStamp+86400 {
		return true
	}
	return false
}

func TimeDifference(oldTimeStamp int64) {
	nowTimeStamp := time.Now().Unix()
	log.Tracef("现在的时间戳：%d，缓存的时间戳：%d，相差 %d 秒", nowTimeStamp, oldTimeStamp, nowTimeStamp-oldTimeStamp)
}
