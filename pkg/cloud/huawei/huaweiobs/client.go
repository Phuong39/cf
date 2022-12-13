package huaweiobs

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"
)

func obsClient(region string) *obs.ObsClient {
	var (
		obsClient *obs.ObsClient
		err       error
	)
	config := cmdutil.GetConfig("huawei")
	if config.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		if region == "all" {
			region = "cn-north-1"
		}
		if config.STSToken == "" {
			obsClient, err = obs.New(config.AccessKeyId, config.AccessKeySecret, "https://obs."+region+".myhuaweicloud.com")
		} else {
			obsClient, err = obs.New(config.AccessKeyId, config.AccessKeySecret, "https://obs."+region+".myhuaweicloud.com", obs.WithSecurityToken(config.STSToken))
		}
		if err == nil {
			log.Traceln("obs Client 连接成功 (obs Client connection successful)")
		} else {
			errutil.HandleErr(err)
		}
		return obsClient
	}
}
