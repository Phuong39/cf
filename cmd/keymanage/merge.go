package keymanage

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

var MergeKeyCmd = &cobra.Command{
	Use:   "merge",
	Short: "保存当前所使用的 Key 对 (Save current using key in Local DataBase)",
	Long:  "保存当前所使用的 Key 对 (Save current using key in Local DataBase)",
	// ToDo: fast add current using key in Local DataBase
	Run: func(cmd *cobra.Command, args []string) {
		if len(HeaderKey) == 0 {
			log.Info("未监测到您配置文件中的密钥 (No Detect your current Key in Provider's Config)")
		}
		for _, key := range HeaderKey {
			log.Info("正在保存当前所使用的 Key 对 (Saving current using key in Local DataBase)...")
			log.Infof("所属平台: %s", key.Platform)
			log.Infof("AK: %s", cmdutil.MaskAK(key.AccessKeyId))
			log.Infof("SK: %s", cmdutil.MaskAK(key.AccessKeySecret))
			log.Infof("STS Token: %s", cmdutil.MaskAK(key.STSToken))
			skip := false
			prompt := &survey.Confirm{
				Message: "是否保存当前所使用的 Key 对 (Save current using key in Local DataBase)?",
			}
			survey.AskOne(prompt, &skip)
			if skip {
				return
			}
			questions := []*survey.Question{
				{
					Name:     "name",
					Prompt:   &survey.Input{Message: "请为当前使用的 Key 输入名称 (Please input name for current using key)"},
					Validate: survey.Required,
				},
				{
					Name: "remark",
					Prompt: &survey.Input{
						Message: "请为当前使用的 Key 输入备注 (可选) " +
							"(Please input remark for current using key)[Optional]",
					},
				},
			}
			answers := struct {
				Name   string `survey:"name"`
				Remark string `survey:"remark"`
			}{}
			survey.Ask(questions, &answers)
			key.Name = answers.Name
			key.Remark = answers.Remark
			KeyDb.Save(key)
			log.Infof("快速保存成功 (fast merge into local db successful)")
		}
	},
}
