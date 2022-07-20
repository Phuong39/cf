package keymanage

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/cmd"
	"github.com/teamssix/cf/cmd/keymanage/keystore"
)

func init() {
	// init keyStore Service
	keystore.Init()
	// init Commands
	KeyManagerRoot.AddCommand(AddKeyCmd)
	KeyManagerRoot.AddCommand(DelKeyCmd)
	KeyManagerRoot.AddCommand(HeadKeyCmd)
	KeyManagerRoot.AddCommand(ListKeyCmd)
	KeyManagerRoot.AddCommand(SwitchKeyCmd)
	// add to root command
	cmd.RootCmd.AddCommand(KeyManagerRoot)
	LoadKeys()
	GetHeader()
}

var KeyManagerRoot = &cobra.Command{
	Use:   "key",
	Short: "AKSK/STS 统一管理/切换模块 (multi AK-SK switch and management)",
	Long:  "AKSK/STS 统一管理/切换模块 (multi AK-SK switch and management)",
}
