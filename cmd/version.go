package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Yeego version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Yeego v0.1.2")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
