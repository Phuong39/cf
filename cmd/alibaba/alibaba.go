package alibaba

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(alibabaCmd)
}

var alibabaCmd = &cobra.Command{
	Use:   "alibaba",
	Short: "执行与阿里云相关的操作 (Perform Alibaba Cloud related operations)",
	Long:  "执行与阿里云相关的操作 (Perform Alibaba Cloud related operations)",
}
