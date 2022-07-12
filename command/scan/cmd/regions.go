package cmd

import (
	"github.com/teamssix/cf/pkg/cloud/tencentcvm"
	"strconv"

	"github.com/teamssix/cf/pkg/cloud/alirds"

	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/aliecs"
)

func init() {
	RootCmd.AddCommand(regionsCmd)
	regionsCmd.AddCommand(aliyunRegionsCmd)
	regionsCmd.AddCommand(tencentRegionsCmd)
	aliyunRegionsCmd.AddCommand(aliyunECSRegionsCmd)
	aliyunRegionsCmd.AddCommand(aliyunRDSRegionsCmd)
	tencentRegionsCmd.AddCommand(tencentCVMRegionsCmd)
}

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "列出所有的区域 (List all regions)",
	Long:  "列出所有的区域 (List all regions)",
}

var aliyunRegionsCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "列出阿里云的区域 (List the regions of alibaba cloud)",
	Long:  "列出阿里云的区域 (List the regions of alibaba cloud)",
}

var tencentRegionsCmd = &cobra.Command{
	Use:   "tencent",
	Short: "列出腾讯云的区域 (List the regions of tencent cloud)",
	Long:  "列出腾讯云的区域 (List the regions of tencent cloud)",
}

var aliyunECSRegionsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "列出阿里云 ECS 的区域 (List the regions of alibaba cloud ECS)",
	Long:  "列出阿里云 ECS 的区域 (List the regions of alibaba cloud ECS)",
	Run: func(cmd *cobra.Command, args []string) {
		regions := aliecs.GetECSRegions()
		var data = make([][]string, len(regions))
		for i, v := range regions {
			SN := strconv.Itoa(i + 1)
			data[i] = []string{SN, v.RegionId, v.LocalName, v.RegionEndpoint}
		}
		var header = []string{"序号 (SN)", "地域 ID (Region Id)", "地理位置 (Local Name)", "区域终端节点 (Region Endpoint)"}
		var td = cloud.TableData{Header: header, Body: data}
		cloud.PrintTable(td, "")
	},
}

var tencentCVMRegionsCmd = &cobra.Command{
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

var aliyunRDSRegionsCmd = &cobra.Command{
	Use:   "rds",
	Short: "列出阿里云 RDS 的区域 (List the regions of alibaba cloud RDS)",
	Long:  "列出阿里云 RDS 的区域 (List the regions of alibaba cloud RDS)",
	Run: func(cmd *cobra.Command, args []string) {
		regions := alirds.GetRDSRegions()
		var data = make([][]string, len(regions))
		for i, v := range regions {
			SN := strconv.Itoa(i + 1)
			data[i] = []string{SN, v.RegionId, v.ZoneId, v.ZoneName, v.LocalName, v.RegionEndpoint}
		}
		var header = []string{"序号 (SN)", "地域 ID (Region Id)", "可用区 ID (Zone ID)", "可用区名称 (Zone Name)", "地理位置 (Local Name)", "区域终端节点 (Region Endpoint)"}
		var td = cloud.TableData{Header: header, Body: data}
		cloud.PrintTable(td, "")
	},
}
