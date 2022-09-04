package database

import (
	"github.com/teamssix/cf/pkg/util/pubutil"
)

func InsertECSCache(ECSCache []pubutil.ECSCache) {
	DeleteECSCache(ECSCache[0].AccessKeyId)
	CacheDb.Create(&ECSCache)
}

func DeleteECSCache(AccessKeyId string) {
	var ECSCache []pubutil.ECSCache
	CacheDb.Where("access_key_id = ?", AccessKeyId).Delete(&ECSCache)
}

func SelectECSCache(provider string) []pubutil.ECSCache {
	var ECSCache []pubutil.ECSCache
	AccessKeyId := SelectConfigInUse(provider).AccessKeyId
	CacheDb.Where("access_key_id = ?", AccessKeyId).Find(&ECSCache)
	return ECSCache
}

func SelectEcsCacheFilter(provider string, region string, specifiedInstanceID string, running bool) []pubutil.ECSCache {
	var ECSCache []pubutil.ECSCache
	AccessKeyId := SelectConfigInUse(provider).AccessKeyId
	switch {
	case region == "all" && specifiedInstanceID == "all" && running == false:
		CacheDb.Where("access_key_id = ? COLLATE NOCASE", AccessKeyId).Find(&ECSCache)
	case region == "all" && specifiedInstanceID == "all" && running == true:
		CacheDb.Where("access_key_id = ? AND status = ? COLLATE NOCASE", AccessKeyId, "Running").Find(&ECSCache)
	case region == "all" && specifiedInstanceID != "all" && running == false:
		CacheDb.Where("access_key_id = ? AND instance_id = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID).Find(&ECSCache)
	case region == "all" && specifiedInstanceID != "all" && running == true:
		CacheDb.Where("access_key_id = ? AND instance_id = ? COLLATE NOCASE AND status = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID, "Running").Find(&ECSCache)
	case region != "all" && specifiedInstanceID == "all" && running == false:
		CacheDb.Where("access_key_id = ? AND region_id = ? COLLATE NOCASE", AccessKeyId, region).Find(&ECSCache)
	case region != "all" && specifiedInstanceID == "all" && running == true:
		CacheDb.Where("access_key_id = ? AND status = ? COLLATE NOCASE AND region_id = ? COLLATE NOCASE", AccessKeyId, "Running", region).Find(&ECSCache)
	case region != "all" && specifiedInstanceID != "all" && running == false:
		CacheDb.Where("access_key_id = ? AND instance_id = ? COLLATE NOCASE AND region_id = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID, region).Find(&ECSCache)
	case region != "all" && specifiedInstanceID != "all" && running == true:
		CacheDb.Where("access_key_id = ? AND instance_id = ? COLLATE NOCASE AND status = ? COLLATE NOCASE AND region_id = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID, "Running", region).Find(&ECSCache)
	}
	return ECSCache
}
