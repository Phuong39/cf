package aws

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/aws/awss3"
)

var (
	s3LsRegion       string
	s3LsFlushCache   bool
	s3LsObjectNumber string
)

func init() {
	awsCmd.AddCommand(s3Cmd)
	s3Cmd.AddCommand(s3LsCmd)

	s3LsCmd.Flags().StringVarP(&s3LsObjectNumber, "number", "n", "all", "指定列出对象的数量 (Specify the number of objects to list)")
	s3LsCmd.Flags().StringVarP(&s3LsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	s3LsCmd.Flags().BoolVar(&s3LsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

}

var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "执行与对象存储相关的操作 (Perform s3-related operations)",
	Long:  "执行与对象存储相关的操作 (Perform s3-related operations)",
}

var s3LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的存储桶 (List all buckets)",
	Long:  "列出所有的存储桶 (List all buckets)",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("s3LsRegion: %s, s3LsFlushCache: %v", s3LsRegion, s3LsFlushCache)
		awss3.PrintBucketsList(s3LsRegion, s3LsFlushCache, s3LsObjectNumber)
	},
}
