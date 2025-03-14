// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/biferdou/illapaca/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Illapa Weather Dashboard v%s\n", config.Version)
	},
}
