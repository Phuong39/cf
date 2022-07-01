package cmdutil

import (
	"cf/pkg/cloud"
	"cf/pkg/util"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

func ConfigureAccessKey() {
	config := GetAliCredential()
	OldAccessKeyId := ""
	OldAccessKeySecret := ""
	OldSTSToken := ""
	AccessKeyId := config.AccessKeyId
	AccessKeySecret := config.AccessKeySecret
	STSToken := config.STSToken
	if AccessKeyId != "" {
		OldAccessKeyId = " [********************" + AccessKeyId[len(AccessKeyId)-4:] + "] "
	}
	if AccessKeySecret != "" {
		OldAccessKeySecret = " [********************" + AccessKeySecret[len(AccessKeySecret)-4:] + "] "
	}
	if STSToken != "" {
		OldSTSToken = " [********************" + STSToken[len(STSToken)-4:] + "] "
	}
	var qs = []*survey.Question{
		{
			Name:     "AccessKeyId",
			Prompt:   &survey.Input{Message: "Access Key Id (必须 Required)" + OldAccessKeyId + ":"},
			Validate: survey.Required,
		},
		{
			Name:     "AccessKeySecret",
			Prompt:   &survey.Input{Message: "Access Key Secret (必须 Required)" + OldAccessKeySecret + ":"},
			Validate: survey.Required,
		},
		{
			Name:   "STSToken",
			Prompt: &survey.Input{Message: "STS Token (可选 Optional)" + OldSTSToken + ":"},
		},
	}
	cred := cloud.Credential{}
	err := survey.Ask(qs, &cred)
	util.HandleErr(err)
	SaveAccessKey(cred)
}

func SaveAccessKey(config cloud.Credential) {
	home, err := GetCFHomeDir()
	util.HandleErr(err)
	if FileExists(home) == false{
		err = os.MkdirAll(home, 0700)
	}
	util.HandleErr(err)
	configJSON, err := json.MarshalIndent(config, "", "    ")
	util.HandleErr(err)
	AliCredentialFilePath := GetAliCredentialFilePath()
	err = ioutil.WriteFile(AliCredentialFilePath, configJSON, 0600)
	util.HandleErr(err)
	log.Infof("配置完成，配置文件路径 (Configure done, Configuration file path): %s ", AliCredentialFilePath)
	createCacheDict()
}

func GetAliCredentialFilePath() string {
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
	}else{
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
