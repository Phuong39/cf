package database

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ssbeatty/sqlite"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"gorm.io/gorm"
	"strings"
)

var ConfigDb *gorm.DB
var ConfigDataBase *GlobalDB

type GlobalDB struct {
	MainDB *gorm.DB
}

func Open(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	util.HandleErr(err)
	return db
}

func InsertConfig(config cloud.Config) {
	ConfigDb.Create(&config)
}

func DeleteConfig() {
	var (
		config       string
		configListId []string
		configList   []cloud.Config
	)
	configList = SelectConfig()
	if len(configList) > 0 {
		configListId = append(configListId, "all")
		for _, v := range configList {
			configListId = append(configListId, v.Provider+"-"+v.Alias+"-"+v.AccessKeyId)
		}
		prompt := &survey.Select{
			Message: "请选择你要删除的访问凭证 (Please select the access key you want to switch): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt, &config)
		util.HandleErr(err)
		if config == "all" {
			ConfigDb.Where("in_use = ?", true).Delete(&configList)
			ConfigDb.Where("in_use = ?", false).Delete(&configList)
			log.Infoln("已删除所有访问凭证 (All access credentials have been deleted)")
		} else {
			ConfigDb.Where("access_key_id = ?", strings.Split(config, "-")[2]).Delete(&configList)
			log.Infof("%s 访问凭证已删除 (%s Access Key deleted)", config, config)
		}
	} else {
		log.Infoln("未找到任何访问凭证 (No access key found)")
	}
}

func UpdateConfigInUse(config cloud.Config) {
	ConfigDb.Model(&cloud.Config{}).Where("provider = ?", config.Provider).Update("InUse", false)
	ConfigDb.Model(&cloud.Config{}).Where("access_key_id = ?", config.AccessKeyId).Update("InUse", true)
}

func UpdateConfigSwitch(provider string) {
	var (
		config       string
		configListId []string
		configList   []cloud.Config
	)
	ConfigDb.Where("provider = ?", provider).Find(&configList)
	if len(configList) > 0 {
		for _, v := range configList {
			configListId = append(configListId, v.Alias+"-"+v.AccessKeyId)
		}
		prompt := &survey.Select{
			Message: "请选择你要切换的访问凭证 (Please select the access key you want to switch): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt, &config)
		util.HandleErr(err)
		ConfigDb.Model(&cloud.Config{}).Where("provider = ?", provider).Update("InUse", false)
		ConfigDb.Model(&cloud.Config{}).Where("access_key_id = ?", strings.Split(config, "-")[1]).Update("InUse", true)
		log.Infof("访问凭证已切换至 %s (Access Key have been switched to %s )", config, config)
	} else {
		log.Infof("未找到 %s 云服务商的访问凭证 (Access credentials for %s provider not found)", provider, provider)
	}
}

func SelectConfig() []cloud.Config {
	var configList []cloud.Config
	ConfigDb.Order("provider").Find(&configList)
	return configList
}

func SelectConfigInUse(provider string) cloud.Config {
	var (
		config     cloud.Config
		configList []cloud.Config
	)
	ConfigDb.Where("provider = ? AND in_use = ?", provider, true).Find(&configList)
	if len(configList) == 0 {
		return config
	} else {
		return configList[0]
	}
}

func init() {
	ConfigDbList := new(GlobalDB)
	ConfigDbList.MainDB = Open(pubutil.GetConfigFilePath())
	ConfigDataBase = ConfigDbList
	err := ConfigDataBase.MainDB.AutoMigrate(&cloud.Config{})
	if err != nil {
		log.Errorln("数据库自动配置失败 (Database AutoMigrate Key Struct failure)")
		util.HandleErr(err)
	}
	ConfigDb = ConfigDataBase.MainDB
}
