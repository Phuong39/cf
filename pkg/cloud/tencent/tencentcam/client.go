package tencentcam

import (
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func CAMClient() *cam.Client {
	tencentConfig := cmdutil.GetConfig("tencent")
	if tencentConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cam.tencentcloudapi.com"
		if tencentConfig.STSToken == "" {
			credential := common.NewCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret)
			client, err := cam.NewClient(credential, "", cpf)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("CAM Client 连接成功 (CAM Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret, tencentConfig.STSToken)
			client, err := cam.NewClient(credential, "", cpf)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("CAM Client 连接成功 (CAM Client connection successful)")
			}
			return client
		}
	}
}
