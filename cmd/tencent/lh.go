package tencent

// 腾讯云lh相关操作

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
	lhCmd.AddCommand(lhSSHCmd)
	lhSSHCmd.AddCommand(lhLsSSHCmd)
	lhSSHCmd.AddCommand(lhGenerateSSHCmd)
	lhSSHCmd.AddCommand(lhDeleteSSHCmd)
	lhCmd.PersistentFlags().BoolVar(&lhFlushCache, "flushCache", false, "刷新缓存，不使用缓存数据 (Refresh the cache without using cached data)")
	lhCmd.Flags().StringVarP(&lhRegion, "region", "r", "all", "指定区域 ID (Set Region ID)")
	lhCmd.Flags().StringVarP(&lhSpecifiedInstanceID, "instanceID", "i", "all", "指定实例 ID (Set Instance ID)")

	lhLsCmd.Flags().BoolVar(&running, "running", false, "只显示正在运行的实例 (Show only running instances)")

	lhExecCmd.Flags().StringVarP(&command, "command", "c", "", "设置待执行的命令 (Set the command you want to execute)")
	lhExecCmd.Flags().StringVarP(&commandFile, "file", "f", "", "设置待执行的命令文件 (Set the command file you want to execute)")
	lhExecCmd.Flags().StringVarP(&scriptType, "scriptType", "s", "auto", "设置执行脚本的类型 (Set the type of script to execute) [sh|bat|ps]")
	lhExecCmd.Flags().StringVar(&lhost, "lhost", "", "设置反弹 shell 的主机 IP (Set the ip of the listening host)")
	lhExecCmd.Flags().StringVar(&lport, "lport", "", "设置反弹 shell 的主机端口 (Set the port of the listening host")
	lhExecCmd.Flags().BoolVarP(&batchCommand, "batchCommand", "b", false, "一键执行三要素，方便 HW (Batch execution of multiple commands used to prove permission acquisition)")
	lhExecCmd.Flags().BoolVarP(&userData, "userData", "u", false, "一键获取实例中的用户数据 (Get the user data on the instance)")
	lhExecCmd.Flags().BoolVarP(&metaDataSTSToken, "metaDataSTSToken", "m", false, "一键获取实例元数据中的临时访问凭证 (Get the STS Token in the instance metadata)")
	lhExecCmd.Flags().IntVarP(&timeOut, "timeOut", "t", 60, "设置命令执行结果的等待时间 (Set the command execution result waiting time)")

	lhGenerateSSHCmd.Flags().StringVarP(&sshKeyName, "keyName", "k", "", "设置新密钥对并绑定,此操作会导致原有登录方式失效(Set a new key pair and bind it. This operation will invalidate the original login method)")
	//lhSSHCmd.Flags().BoolVarP(&deleteSshKey, "deleteKey", "d", false, "删除绑定的密钥对并替换回原有密钥对(Delete the bound key pair and replace it with the original key pair)")
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

var lhSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "与ssh有关的操作(SSH related operations)",
	Long:  "与ssh有关的操作(SSH related operations)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var lhLsSSHCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出目前所有密钥与相关实例信息(List all current keys and related instance information)",
	Long:  "列出目前所有密钥与相关实例信息(List all current keys and related instance information)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentlh2.GetSSHKeysListInfo(lhFlushCache)
	},
}

var lhGenerateSSHCmd = &cobra.Command{
	Use:   "generate",
	Short: "创建并绑定新密钥,该操作生效存在一定延迟，可通过ls --flushCache确定是否绑定成功(Create and bind a new key. There is a certain delay in the effectiveness of this operation. You can determine whether the binding is successful through LS --flushcache)",
	Long:  "创建并绑定新密钥,该操作生效存在一定延迟，可通过ls --flushCache确定是否绑定成功(Create and bind a new key. There is a certain delay in the effectiveness of this operation. You can determine whether the binding is successful through LS --flushcache)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentlh2.GenerateSSHKeyOnInstance(lhRegion, sshKeyName, lhSpecifiedInstanceID, lhFlushCache)
	},
}

var lhDeleteSSHCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除通过generate方法生成的密钥,该操作生效存在一定延迟，可通过ls --flushCache确定是否解绑成功(Delete the key generated by the generate method. There is a certain delay for this operation to take effect. You can determine whether the unbinding is successful through LS --flushcache)",
	Long:  "删除通过generate方法生成的密钥,该操作生效存在一定延迟，可通过ls --flushCache确定是否解绑成功(Delete the key generated by the generate method. There is a certain delay for this operation to take effect. You can determine whether the unbinding is successful through LS --flushcache)",
	Run: func(cmd *cobra.Command, args []string) {
		tencentlh2.DeleteSSHKeyOnInstance(lhRegion, lhSpecifiedInstanceID, lhFlushCache)
	},
}
