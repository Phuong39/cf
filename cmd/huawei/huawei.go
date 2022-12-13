package huawei

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(huaweiCmd)
}

var huaweiCmd = &cobra.Command{
	Use:   "huawei",
	Short: "执行与华为云相关的操作 (Perform Huawei Cloud related operations)",
	Long:  "执行与华为云相关的操作 (Perform Huawei Cloud related operations)",
}
