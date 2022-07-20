package keymanage

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

var ListKeyCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有已保存的 AK/SK (List all AK/SK)",
	Long:  "列出所有已保存的 AK/SK (List all AK/SK)",
	// ToDo: List keys
	Run: func(cmd *cobra.Command, args []string) {
		Data := cloud.TableData{
			Header: []string{"Name", "Platform", "AK", "SK", "STS", "Remark"},
		}
		for _, key := range KeyChain {
			Data.Body = append(Data.Body, []string{key.Name, key.Platform,
				cmdutil.MaskAK(key.AccessKeyId),
				cmdutil.MaskAK(key.AccessKeySecret),
				cmdutil.MaskAK(key.STSToken),
				key.Remark,
			})
		}
		cloud.PrintTable(Data, "")
	},
}
