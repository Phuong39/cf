package alibaba

import (
	"github.com/spf13/cobra"
	alioss2 "github.com/teamssix/cf/pkg/cloud/alibaba/alioss"
)

var (
	osslsRegion           string
	osslsBucket           string
	osslsFlushCache       bool
	ossdownloadBucket     string
	ossdownloadObject     string
	ossdownloadOutputPath string
)

func init() {
	alibabaCmd.AddCommand(ossCmd)

	ossCmd.AddCommand(osslsCmd)
	osslsCmd.Flags().StringVarP(&osslsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	osslsCmd.Flags().StringVarP(&osslsBucket, "bucket", "b", "all", "列出指定 Bucket 下的对象 (List objects in Bucket)")
	ossCmd.PersistentFlags().BoolVar(&osslsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	ossCmd.AddCommand(ossdownloadCmd)
	ossdownloadCmd.Flags().StringVarP(&ossdownloadBucket, "bucket", "b", "", "指定存储桶 (Set Bucket)")
	ossdownloadCmd.Flags().StringVarP(&ossdownloadObject, "objectKey", "k", "all", "指定对象 (Set object key)")
	ossdownloadCmd.Flags().StringVarP(&ossdownloadOutputPath, "outputPath", "o", "./", "指定导出路径 (Set output path)")
	_ = ossdownloadCmd.MarkFlagRequired("bucket")
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
			alioss2.PrintBucketsList(osslsRegion, osslsFlushCache)
		} else {
			alioss2.PrintObjectsList(osslsBucket)
		}
	},
}

var ossdownloadCmd = &cobra.Command{
	Use:   "get",
	Short: "下载指定的对象 (Download objects)",
	Long:  "下载指定的对象 (Download objects)",
	Run: func(cmd *cobra.Command, args []string) {
		alioss2.DownloadObjects(ossdownloadBucket, ossdownloadObject, ossdownloadOutputPath)
	},
}
