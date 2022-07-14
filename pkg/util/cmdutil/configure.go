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

func ConfigureAccessKey(cf string) {
	config := GetAllCredential()
	//config := GetAliCredential()
	OldAccessKeyId := ""
	OldAccessKeySecret := ""
	OldSTSToken := ""
	OldToken := ""
	OldTmpSecretId := ""
	OldTmpSecretKey := ""
	if cf == "tencent" {
		SecretId := config.Tencent.SecretId
		SecretKey := config.Tencent.SecretKey
		if SecretId != "" {
			OldAccessKeyId = fmt.Sprintf(" [%s] ", maskAK(SecretId))
		}
		if SecretKey != "" {
			OldAccessKeySecret = fmt.Sprintf(" [%s] ", maskAK(SecretKey))
		}
		var qs = []*survey.Question{
			{
				Name:   "SecretId",
				Prompt: &survey.Input{Message: "Secret Key Id (可选 Optional)" + OldAccessKeyId + ":"},
				//Validate: survey.Required,
			},
			{
				Name:   "SecretKey",
				Prompt: &survey.Password{Message: "Secret Key (可选 Optional)" + OldAccessKeySecret + ":"},
				//Validate: survey.Required,
			},
			{
				Name:   "Token",
				Prompt: &survey.Input{Message: "Token (可选 Optional)" + OldToken + ":"},
			},
			{
				Name:   "TmpSecretId",
				Prompt: &survey.Input{Message: "Tmp Secret Id (可选 Optional)" + OldTmpSecretId + ":"},
			},
			{
				Name:   "TmpSecretKey",
				Prompt: &survey.Password{Message: "Tmp Secret Key (可选 Optional)" + OldTmpSecretKey + ":"},
			},
		}
		cred := cloud.Credential{}
		err := survey.Ask(qs, &cred.Tencent)
		util.HandleErr(err)
		SaveAccessKey(cred)
	} else if cf == "alibaba" {
		AccessKeyId := config.Alibaba.AccessKeyId
		AccessKeySecret := config.Alibaba.AccessKeySecret
		STSToken := config.Alibaba.STSToken
		if AccessKeyId != "" {
			OldAccessKeyId = fmt.Sprintf(" [%s] ", maskAK(AccessKeyId))
		}
		if AccessKeySecret != "" {
			OldAccessKeySecret = fmt.Sprintf(" [%s] ", maskAK(AccessKeySecret))
		}
		if STSToken != "" {
			OldSTSToken = fmt.Sprintf(" [%s] ", maskAK(STSToken))
		}
		var qs = []*survey.Question{
			{
				Name:   "AccessKeyId",
				Prompt: &survey.Input{Message: "Access Key Id (可选 Optional)" + OldAccessKeyId + ":"},
				//Validate: survey.Required,
			},
			{
				Name:   "AccessKeySecret",
				Prompt: &survey.Password{Message: "Access Key Secret (可选 Optional)" + OldAccessKeySecret + ":"},
				//Validate: survey.Required,
			},
			{
				Name:   "STSToken",
				Prompt: &survey.Input{Message: "STS Token (可选 Optional)" + OldSTSToken + ":"},
			},
		}
		cred := cloud.Credential{}
		err := survey.Ask(qs, &cred.Alibaba)
		util.HandleErr(err)
		SaveAccessKey(cred)
	} else {
		log.Fatal("请检查输入的云厂商名称是否正确！(Please check cloud name!)")
	}
}

func SaveAccessKey(config cloud.Credential) {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	if FileExists(home) == false {
		err = os.MkdirAll(home, 0700)
	}
	util.HandleErr(err)
	configJSON, err := json.MarshalIndent(config, "", "    ")
	util.HandleErr(err)
	AllCredentialFilePath := GetAllCredentialFilePath()
	err = ioutil.WriteFile(AllCredentialFilePath, configJSON, 0600)
	util.HandleErr(err)
	log.Infof("配置完成，配置文件路径 (Configure done, Configuration file path): %s ", AllCredentialFilePath)
	createCacheDict()
}

//func SaveAccessKey(config cloud.Credential) {
//	home, err := GetCFHomeDir()
//	util.HandleErr(err)
//	if FileExists(home) == false {
//		err = os.MkdirAll(home, 0700)
//	}
//	util.HandleErr(err)
//	configJSON, err := json.MarshalIndent(config, "", "    ")
//	util.HandleErr(err)
//	AliCredentialFilePath := GetAliCredentialFilePath()
//	err = ioutil.WriteFile(AliCredentialFilePath, configJSON, 0600)
//	util.HandleErr(err)
//	log.Infof("配置完成，配置文件路径 (Configure done, Configuration file path): %s ", AliCredentialFilePath)
//	createCacheDict()
//}

func GetAliCredentialFilePath() string {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	AliCredential := filepath.Join(home, "config.json")
	return AliCredential
}

func GetAllCredentialFilePath() string {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	AliCredential := filepath.Join(home, "config.json")
	return AliCredential
}

func GetAliCredential() cloud.Credential {
	AliCredentialFilePath := GetAliCredentialFilePath()
	var credentials cloud.Credential
	if _, err := os.Stat(AliCredentialFilePath); errors.Is(err, os.ErrNotExist) {
		return credentials
	} else {
		file, err := ioutil.ReadFile(AliCredentialFilePath)
		if err != nil {
			util.HandleErr(err)
		}
		err = json.Unmarshal(file, &credentials)
		if err != nil {
			util.HandleErr(err)
		}
		return credentials
	}
}

func GetAllCredential() cloud.Credential {
	AllCredentialFilePath := GetAllCredentialFilePath()
	var credentials cloud.Credential
	if _, err := os.Stat(AllCredentialFilePath); errors.Is(err, os.ErrNotExist) {
		return credentials
	} else {
		file, err := ioutil.ReadFile(AllCredentialFilePath)
		if err != nil {
			util.HandleErr(err)
		}
		err = json.Unmarshal(file, &credentials)
		if err != nil {
			util.HandleErr(err)
		}
		return credentials
	}
}

func maskAK(ak string) string {
	prefix := ak[:2]
	suffix := ak[len(ak)-6:]
	return prefix + strings.Repeat("*", 18) + suffix
}
