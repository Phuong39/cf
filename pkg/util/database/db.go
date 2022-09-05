package database

import (
	log "github.com/sirupsen/logrus"
	"github.com/ssbeatty/sqlite"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"gorm.io/gorm"
)

var CacheDb *gorm.DB
var CacheDataBase *GlobalDB

type GlobalDB struct {
	MainDB *gorm.DB
}

func Open(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	errutil.HandleErr(err)
	return db
}

func init() {
	var err error
	CacheDbList := new(GlobalDB)
	CacheDbList.MainDB = Open(pubutil.GetConfigFilePath())
	CacheDataBase = CacheDbList
	err = CacheDataBase.MainDB.AutoMigrate(&cloud.Config{})
	errutil.HandleErr(err)
	err = CacheDataBase.MainDB.AutoMigrate(&pubutil.TimestampCache{})
	errutil.HandleErr(err)
	err = CacheDataBase.MainDB.AutoMigrate(&pubutil.OSSCache{})
	errutil.HandleErr(err)
	err = CacheDataBase.MainDB.AutoMigrate(&pubutil.ECSCache{})
	errutil.HandleErr(err)
	err = CacheDataBase.MainDB.AutoMigrate(&pubutil.RDSCache{})
	errutil.HandleErr(err)
	err = CacheDataBase.MainDB.AutoMigrate(&pubutil.TakeoverConsoleCache{})
	if err != nil {
		log.Errorln("数据库自动配置失败 (Database AutoMigrate Key Struct failure)")
		errutil.HandleErr(err)
	}
	CacheDb = CacheDataBase.MainDB
}
