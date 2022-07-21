package tencentcwp

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	"os"
)

func CWPClient(region string) *cwp.Client {
	tencentConfig := cmdutil.GetConfig("tencent")
	if tencentConfig.AccessKeyId == "" {
		log.Warnln("需要先配置访问凭证 (Access Key need to be configured first)")
		os.Exit(0)
		return nil
	} else {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cwp.tencentcloudapi.com"
		if tencentConfig.STSToken == "" {
			credential := common.NewCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeySecret)
			client, err := cwp.NewClient(credential, region, cpf)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("CWP Client 连接成功 (CWP Client connection successful)")
			}
			return client
		} else {
			credential := common.NewTokenCredential(tencentConfig.AccessKeyId, tencentConfig.AccessKeyId, tencentConfig.STSToken)
			client, err := cwp.NewClient(credential, region, cpf)
			util.HandleErr(err)
			if err == nil {
				log.Traceln("CWP Client 连接成功 (CWP Client connection successful)")
			}
			return client
		}
	}
}

func DescribeMachineCWPStatus(MachineType string, Quuid string) (*string, *string) {
	client := CWPClient("")
	request := cwp.NewDescribeMachinesRequest()
	request.MachineType = common.StringPtr(MachineType)
	request.Filters = []*cwp.Filter{
		{
			Name:   common.StringPtr("Quuid"),
			Values: common.StringPtrs([]string{Quuid}),
		},
	}
	request.MachineRegion = common.StringPtr("all-regions")
	response, err := client.DescribeMachines(request)
	util.HandleErr(err)
	return response.Response.Machines[0].MachineStatus, response.Response.Machines[0].Uuid

}
