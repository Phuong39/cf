package alibaba

import (
	"github.com/spf13/cobra"
	aliram2 "github.com/teamssix/cf/pkg/cloud/alibaba/aliram"
)

func init() {
	alibabaCmd.AddCommand(consoleCmd)
	consoleCmd.AddCommand(cancelConsoleCmd)
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "一键接管控制台 (Takeover console)",
	Long:  "一键接管控制台 (Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		aliram2.TakeoverConsole()
	},
}

var cancelConsoleCmd = &cobra.Command{
	Use:   "cancel",
	Short: "取消接管控制台 (Cancel Takeover console)",
	Long:  "取消接管控制台 (Cancel Takeover console)",
	Run: func(cmd *cobra.Command, args []string) {
		aliram2.CancelTakeoverConsole()
	},
}
