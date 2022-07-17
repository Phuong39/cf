package alibaba

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliecs"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alioss"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alirds"
)

var (
	lsRegion     string
	lsFlushCache bool
)

func init() {
	alibabaCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&lsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	lsCmd.PersistentFlags().BoolVar(&lsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出当前凭证下的云服务资源 (List all resources)",
	Long:  `列出当前凭证下的云服务资源 (List all resources)`,
	Run: func(cmd *cobra.Command, args []string) {
		alioss.PrintBucketsList(lsRegion, lsFlushCache)
		fmt.Println("")
		aliecs.PrintInstancesList(lsRegion, false, "all", lsFlushCache)
		fmt.Println("")
		alirds.PrintDBInstancesList(lsRegion, false, "all", "all", lsFlushCache)
	},
}
