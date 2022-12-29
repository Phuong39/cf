package database

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"os"
	"sort"
	"strings"
)

func InsertConfig(config cloud.Config) {
	if config.AccessKeyId == "" {
		log.Warnln("当访问密钥 ID 为空的时候将不会被存储 (When the Access Key ID is empty it will not be stored.)")
	} else {
		var configAccessKeyIDList []string
		configList := SelectConfig()
		for _, v := range configList {
			configAccessKeyIDList = append(configAccessKeyIDList, v.AccessKeyId)
		}
		sort.Strings(configAccessKeyIDList)
		index := sort.SearchStrings(configAccessKeyIDList, config.AccessKeyId)

		if index < len(configAccessKeyIDList) && configAccessKeyIDList[index] == config.AccessKeyId {
			log.Warnf("已配置过 %s 访问密钥 (The %s Access Key has been configured.)", pubutil.MaskAK(config.AccessKeyId), pubutil.MaskAK(config.AccessKeyId))
		} else {
			CacheDb.Create(&config)
			log.Infof("%s 访问密钥配置完成 (%s Access Key configuration complete.)", pubutil.MaskAK(config.AccessKeyId), pubutil.MaskAK(config.AccessKeyId))
		}
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
		configListId = append(configListId, "全部访问密钥 (All access keys)")
		for _, v := range configList {
			configListId = append(configListId, v.Provider+"\t"+v.Alias+"\t"+v.AccessKeyId)
		}
		configListId = append(configListId, "退出 (Exit)")
		sort.Strings(configListId)
		prompt := &survey.Select{
			Message: "请选择您要删除的访问密钥 (Please select the access key you want to switch): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt, &config)
		errutil.HandleErr(err)
		if config == "全部访问密钥 (All access keys)" {
			var isTrue bool
			prompt := &survey.Confirm{
				Message: "此操作不可逆，您确定要删除全部的访问密钥吗？(This operation is not reversible, are you sure you want to delete all access keys?)",
				Default: false,
			}
			err := survey.AskOne(prompt, &isTrue)
			errutil.HandleErr(err)
			if isTrue {
				CacheDb.Where("in_use = ?", true).Delete(&configList)
				CacheDb.Where("in_use = ?", false).Delete(&configList)
				log.Infoln("已删除所有访问密钥 (All access keys have been deleted.)")
			} else {
				log.Infoln("已取消删除所有访问密钥 (Canceled delete all access keys.)")
			}
		} else if config == "退出 (Exit)" {
			os.Exit(0)
		} else {
			var isTrue bool
			prompt := &survey.Confirm{
				Message: "此操作不可逆，您确定要删除选中的访问密钥吗？(This operation is not reversible, are you sure you want to delete the selected access key?)",
				Default: false,
			}
			err := survey.AskOne(prompt, &isTrue)
			errutil.HandleErr(err)
			if isTrue {
				accessKeyId := strings.Split(config, "\t")[2]
				CacheDb.Where("access_key_id = ?", accessKeyId).Delete(&configList)
				log.Infof("%s 访问密钥已删除 (%s Access Key deleted)", pubutil.MaskAK(accessKeyId), pubutil.MaskAK(accessKeyId))
			} else {
				log.Infoln("已取消删除选中的访问密钥 (Canceled delete the selected access key.)")
			}
		}
	} else {
		log.Infoln("未找到任何访问密钥 (No access key found)")
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
			configListId = append(configListId, v.Provider+"\t"+v.Alias+"\t"+v.AccessKeyId)
		}
		sort.Strings(configListId)
		prompt := &survey.Select{
			Message: "请选择您要切换的访问密钥 (Please select the access key you want to switch): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt, &config)
		errutil.HandleErr(err)
		accessKeyID := strings.Split(config, "\t")[2]
		CacheDb.Model(&cloud.Config{}).Where("provider = ?", provider).Update("InUse", false)
		CacheDb.Model(&cloud.Config{}).Where("access_key_id = ?", accessKeyID).Update("InUse", true)
		log.Infof("访问密钥已切换至 %s (Access Key have been switched to %s )", pubutil.MaskAK(accessKeyID), pubutil.MaskAK(accessKeyID))
	} else {
		log.Infof("未找到 %s 云服务商的访问密钥 (access keys for %s provider not found)", provider, provider)
	}
}

func UpdateConfigModify() {
	var (
		config           string
		selectColumn     string
		mfValue          string
		configListId     []string
		configList       []cloud.Config
		selectColumnList = []string{"别名 (Alias)", "访问密钥 ID (Access Key Id)", "访问密钥密钥 (Secret Key)", "临时访问密钥令牌 (STS Token)"}
	)
	configList = SelectConfig()
	if len(configList) > 0 {
		for _, v := range configList {
			configListId = append(configListId, v.Provider+"\t"+v.Alias+"\t"+v.AccessKeyId)
		}
		sort.Strings(configListId)
		prompt1 := &survey.Select{
			Message: "请选择您要修改的访问密钥 (Please select the access key you want to modify): ",
			Options: configListId,
		}
		err := survey.AskOne(prompt1, &config)
		errutil.HandleErr(err)
		configSplit := strings.Split(config, "\t")
		sort.Strings(selectColumnList)
		prompt2 := &survey.Select{
			Message: "请选择您要修改的属性 (Please select the type you want to modify): ",
			Options: selectColumnList,
		}
		err = survey.AskOne(prompt2, &selectColumn)
		errutil.HandleErr(err)
		switch {
		case selectColumn == "别名 (Alias)":
			selectColumn = "alias"
		case selectColumn == "访问密钥 ID (Access Key Id)":
			selectColumn = "access_key_id"
		case selectColumn == "访问密钥密钥 (Secret Key)":
			selectColumn = "access_key_secret"
		case selectColumn == "临时访问密钥令牌 (STS Token)":
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
		log.Infoln("未找到任何访问密钥 (No access key found)")
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
