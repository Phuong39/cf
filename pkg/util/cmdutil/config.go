package cmdutil

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/global"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
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
		Message: "选择你要设置的云服务商 (Select a cloud provider): ",
		Options: cloudProviderList,
	}
	err := survey.AskOne(prompt, &cloudProvider)
	errutil.HandleErr(err)
	return cloudConfigList, cloudProviderList, cloudProvider
}

func ReturnCloudProviderList() ([]string, []string) {
	var (
		cloudConfigList   []string
		cloudProviderList []string
		CloudProviderMap  = global.CloudProviderMap
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
	errutil.HandleErr(err)
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

func ConfigLs(selectAll bool) {
	var (
		STSToken          string
		CommonTableHeader = []string{"别名 (Alias)", "访问凭证 ID (Access Key Id)", "访问凭证密钥 (Secret Key)", "临时访问凭证令牌 (STS Token)", "云服务提供商 (Provider)", "是否在使用 (In Use)"}
	)
	configList := database.SelectConfig()
	if selectAll {
		for _, v := range configList {
			color.Tag("info").Print("\n别名 (Alias): ")
			fmt.Println(v.Alias)
			color.Tag("info").Print("访问凭证 ID (Access Key Id): ")
			fmt.Println(v.AccessKeyId)
			color.Tag("info").Print("访问凭证密钥 (Secret Key): ")
			fmt.Println(v.AccessKeySecret)
			color.Tag("info").Print("临时访问凭证令牌 (STS Token): ")
			fmt.Println(v.STSToken)
			color.Tag("info").Print("云服务提供商 (Provider): ")
			fmt.Println(v.Provider)
			color.Tag("info").Print("是否在使用 (In Use): ")
			fmt.Println(v.InUse)
		}
	} else {
		Data := cloud.TableData{
			Header: CommonTableHeader,
		}
		if len(configList) == 0 {
			log.Info("未找到任何密钥 (No key found)")
		} else {
			for _, v := range configList {
				if len(v.STSToken) > 10 {
					STSToken = MaskAK(v.STSToken)
				} else {
					STSToken = v.STSToken
				}
				Data.Body = append(Data.Body, []string{
					v.Alias,
					v.AccessKeyId,
					v.AccessKeySecret,
					STSToken,
					v.Provider,
					strconv.FormatBool(v.InUse),
				})
			}
			cloud.PrintTable(Data, "当前存储的访问凭证信息")
		}
	}
}

func ConfigMf() {
	database.UpdateConfigModify()
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
