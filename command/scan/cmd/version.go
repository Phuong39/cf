package cmd

import (
	"githubu.com/teamssix/cf/pkg/cloud"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "输出 cf 的版本和更新时间 (Print the version number and update time of cf)",
	Long:  "输出 cf 的版本和更新时间 (Print the version number and update time of cf)",
	Run: func(cmd *cobra.Command, args []string) {
		data := [][]string{
			[]string{"0.0.1", "2022.7.1"},
		}
		var header = []string{"当前版本 (Version)", "更新时间 (Update Time)"}
		var td = cloud.TableData{Header: header, Body: data}
		cloud.PrintTable(td,"")
	},
}
