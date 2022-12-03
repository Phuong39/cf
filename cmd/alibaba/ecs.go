package alibaba

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	aliecs2 "github.com/teamssix/cf/pkg/cloud/alibaba/aliecs"
)

var (
	timeOut int

	running           bool
	userData          bool
	batchCommand      bool
	ecsFlushCache     bool
	ecsLsAllRegions   bool
	ecsExecAllRegions bool
	metaDataSTSToken  bool

	lhost                      string
	lport                      string
	command                    string
	scriptType                 string
	commandFile                string
	ecsLsRegion                string
	ecsExecRegion              string
	ecsLsSpecifiedInstanceID   string
	ecsExecSpecifiedInstanceID string
)

func init() {
	alibabaCmd.AddCommand(ecsCmd)
	ecsCmd.AddCommand(ecsLsCmd)
	ecsCmd.AddCommand(ecsExecCmd)

	ecsCmd.PersistentFlags().BoolVar(&ecsFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	ecsLsCmd.Flags().StringVarP(&ecsLsRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")
	ecsLsCmd.Flags().StringVarP(&ecsLsSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Specify instance ID)")
	ecsLsCmd.Flags().BoolVar(&running, "running", false, "只显示正在运行的实例 (Show only running instances)")
	ecsLsCmd.Flags().BoolVarP(&ecsLsAllRegions, "allRegions", "a", false, "使用所有区域，包括私有区域 (Use all regions, including private regions)")

	ecsExecCmd.Flags().StringVarP(&ecsExecSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Specify Instance ID)")
	ecsExecCmd.Flags().StringVarP(&command, "command", "c", "", "设置待执行的命令 (Set the command you want to execute)")
	ecsExecCmd.Flags().StringVarP(&commandFile, "file", "f", "", "设置待执行的命令文件 (Set the command file you want to execute)")
	ecsExecCmd.Flags().StringVarP(&scriptType, "scriptType", "s", "auto", "设置执行脚本的类型 (Specify the type of script to execute) [sh|bat|ps]")
	ecsExecCmd.Flags().StringVar(&lhost, "lhost", "", "设置反弹 shell 的主机 IP (Set the ip of the listening host)")
	ecsExecCmd.Flags().StringVar(&lport, "lport", "", "设置反弹 shell 的主机端口 (Set the port of the listening host)")
	ecsExecCmd.Flags().BoolVarP(&batchCommand, "batchCommand", "b", false, "一键执行三要素，方便 HW (Batch execution of multiple commands used to prove permission acquisition)")
	ecsExecCmd.Flags().BoolVarP(&userData, "userData", "u", false, "一键获取实例中的用户数据 (Get the user data on the instance)")
	ecsExecCmd.Flags().BoolVarP(&metaDataSTSToken, "metaDataSTSToken", "m", false, "一键获取实例元数据中的临时访问密钥 (Get the STS Token in the instance metadata)")
	ecsExecCmd.Flags().IntVarP(&timeOut, "timeOut", "t", 60, "设置命令执行结果的等待时间 (Set the command execution result waiting time)")
	ecsExecCmd.Flags().BoolVarP(&ecsExecAllRegions, "allRegions", "a", false, "使用所有区域，包括私有区域 (Use all regions, including private regions)")
	ecsExecCmd.Flags().StringVarP(&ecsExecRegion, "region", "r", "all", "指定区域 ID (Specify region ID)")

}

var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "执行与弹性计算服务相关的操作 (Perform ecs-related operations)",
	Long:  "执行与弹性计算服务相关的操作 (Perform ecs-related operations)",
}

var ecsLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的实例 (List all instances)",
	Long:  "列出所有的实例 (List all instances)",
	Run: func(cmd *cobra.Command, args []string) {
		aliecs2.PrintInstancesList(ecsLsRegion, running, ecsLsSpecifiedInstanceID, ecsFlushCache, ecsLsAllRegions)
	},
}

var ecsExecCmd = &cobra.Command{
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
			log.Warnf("还未指定要执行的命令 (The command to be executed has not been specified yet)\n")
			cmd.Help()
		} else {
			aliecs2.ECSExec(command, commandFile, scriptType, ecsExecSpecifiedInstanceID, ecsExecRegion, batchCommand, userData, metaDataSTSToken, ecsFlushCache, lhost, lport, timeOut, ecsExecAllRegions)
		}
	},
}
