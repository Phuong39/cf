package aliecs

import (
	"os"

	"github.com/teamssix/cf/pkg/util/errutil"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func ECSClient(region string) *ecs.Client {
	aliconfig := cmdutil.GetConfig("alibaba")
	if aliconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
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

func GetECSRegions(fullRegions bool) []ecs.Region {
	var ecsRegions []ecs.Region
	client := ECSClient("cn-hangzhou")
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"
	response, err := client.DescribeRegions(request)
	errutil.HandleErr(err)
	ecsRegions = response.Regions.Region
	if fullRegions {
		var privateRegions = map[string]string{
			"cn-shanghai-internal-test-1": "ecs-cn-hangzhou.aliyuncs.com",
			"cn-beijing-gov-1":            "ecs.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "ecs-cn-hangzhou.aliyuncs.com",
			"cn-shanghai-inner":           "ecs.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "ecs-cn-hangzhou.aliyuncs.com",
			"cn-haidian-cm12-c01":         "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "ecs-cn-hangzhou.aliyuncs.com",
			"cn-north-2-gov-1":            "ecs.aliyuncs.com",
			"cn-yushanfang":               "ecs.aliyuncs.com",
			"cn-hongkong-finance-pop":     "ecs.aliyuncs.com",
			"cn-shanghai-finance-1":       "ecs-cn-hangzhou.aliyuncs.com",
			"cn-beijing-finance-pop":      "ecs.aliyuncs.com",
			"cn-wuhan":                    "ecs.aliyuncs.com",
			"cn-zhangbei":                 "ecs.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "ecs.cn-qingdao-nebula.aliyuncs.com",
			"rus-west-1-pop":              "ecs.aliyuncs.com",
			"cn-shanghai-et15-b01":        "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "ecs-cn-hangzhou.aliyuncs.com",
			"eu-west-1-oxs":               "ecs.cn-shenzhen-cloudstone.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "ecs-cn-hangzhou.aliyuncs.com",
			"cn-shenzhen-finance-1":       "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-test-306":        "ecs-cn-hangzhou.aliyuncs.com",
			"cn-huhehaote-nebula-1":       "ecs.cn-qingdao-nebula.aliyuncs.com",
			"cn-shanghai-et2-b01":         "ecs-cn-hangzhou.aliyuncs.com",
			"cn-hangzhou-finance":         "ecs.aliyuncs.com",
			"cn-beijing-nu16-b01":         "ecs-cn-hangzhou.aliyuncs.com",
			"cn-edge-1":                   "ecs.cn-qingdao-nebula.aliyuncs.com",
			"cn-fujian":                   "ecs-cn-hangzhou.aliyuncs.com",
			"ap-northeast-2-pop":          "ecs.aliyuncs.com",
			"cn-shenzhen-inner":           "ecs.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "ecs.cn-zhangjiakou.aliyuncs.com",
		}
		for k, v := range privateRegions {
			ecsRegion := ecs.Region{
				Status:         "",
				RegionEndpoint: v,
				LocalName:      "",
				RegionId:       k,
			}
			ecsRegions = append(ecsRegions, ecsRegion)
		}
	}
	return ecsRegions
}
