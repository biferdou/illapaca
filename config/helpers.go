package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

// Save the current config to disk
func SaveConfig() error {
	viper.Set("api_key", AppConfig.APIKey)
	viper.Set("default_location", AppConfig.DefaultLocation)
	viper.Set("units", AppConfig.Units)
	viper.Set("favorite_locations", AppConfig.FavoriteLocations)
	viper.Set("alert_thresholds.high_temp", AppConfig.AlertThresholds.HighTemp)
	viper.Set("alert_thresholds.low_temp", AppConfig.AlertThresholds.LowTemp)
	viper.Set("alert_thresholds.precipitation", AppConfig.AlertThresholds.Precipitation)
	viper.Set("alert_thresholds.wind_speed", AppConfig.AlertThresholds.WindSpeed)

	return viper.WriteConfig()
}

// List favorite locations
func ListFavoriteLocations() {
	if len(AppConfig.FavoriteLocations) == 0 {
		println("No favorite locations saved")
		return
	}

	fmt.Println("Favorite locations:")
	for i, loc := range AppConfig.FavoriteLocations {
		if loc == AppConfig.DefaultLocation {
			fmt.Printf("%d. %s (default)\n", i+1, loc)
		} else {
			fmt.Printf("%d. %s\n", i+1, loc)
		}
	}
}

// Add a location to the list of favorite locations
func SaveFavoriteLocation(location string) error {
	// Check if the location is already in the list
	for _, loc := range AppConfig.FavoriteLocations {
		if loc == location {
			return fmt.Errorf("%s is already a favorite location", location)
		}
	}

	// Add the location to the list
	AppConfig.FavoriteLocations = append(AppConfig.FavoriteLocations, location)

	// Save the config
	return SaveConfig()
}

// Remove a favorite location by index
func RemoveFavoriteLocationByIndex(index int) error {
	if index < 1 || index > len(AppConfig.FavoriteLocations) {
		return fmt.Errorf("index out of range: %d", index)
	}

	// Adjust for 0-based indexing
	index--

	// Get location name for output
	removedLocation := AppConfig.FavoriteLocations[index]

	// Remove from slice
	AppConfig.FavoriteLocations = append(
		AppConfig.FavoriteLocations[:index],
		AppConfig.FavoriteLocations[index+1:]...,
	)

	// Save the config
	if err := SaveConfig(); err != nil {
		return err
	}

	fmt.Printf("Removed %s from favorite locations\n", removedLocation)
	return nil
}

// Remove a favorite location by name
func RemoveFavoriteLocationByName(location string) error {
	var newFavoriteLocations []string
	found := false

	for _, loc := range AppConfig.FavoriteLocations {
		if loc != location {
			newFavoriteLocations = append(newFavoriteLocations, loc)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("%s is not a favorite location", location)
	}

	AppConfig.FavoriteLocations = newFavoriteLocations

	// Save the config
	return SaveConfig()
}

// Set the default location
func SetDefaultLocation(arg string) error {
	// Check if the argument is an index
	index, err := strconv.Atoi(arg)
	if err == nil && index > 0 && index <= len(AppConfig.FavoriteLocations) {
		// It's an index
		defaultLocation := AppConfig.FavoriteLocations[index-1]

		// Save the config
		AppConfig.DefaultLocation = defaultLocation
		if err := SaveConfig(); err != nil {
			return err
		}

		fmt.Printf("Set %s as the default location\n", defaultLocation)
		return nil
	}

	// It's a location name
	location := arg

	// Save config
	AppConfig.DefaultLocation = location
	if err := SaveConfig(); err != nil {
		return err
	}

	fmt.Printf("Set %s as the default location\n", location)
	return nil
}

// Set an alert thresholds
func SetAlertThresholds(highTemp, lowTemp, precipitation, windSpeed float64) error {
	if highTemp != 0 {
		AppConfig.AlertThresholds.HighTemp = highTemp
	}

	if lowTemp != 0 {
		AppConfig.AlertThresholds.LowTemp = lowTemp
	}

	if precipitation != 0 {
		AppConfig.AlertThresholds.Precipitation = precipitation
	}

	if windSpeed != 0 {
		AppConfig.AlertThresholds.WindSpeed = windSpeed
	}

	// Save the config
	return SaveConfig()
}
