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
	ossDownloadFlushCache bool
)

func init() {
	alibabaCmd.AddCommand(ossCmd)

	ossCmd.AddCommand(ossLsCmd)
	ossLsCmd.Flags().StringVarP(&ossLsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	ossLsCmd.Flags().StringVarP(&ossLsBucket, "bucket", "b", "all", "列出指定 Bucket 下的对象 (List objects in Bucket)")
	ossLsCmd.Flags().BoolVar(&ossLsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	ossCmd.AddCommand(ossDownloadCmd)
	ossDownloadCmd.Flags().StringVarP(&ossDownloadBucket, "bucket", "b", "all", "指定存储桶 (Set Bucket)")
	ossDownloadCmd.Flags().StringVarP(&ossDownloadObject, "objectKey", "k", "all", "指定对象 (Set object key)")
	ossDownloadCmd.Flags().StringVarP(&ossDownloadOutputPath, "outputPath", "o", "./result", "指定导出路径 (Set output path)")
	ossDownloadCmd.Flags().BoolVar(&ossDownloadFlushCache, "flushCache", true, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
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
			log.Debugf("ossLsBucket: %s", ossLsBucket)
			alioss.PrintObjectsList(ossLsBucket)
		}
	},
}

var ossDownloadCmd = &cobra.Command{
	Use:   "get",
	Short: "下载指定的对象 (Download objects)",
	Long:  "下载指定的对象 (Download objects)",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("ossDownloadBucket: %s, ossDownloadObject: %s, ossDownloadOutputPath: %s, ossDownloadFlushCache: %v", ossDownloadBucket, ossDownloadObject, ossDownloadOutputPath, ossDownloadFlushCache)
		alioss.DownloadObjects(ossDownloadBucket, ossDownloadObject, ossDownloadOutputPath, ossDownloadFlushCache)
	},
}
