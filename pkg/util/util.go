package util

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

const AppDirName = ".cf"

const (
	CFCredentialPathEnvVar = "CF_CREDENTIAL_PATH"
	CFHomeEnvVar           = "CF_HOME"
	CFCloudTokenEnvVar     = "CF_CLOUD_TOKEN"
)

type error interface {
	Error() string
}

var errorMessages = map[string]string{
	"InvalidAccessKeyId.NotFound":"当前访问凭证无效 (Current access key are invalid)",
	"Message: The specified parameter \"SecurityToken.Expired\" is not valid.":"当前临时访问凭证已过期 (Current SecurityToken has expired)",
	"Message: The Access Key is disabled.":"当前访问凭证已被禁用 (The Access Key is disabled)",
}

func HandleErr(e error) {
	if e != nil {
		for k,v := range errorMessages{
			if strings.Contains(e.Error(),k){
				log.Errorf(v)
				os.Exit(0)
			}
		}
		log.Errorln(e)
	}
}