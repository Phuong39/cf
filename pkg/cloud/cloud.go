package cloud

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Config struct {
	Alias           string
	AccessKeyId     string
	AccessKeySecret string
	STSToken        string
	Provider        string
	InUse           bool
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
