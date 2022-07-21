package tencent

// 腾讯云lh相关操作

import (
	"github.com/spf13/cobra"
	tencentlh2 "github.com/teamssix/cf/pkg/cloud/tencent/tencentlh"
)

var (
	lhFlushCache          bool
	lhRegion              string
	lhSpecifiedInstanceID string
)

func init() {
	tencentCmd.AddCommand(lhCmd)
	lhCmd.AddCommand(lhLsCmd)
	lhCmd.PersistentFlags().BoolVar(&lhFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
	lhCmd.Flags().StringVarP(&lhRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	lhCmd.Flags().StringVarP(&lhSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Set Instance ID)")

	lhLsCmd.Flags().BoolVar(&running, "running", false, "只显示正在运行的实例 (Show only running instances)")
}

var lhCmd = &cobra.Command{
	Use:   "lh",
	Short: "执行与轻量计算服务相关的操作 (Perform lh-related operations)",
	Long:  "执行与轻量计算服务相关的操作 (Perform lh-related operations)",
}

var lhLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的实例 (List all instances)",
	Long:  "列出所有的实例 (List all instances)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentlh2.PrintInstancesList(lhRegion, running, lhSpecifiedInstanceID, lhFlushCache)
	},
}
