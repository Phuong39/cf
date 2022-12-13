package cmdutil

import (
	"github.com/bitly/go-simplejson"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/pubutil"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func findAlibabaConfig() []cloud.Config {
	var credList []cloud.Config
	// 1. credential file
	alibabaConfigFile := filepath.Join(pubutil.GetUserDir(), "/.aliyun/config.json")
	isTrue, content := pubutil.ReadFile(alibabaConfigFile)
	if isTrue {
		contentJson, _ := simplejson.NewJson([]byte(content))
		contentJsonArray, _ := contentJson.Get("profiles").Array()
		for _, v := range contentJsonArray {
			cred := cloud.Config{}
			contentResult, _ := v.(map[string]interface{})
			cred.Alias = "local_" + contentResult["name"].(string)
			cred.AccessKeyId = contentResult["access_key_id"].(string)
			cred.AccessKeySecret = contentResult["access_key_secret"].(string)
			cred.STSToken = contentResult["sts_token"].(string)
			cred.Provider = alibaba
			if cred.AccessKeyId != "" {
				credList = append(credList, cred)
			}
		}
	}
	// 2. environment variables
	cred := cloud.Config{}
	cred.Provider = alibaba
	cred.Alias = "local_env"
	cred.AccessKeyId = os.Getenv("ALIBABACLOUD_ACCESS_KEY_ID")
	cred.AccessKeySecret = os.Getenv("ALIBABACLOUD_ACCESS_KEY_SECRET")
	cred.STSToken = os.Getenv("SECURITY_TOKEN")
	if cred.AccessKeyId != "" {
		credList = append(credList, cred)
	}
	return credList
}

func findTencentConfig() []cloud.Config {
	var credList []cloud.Config
	// 1. credential file
	tencentConfigPath := filepath.Join(pubutil.GetUserDir(), "/.tccli")
	tencentConfigFiles, _ := ioutil.ReadDir(tencentConfigPath)
	for _, f := range tencentConfigFiles {
		tencentConfigName := f.Name()
		if path.Ext(tencentConfigName) == ".credential" {
			tencentConfigFile := filepath.Join(tencentConfigPath, tencentConfigName)
			isTrue, content := pubutil.ReadFile(tencentConfigFile)
			if isTrue {
				contentJson, _ := simplejson.NewJson([]byte(content))
				cred := cloud.Config{}
				cred.Alias = "local_" + strings.TrimSuffix(tencentConfigName, ".credential")
				cred.AccessKeyId = contentJson.Get("secretId").MustString()
				cred.AccessKeySecret = contentJson.Get("secretKey").MustString()
				cred.Provider = tencent
				if cred.AccessKeyId != "" {
					credList = append(credList, cred)
				}
			}
		}
	}
	// 2. environment variables
	cred := cloud.Config{}
	cred.Provider = tencent
	cred.Alias = "local_env"
	cred.AccessKeyId = os.Getenv("TENCENTCLOUD_SECRET_ID")
	cred.AccessKeySecret = os.Getenv("TENCENTCLOUD_SECRET_KEY")
	if cred.AccessKeyId != "" {
		credList = append(credList, cred)
	}
	return credList
}

func findAWSConfig() []cloud.Config {
	var credList []cloud.Config
	// 1. credential file
	awsConfigFile := filepath.Join(pubutil.GetUserDir(), "/.aws/credentials")
	isTrue, content := pubutil.ReadFile(awsConfigFile)
	if isTrue {
		for _, v := range strings.Split(content, "[") {
			cred := cloud.Config{}
			if len(pubutil.StringClean(v)) != 0 {
				for _, j := range strings.Split(v, "\n") {
					if strings.Contains(j, "]") {
						cred.Alias = "local_" + strings.Replace(j, "]", "", -1)
					} else if strings.Contains(j, "aws_access_key_id") {
						cred.AccessKeyId = pubutil.StringClean(strings.Split(j, "=")[1])
					} else if strings.Contains(j, "aws_secret_access_key") {
						cred.AccessKeySecret = pubutil.StringClean(strings.Split(j, "=")[1])
					} else if strings.Contains(j, "aws_session_token") {
						cred.STSToken = pubutil.StringClean(strings.Split(j, "=")[1])
					}
				}
				cred.Provider = aws
				if cred.AccessKeyId != "" {
					credList = append(credList, cred)
				}
			}
		}
	}
	// 2. environment variables
	cred := cloud.Config{}
	cred.Provider = aws
	cred.Alias = "local_env"
	cred.AccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	cred.AccessKeySecret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	cred.STSToken = os.Getenv("AWS_SESSION_TOKEN")
	if cred.AccessKeyId != "" {
		credList = append(credList, cred)
	}
	return credList
}

func findHuaweiConfig() []cloud.Config {
	var credList []cloud.Config
	// 1. credential file
	huaweiConfigFile := filepath.Join(pubutil.GetUserDir(), "/.huaweicloud/credentials")
	isTrue, content := pubutil.ReadFile(huaweiConfigFile)
	if isTrue {
		for _, v := range strings.Split(content, "[") {
			cred := cloud.Config{}
			if len(pubutil.StringClean(v)) != 0 {
				for _, j := range strings.Split(v, "\n") {
					if strings.Contains(j, "]") {
						cred.Alias = "local_" + strings.Replace(j, "]", "", -1)
					} else if strings.Contains(j, "ak") {
						cred.AccessKeyId = pubutil.StringClean(strings.Split(j, "=")[1])
					} else if strings.Contains(j, "sk") {
						cred.AccessKeySecret = pubutil.StringClean(strings.Split(j, "=")[1])
					} else if strings.Contains(j, "security_token") {
						cred.STSToken = pubutil.StringClean(strings.Split(j, "=")[1])
					}
				}
				cred.Provider = huawei
				if cred.AccessKeyId != "" {
					credList = append(credList, cred)
				}
			}
		}
	}
	// 2. environment variables
	cred := cloud.Config{}
	cred.Provider = huawei
	cred.Alias = "local_env_sdk"
	cred.AccessKeyId = os.Getenv("HUAWEICLOUD_SDK_AK")
	cred.AccessKeySecret = os.Getenv("HUAWEICLOUD_SDK_SK")
	cred.STSToken = os.Getenv("HUAWEICLOUD_SDK_SECURITY_TOKEN")
	if cred.AccessKeyId != "" {
		credList = append(credList, cred)
	}

	cred = cloud.Config{}
	cred.Provider = huawei
	cred.Alias = "local_env_obs"
	cred.AccessKeyId = os.Getenv("OBS_ACCESS_KEY_ID")
	cred.AccessKeySecret = os.Getenv("OBS_SECRET_ACCESS_KEY")
	cred.STSToken = os.Getenv("OBS_SECURITY_TOKEN")
	if cred.AccessKeyId != "" {
		credList = append(credList, cred)
	}
	return credList
}
