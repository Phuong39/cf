package cloud

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

//type Credential struct {
//	AccessKeyId     string `json:"AccessKeyId"`
//	AccessKeySecret string `json:"AccessKeySecret"`
//	STSToken        string `json:"STSToken"`
//}

//增加腾讯云厂商结构支持，未来方便拓展
type Credential struct {
	Tencent struct {
		SecretId     string `json:"SecretId"`
		SecretKey    string `json:"SecretKey"`
		Token        string `json:"Token"`
		TmpSecretId  string `json:"TmpSecretId"`
		TmpSecretKey string `json:"TmpSecretKey"`
	} `json:"tencent"`
	Alibaba struct {
		AccessKeyId     string `json:"AccessKeyId"`
		AccessKeySecret string `json:"AccessKeySecret"`
		STSToken        string `json:"STSToken"`
	} `json:"alibaba"`
}

type Bucket = Resource

type TableData struct {
	Header []string
	Body   [][]string
}

type Resource struct {
	Name       string
	Region     string
	Properties *Property
}

type Property map[string]string

var regionMaps = map[string][]string{
	"default": {"cn-beijing"},
}

func GetGlobalRegions() []string {
	return GetRegions("default")
}

func GetRegions(key string) []string {
	if val, ok := regionMaps[key]; ok {
		return val
	}
	return regionMaps["default"]
}

func PrintTable(data TableData, Caption string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(data.Header)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	var TableHeaderColor = make([]tablewriter.Colors, len(data.Header))
	for i := range TableHeaderColor {
		TableHeaderColor[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(TableHeaderColor...)
	if Caption != "" {
		table.SetCaption(true, Caption)
	}
	table.AppendBulk(data.Body)
	table.Render()
}
