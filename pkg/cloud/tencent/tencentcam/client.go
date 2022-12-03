package tencentcam

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
	"os"
)

func CAMClient() *cam.Client {
	tencentConfig := cmdutil.GetConfig("tencent")
	if tencentConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cam.tencentcloudapi.com"
		if tencentConfig.STSToken == "" {
			credential := common.NewCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret)
			client, err := cam.NewClient(credential, "", cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("CAM Client 连接成功 (CAM Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret, tencentConfig.STSToken)
			client, err := cam.NewClient(credential, "", cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("CAM Client 连接成功 (CAM Client connection successful)")
			}
			return client
		}
	}
}

func STSClient() *sts.Client {
	tencentConfig := cmdutil.GetConfig("tencent")
	if tencentConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "sts.tencentcloudapi.com"
		if tencentConfig.STSToken == "" {
			credential := common.NewCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret)
			client, err := sts.NewClient(credential, "ap-beijing", cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("STS Client 连接成功 (STS Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret, tencentConfig.STSToken)
			client, err := sts.NewClient(credential, "ap-beijing", cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("STS Client 连接成功 (STS Client connection successful)")
			}
			return client
		}
	}
}
