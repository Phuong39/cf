package keymanage

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/cmd/keymanage/keystore"
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

var KeyChain = []Key{}

var HeaderKey = []Key{}

func LoadKeys() error {
	// ToDo: read cf key storage file.
	err := keystore.KeyConfig.UnmarshalKey("keys", &KeyChain)
	if err != nil {
		log.Error("加载密钥配置出错 (Loading Key Config Error)", err)
	}
	return err
}

func GetHeader() {
	_, cloudProviderList := cmdutil.ReturnCloudProviderList()
	for i, provider := range cloudProviderList {
		config := cmdutil.GetConfig(provider)
		AccessKeyId := config.AccessKeyId
		if AccessKeyId == "" {
			log.Infof("当前未配置平台 %s 访问密钥 (No access key configured)", provider)
		} else {
			HeaderKey[i] = Key{
				Name:     "Current(当前)",
				Platform: provider,
				Config:   &config,
				Remark:   "当前配置文件中所设置的访问密钥",
			}
		}
	}
}

func SaveKeys() error {
	// ToDo: Set cf key storage file.
	err := keystore.KeyConfig.SafeWriteConfig()
	if err != nil {
		log.Error("写入密钥配置出错 (Writing Key Config Error)", err)
	}
	return err
}
