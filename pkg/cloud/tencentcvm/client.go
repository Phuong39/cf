package tencentcvm

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"os"
)

//腾讯云主机结构
type Instances struct {
	InstanceId       string
	OSName           string
	InstanceType     string
	InstanceState    string
	PrivateIpAddress string
	PublicIpAddress  string
	Zone             string
}

func CVMClient(region string) *cvm.Client {
	tencentconfig := cmdutil.GetAllCredential()
	if tencentconfig.Tencent.TmpSecretId == "" && tencentconfig.Tencent.TmpSecretKey == "" && tencentconfig.Tencent.Token == "" {
		if tencentconfig.Tencent.SecretId == "" && tencentconfig.Tencent.SecretKey == "" {
			log.Warnln("需要先配置访问凭证 (Secret Id need to be configured first)")
			os.Exit(0)
			return nil
		}
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
		credential := common.NewCredential(tencentconfig.Tencent.SecretId, tencentconfig.Tencent.SecretKey)
		client, err := cvm.NewClient(credential, region, cpf)
		util.HandleErr(err)
		if err == nil {
			log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
		}
		return client
	} else {
		//临时密钥
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
		credential := common.NewTokenCredential(tencentconfig.Tencent.SecretId, tencentconfig.Tencent.SecretKey, tencentconfig.Tencent.Token)
		client, err := cvm.NewClient(credential, region, cpf)
		util.HandleErr(err)
		if err == nil {
			log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
		}
		return client
	}
}

//执行命令接口
func TATClient(region string) *tat.Client {
	tencentconfig := cmdutil.GetAllCredential()
	if tencentconfig.Tencent.TmpSecretId == "" && tencentconfig.Tencent.TmpSecretKey == "" && tencentconfig.Tencent.Token == "" {
		if tencentconfig.Tencent.SecretId == "" && tencentconfig.Tencent.SecretKey == "" {
			log.Warnln("需要先配置访问凭证 (Secret Id need to be configured first)")
			os.Exit(0)
			return nil
		}
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "tat.tencentcloudapi.com"
		credential := common.NewCredential(tencentconfig.Tencent.SecretId, tencentconfig.Tencent.SecretKey)
		client, err := tat.NewClient(credential, region, cpf)
		util.HandleErr(err)
		if err == nil {
			log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
		}
		return client
	} else {
		//临时密钥
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "tat.tencentcloudapi.com"
		credential := common.NewTokenCredential(tencentconfig.Tencent.SecretId, tencentconfig.Tencent.SecretKey, tencentconfig.Tencent.Token)
		client, err := tat.NewClient(credential, region, cpf)
		util.HandleErr(err)
		if err == nil {
			log.Traceln("ECS Client 连接成功 (ECS Client connection successful)")
		}
		return client
	}
}

func GetCVMRegions() []*cvm.RegionInfo {
	client := CVMClient(regions.Nanjing)
	request := cvm.NewDescribeRegionsRequest()
	request.SetScheme("https")
	response, err := client.DescribeRegions(request)
	util.HandleErr(err)
	return response.Response.RegionSet
}
