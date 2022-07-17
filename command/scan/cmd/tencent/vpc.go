package tencent

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentvpc"
)

var (
	vpclsRegion     string
	vpclsFlushCache bool
	securityGroupId string
)

func init() {
	tencentCmd.AddCommand(vpcCmd)
	vpcCmd.AddCommand(vpclsCmd)
	vpclsCmd.Flags().StringVarP(&vpclsRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	vpclsCmd.Flags().StringVarP(&securityGroupId, "securityGroupId", "i", "all", "指定安全组实例 ID (Set Security Group Id)")
	vpclsCmd.PersistentFlags().BoolVar(&vpclsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
}

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "执行与 VPC 相关的操作 (Perform vpc-related operations)",
	Long:  "执行与 VPC 相关的操作 (Perform vpc-related operations)",
}

var vpclsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出当前凭证下的VPC安全组策略 (List all vpc security group policy)",
	Long:  `列出当前凭证下的VPC安全组策略 (List all vpc security group policy)`,
	Run: func(cmd *cobra.Command, args []string) {
		tencentvpc.PrintVPCSecurityGroupPoliciesList(vpclsRegion, &securityGroupId, vpclsFlushCache)
	},
}
