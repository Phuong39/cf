package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/tencentvpc"
)

var (
	//lsRegion     string
	//lsFlushCache bool
	securityGroupId string
	//cloudName       string
)

func init() {
	RootCmd.AddCommand(vpcCmd)
	//vpcCmd.Flags().StringVarP(&cloudName, "cloudname", "cloudname", "tencent", "指定云厂商 (Set Cloud Vendors)")
	vpcCmd.Flags().StringVarP(&lsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	vpcCmd.Flags().StringVarP(&securityGroupId, "sgid", "s", "all", "指定安全组实例 ID (Set Security Group Id)")
	vpcCmd.PersistentFlags().BoolVar(&lsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
}

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "列出当前凭证下的VPC安全组策略 (List all vpc security group policy)",
	Long:  `列出当前凭证下的VPC安全组策略 (List all vpc security group policy)`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "tencent" {
			tencentvpc.PrintVPCSecurityGroupList(lsRegion, &securityGroupId, lsFlushCache)
			fmt.Println("")
			tencentvpc.PrintVPCSecurityGroupPoliciesList(lsRegion, &securityGroupId, lsFlushCache)
			fmt.Println("")
		}
	},
}
