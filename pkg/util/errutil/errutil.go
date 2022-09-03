package errutil

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type error interface {
	Error() string
}

var errorMessages = map[string]string{
	"InvalidAccessKeyId.NotFound": "当前访问凭证无效 (Current access key are invalid)",
	"Message: The specified parameter \"SecurityToken.Expired\" is not valid.": "当前临时访问凭证已过期 (Current SecurityToken has expired)",
	"ErrorCode: InvalidSecurityToken.Expired":                                  "当前临时访问凭证已过期 (Current SecurityToken has expired)",
	"Message: The Access Key is disabled.":                                     "当前访问凭证已被禁用 (The Access Key is disabled)",
	"ErrorCode: Forbidden.RAM":                                                 "当前访问凭证没有执行命令的权限 (Current Access Key do not have permission to execute commands)",
	"ErrorCode: NoPermission":                                                  "当前访问凭证没有接管控制台的权限 (Current Access Key do not have permission to take over the console)",
	"ErrorCode=NoSuchKey":                                                      "存储桶中没有这个对象 (There is no such key in the bucket)",
	"Code=ResourceNotFound, Message=未查询到对应机器":                                  "指定资源不存在 (ResourceNotFound)",
	//"Code=UnauthorizedOperation":                                               "当前 AK 权限不足 (Insufficient Access Key permissions)",
	"you are not authorized to perform operation (tat:CreateCommand)": "当前 AK 不具备执行命令的权限 (This Access Key does not have permission to execute commands)",
}

var errorMessagesNoExit = map[string]string{
	"ErrorCode: Forbidden.RAM":                                           "当前访问凭证没有执行命令的权限 (Current Access Key do not have permission to execute commands)",
	"ErrorCode: Forbidden":                                               "当前访问凭证没有 RDS 的读取权限 Current Access Key do not have read access to RDS",
	"You are forbidden to list buckets.":                                 "当前凭证不具备 OSS 的读取权限，无法获取 OSS 数据。 (OSS data is not available because the current credential does not have read access to OSS.)",
	"ErrorCode: EntityAlreadyExists.User.Policy":                         "已接管过控制台，无需重复接管 (Console has been taken over)",
	"ErrorCode: EntityAlreadyExists.User":                                "已接管过控制台，无需重复接管 (Console has been taken over)",
	"ErrorCode: EntityNotExist.User":                                     "已取消接管控制台，无需重复取消 (Console has been de-taken over)",
	"Code=ResourceNotFound, Message=指定资源":                                "指定资源不存在 (ResourceNotFound)",
	"InvalidParameter.SubUserNameInUse":                                  "已接管过控制台，无需重复接管 (Console has been taken over)",
	"you are not authorized to perform operation (cwp:DescribeMachines)": "当前 AK 没有 CWP 权限",
}

var errorMessagesExit = map[string]string{
	"ErrorCode: Forbidden.RAM": "当前访问凭证没有执行命令的权限 (Current Access Key do not have permission to execute commands)",
	"ErrorCode: NoPermission":  "当前访问凭证没有接管控制台的权限 (Current Access Key do not have permission to take over the console)",
	//"Message=操作未授权，请检查CAM策略。":  "当前 AK 权限不足 (Insufficient Access Key permissions)",
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
		for k, v := range errorMessagesExit {
			if strings.Contains(e.Error(), k) {
				log.Errorln(v)
				os.Exit(0)
			}
		}
	}
}
