package aliram

import (
	"os"

	"github.com/teamssix/cf/pkg/util/errutil"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func RAMClient() *ram.Client {
	// 判断是否已经配置了访问密钥
	aliconfig := cmdutil.GetConfig("alibaba")
	if aliconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		// 判断是否已经配置了 STS Token
		config := sdk.NewConfig()
		if aliconfig.STSToken == "" {
			credential := credentials.NewAccessKeyCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret)
			client, err := ram.NewClientWithOptions("cn-beijing", config, credential)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RAM Client connection successful)")
			} else {
				errutil.HandleErr(err)
			}
			return client
		} else {
			// 使用 STS Token 连接
			credential := credentials.NewStsTokenCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret, aliconfig.STSToken)
			client, err := ram.NewClientWithOptions("cn-beijing", config, credential)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RAM Client connection successful)")
			} else {
				errutil.HandleErr(err)
			}
			return client
		}
	}
}

func STSClient() *sts.Client {
	// 判断是否已经配置了访问密钥
	aliconfig := cmdutil.GetConfig("alibaba")
	if aliconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		// 判断是否已经配置了 STS Token
		config := sdk.NewConfig()
		if aliconfig.STSToken == "" {
			credential := credentials.NewAccessKeyCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret)
			client, err := sts.NewClientWithOptions("cn-beijing", config, credential)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RAM Client connection successful)")
			}
			return client
		} else {
			// 使用 STS Token 连接
			credential := credentials.NewStsTokenCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret, aliconfig.STSToken)
			client, err := sts.NewClientWithOptions("cn-beijing", config, credential)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RAM Client connection successful)")
			}
			return client
		}
	}
}
