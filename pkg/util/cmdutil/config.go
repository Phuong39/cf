package cmdutil

import (
	"fmt"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util"
)

func ConfigureAccessKey() {
	cloudConfigList, cloudProviderList, cloudProvider := selectProvider()
	for i, j := range cloudProviderList {
		if j == cloudProvider {
			config := GetConfig(cloudConfigList[i])
			inputAccessKey(config, cloudConfigList[i])
		}
	}
}

func selectProvider() ([]string, []string, string) {
	var cloudProvider string
	cloudConfigList, cloudProviderList := ReturnCloudProviderList()
	prompt := &survey.Select{
		Message: "选择你要配置的云服务商 (Select a cloud provider): ",
		Options: cloudProviderList,
	}
	err := survey.AskOne(prompt, &cloudProvider)
	util.HandleErr(err)
	return cloudConfigList, cloudProviderList, cloudProvider
}

func ReturnCloudProviderList() ([]string, []string) {
	var (
		cloudConfigList   []string
		cloudProviderList []string
		CloudProviderMap  = map[string]string{
			"alibaba": "阿里云 (Alibaba Cloud)",
			"tencent": "腾讯云 (Tencent Cloud)",
		}
	)
	for k, v := range CloudProviderMap {
		cloudConfigList = append(cloudConfigList, k)
		cloudProviderList = append(cloudProviderList, v)
	}
	return cloudConfigList, cloudProviderList
}

func inputAccessKey(config cloud.Config, provider string) {
	OldAlias := ""
	OldAccessKeyId := ""
	OldAccessKeySecret := ""
	OldSTSToken := ""
	Alias := config.Alias
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
			Name:   "Alias",
			Prompt: &survey.Input{Message: "输入访问凭证别名 (Input Access Key Alias) (必须 Required)" + OldAlias + ":"},
		},
		{
			Name:   "AccessKeyId",
			Prompt: &survey.Input{Message: "输入访问凭证 ID (Input Access Key Id) (必须 Required)" + OldAccessKeyId + ":"},
		},
		{
			Name:   "AccessKeySecret",
			Prompt: &survey.Password{Message: "输入访问凭证密钥 (Input Access Key Secret) (必须 Required)" + OldAccessKeySecret + ":"},
		},
		{
			Name:   "STSToken",
			Prompt: &survey.Input{Message: "输入临时凭证的 Token (Input STS Token) (可选 Optional)" + OldSTSToken + ":"},
		},
	}

	cred := cloud.Config{}
	err := survey.Ask(qs, &cred)
	cred.Alias = strings.TrimSpace(cred.Alias)
	cred.AccessKeyId = strings.TrimSpace(cred.AccessKeyId)
	cred.AccessKeySecret = strings.TrimSpace(cred.AccessKeySecret)
	cred.STSToken = strings.TrimSpace(cred.STSToken)
	cred.Provider = provider
	if cred.Alias == "" {
		cred.Alias = Alias
	}
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
	SaveAccessKey(cred)
}

func SaveAccessKey(config cloud.Config) {
	configFilePath := pubutil.GetConfigFilePath()
	database.InsertConfig(config)
	database.UpdateConfigInUse(config)
	log.Debugf("配置文件路径 (Configuration file path): %s ", configFilePath)
	pubutil.CreateFolder(ReturnCacheDict())
}

func GetConfig(provider string) cloud.Config {
	return database.SelectConfigInUse(provider)
}

func ConfigLs() {
	var (
		STSToken          string
		CommonTableHeader = []string{"别名 (alias)", "访问凭证 ID (access_key_id)", "访问凭证密钥 (access_key_secret)", "临时访问凭证令牌 (sts_token)", "云服务提供商 (provider)", "是否在使用 (in_use)"}
	)
	configList := database.SelectConfig()
	Data := cloud.TableData{
		Header: CommonTableHeader,
	}
	if len(configList) == 0 {
		log.Info("未找到任何密钥 (No key found)")
	} else {
		for _, k := range configList {
			if len(STSToken) > 10 {
				STSToken = MaskAK(STSToken)
			} else {
				STSToken = k.STSToken
			}
			Data.Body = append(Data.Body, []string{
				k.Alias,
				k.AccessKeyId,
				k.AccessKeySecret,
				STSToken,
				k.Provider,
				strconv.FormatBool(k.InUse),
			})
		}
		cloud.PrintTable(Data, "当前存储的访问凭证信息")
	}
}

func ConfigSw() {
	cloudConfigList, cloudProviderList, cloudProvider := selectProvider()
	for i, j := range cloudProviderList {
		if j == cloudProvider {
			database.UpdateConfigSwitch(cloudConfigList[i])
		}
	}
}

func ConfigDel() {
	database.DeleteConfig()
}

func MaskAK(ak string) string {
	prefix := ak[:2]
	suffix := ak[len(ak)-6:]
	return prefix + strings.Repeat("*", 18) + suffix
}
