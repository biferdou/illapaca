package cmd

import (
	"fmt"
	"strconv"

	"github.com/biferdou/illapaca/api"
	"github.com/biferdou/illapaca/config"
	"github.com/spf13/cobra"
)

var favoriteCmd = &cobra.Command{
	Use:   "favorite [command]",
	Short: "Manage favorite locations",
	Long: `Manage your favorite locations. Available commands:
  list    - List all favorite locations
  add     - Add a location to favorites
  remove  - Remove a location from favorites
  set-default - Set a location as default`,
}

var favoriteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List favorite locations",
	Run: func(cmd *cobra.Command, args []string) {
		config.ListFavoriteLocations()
	},
}

var favoriteAddCmd = &cobra.Command{
	Use:   "add [location]",
	Short: "Add a location to favorites",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := args[0]

		// Verify location by fetching its weather
		_, err := api.FetchWeather(location, 1)
		if err != nil {
			fmt.Printf("Error: invalid location or API error - %v\n", err)
			return
		}

		err = config.SaveFavoriteLocation(location)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Added %s to favorites\n", location)
	},
}

var favoriteRemoveCmd = &cobra.Command{
	Use:   "remove [location or index]",
	Short: "Remove a location from favorites",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]

		// Check if argument is an index
		index, err := strconv.Atoi(arg)
		if err == nil {
			err = config.RemoveFavoriteLocationByIndex(index)
		} else {
			err = config.RemoveFavoriteLocationByName(arg)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Removed location from favorites\n")
	},
}

var favoriteSetDefaultCmd = &cobra.Command{
	Use:   "set-default [location or index]",
	Short: "Set a location as default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]

		err := config.SetDefaultLocation(arg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

func init() {
	favoriteCmd.AddCommand(favoriteListCmd)
	favoriteCmd.AddCommand(favoriteAddCmd)
	favoriteCmd.AddCommand(favoriteRemoveCmd)
	favoriteCmd.AddCommand(favoriteSetDefaultCmd)
}
