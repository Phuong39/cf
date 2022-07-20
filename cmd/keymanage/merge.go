package keymanage

import (
	"github.com/spf13/cobra"
)

var MergeKeyCmd = &cobra.Command{
	Use:   "merge",
	Short: "保存当前所使用的 Key 对 (Save current using key in Local DataBase)",
	Long:  "保存当前所使用的 Key 对 (Save current using key in Local DataBase)",
	// ToDo: fast add current using key in Local DataBase
}
