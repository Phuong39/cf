package tencent

import (
	"strconv"

	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcvm"

	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
)

func init() {
	tencentCmd.AddCommand(regionsCmd)
	regionsCmd.AddCommand(CVMRegionsCmd)
}

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "列出可用区域 (List available regions)",
	Long:  "列出可用区域 (List available regions)",
}

var CVMRegionsCmd = &cobra.Command{
	Use:   "cvm",
	Short: "列出腾讯云 CVM 的区域 (List the regions of tencent cloud CVM)",
	Long:  "列出腾讯云 CVM 的区域 (List the regions of tencent cloud CVM)",
	Run: func(cmd *cobra.Command, args []string) {
		regions := tencentcvm.GetCVMRegions()
		var data = make([][]string, len(regions))
		for i, v := range regions {
			SN := strconv.Itoa(i + 1)
			data[i] = []string{SN, *v.Region, *v.RegionName, *v.RegionState}
		}
		var header = []string{"序号 (SN)", "地域名称 (Region)", "地域描述 (Region Name)", "地域是否可用状态 (Region State)"}
		var td = cloud.TableData{Header: header, Body: data}
		cloud.PrintTable(td, "")
	},
}
