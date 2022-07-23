package keymanage

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
)

var HeadKeyCmd = &cobra.Command{
	Use:     "head",
	Short:   "展示当前所使用的 Key 对 (Show current using key in All Provider config Path)",
	Long:    "展示当前所使用的 Key 对 (Show current using key in All Provider config Path)",
	Aliases: []string{"h"}, // Short Command
	Run: func(cmd *cobra.Command, args []string) {
		Data := cloud.TableData{
			Header: CommonTableHeader,
		}
		for _, key := range HeaderKey {
			Data.Body = append(Data.Body, []string{key.Name, key.Platform,
				key.AccessKeyId,
				key.AccessKeySecret,
				key.STSToken,
				key.Remark,
			})
		}
		cloud.PrintTable(Data, "")
	},
}
