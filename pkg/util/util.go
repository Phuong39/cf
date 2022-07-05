package util

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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
	"InvalidAccessKeyId.NotFound": "当前访问凭证无效 (Current access key are invalid)",
	"Message: The specified parameter \"SecurityToken.Expired\" is not valid.": "当前临时访问凭证已过期 (Current SecurityToken has expired)",
	"Message: The Access Key is disabled.":                                     "当前访问凭证已被禁用 (The Access Key is disabled)",
	"ErrorCode: Forbidden.RAM":                                                 "当前访问凭证权限不足 (No permission)",
}

var errorMessagesNoExit = map[string]string{
	"ErrorCode: Forbidden.RAM":           "当前访问凭证权限不足 (No permission)",
	"ErrorCode: Forbidden":               "当前访问凭证权限不足 (No permission)",
	"You are forbidden to list buckets.": "当前凭证不具备 OSS 的读取权限，无法获取 OSS 数据。 (OSS data is not available because the current credential does not have read access to OSS.)",
}

func HandleErr(e error) {
	if e != nil {
		for k, v := range errorMessages {
			if strings.Contains(e.Error(), k) {
				log.Errorln(v)
				os.Exit(0)
			}
		}
		log.Errorln(e)
	}
}

func HandleErrNoExit(e error) {
	if e != nil {
		for k, v := range errorMessagesNoExit {
			if strings.Contains(e.Error(), k) {
				log.Debugln(v)
			}
		}
	}
}
