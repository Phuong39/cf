package keymanage

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

// Key struct store the KeyPairs.
type Key struct {
	Name          string // KeyName to help Red-teamer the Key.
	Platform      string // Which Cloud Service provider
	*cloud.Config        // AK-SK config
	Remark        string // remarks.
}

var HeaderKey = []Key{}

var CommonTableHeader = []string{
	"名称 (Name)", "所属平台 (Platform)",
	"AccessKeyId", "AccessKeySecret",
	"STSToken", "备注 (Remark)",
}

// GetHeader Get the Current Key config for all Cloud Service Provider.
func GetHeader() {
	cloudConfigList, _ := cmdutil.ReturnCloudProviderList()
	for _, provider := range cloudConfigList {
		config := cmdutil.GetConfig(provider)
		AccessKeyId := config.AccessKeyId
		if AccessKeyId == "" {
			log.Infof("当前未配置平台 %s 访问密钥 (No access key configured)", provider)
		} else {
			HeaderKey = append(HeaderKey, Key{
				Name:     "Current(当前)",
				Platform: provider,
				Config:   &config,
				Remark:   "当前配置文件中所设置的访问密钥",
			})
		}
	}
}
