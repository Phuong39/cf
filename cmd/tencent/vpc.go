package tencent

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentvpc"
)

var (
	vpcRegion                 string
	vpclsFlushCache           bool
	op                        string
	vpclsSecurityGroupId      string
	vpccontrolSecurityGroupId string
)

func init() {
	tencentCmd.AddCommand(vpcCmd)
	vpcCmd.AddCommand(vpclsCmd)
	vpcCmd.AddCommand(vpccontrolCmd)
	vpcCmd.PersistentFlags().BoolVar(&vpclsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	vpclsCmd.Flags().StringVarP(&vpcRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	vpclsCmd.Flags().StringVarP(&vpclsSecurityGroupId, "securityGroupId", "i", "all", "指定安全组 ID (Set Security Group Id)")
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
		tencentvpc.PrintVPCSecurityGroupPoliciesList(vpcRegion, &vpclsSecurityGroupId, vpclsFlushCache)
	},
}

var vpccontrolCmd = &cobra.Command{
	Use:   "control",
	Short: "添加或删除当前凭证下的VPC安全组策略 (Add/Del current vpc security group policy rule)",
	Long:  "添加或删除当前凭证下的VPC安全组策略 (Add/Del current vpc security group policy rule)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentvpc.VPCControl()
	},
}
