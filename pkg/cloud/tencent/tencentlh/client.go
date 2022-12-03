package tencentlh

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"os"
)

func LHClient(region string) *lh.Client {
	tencentConfig := cmdutil.GetConfig("tencent")
	if tencentConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问密钥 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
		if tencentConfig.STSToken == "" {
			credential := common.NewCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret)
			client, err := lh.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("LH Client 连接成功 (LH Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret, tencentConfig.STSToken)
			client, err := lh.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("LH Client 连接成功 (LH Client connection successful)")
			}
			return client
		}
	}
}

func GetLHRegions() []*lh.RegionInfo {
	client := LHClient("ap-guangzhou")
	request := lh.NewDescribeRegionsRequest()
	response, err := client.DescribeRegions(request)
	errutil.HandleErr(err)
	return response.Response.RegionSet
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
			credential := common.NewTokenCredential(tencentconfig.AccessKeyId, tencentconfig.AccessKeyId, tencentconfig.STSToken)
			client, err := tat.NewClient(credential, region, cpf)
			errutil.HandleErr(err)
			if err == nil {
				log.Traceln("TAT Client 连接成功 (CVM Client connection successful)")
			}
			return client
		}
	}
}
