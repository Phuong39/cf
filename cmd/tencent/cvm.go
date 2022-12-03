package tencent

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	tencentcvm2 "github.com/teamssix/cf/pkg/cloud/tencent/tencentcvm"
)

var (
	timeOut int

	running          bool
	userData         bool
	batchCommand     bool
	cvmFlushCache    bool
	metaDataSTSToken bool

	lhost                      string
	lport                      string
	command                    string
	scriptType                 string
	commandFile                string
	cvmLsRegion                string
	cvmLsSpecifiedInstanceID   string
	cvmExecSpecifiedInstanceID string
)

func init() {
	tencentCmd.AddCommand(cvmCmd)
	cvmCmd.AddCommand(cvmLsCmd)
	cvmCmd.AddCommand(cvmExecCmd)
	cvmCmd.PersistentFlags().BoolVar(&cvmFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")

	cvmLsCmd.Flags().BoolVar(&running, "running", false, "只显示正在运行的实例 (Show only running instances)")
	cvmLsCmd.Flags().StringVarP(&cvmLsRegion, "region", "r", "all", "指定区域 ID (Specify Region ID)")
	cvmLsCmd.Flags().StringVarP(&cvmLsSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Specify Instance ID)")

	cvmExecCmd.Flags().StringVarP(&command, "command", "c", "", "设置待执行的命令 (Set the command you want to execute)")
	cvmExecCmd.Flags().StringVarP(&cvmExecSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Specify Instance ID)")
	cvmExecCmd.Flags().StringVarP(&commandFile, "file", "f", "", "设置待执行的命令文件 (Set the command file you want to execute)")
	cvmExecCmd.Flags().StringVarP(&scriptType, "scriptType", "s", "auto", "设置执行脚本的类型 (Set the type of script to execute) [sh|bat|ps]")
	cvmExecCmd.Flags().StringVar(&lhost, "lhost", "", "设置反弹 shell 的主机 IP (Set the ip of the listening host)")
	cvmExecCmd.Flags().StringVar(&lport, "lport", "", "设置反弹 shell 的主机端口 (Set the port of the listening host")
	cvmExecCmd.Flags().BoolVarP(&batchCommand, "batchCommand", "b", false, "一键执行三要素，方便 HW (Batch execution of multiple commands used to prove permission acquisition)")
	cvmExecCmd.Flags().BoolVarP(&userData, "userData", "u", false, "一键获取实例中的用户数据 (Get the user data on the instance)")
	cvmExecCmd.Flags().BoolVarP(&metaDataSTSToken, "metaDataSTSToken", "m", false, "一键获取实例元数据中的临时访问密钥 (Get the STS Token in the instance metadata)")
	cvmExecCmd.Flags().IntVarP(&timeOut, "timeOut", "t", 60, "设置命令执行结果的等待时间 (Set the command execution result waiting time)")
}

var cvmCmd = &cobra.Command{
	Use:   "cvm",
	Short: "执行与弹性计算服务相关的操作 (Perform cvm-related operations)",
	Long:  "执行与弹性计算服务相关的操作 (Perform cvm-related operations)",
}

var cvmLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有的实例 (List all instances)",
	Long:  "列出所有的实例 (List all instances)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentcvm2.PrintInstancesList(cvmLsRegion, running, cvmLsSpecifiedInstanceID, cvmFlushCache)
	},
}

var cvmExecCmd = &cobra.Command{
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
			tencentcvm2.CVMExec(command, commandFile, scriptType, cvmExecSpecifiedInstanceID, "all", batchCommand, userData, metaDataSTSToken, cvmFlushCache, lhost, lport, timeOut)
		}
	},
}
