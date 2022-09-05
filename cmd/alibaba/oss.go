package alibaba

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alioss"
)

var (
	ossLsRegion           string
	ossLsBucket           string
	ossLsFlushCache       bool
	ossDownloadBucket     string
	ossDownloadObject     string
	ossDownloadOutputPath string
)

func init() {
	alibabaCmd.AddCommand(ossCmd)
	ossCmd.AddCommand(ossLsCmd)
	ossCmd.AddCommand(ossObjCmd)
	ossLsCmd.Flags().StringVarP(&ossLsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	ossLsCmd.Flags().BoolVar(&ossLsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	ossObjCmd.AddCommand(ossObjLsCmd)
	ossObjCmd.AddCommand(ossObjGetCmd)
	ossObjGetCmd.Flags().StringVarP(&ossDownloadBucket, "bucket", "b", "all", "指定存储桶 (Specify Bucket)")
	ossObjGetCmd.Flags().StringVarP(&ossDownloadObject, "objectKey", "k", "all", "指定对象 (Specify object key)")
	ossObjGetCmd.Flags().StringVarP(&ossDownloadOutputPath, "outputPath", "o", "./result", "指定导出路径 (Specify output path)")
}

var ossCmd = &cobra.Command{
	Use:   "oss",
	Short: "执行与对象存储相关的操作 (Perform oss-related operations)",
	Long:  "执行与对象存储相关的操作 (Perform oss-related operations)",
}

var ossLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的存储桶 (List all buckets)",
	Long:  "列出所有的存储桶 (List all buckets)",
	Run: func(cmd *cobra.Command, args []string) {
		if ossLsBucket == "all" {
			log.Debugf("ossLsRegion: %s, ossLsFlushCache: %v", ossLsRegion, ossLsFlushCache)
			alioss.PrintBucketsList(ossLsRegion, ossLsFlushCache)
		} else {

		}
	},
}

var ossObjCmd = &cobra.Command{
	Use:   "obj",
	Short: "执行与对象相关的操作 (Perform oss-related operations)",
	Long:  "执行与对象相关的操作 (Perform oss-related operations)",
}

var ossObjGetCmd = &cobra.Command{
	Use:   "get",
	Short: "下载指定的对象 (Download objects)",
	Long:  "下载指定的对象 (Download objects)",
	Run: func(cmd *cobra.Command, args []string) {
		alioss.DownloadObjects(ossDownloadBucket, ossDownloadObject, ossDownloadOutputPath)
	},
}

var ossObjLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出指定存储桶的对象 (List objects in the specified bucket)",
	Long:  "列出指定存储桶的对象 (List objects in the specified bucket)",
	Run: func(cmd *cobra.Command, args []string) {
		alioss.PrintObjectsList()
	},
}
