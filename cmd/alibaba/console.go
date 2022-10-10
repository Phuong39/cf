package alibaba

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliconsole"
	"github.com/teamssix/cf/pkg/cloud/cloudpub"
)

var (
	userName string
	password string
)

func init() {
	alibabaCmd.AddCommand(consoleCmd)
	consoleCmd.AddCommand(cancelConsoleCmd)
	consoleCmd.AddCommand(lsConsoleCmd)

	consoleCmd.Flags().StringVarP(&userName, "userName", "u", "crossfire", "指定用户名 (Specify user name)")
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "一键接管控制台 (Takeover console)",
	Long:  "一键接管控制台 (Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		aliconsole.TakeoverConsole(userName)
	},
}

var cancelConsoleCmd = &cobra.Command{
	Use:   "cancel",
	Short: "取消接管控制台 (Cancel Takeover console)",
	Long:  "取消接管控制台 (Cancel Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		aliconsole.CancelTakeoverConsole()
	},
}

var lsConsoleCmd = &cobra.Command{
	Use:   "ls",
	Short: "查看接管控制台的信息 (View Takeover console information)",
	Long:  "查看接管控制台的信息 (View Takeover console information)",
	Run: func(cmd *cobra.Command, args []string) {
		cloudpub.LsTakeoverConsole("alibaba")
	},
}
