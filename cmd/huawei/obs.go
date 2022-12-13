package huawei

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/huawei/huaweiobs"
)

var (
	obsLsRegion       string
	obsLsFlushCache   bool
	obsLsObjectNumber string
)

func init() {
	huaweiCmd.AddCommand(obsCmd)
	obsCmd.AddCommand(obsLsCmd)

	obsLsCmd.Flags().StringVarP(&obsLsObjectNumber, "number", "n", "all", "指定列出对象的数量 (Specify the number of objects to list)")
	obsLsCmd.Flags().StringVarP(&obsLsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	obsLsCmd.Flags().BoolVar(&obsLsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

}

var obsCmd = &cobra.Command{
	Use:   "obs",
	Short: "执行与对象存储相关的操作 (Perform obs-related operations)",
	Long:  "执行与对象存储相关的操作 (Perform obs-related operations)",
}

var obsLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的存储桶 (List all buckets)",
	Long:  "列出所有的存储桶 (List all buckets)",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("obsLsRegion: %s, obsLsFlushCache: %v", obsLsRegion, obsLsFlushCache)
		huaweiobs.PrintBucketsList(obsLsRegion, obsLsFlushCache, obsLsObjectNumber)
	},
}
