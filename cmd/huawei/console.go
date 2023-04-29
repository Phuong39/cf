package huawei

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/cloudpub"
	"github.com/teamssix/cf/pkg/cloud/huawei/huaweiconsole"
)

var (
	userName string
	password string
)

func init() {
	huaweiCmd.AddCommand(consoleCmd)
	consoleCmd.AddCommand(cancelConsoleCmd)
	consoleCmd.AddCommand(lsConsoleCmd)

	consoleCmd.Flags().StringVarP(&userName, "userName", "u", "crossfire", "指定用户名 (Specify user name)")
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "一键接管控制台 (Takeover console)",
	Long:  "一键接管控制台 (Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		huaweiconsole.TakeoverConsole(userName)
	},
}

var cancelConsoleCmd = &cobra.Command{
	Use:   "cancel",
	Short: "取消接管控制台 (Cancel Takeover console)",
	Long:  "取消接管控制台 (Cancel Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		huaweiconsole.CancelTakeoverConsole()
	},
}

var lsConsoleCmd = &cobra.Command{
	Use:   "ls",
	Short: "查看接管控制台的信息 (View Takeover console information)",
	Long:  "查看接管控制台的信息 (View Takeover console information)",
	Run: func(cmd *cobra.Command, args []string) {
		cloudpub.LsTakeoverConsole("huawei")
	},
}
