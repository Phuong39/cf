package keymanage

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
)

var ListKeyCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有已保存的 AK/SK (List all AK/SK)",
	Long:  "列出所有已保存的 AK/SK (List all AK/SK)",
	// ToDo: List keys
	Run: func(cmd *cobra.Command, args []string) {
		Data := cloud.TableData{
			Header: CommonTableHeader,
		}
		KeyChains := []Key{} // Get all Keys
		result := KeyDb.Find(&KeyChains)
		if result.RowsAffected == 0 {
			Data.Body = append(Data.Body,
				[]string{"", "", "", "", "", ""})
			// no handle the Result.Error
			log.Info("没有在本地数据库中找到任何密钥 (No key found in local database)")
		} else {
			for _, key := range KeyChains {
				Data.Body = append(Data.Body, []string{key.Name, key.Platform,
					key.AccessKeyId,
					key.AccessKeySecret,
					key.STSToken,
					key.Remark,
				})
			}
			cloud.PrintTable(Data, "")
		}
	},
}
