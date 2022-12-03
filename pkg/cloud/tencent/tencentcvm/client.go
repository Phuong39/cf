package tencentcvm

import (
	"github.com/teamssix/cf/pkg/util/errutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
)

func CVMClient(region string) *cvm.Client {
	tencentconfig := cmdutil.GetConfig("tencent")
	if tencentconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
		if tencentconfig.STSToken == "" {
			credential := common.NewCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeySecret)
			client, err := cvm.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("CVM Client 连接成功 (CVM Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeySecret, tencentconfig.STSToken)
			client, err := cvm.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("CVM Client 连接成功 (CVM Client connection successful)")
			}
			return client
		}
	}
}

func TATClient(region string) *tat.Client {
	tencentconfig := cmdutil.GetConfig("tencent")
	if tencentconfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "tat.tencentcloudapi.com"
		if tencentconfig.STSToken == "" {
			credential := common.NewCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeySecret)
			client, err := tat.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("TAT Client 连接成功 (CVM Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeySecret, tencentconfig.STSToken)
			client, err := tat.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("TAT Client 连接成功 (CVM Client connection successful)")
			}
			return client
		}
	}
}

func GetCVMRegions() []*cvm.RegionInfo {
	client := CVMClient(regions.Nanjing)
	request := cvm.NewDescribeRegionsRequest()
	request.SetScheme("https")
	response, err := client.DescribeRegions(request)
	errutil.HandleErr(err)
	return response.Response.RegionSet
}
