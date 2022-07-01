package cmd

import (
	"githubu.com/teamssix/cf/pkg/cloud"
	"githubu.com/teamssix/cf/pkg/cloud/aliecs"
	"github.com/spf13/cobra"
	"strconv"
)

func init() {
	RootCmd.AddCommand(regionsCmd)
	regionsCmd.AddCommand(aliyunRegionsCmd)
}

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "列出所有的区域 (List all regions)",
	Long:  "列出所有的区域 (List all regions)",
}

var aliyunRegionsCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "列出阿里云的区域 (List alibaba cloud regions)",
	Long:  "列出阿里云的区域 (List alibaba cloud regions)",
	Run: func(cmd *cobra.Command, args []string) {
		regions := aliecs.GetECSRegions()
		var data = make([][]string, len(regions))
		for i, v := range regions {
			SN := strconv.Itoa(i + 1)
			data[i] = []string{SN, v.RegionId, v.LocalName, v.RegionEndpoint}
		}
		var header = []string{"序号 (SN)", "区域 ID (Region Id)", "地理位置 (Local Name)", "区域终端节点 (Region Endpoint)"}
		var td = cloud.TableData{Header: header, Body: data}
		cloud.PrintTable(td, "")
	},
}
