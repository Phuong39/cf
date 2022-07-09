package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alioss"
)

var (
	osslsRegion     string
	osslsBucket     string
	osslsFlushCache bool
)

func init() {
	RootCmd.AddCommand(ossCmd)
	ossCmd.AddCommand(osslsCmd)
	osslsCmd.Flags().StringVarP(&osslsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	osslsCmd.Flags().StringVarP(&osslsBucket, "bucket", "b", "all", "列出指定 Bucket 下的对象 (List objects in Bucket)")
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
		if osslsBucket == "all" {
			alioss.PrintBucketsList(osslsRegion, osslsFlushCache)
		} else {
			alioss.PrintObjectsList(osslsBucket)
		}
	},
}
