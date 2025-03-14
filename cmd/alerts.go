package cmd

import (
	"fmt"

	"github.com/biferdou/illapaca/config"
	"github.com/spf13/cobra"
)

var alertsCmd = &cobra.Command{
	Use:   "alerts [command]",
	Short: "Manage weather alerts",
	Long: `Manage weather alert thresholds. Available commands:
  show    - Show current alert thresholds
  set     - Set alert thresholds`,
}

var alertsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current alert thresholds",
	Run: func(cmd *cobra.Command, args []string) {
		config.ShowAlertThresholds()
	},
}

var alertsSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set alert thresholds",
	Run: func(cmd *cobra.Command, args []string) {
		highTemp, _ := cmd.Flags().GetFloat64("high-temp")
		lowTemp, _ := cmd.Flags().GetFloat64("low-temp")
		precip, _ := cmd.Flags().GetFloat64("precipitation")
		wind, _ := cmd.Flags().GetFloat64("wind-speed")

		err := config.SetAlertThresholds(highTemp, lowTemp, precip, wind)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("Alert thresholds updated successfully")
	},
}

func init() {
	alertsCmd.AddCommand(alertsShowCmd)
	alertsCmd.AddCommand(alertsSetCmd)

	alertsSetCmd.Flags().Float64("high-temp", 0, "High temperature threshold (°C)")
	alertsSetCmd.Flags().Float64("low-temp", 0, "Low temperature threshold (°C)")
	alertsSetCmd.Flags().Float64("precipitation", 0, "Precipitation chance threshold (%)")
	alertsSetCmd.Flags().Float64("wind-speed", 0, "Wind speed threshold (km/h)")
}
