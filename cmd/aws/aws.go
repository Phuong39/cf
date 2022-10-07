package aws

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(awsCmd)
}

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "执行与 AWS 相关的操作 (Perform AWS related operations)",
	Long:  "执行与 AWS 相关的操作 (Perform AWS related operations)",
}
