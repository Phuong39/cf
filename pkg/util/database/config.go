package database

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"
	"sort"
	"strings"
)

func InsertConfig(config cloud.Config) {
	var configAccessKeyIDList []string
	configList := SelectConfig()
	for _, v := range configList {
		configAccessKeyIDList = append(configAccessKeyIDList, v.AccessKeyId)
	}
	sort.Strings(configAccessKeyIDList)
	index := sort.SearchStrings(configAccessKeyIDList, config.AccessKeyId)

	if index < len(configAccessKeyIDList) && configAccessKeyIDList[index] == config.AccessKeyId {
		log.Warnln("已配置过该 Access Key (The Access Key has been configured.)")
	} else {
		CacheDb.Create(&config)
		log.Infoln("访问凭证配置完成 (Access Key configuration complete.)")
	}
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

func UpdateConfigModify() {
	var (
		config           string
		selectColumn     string
		mfValue          string
		configListId     []string
		configList       []cloud.Config
		selectColumnList = []string{"别名 (Alias)", "访问凭证 ID (Access Key Id)", "访问凭证密钥 (Secret Key)", "临时访问凭证令牌 (STS Token)"}
	)
	configList = SelectConfig()
	if len(configList) > 0 {
		for _, v := range configList {
			configListId = append(configListId, v.Provider+"-"+v.Alias+"-"+v.AccessKeyId)
		}

		prompt1 := &survey.Select{
			Message: "请选择你要修改的访问凭证 (Please select the access key you want to modify): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt1, &config)
		errutil.HandleErr(err)
		configSplit := strings.Split(config, "-")

		prompt2 := &survey.Select{
			Message: "请选择你要修改的属性 (Please select the type you want to modify): ",
			Options: selectColumnList,
		}
		err = survey.AskOne(prompt2, &selectColumn)
		errutil.HandleErr(err)
		switch {
		case selectColumn == "别名 (Alias)":
			selectColumn = "alias"
		case selectColumn == "访问凭证 ID (Access Key Id)":
			selectColumn = "access_key_id"
		case selectColumn == "访问凭证密钥 (Secret Key)":
			selectColumn = "access_key_secret"
		case selectColumn == "临时访问凭证令牌 (STS Token)":
			selectColumn = "sts_token"
		}

		var qs = []*survey.Question{
			{
				Name:   "ques",
				Prompt: &survey.Input{Message: "请输入修改后的值 (Please enter the modified value): "},
			},
		}
		err = survey.Ask(qs, &mfValue)
		errutil.HandleErr(err)

		CacheDb.Model(&cloud.Config{}).Where("provider = ? AND access_key_id = ?", configSplit[0], configSplit[2]).Update(selectColumn, mfValue)
		log.Infof("修改成功 (Successfully modified)")
	} else {
		log.Infoln("未找到任何访问凭证 (No access key found)")
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
