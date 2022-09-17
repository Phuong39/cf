package alibaba

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alioss"
)

var (
	ossLsRegion           string
	ossLsFlushCache       bool
	ossDownloadBucket     string
	ossDownloadObject     string
	ossDownloadOutputPath string
	ossLsObjectNumber     string
)

func init() {
	alibabaCmd.AddCommand(ossCmd)
	ossCmd.AddCommand(ossLsCmd)
	ossCmd.AddCommand(ossObjCmd)
	ossLsCmd.Flags().StringVarP(&ossLsObjectNumber, "number", "n", "all", "指定列出对象的数量 (Specify the number of objects to list)")
	ossLsCmd.Flags().StringVarP(&ossLsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	ossLsCmd.Flags().BoolVar(&ossLsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	ossObjCmd.AddCommand(ossObjLsCmd)
	ossObjCmd.AddCommand(ossObjGetCmd)
	ossObjGetCmd.Flags().StringVarP(&ossDownloadBucket, "bucket", "b", "all", "指定存储桶 (Specify Bucket)")
	ossObjGetCmd.Flags().StringVarP(&ossDownloadObject, "objectKey", "k", "all", "指定对象 (Specify object key)")
	ossObjGetCmd.Flags().StringVarP(&ossDownloadOutputPath, "outputPath", "o", "./result", "指定导出路径 (Specify output path)")

	ossObjLsCmd.Flags().StringVarP(&ossLsObjectNumber, "number", "n", "all", "指定列出对象的数量 (Specify the number of objects to list)")
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
		log.Debugf("ossLsRegion: %s, ossLsFlushCache: %v", ossLsRegion, ossLsFlushCache)
		alioss.PrintBucketsList(ossLsRegion, ossLsFlushCache, ossLsObjectNumber)
	},
}

var ossObjCmd = &cobra.Command{
	Use:   "obj",
	Short: "执行与对象相关的操作 (Perform objects-related operations)",
	Long:  "执行与对象相关的操作 (Perform objects-related operations)",
}

var ossObjGetCmd = &cobra.Command{
	Use:   "get",
	Short: "下载存储桶里的对象 (Downloading objects from the bucket)",
	Long:  "下载存储桶里的对象 (Downloading objects from the bucket)",
	Run: func(cmd *cobra.Command, args []string) {
		alioss.DownloadObjects(ossDownloadBucket, ossDownloadObject, ossDownloadOutputPath)
	},
}

var ossObjLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出存储桶里的对象 (List objects in the bucket)",
	Long:  "列出存储桶里的对象 (List objects in the bucket)",
	Run: func(cmd *cobra.Command, args []string) {
		alioss.PrintObjectsList(ossLsObjectNumber)
	},
}
