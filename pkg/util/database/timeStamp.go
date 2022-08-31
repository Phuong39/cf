package database

import (
	"github.com/teamssix/cf/pkg/util/pubutil"
)

func InsertTimestamp(TimestampCache pubutil.TimestampCache) {
	var TimestampCacheList []pubutil.TimestampCache
	TimestampType := TimestampCache.TimestampType
	if SelectTimestampType(TimestampType) != 0 {
		CacheDb.Where("timestamp_type = ?", TimestampType).Delete(&TimestampCacheList)
	}
	CacheDb.Create(&TimestampCache)
}

func SelectTimestampType(TimestampType string) int64 {
	var (
		TimestampCache     pubutil.TimestampCache
		TimestampCacheList []pubutil.TimestampCache
	)
	CacheDb.Where("timestamp_type = ?", TimestampType).Find(&TimestampCacheList)
	if len(TimestampCacheList) == 0 {
		return TimestampCache.Timestamp
	} else {
		return TimestampCacheList[0].Timestamp
	}
}
