package alibaba

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliram"
)

func init() {
	alibabaCmd.AddCommand(permissionsCmd)
}

var permissionsCmd = &cobra.Command{
	Use:   "permissions",
	Short: "列出当前凭证下所拥有的权限 (List access key permissions)",
	Long:  `列出当前凭证下所拥有的权限 (List access key permissions)`,
	Run: func(cmd *cobra.Command, args []string) {
		aliram.ListPermissions()
	},
}
