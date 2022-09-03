package tencentvpc

import (
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func VPCClient(region string) *vpc.Client {
	tencentconfig := cmdutil.GetConfig("tencent")
	if tencentconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
		if tencentconfig.STSToken == "" {
			credential := common.NewCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeySecret)
			client, err := vpc.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("VPC Client 连接成功 (VPC Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeySecret, tencentconfig.STSToken)
			client, err := vpc.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("VPC Client 连接成功 (VPC Client connection successful)")
			}
			return client
		}
	}
}
