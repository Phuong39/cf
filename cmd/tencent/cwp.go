package tencent

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	tencentcwp2 "github.com/teamssix/cf/pkg/cloud/tencent/tencentcwp"
)

var (
	UUID string
)

func init() {
	tencentCmd.AddCommand(cwpUninstall)
	cwpUninstall.Flags().StringVarP(&UUID, "UUID", "u", "", "指定云镜 UUID (Specify Agent UUID)")
}

var cwpUninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "一键卸载云镜 (Uninstall Agent)",
	Long:  "一键卸载云镜  (Uninstall Agent)",
	Run: func(cmd *cobra.Command, args []string) {
		if UUID == "" {
			log.Warnf("还未指定要卸载云镜的 UUID (The agent-UUID to be uninstall has not been specified yet)\n")
			cmd.Help()
		} else {
			tencentcwp2.UninstallAgent(UUID)
		}

	},
}
