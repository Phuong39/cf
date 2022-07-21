package tencentlh

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"os"
)

func LHClient(region string) *lh.Client {
	tencentConfig := cmdutil.GetConfig("tencent")
	if tencentConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
		if tencentConfig.STSToken == "" {
			credential := common.NewCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret)
			client, err := lh.NewClient(credential, region, cpf)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("LH Client 连接成功 (LH Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeyId, tencentConfig.STSToken)
			client, err := lh.NewClient(credential, region, cpf)
			util.HandleErr(err)
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
	util.HandleErr(err)
	return response.Response.RegionSet
}
