// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version information
const Version = "1.0.0"

// Config variables
var (
	CfgFile   string
	AppConfig Config
)

// Config struct for app configuration
type Config struct {
	APIKey            string
	DefaultLocation   string
	Units             string
	FavoriteLocations []string
	AlertThresholds   AlertThresholds
}

// AlertThresholds for weather alerts
type AlertThresholds struct {
	HighTemp      float64
	LowTemp       float64
	Precipitation float64
	WindSpeed     float64
}

// InitConfig initializes the configuration
func InitConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".illapa")
	}

	viper.SetEnvPrefix("ILLAPA")
	viper.AutomaticEnv()

	// Load .env file if it exists
	godotenv.Load()

	// Set default values
	viper.SetDefault("units", "metric") // OpenWeatherMap supports metric, imperial, standard
	viper.SetDefault("favorite_locations", []string{})
	viper.SetDefault("alert_thresholds", map[string]float64{
		"high_temp":     35.0,
		"low_temp":      0.0,
		"precipitation": 70.0,
		"wind_speed":    30.0,
	})

	// Read config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Parse config
	AppConfig = Config{
		APIKey:            viper.GetString("api_key"),
		DefaultLocation:   viper.GetString("default_location"),
		Units:             viper.GetString("units"),
		FavoriteLocations: viper.GetStringSlice("favorite_locations"),
		AlertThresholds: AlertThresholds{
			HighTemp:      viper.GetFloat64("alert_thresholds.high_temp"),
			LowTemp:       viper.GetFloat64("alert_thresholds.low_temp"),
			Precipitation: viper.GetFloat64("alert_thresholds.precipitation"),
			WindSpeed:     viper.GetFloat64("alert_thresholds.wind_speed"),
		},
	}

	// Override with environment variables if they exist
	if os.Getenv("ILLAPA_API_KEY") != "" {
		AppConfig.APIKey = os.Getenv("ILLAPA_API_KEY")
	}
}

// BindFlags binds command flags to viper
func BindFlags(cmd *cobra.Command) {
	viper.BindPFlag("api_key", cmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("units", cmd.PersistentFlags().Lookup("units"))
}

// Rest of functions remain the same as they don't depend on the API structure
