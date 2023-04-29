package tencent

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	tencentlh2 "github.com/teamssix/cf/pkg/cloud/tencent/tencentlh"
)

var (
	lhFlushCache          bool
	lhRegion              string
	lhSpecifiedInstanceID string

	sshKeyName   string
	deleteSshKey bool
)

func init() {
	tencentCmd.AddCommand(lhCmd)
	lhCmd.AddCommand(lhLsCmd)
	lhCmd.AddCommand(lhExecCmd)
	lhCmd.PersistentFlags().BoolVar(&lhFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
	lhCmd.Flags().StringVarP(&lhSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Specify Instance ID)")

	lhLsCmd.Flags().BoolVar(&running, "running", false, "只显示正在运行的实例 (Show only running instances)")
	lhLsCmd.Flags().StringVarP(&lhRegion, "region", "r", "all", "指定区域 ID (Specify Region ID)")
	lhExecCmd.Flags().StringVarP(&command, "command", "c", "", "设置待执行的命令 (Set the command you want to execute)")
	lhExecCmd.Flags().StringVarP(&commandFile, "file", "f", "", "设置待执行的命令文件 (Set the command file you want to execute)")
	lhExecCmd.Flags().StringVarP(&scriptType, "scriptType", "s", "auto", "设置执行脚本的类型 (Set the type of script to execute) [sh|bat|ps]")
	lhExecCmd.Flags().StringVar(&lhost, "lhost", "", "设置反弹 shell 的主机 IP (Set the ip of the listening host)")
	lhExecCmd.Flags().StringVar(&lport, "lport", "", "设置反弹 shell 的主机端口 (Set the port of the listening host")
	lhExecCmd.Flags().BoolVarP(&batchCommand, "batchCommand", "b", false, "一键执行三要素，方便 HW (Batch execution of multiple commands used to prove permission acquisition)")
	lhExecCmd.Flags().BoolVarP(&userData, "userData", "u", false, "一键获取实例中的用户数据 (Get the user data on the instance)")
	lhExecCmd.Flags().BoolVarP(&metaDataSTSToken, "metaDataSTSToken", "m", false, "一键获取实例元数据中的临时访问密钥 (Get the STS Token in the instance metadata)")
	lhExecCmd.Flags().IntVarP(&timeOut, "timeOut", "t", 60, "设置命令执行结果的等待时间 (Set the command execution result waiting time)")
}

var lhCmd = &cobra.Command{
	Use:   "lh",
	Short: "执行与轻量计算服务相关的操作 (Perform lh-related operations)",
	Long:  "执行与轻量计算服务相关的操作 (Perform lh-related operations)",
}

var lhLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的实例 (List all instances)",
	Long:  "列出所有的实例 (List all instances)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentlh2.PrintInstancesList(lhRegion, running, lhSpecifiedInstanceID, lhFlushCache)
	},
}

var lhExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "在实例上执行命令 (Execute the command on the instance)",
	Long:  "在实例上执行命令 (Execute the command on the instance)",
	Run: func(cmd *cobra.Command, args []string) {
		if lhost != "" && lport == "" {
			log.Warnln("未指定反弹 shell 的主机端口 (The port of the listening host is not set)")
			cmd.Help()
		} else if lhost == "" && lport != "" {
			log.Warnln("未指定反弹 shell 的主机 IP (The ip of the listening host is not set)")
			cmd.Help()
		} else if command == "" && batchCommand == false && userData == false && metaDataSTSToken == false && commandFile == "" && lhost == "" && lport == "" {
			log.Warnln("还未指定要执行的命令 (The command to be executed has not been specified yet)")
			cmd.Help()
		} else {
			tencentlh2.LhExec(command, commandFile, scriptType, lhSpecifiedInstanceID, lhRegion, batchCommand, userData, metaDataSTSToken, lhFlushCache, lhost, lport, timeOut)
		}
	},
}
