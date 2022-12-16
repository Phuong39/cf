package cmdutil

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/teamssix/cf/pkg/util/global"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
)

const (
	alibaba = "alibaba"
	tencent = "tencent"
	aws     = "aws"
	huawei  = "huawei"
)

func ConfigureAccessKey() {
	cloudConfigList, cloudProviderList, cloudProvider := selectProvider()
	for i, j := range cloudProviderList {
		if j == cloudProvider {
			var credList []cloud.Config
			switch cloudConfigList[i] {
			case alibaba:
				credList = append(credList, findAlibabaConfig()...)
			case tencent:
				credList = append(credList, findTencentConfig()...)
			case aws:
				credList = append(credList, findAWSConfig()...)
			case huawei:
				credList = append(credList, findHuaweiConfig()...)
			}
			if len(credList) != 0 {
				var (
					isTrue     bool
					selectedAK string
				)
				prompt := &survey.Confirm{
					Message: "在当前系统中发现访问密钥，是否导入？(Access keys were found in the current system, are they import?)",
					Default: false,
				}
				err := survey.AskOne(prompt, &isTrue)
				errutil.HandleErr(err)
				if isTrue {
					var accessKeyList []string
					if len(credList) > 1 {
						accessKeyList = append(accessKeyList, "全部访问密钥 (All access keys)")
					}
					for i, v := range credList {
						i = i + 1
						accessKeyList = append(accessKeyList, strconv.Itoa(i)+"\t"+v.Provider+"\t"+v.Alias+"\t"+v.AccessKeyId)
					}
					accessKeyList = append(accessKeyList, "退出 (Exit)")
					sort.Strings(accessKeyList)
					prompt := &survey.Select{
						Message: "选择您要导入的访问密钥 (Select the access key you want to import): ",
						Options: accessKeyList,
					}
					err := survey.AskOne(prompt, &selectedAK)
					errutil.HandleErr(err)
					if selectedAK == "全部访问密钥 (All access keys)" {
						log.Infoln("在导入全部的访问密钥后，您可以通过 \"cf config sw\" 来切换访问密钥。 (After importing all access keys, you can switch access key via \"cf config sw\".)")
						for _, v := range credList {
							SaveAccessKey(v)
						}
					} else if selectedAK == "退出 (Exit)" {
						log.Debugln("正在退出…… (Exiting...)")
					} else {
						for _, v := range credList {
							if v.AccessKeyId == strings.Split(selectedAK, "\t")[3] {
								SaveAccessKey(v)
							}
						}
					}
				} else {
					log.Infoln("已取消自动导入，请输入您要添加的访问密钥 (Automatic import has been cancelled, please enter the access key you want to add.)")
					config := GetConfig(cloudConfigList[i])
					inputAccessKey(config, cloudConfigList[i])
				}
			} else {
				config := GetConfig(cloudConfigList[i])
				inputAccessKey(config, cloudConfigList[i])
			}
		}
	}
}

func selectProvider() ([]string, []string, string) {
	var cloudProvider string
	cloudConfigList, cloudProviderList := ReturnCloudProviderList()
	prompt := &survey.Select{
		Message: "选择您要设置的云服务商 (Select a cloud provider): ",
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
		OldAccessKeyId = fmt.Sprintf(" [%s] ", pubutil.MaskAK(AccessKeyId))
	}
	if AccessKeySecret != "" {
		OldAccessKeySecret = fmt.Sprintf(" [%s] ", pubutil.MaskAK(AccessKeySecret))
	}
	if STSToken != "" {
		OldSTSToken = fmt.Sprintf(" [%s] ", pubutil.MaskAK(STSToken))
	}
	var qs = []*survey.Question{
		{
			Name:   "Alias",
			Prompt: &survey.Input{Message: "输入访问密钥别名 (Input Access Key Alias) (必须 Required)" + OldAlias + ":"},
		},
		{
			Name:   "AccessKeyId",
			Prompt: &survey.Input{Message: "输入访问密钥 ID (Input Access Key Id) (必须 Required)" + OldAccessKeyId + ":"},
			Validate: func(val interface{}) error {
				str := val.(string)
				if len(strings.TrimSpace(str)) < 7 {
					log.Warnln("访问凭证似乎输入有误 (This access credential appears to be incorrect.)")
				}
				return nil
			},
		},
		{
			Name:   "AccessKeySecret",
			Prompt: &survey.Password{Message: "输入访问密钥密钥 (Input Access Key Secret) (必须 Required)" + OldAccessKeySecret + ":"},
			Validate: func(val interface{}) error {
				str := val.(string)
				if len(strings.TrimSpace(str)) < 7 {
					log.Warnln("访问凭证似乎输入有误 (This access credential appears to be incorrect.)")
				}
				return nil
			},
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
		CommonTableHeader = []string{"云服务提供商 (Provider)", " 别名 (Alias)", "访问密钥 ID (Access Key Id)", "访问密钥密钥 (Secret Key)", "临时访问密钥令牌 (STS Token)", "是否在使用 (In Use)"}
	)
	configList := database.SelectConfig()
	if selectAll {
		for _, v := range configList {
			color.Tag("info").Print("\n云服务提供商 (Provider): ")
			fmt.Println(v.Provider)
			color.Tag("info").Print("别名 (Alias): ")
			fmt.Println(v.Alias)
			color.Tag("info").Print("访问密钥 ID (Access Key Id): ")
			fmt.Println(v.AccessKeyId)
			color.Tag("info").Print("访问密钥密钥 (Secret Key): ")
			fmt.Println(v.AccessKeySecret)
			color.Tag("info").Print("临时访问密钥令牌 (STS Token): ")
			fmt.Println(v.STSToken)
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
					STSToken = pubutil.MaskAK(v.STSToken)
				} else {
					STSToken = v.STSToken
				}
				Data.Body = append(Data.Body, []string{
					v.Provider,
					v.Alias,
					v.AccessKeyId,
					v.AccessKeySecret,
					STSToken,
					strconv.FormatBool(v.InUse),
				})
			}
			cloud.PrintTable(Data, "当前存储的访问密钥信息 (Current stored access key information)")
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

func ScanAccessKey(selectAll bool) {
	var credList []cloud.Config
	credList = append(credList, findAlibabaConfig()...)
	credList = append(credList, findTencentConfig()...)
	credList = append(credList, findAWSConfig()...)
	credList = append(credList, findHuaweiConfig()...)
	if len(credList) == 0 {
		log.Infoln("在当前系统中未扫描到任何访问密钥 (No access keys were scanned in the current environment.)")
	} else {
		if selectAll {
			for _, v := range credList {
				color.Tag("info").Print("\n云服务提供商 (Provider): ")
				fmt.Println(v.Provider)
				color.Tag("info").Print("别名 (Alias): ")
				fmt.Println(v.Alias)
				color.Tag("info").Print("访问密钥 ID (Access Key Id): ")
				fmt.Println(v.AccessKeyId)
				color.Tag("info").Print("访问密钥密钥 (Secret Key): ")
				fmt.Println(v.AccessKeySecret)
				color.Tag("info").Print("临时访问密钥令牌 (STS Token): ")
				fmt.Println(v.STSToken)
			}
		} else {
			Data := cloud.TableData{
				Header: []string{"云服务提供商 (Provider)", "别名 (Alias)", "访问密钥 ID (Access Key Id)", "访问密钥密钥 (Secret Key)", "临时访问密钥令牌 (STS Token)"},
			}
			var STSToken string
			for _, v := range credList {
				if len(v.STSToken) > 10 {
					STSToken = pubutil.MaskAK(v.STSToken)
				} else {
					STSToken = v.STSToken
				}
				Data.Body = append(Data.Body, []string{
					v.Provider,
					v.Alias,
					v.AccessKeyId,
					v.AccessKeySecret,
					STSToken,
				})
			}
			cloud.PrintTable(Data, "扫描到的访问密钥信息 (Scanned access key information)")
		}
	}
}
