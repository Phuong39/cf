package aliecs

import (
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func ECSClient(region string) *ecs.Client {
	aliconfig := cmdutil.GetConfig("alibaba")
	if aliconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		config := sdk.NewConfig()
		if aliconfig.STSToken == "" {
			credential := credentials.NewAccessKeyCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret)
			client, err := ecs.NewClientWithOptions(region, config, credential)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
			}
			return client
		} else {
			credential := credentials.NewStsTokenCredential(aliconfig.AccessKeyId, aliconfig.AccessKeySecret, aliconfig.STSToken)
			client, err := ecs.NewClientWithOptions(region, config, credential)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
			}
			return client
		}
	}
}

func GetECSRegions() []ecs.Region {
	client := ECSClient("cn-hangzhou")
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"
	response, err := client.DescribeRegions(request)
	errutil.HandleErr(err)
	return response.Regions.Region
}
