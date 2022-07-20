package keymanage

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/cmd"
	"github.com/teamssix/cf/cmd/keymanage/keyop"
	"github.com/teamssix/cf/cmd/keymanage/keystore"
)

func init() {
	// init keyStore Service
	keystore.Init()
	// init Commands
	KeyManagerRoot.AddCommand(keyop.AddKeyCmd)
	KeyManagerRoot.AddCommand(keyop.DelKeyCmd)
	KeyManagerRoot.AddCommand(keyop.HeadKeyCmd)
	KeyManagerRoot.AddCommand(keyop.ListKeyCmd)
	KeyManagerRoot.AddCommand(keyop.SwitchKeyCmd)
	cmd.RootCmd.AddCommand(KeyManagerRoot)
}

var KeyManagerRoot = &cobra.Command{
	Use:   "key",
	Short: "AKSK/STS 统一管理/切换模块 (multi AK-SK switch and management)",
	Long:  "AKSK/STS 统一管理/切换模块 (multi AK-SK switch and management)",
}
