package keystore

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var KeyConfig *viper.Viper

func Init() {
	KeyConfig = viper.New()
	KeyConfig.SetConfigName("keymanage")
	KeyConfig.SetConfigType("yaml")
	KeyConfig.AddConfigPath("$HOME/.cf/config")
	KeyConfig.SetDefault("enable", true)
	KeyConfig.SafeWriteConfig() // Create but no replace config file when not Exist
	err := KeyConfig.ReadInConfig()
	if err != nil {
		logrus.WithField("config", "KeyConfig").WithError(err).Panicf("unable to read config file")
	}

	KeyConfig.WatchConfig() // 自动更新配置
	KeyConfig.OnConfigChange(func(e fsnotify.Event) {
		err := KeyConfig.ReadInConfig()
		if err == nil {
			logrus.WithField("config", "KeyConfig").Info("config updated")
		}
	})
}
