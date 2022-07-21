package keymanage

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/cmd"
	"github.com/teamssix/cf/cmd/database"
	"gorm.io/gorm"
)

var KeyDb *gorm.DB

func init() {
	// init Commands
	KeyManagerRoot.AddCommand(AddKeyCmd)
	KeyManagerRoot.AddCommand(DelKeyCmd)
	KeyManagerRoot.AddCommand(HeadKeyCmd)
	KeyManagerRoot.AddCommand(ListKeyCmd)
	KeyManagerRoot.AddCommand(SwitchKeyCmd)
	// add to root command
	cmd.RootCmd.AddCommand(KeyManagerRoot)
	// Do Some prepare Loading keys in Config and local db
	GetHeader() // get current config in aliyun/tencent/... config PATH
	err := database.GlobalDataBase.MainDB.AutoMigrate(&Key{})
	if err != nil {
		log.Panic("数据库自动配置失败 ( Database AutoMigrate Key Struct failure )", err)
	}
	KeyDb = database.GlobalDataBase.MainDB
}

var KeyManagerRoot = &cobra.Command{
	Use:   "key",
	Short: "AKSK/STS 统一管理/切换模块 (multi AK-SK switch and management)",
	Long:  "AKSK/STS 统一管理/切换模块 (multi AK-SK switch and management)",
}
