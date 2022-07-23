package keymanage

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

var SwitchKeyCmd = &cobra.Command{
	Use:   "switch",
	Short: "切换当前使用的 Key 对 (Switch current using key in framework)",
	Long: "切换当前使用的 Key 对, 可以不传入参数也可以后跟一个存在于数据库的 AK 来进行快速切换 " +
		"(Switch current using key in framework, " +
		"you can not pass any parameter to switch, or you can pass one existing AK to switch)",
	Aliases: []string{"s", "checkout"}, // Short Command
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("SwitchKeyCmd", args)
		// Sure the AccessKeyId
		var InputAccessKeyId string
		if len(args) == 1 {
			InputAccessKeyId = args[0]
		} else {
			Promot := &survey.Input{
				Message: "输入存放的 Key 对的名称[支持LIKE正则] (Please input name of key pair)[Support SQL LIKE]",
			}
			var InputName string
			survey.AskOne(Promot, &InputName)
			var keys = []Key{}
			KeyDb.Where("name LIKE ? ", InputName).Find(&keys)
			if len(keys) == 0 {
				log.Error("没有找到对应的 Key 对 (No key found)")
				return
			}
			// PrintKeysTable(keys)
			keyAccessKeyIdList := []string{}
			for _, key := range keys {
				keyAccessKeyIdList = append(keyAccessKeyIdList, key.AccessKeyId)
			}
			Promot2 := &survey.Select{
				Message: "请输入要切换的 Key 对的 Access Key Id (Please input Access Key Id)",
				Default: keyAccessKeyIdList[0],
				Options: keyAccessKeyIdList,
			}
			PrintKeysTable(keys)
			survey.AskOne(Promot2, &InputAccessKeyId)
		}

		var key = Key{}
		Result := KeyDb.Where("access_key_id = ? ", InputAccessKeyId).First(&key)
		if Result.RowsAffected == 0 {
			log.Error("没有找到对应的 Key 对 (No key found)")
			return
		} else {
			cmdutil.SaveAccessKey(key.Config, key.Platform)
		}
	},
}
