package aws

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/cloud/aws/awsec2"
	"strconv"
)

var ec2RegionsAllRegions bool

func init() {
	awsCmd.AddCommand(regionsCmd)
	regionsCmd.AddCommand(ec2RegionsCmd)
}

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "列出可用区域 (List available regions)",
	Long:  "列出可用区域 (List available regions)",
}

var ec2RegionsCmd = &cobra.Command{
	Use:   "ec2",
	Short: "列出 aws ec2 的区域 (List the regions of aws ec2)",
	Long:  "列出 aws ec2 的区域 (List the regions of aws ec2)",
	Run: func(cmd *cobra.Command, args []string) {
		awsec2.GetEC2Regions()
		regions := awsec2.GetEC2Regions()
		if len(regions) > 0 {
			var data = make([][]string, len(regions))
			for i, v := range regions {
				SN := strconv.Itoa(i + 1)
				data[i] = []string{SN, *v.RegionName, *v.Endpoint}
			}
			var header = []string{"序号 (SN)", "区域名称 (Region Name)", "区域终端节点 (Region Endpoint)"}
			var td = cloud.TableData{Header: header, Body: data}
			cloud.PrintTable(td, "")
		} else {
			log.Infoln("未找到区域 (Regions not found)")
		}
	},
}
