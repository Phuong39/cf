package alibaba

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliecs"
	"github.com/teamssix/cf/pkg/cloud/alibaba/alirds"
)

var ecsRegionsAllRegions bool

func init() {
	alibabaCmd.AddCommand(regionsCmd)
	regionsCmd.AddCommand(ecsRegionsCmd)
	regionsCmd.AddCommand(rdsRegionsCmd)
	ecsRegionsCmd.Flags().BoolVarP(&ecsRegionsAllRegions, "allRegions", "a", false, "列出所有区域，包括私有区域 (List all regions, including private regions)")
}

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "列出可用区域 (List available regions)",
	Long:  "列出可用区域 (List available regions)",
}

var ecsRegionsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "列出阿里云 ECS 的区域 (List the regions of alibaba cloud ECS)",
	Long:  "列出阿里云 ECS 的区域 (List the regions of alibaba cloud ECS)",
	Run: func(cmd *cobra.Command, args []string) {
		regions := aliecs.GetECSRegions(ecsRegionsAllRegions)
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

var rdsRegionsCmd = &cobra.Command{
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
		var header = []string{"序号 (SN)", "区域 ID (Region Id)", "可用区 ID (Zone ID)", "可用区名称 (Zone Name)", "地理位置 (Local Name)", "区域终端节点 (Region Endpoint)"}
		var td = cloud.TableData{Header: header, Body: data}
		cloud.PrintTable(td, "")
	},
}
