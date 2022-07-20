package keymanage

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/cmd/keymanage/keystore"
	"github.com/teamssix/cf/pkg/cloud"
)

// Key struct store the KeyPairs.
type Key struct {
	Name          string // KeyName to help Red-teamer the Key.
	Platform      string // Which Cloud Service provider
	*cloud.Config        // AK-SK config
	Remark        string // remarks.
}

var KeyChain = []Key{}

func LoadKeys() error {
	// ToDo: read cf key storage file.
	err := keystore.KeyConfig.UnmarshalKey("keys", &KeyChain)
	if err != nil {
		log.Error("加载密钥配置出错 (Loading Key Config Error)", err)
	}
	return err
}

func SaveKeys() error {
	// ToDo: Set cf key storage file.
	err := keystore.KeyConfig.SafeWriteConfig()
	if err != nil {
		log.Error("写入密钥配置出错 (Writing Key Config Error)", err)
	}
	return err
}
