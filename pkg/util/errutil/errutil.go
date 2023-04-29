package errutil

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type error interface {
	Error() string
}

var errorMessages = map[string]string{
	"InvalidAccessKeyId.NotFound": "当前访问密钥无效 (Current access key are invalid)",
	"Message: The specified parameter \"SecurityToken.Expired\" is not valid.": "当前临时访问密钥已过期 (Current SecurityToken has expired)",
	"ErrorCode: InvalidSecurityToken.Expired":                                  "当前临时访问密钥已过期 (Current SecurityToken has expired)",
	"Message: The Access Key is disabled.":                                     "当前访问密钥已被禁用 (The Access Key is disabled)",
	"ErrorCode: Forbidden.RAM":                                                 "当前访问密钥没有执行命令的权限 (Current Access Key do not have permission to execute commands)",
	"ErrorCode: NoPermission":                                                  "当前访问密钥没有接管控制台的权限 (Current Access Key do not have permission to take over the console)",
	"ErrorCode=NoSuchKey":                                                      "存储桶中没有这个对象 (There is no such key in the bucket)",
	"Code=ResourceNotFound, Message=未查询到对应机器":                                  "指定资源不存在 (Resource not found)",
	//"Code=UnauthorizedOperation":                                               "当前 AK 权限不足 (Insufficient Access Key permissions)",
	"you are not authorized to perform operation (tat:CreateCommand)": "当前 AK 不具备执行命令的权限 (This Access Key does not have permission to execute commands)",
	"network is unreachable":       "当前网络连接异常 (Network is unreachable)",
	"InvalidSecurityToken.Expired": "临时令牌已过期 (STS token has expired)",
	"InvalidAccessKeyId.Inactive":  "当前 AK 已被禁用 (The current AccessKeyId is inactive)",
	"interrupt":                    "程序已退出 (Program exited.)",
	"ErrorCode=AccessDenied, ErrorMessage=\"The bucket you access does not belong to you.\"": "获取 Bucket 信息失败，访问被拒绝 (Failed to get Bucket information, access is denied.)",
	"ExpiredToken":                                                        "当前访问密钥已过期 (Current token has expired)",
	"read: connection reset by peer":                                      "网络连接出现错误，请检查您的网络环境是否正常 (There is an error in your network connection, please check if your network environment is normal.)",
	"Code=ResourceUnavailable.AgentNotInstalled":                          "Agent 未安装 (Agent not installed)",
	"Incorrect IAM authentication information":                            "当前 AK 信息无效 (Current AccessKey information is invalid)",
	"The API does not exist or has not been published in the environment": "当前用户已存在，请指定其他用户名 (User already exists, please specify another user name)",
	"Status=403 Forbidden, Code=AccessDenied":                             "当前权限不足 (Insufficient permissions)",
}

var errorMessagesNoExit = map[string]string{
	"ErrorCode: Forbidden.RAM": "当前访问密钥没有执行命令的权限 (Current Access Key do not have permission to execute commands)",
	//"ErrorCode: Forbidden":                                               " 当前访问密钥没有 RDS 的读取权限 (Current Access Key do not have read access to RDS"),
	"You are forbidden to list buckets.":                                 "当前凭证不具备 OSS 的读取权限，无法获取 OSS 数据。 (OSS data is not available because the current credential does not have read access to OSS.)",
	"ErrorCode: EntityAlreadyExists.User.Policy":                         "已接管过控制台，无需重复接管 (Console has been taken over)",
	"ErrorCode: EntityAlreadyExists.User":                                "已接管过控制台，无需重复接管 (Console has been taken over)",
	"ErrorCode: EntityNotExist.User":                                     "已取消接管控制台，无需重复取消 (Console has been de-taken over)",
	"Code=ResourceNotFound, Message=指定资源":                                "指定资源不存在 (ResourceNotFound)",
	"InvalidParameter.SubUserNameInUse":                                  "已接管过控制台，无需重复接管 (Console has been taken over)",
	"you are not authorized to perform operation (cwp:DescribeMachines)": "当前 AK 没有 CWP 权限",
}

var errorMessagesExit = map[string]string{
	"ErrorCode: Forbidden.RAM":     "当前访问密钥没有执行命令的权限 (Current Access Key do not have permission to execute commands)",
	"ErrorCode: NoPermission":      "当前访问密钥没有接管控制台的权限 (Current Access Key do not have permission to take over the console)",
	"network is unreachable":       "当前网络连接异常 (Network is unreachable)",
	"InvalidSecurityToken.Expired": "临时令牌已过期 (STS token has expired)",
	"InvalidAccessKeyId.Inactive":  "当前 AK 已被禁用 (The current AccessKeyId is inactive)",
	//"Message=操作未授权，请检查CAM策略。":  "当前 AK 权限不足 (Insufficient Access Key permissions)",
	"Code=AuthFailure.SecretIdNotFound": "SecretId 不存在，请输入正确的密钥 (SecretId does not exist, please enter the correct key.)",
	"Code=AuthFailure.SignatureFailure": "请求签名验证失败，请检查您的访问密钥是否正确 (Request signature verification failed, please check if your access key is correct.)",
	"read: connection reset by peer":    "网络连接出现错误，请检查您的网络环境是否正常 (There is an error in your network connection, please check if your network environment is normal.)",
	"InvalidAccessKeyId.NotFound":       "当前访问密钥无效 (Current access key are invalid)",
}

func HandleErr(e error) {
	if e != nil {
		log.Traceln(e.Error())
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
		log.Traceln(e.Error())
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
