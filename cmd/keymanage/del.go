package keymanage

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DelKeyCmd = &cobra.Command{
	Use:   "del",
	Short: "删除 key (delete key)",
	Long:  "删除 key (delete key)",
	Run: func(cmd *cobra.Command, args []string) {
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
		Result := KeyDb.Where("access_key_id = ? ", InputAccessKeyId).Delete(&Key{})
		if Result.RowsAffected == 0 {
			log.Error("没有找到对应的 Key 对 (No key found)")
			return
		} else {
			log.Info("删除成功 (Delete success)")
		}
	},
}
