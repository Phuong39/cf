package database

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"
	"strings"
)

func InsertConfig(config cloud.Config) {
	CacheDb.Create(&config)
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
		configListId = append(configListId, "exit")
		prompt := &survey.Select{
			Message: "请选择你要删除的访问凭证 (Please select the access key you want to switch): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt, &config)
		errutil.HandleErr(err)
		if config == "all" {
			CacheDb.Where("in_use = ?", true).Delete(&configList)
			CacheDb.Where("in_use = ?", false).Delete(&configList)
			log.Infoln("已删除所有访问凭证 (All access credentials have been deleted)")
		} else if config == "exit" {
			os.Exit(0)
		} else {
			CacheDb.Where("access_key_id = ?", strings.Split(config, "-")[2]).Delete(&configList)
			log.Infof("%s 访问凭证已删除 (%s Access Key deleted)", config, config)
		}
	} else {
		log.Infoln("未找到任何访问凭证 (No access key found)")
	}
}

func UpdateConfigInUse(config cloud.Config) {
	CacheDb.Model(&cloud.Config{}).Where("provider = ?", config.Provider).Update("InUse", false)
	CacheDb.Model(&cloud.Config{}).Where("access_key_id = ?", config.AccessKeyId).Update("InUse", true)
}

func UpdateConfigSwitch(provider string) {
	var (
		config       string
		configListId []string
		configList   []cloud.Config
	)
	CacheDb.Where("provider = ?", provider).Find(&configList)
	if len(configList) > 0 {
		for _, v := range configList {
			configListId = append(configListId, v.Alias+"-"+v.AccessKeyId)
		}
		prompt := &survey.Select{
			Message: "请选择你要切换的访问凭证 (Please select the access key you want to switch): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt, &config)
		errutil.HandleErr(err)
		CacheDb.Model(&cloud.Config{}).Where("provider = ?", provider).Update("InUse", false)
		CacheDb.Model(&cloud.Config{}).Where("access_key_id = ?", strings.Split(config, "-")[1]).Update("InUse", true)
		log.Infof("访问凭证已切换至 %s (Access Key have been switched to %s )", config, config)
	} else {
		log.Infof("未找到 %s 云服务商的访问凭证 (Access credentials for %s provider not found)", provider, provider)
	}
}

func SelectConfig() []cloud.Config {
	var configList []cloud.Config
	CacheDb.Order("provider").Find(&configList)
	return configList
}

func SelectConfigInUse(provider string) cloud.Config {
	var (
		config     cloud.Config
		configList []cloud.Config
	)
	CacheDb.Where("provider = ? AND in_use = ?", provider, true).Find(&configList)
	if len(configList) == 0 {
		return config
	} else {
		return configList[0]
	}
}
