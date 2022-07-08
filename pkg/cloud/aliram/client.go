package aliram

import (
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func RAMClient() *ram.Client {
	aliconfig := cmdutil.GetAliCredential()
	if aliconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		config := sdk.NewConfig()
		if aliconfig.STSToken == "" {
			credential := credentials.NewAccessKeyCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret)
			client, err := ram.NewClientWithOptions("cn-beijing", config, credential)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RDS Client connection successful)")
			}
			return client
		} else {
			credential := credentials.NewStsTokenCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret, aliconfig.STSToken)
			client, err := ram.NewClientWithOptions("cn-beijing", config, credential)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RDS Client connection successful)")
			}
			return client
		}
	}
}

func STSClient() *sts.Client {
	aliconfig := cmdutil.GetAliCredential()
	if aliconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		config := sdk.NewConfig()
		if aliconfig.STSToken == "" {
			credential := credentials.NewAccessKeyCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret)
			client, err := sts.NewClientWithOptions("cn-beijing", config, credential)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RDS Client connection successful)")
			}
			return client
		} else {
			credential := credentials.NewStsTokenCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret, aliconfig.STSToken)
			client, err := sts.NewClientWithOptions("cn-beijing", config, credential)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("RAM Client 连接成功 (RDS Client connection successful)")
			}
			return client
		}
	}
}
