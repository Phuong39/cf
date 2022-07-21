package keymanage

import (
	"github.com/spf13/cobra"
)

var ListKeyCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出所有已保存的 AK/SK (List all AK/SK)",
	Long:  "列出所有已保存的 AK/SK (List all AK/SK)",
	// ToDo: List keys
	Run: func(cmd *cobra.Command, args []string) {
		KeyChains := []Key{} // Get all Keys
		KeyDb.Find(&KeyChains)
		PrintKeysTable(KeyChains)
	},
}
