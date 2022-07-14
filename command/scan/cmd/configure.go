package cmd

import (
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

func init() {
	RootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(getconfigCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "配置云服务商的访问密钥 (Configure cloud provider access key)",
	Long:  `配置云服务商的访问密钥 (Configure cloud provider access key)`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdutil.ConfigureAccessKey(args[0])
		//cmdutil.ConfigureAccessKey()
	},
}

var getconfigCmd = &cobra.Command{
	Use:   "ls",
	Short: "获取当前配置的访问凭证 (Get the currently configured access key)",
	Long:  `获取当前配置的访问凭证 (Get the currently configured access key)`,
	Run: func(cmd *cobra.Command, args []string) {
		config := cmdutil.GetAliCredential()
		AccessKeyId := config.Alibaba.AccessKeyId
		AccessKeySecret := config.Alibaba.AccessKeySecret
		STSToken := config.Alibaba.STSToken
		if AccessKeyId == "" {
			log.Infoln("当前未配置访问密钥 (No access key configured)")
		} else {
			color.Printf(`<lightGreen>访问凭证 ID (Access key id):</> %s
<lightGreen>访问凭证密钥 (Access key secret):</> %s
<lightGreen>临时访问凭证令牌 (STS token):</> %s
<lightGreen>配置文件路径 (Configuration file path):</> %s
`, AccessKeyId, AccessKeySecret, STSToken, cmdutil.GetAliCredentialFilePath())
		}
	},
}
