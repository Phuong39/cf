package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(ConfigLs)
	configCmd.AddCommand(ConfigSw)
	configCmd.AddCommand(ConfigDel)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置云服务商的访问密钥 (Configure cloud provider access key)",
	Long:  `配置云服务商的访问密钥 (Configure cloud provider access key)`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdutil.ConfigureAccessKey()
	},
}

var ConfigLs = &cobra.Command{
	Use:   "ls",
	Short: "获取已配置过的访问凭证 (Get configured access key)",
	Long:  `获取已配置过的访问凭证 (Get configured access key)`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdutil.ConfigLs()
	},
}

var ConfigSw = &cobra.Command{
	Use:   "sw",
	Short: "切换访问凭证 (Switch access key)",
	Long:  `切换访问凭证 (Switch access key)`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdutil.ConfigSw()
	},
}

var ConfigDel = &cobra.Command{
	Use:   "del",
	Short: "删除访问凭证 (Delete access key)",
	Long:  `删除访问凭证 (Delete access key)`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdutil.ConfigDel()
	},
}
