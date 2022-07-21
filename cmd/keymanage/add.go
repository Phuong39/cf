package keymanage

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"strings"
)

var AddKeyCmd = &cobra.Command{
	Use: "add",
	// ToDo: Add keys
	Short: "添加密钥 (Add Key)",
	Long:  "添加密钥到数据库 (Add Key)",
	Run: func(cmd *cobra.Command, args []string) {
		cloudConfigList, _ := cmdutil.ReturnCloudProviderList()
		cred := Key{}

		var qs = []*survey.Question{
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
			{
				Name: "platform",
				Prompt: &survey.Select{
					Message: "Key 所属的云服务平台",
					Options: cloudConfigList,
				},
			},
			{
				Name:     "AccessKeyId",
				Prompt:   &survey.Input{Message: "Access Key Id (必须 Required):"},
				Validate: survey.Required,
			},
			{
				Name:     "AccessKeySecret",
				Prompt:   &survey.Password{Message: "Access Key Secret (必须 Required):"},
				Validate: survey.Required,
			},
			{
				Name:   "STSToken",
				Prompt: &survey.Input{Message: "STS Token (可选 Optional):"},
			},
		}

		// Generate the new config struct named cred to receive the inputted values.
		survey.Ask(qs, &cred)
		cred.AccessKeyId = strings.TrimSpace(cred.AccessKeyId)
		cred.AccessKeySecret = strings.TrimSpace(cred.AccessKeySecret)
		cred.STSToken = strings.TrimSpace(cred.STSToken)

		// Make user to check
		PrintSaving(cred)
		promot := &survey.Confirm{
			Message: "以上信息是否正确 (make sure correctness) "}
		sure := true // Break out
		survey.AskOne(promot, sure)
		if sure {
			KeyDb.Save(cred)
		}
	},
}

func PrintSaving(key Key) {
	Data := cloud.TableData{
		Header: CommonTableHeader,
	}
	Data.Body = append(Data.Body, []string{
		key.Name, key.Platform,
		key.AccessKeyId, key.AccessKeySecret,
		key.STSToken, key.Remark,
	})
	cloud.PrintTable(Data, "")
}
