package util

import (
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"time"

	log "github.com/sirupsen/logrus"
)

func WriteTimestamp(TimestampType string) {
	var TimestampCache pubutil.TimestampCache
	log.Tracef("写入 %s 时间戳 (Write %s Timestamp)", TimestampType, TimestampType)
	Timestamp := time.Now().Unix()
	TimestampCache.TimestampType = TimestampType
	TimestampCache.Timestamp = Timestamp
	database.InsertTimestamp(TimestampCache)
}

func ReadTimestamp(TimestampType string) int64 {
	log.Tracef("读取 %s 时间戳 (Reading %s Timestamp)", TimestampType, TimestampType)
	Timestamp := database.SelectTimestampType(TimestampType)
	return Timestamp
}

func ReturnVersionTimestampFile() string {
	cacheType := "version"
	return cacheType
}

func ReturnTimestampType(provider string, TimestampType string) string {
	cacheType := provider + "-" + TimestampType + "-" + pubutil.MaskAK(database.SelectConfigInUse(provider).AccessKeyId)
	return cacheType
}

func IsFlushCache(oldTimestamp int64) bool {
	nowTimestamp := time.Now().Unix()
	if nowTimestamp > oldTimestamp+86400 {
		return true
	}
	return false
}

func TimeDifference(oldTimestamp int64) {
	nowTimestamp := time.Now().Unix()
	log.Tracef("现在的时间戳：%d，缓存的时间戳：%d，相差 %d 秒", nowTimestamp, oldTimestamp, nowTimestamp-oldTimestamp)
}
