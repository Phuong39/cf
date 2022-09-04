package tencent

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/console"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentconsole"
)

func init() {
	tencentCmd.AddCommand(consoleCmd)
	consoleCmd.AddCommand(cancelConsoleCmd)
	consoleCmd.AddCommand(lsConsoleCmd)
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "一键接管控制台 (Takeover console)",
	Long:  "一键接管控制台 (Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentconsole.TakeoverConsole()
	},
}

var cancelConsoleCmd = &cobra.Command{
	Use:   "cancel",
	Short: "取消接管控制台 (Cancel Takeover console)",
	Long:  "取消接管控制台 (Cancel Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentconsole.CancelTakeoverConsole()
	},
}

var lsConsoleCmd = &cobra.Command{
	Use:   "ls",
	Short: "查看接管控制台的信息 (View Takeover console information)",
	Long:  "查看接管控制台的信息 (View Takeover console information)",
	Run: func(cmd *cobra.Command, args []string) {
		console.LsTakeoverConsole("tencent")
	},
}
