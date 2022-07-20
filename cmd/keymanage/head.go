package keymanage

import (
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"
)

var HeadKeyCmd = &cobra.Command{
	Use: "head",
	// ToDo: Current using keys
	Run: func(cmd *cobra.Command, args []string) {
		Data := cloud.TableData{
			Header: []string{"Name", "Platform", "AK", "SK", "STS", "Remark"},
		}
		for _, key := range HeaderKey {
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
