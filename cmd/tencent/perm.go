package tencent

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcam"
)

func init() {
	tencentCmd.AddCommand(permCmd)
}

var permCmd = &cobra.Command{
	Use:   "perm",
	Short: "列出当前凭证下所拥有的权限 (List access key permissions)",
	Long:  `列出当前凭证下所拥有的权限 (List access key permissions)`,
	Run: func(cmd *cobra.Command, args []string) {
		tencentcam.ListPermissions()
	},
}
