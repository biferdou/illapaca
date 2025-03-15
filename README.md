# Illapaca

A terminal-based weather dashboard with clean, simple visualizations.

## Features

- 🌡️ Current weather conditions with color-coded information
- 🔮 Multi-day weather forecast
- 📊 Temperature trend visualization
- 🌧️ Precipitation chance charts
- 🔄 Location comparison
- ⚠️ Customizable weather alerts
- 📍 Favorite locations management

## Installation

### Prerequisites

- Go 1.24 or later

### From Source

```bash
# Clone the repository
git clone https://github.com/biferdou/illapaca.git
cd illapaca

# Build the binary
go build -o illapaca

# Move to a directory in your PATH (optional)
sudo mv illapaca /usr/local/bin/
```

### Configuration

On first run, Illapaca will create a default configuration file at `~/.illapaca.yaml`. You'll need to add your API key. You can get a free API key from [WeatherAPI.com](https://www.weatherapi.com/).

```bash
# Set your API key
illapaca --api-key=YOUR_API_KEY
```

## Usage

### Quick Examples

```bash
# Show current weather for a location
illapaca current "New York"

# Show forecast for the next 5 days
illapaca forecast "Tokyo" --days=5

# Display the full dashboard
illapaca dashboard "London"

# Compare weather between two locations
illapaca compare "Paris" "Rome"
```

## Commands

### Basic Commands

- `current`: Show current weather conditions
- `forecast`: Show weather forecast for next few days
- `dashboard`: Show complete weather dashboard
- `compare`: Compare weather between two locations

### Location Management

- `favorite list`: List all favorite locations
- `favorite add [location]`: Add a location to favorites
- `favorite remove [location/index]`: Remove a location from favorites
- `favorite set-default [location/index]`: Set a location as default

### Alert Management

- `alerts show`: Show current alert thresholds
- `alerts set`: Set alert thresholds

## Configuration

You can configure Illapaca with the following options:

- API key for weather service
- Default location
- Units (metric/imperial)
- Favorite locations
- Alert thresholds

Example configuration file:

```yaml
api_key: your_api_key_here
default_location: "New York"
units: metric
favorite_locations:
  - "New York"
  - "London"
  - "Tokyo"
alert_thresholds:
  high_temp: 30.0
  low_temp: 0.0
  precipitation: 70.0
  wind_speed: 30.0
```

## Development

### Dependencies

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [github.com/spf13/viper](https://github.com/spf13/viper) - Configuration management
- [github.com/fatih/color](https://github.com/fatih/color) - Terminal color output
- [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) - ASCII table rendering

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Weather data provided by [WeatherAPI.com](https://www.weatherapi.com/)
