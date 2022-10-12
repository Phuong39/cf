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
	lsAllRegions bool
)

func init() {
	alibabaCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&lsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	lsCmd.PersistentFlags().BoolVar(&lsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
	lsCmd.Flags().BoolVarP(&lsAllRegions, "allRegions", "a", false, "使用所有区域，包括私有区域 (Use all regions, including private regions)")
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "一键列出当前凭证下的 OSS、ECS、RDS 资源 (List OSS, ECS, RDS resources)",
	Long:  `一键列出当前凭证下的 OSS、ECS、RDS 资源 (List OSS, ECS, RDS resources)`,
	Run: func(cmd *cobra.Command, args []string) {
		alioss.PrintBucketsList(lsRegion, lsFlushCache, "all", "all")
		fmt.Println("")
		aliecs.PrintInstancesList(lsRegion, false, "all", lsFlushCache, lsAllRegions)
		fmt.Println("")
		alirds.PrintDBInstancesList(lsRegion, false, "all", "all", lsFlushCache)
	},
}
