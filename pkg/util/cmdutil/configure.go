package cmdutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
)

func ConfigureAccessKey() {
	var cloudProvider string
	cloudConfigList, cloudProviderList := ReturnCloudProviderList()
	prompt := &survey.Select{
		Message: "选择你要配置的云服务商 (Select a cloud provider): ",
		Options: cloudProviderList,
	}
	err := survey.AskOne(prompt, &cloudProvider)
	util.HandleErr(err)
	for i, j := range cloudProviderList {
		if j == cloudProvider {
			config := GetConfig(cloudConfigList[i])
			inputAccessKey(config, cloudConfigList[i])
		}
	}
}

func ReturnCloudProviderList() ([]string, []string) {
	var (
		cloudConfigList   []string
		cloudProviderList []string
		CloudProviderMap  = map[string]string{"alibaba": "阿里云 (Alibaba Cloud)", "tencent": "腾讯云 (Tencent Cloud)"}
	)
	for k, v := range CloudProviderMap {
		cloudConfigList = append(cloudConfigList, k)
		cloudProviderList = append(cloudProviderList, v)
	}
	return cloudConfigList, cloudProviderList
}

func inputAccessKey(config cloud.Config, provider string) {
	OldAccessKeyId := ""
	OldAccessKeySecret := ""
	OldSTSToken := ""
	AccessKeyId := config.AccessKeyId
	AccessKeySecret := config.AccessKeySecret
	STSToken := config.STSToken
	if AccessKeyId != "" {
		OldAccessKeyId = fmt.Sprintf(" [%s] ", MaskAK(AccessKeyId))
	}
	if AccessKeySecret != "" {
		OldAccessKeySecret = fmt.Sprintf(" [%s] ", MaskAK(AccessKeySecret))
	}
	if STSToken != "" {
		OldSTSToken = fmt.Sprintf(" [%s] ", MaskAK(STSToken))
	}
	var qs = []*survey.Question{
		{
			Name:   "AccessKeyId",
			Prompt: &survey.Input{Message: "Access Key Id (必须 Required)" + OldAccessKeyId + ":"},
		},
		{
			Name:   "AccessKeySecret",
			Prompt: &survey.Password{Message: "Access Key Secret (必须 Required)" + OldAccessKeySecret + ":"},
		},
		{
			Name:   "STSToken",
			Prompt: &survey.Input{Message: "STS Token (可选 Optional)" + OldSTSToken + ":"},
		},
	}
	cred := cloud.Config{}
	err := survey.Ask(qs, &cred)
	cred.AccessKeyId = strings.TrimSpace(cred.AccessKeyId)
	cred.AccessKeySecret = strings.TrimSpace(cred.AccessKeySecret)
	cred.STSToken = strings.TrimSpace(cred.STSToken)
	if cred.AccessKeyId == "" {
		cred.AccessKeyId = AccessKeyId
	}
	if cred.AccessKeySecret == "" {
		cred.AccessKeySecret = AccessKeySecret
	}
	if cred.STSToken == "" && strings.Contains(cred.AccessKeyId, "STS.") {
		cred.STSToken = STSToken
	}
	util.HandleErr(err)
	SaveAccessKey(cred, provider)
}

func SaveAccessKey(config cloud.Config, provider string) {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	if FileExists(home) == false {
		err = os.MkdirAll(home, 0700)
	}
	util.HandleErr(err)
	configJSON, err := json.MarshalIndent(config, "", "    ")
	util.HandleErr(err)
	configFilePath := GetConfigFilePath(provider)
	err = ioutil.WriteFile(configFilePath, configJSON, 0600)
	util.HandleErr(err)
	log.Infof("配置文件路径 (Configuration file path): %s ", configFilePath)
	createCacheDict()
}

func GetConfigFilePath(provider string) string {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	configHomeFile := filepath.Join(home, "config")
	if FileExists(configHomeFile) == false {
		err = os.MkdirAll(configHomeFile, 0700)
		util.HandleErr(err)
	}
	configFilePath := filepath.Join(configHomeFile, provider+"Config.json")
	return configFilePath
}

func GetConfig(provider string) cloud.Config {
	configFilePath := GetConfigFilePath(provider)
	var config cloud.Config
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return config
	} else {
		file, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			util.HandleErr(err)
		}
		err = json.Unmarshal(file, &config)
		if err != nil {
			util.HandleErr(err)
		}
		return config
	}
}

func MaskAK(ak string) string {
	prefix := ak[:2]
	suffix := ak[len(ak)-6:]
	return prefix + strings.Repeat("*", 18) + suffix
}
