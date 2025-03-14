package cmd

import (
	"github.com/biferdou/illapaca/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "illapaca",
	Short: "Illapaca - Terminal Weather Dashboard",
	Long: `Illapaca is a terminal-based weather dashboard that provides
current conditions, forecasts, and visualizations right in your terminal.

Named after the Inca god of weather, this tool offers quick access to
weather information for any location with a visually appealing interface.`,
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.illapa.yaml)")
	rootCmd.PersistentFlags().String("api-key", "", "API key for weather service")
	rootCmd.PersistentFlags().String("units", "metric", "Units to display (metric or imperial)")

	config.BindFlags(rootCmd)

	// Add all subcommands
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(forecastCmd)
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(favoriteCmd)
	rootCmd.AddCommand(alertsCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(versionCmd)
}

// Helper to get location from args or default
func getLocation(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return config.AppConfig.DefaultLocation
}
