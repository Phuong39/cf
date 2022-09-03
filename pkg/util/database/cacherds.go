package database

import (
	"github.com/teamssix/cf/pkg/util/pubutil"
)

func InsertRDSCache(RDSCache []pubutil.RDSCache) {
	DeleteRDSCache(RDSCache[0].AccessKeyId)
	CacheDb.Create(&RDSCache)
}

func DeleteRDSCache(AccessKeyId string) {
	var RDSCache []pubutil.RDSCache
	CacheDb.Where("access_key_id = ?", AccessKeyId).Delete(&RDSCache)
}

func SelectRDSCache(provider string) []pubutil.RDSCache {
	var RDSCache []pubutil.RDSCache
	AccessKeyId := SelectConfigInUse(provider).AccessKeyId
	CacheDb.Where("access_key_id = ?", AccessKeyId).Find(&RDSCache)
	return RDSCache
}

func SelectRDSCacheFilter(provider string, region string, specifiedInstanceID string, engine string) []pubutil.RDSCache {
	var RDSCache []pubutil.RDSCache
	AccessKeyId := SelectConfigInUse(provider).AccessKeyId
	switch {
	case region == "all" && specifiedInstanceID == "all" && engine == "all":
		CacheDb.Where("access_key_id = ? COLLATE NOCASE", AccessKeyId).Find(&RDSCache)
	case region == "all" && specifiedInstanceID == "all" && engine != "all":
		CacheDb.Where("access_key_id = ? AND engine = ? COLLATE NOCASE", AccessKeyId, engine).Find(&RDSCache)
	case region == "all" && specifiedInstanceID != "all" && engine == "all":
		CacheDb.Where("access_key_id = ? AND db_instance_id = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID).Find(&RDSCache)
	case region == "all" && specifiedInstanceID != "all" && engine != "all":
		CacheDb.Where("access_key_id = ? AND db_instance_id = ? COLLATE NOCASE AND engine = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID, engine).Find(&RDSCache)
	case region != "all" && specifiedInstanceID == "all" && engine == "all":
		CacheDb.Where("access_key_id = ? AND region_id = ? COLLATE NOCASE", AccessKeyId, region).Find(&RDSCache)
	case region != "all" && specifiedInstanceID == "all" && engine != "all":
		CacheDb.Where("access_key_id = ? AND engine = ? COLLATE NOCASE AND region_id = ? COLLATE NOCASE", AccessKeyId, engine, region).Find(&RDSCache)
	case region != "all" && specifiedInstanceID != "all" && engine == "all":
		CacheDb.Where("access_key_id = ? AND db_instance_id = ? COLLATE NOCASE AND region_id = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID, region).Find(&RDSCache)
	case region != "all" && specifiedInstanceID != "all" && engine != "all":
		CacheDb.Where("access_key_id = ? AND db_instance_id = ? COLLATE NOCASE AND engine = ? COLLATE NOCASE AND region_id = ? COLLATE NOCASE", AccessKeyId, specifiedInstanceID, engine, region).Find(&RDSCache)
	}
	return RDSCache
}
