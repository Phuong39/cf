package cmd

import (
	"github.com/spf13/cobra"
	"githubu.com/teamssix/cf/pkg/cloud/alioss"
)

var (
	osslsregion     string
	osslsFlushCache bool
)

func init() {
	RootCmd.AddCommand(ossCmd)
	ossCmd.AddCommand(osslsCmd)
	osslsCmd.Flags().StringVarP(&osslsregion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	ossCmd.PersistentFlags().BoolVar(&osslsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
}

var ossCmd = &cobra.Command{
	Use:   "oss",
	Short: "执行与对象存储相关的操作 (Perform oss-related operations)",
	Long:  "执行与对象存储相关的操作 (Perform oss-related operations)",
}

var osslsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的存储桶 (List all buckets)",
	Long:  "列出所有的存储桶 (List all buckets)",
	Run: func(cmd *cobra.Command, args []string) {
		alioss.PrintBucketsList(osslsregion, osslsFlushCache)
	},
}
