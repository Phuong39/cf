package aws

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud/aws/awsec2"
)

var (
	timeOut int

	running           bool
	userData          bool
	batchCommand      bool
	ec2FlushCache     bool
	ec2LsAllRegions   bool
	ec2ExecAllRegions bool
	metaDataSTSToken  bool

	lhost                      string
	lport                      string
	command                    string
	scriptType                 string
	commandFile                string
	ec2LsRegion                string
	ec2ExecRegion              string
	ec2LsSpecifiedInstanceID   string
	ec2Exec2pecifiedInstanceID string
)

func init() {
	awsCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(ec2LsCmd)

	ec2Cmd.PersistentFlags().BoolVar(&ec2FlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	ec2LsCmd.Flags().StringVarP(&ec2LsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	ec2LsCmd.Flags().StringVarP(&ec2LsSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Specify instance ID)")
	ec2LsCmd.Flags().BoolVar(&running, "running", false, "只显示正在运行的实例 (Show only running instances)")
	ec2LsCmd.Flags().BoolVarP(&ec2LsAllRegions, "allRegions", "a", false, "使用所有区域，包括私有区域 (Use all regions, including private regions)")
}

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "执行与弹性计算服务相关的操作 (Perform ec2-related operations)",
	Long:  "执行与弹性计算服务相关的操作 (Perform ec2-related operations)",
}

var ec2LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的实例 (List all instances)",
	Long:  "列出所有的实例 (List all instances)",
	Run: func(cmd *cobra.Command, args []string) {
		awsec2.PrintInstancesList(ec2LsRegion, running, ec2LsSpecifiedInstanceID, ec2FlushCache, ec2LsAllRegions)
	},
}
