package ui

import (
	"strings"
)

// WeatherIcons maps weather condition text to appropriate emoji
var weatherIcons = map[string]string{
	"sunny":         "☀️",
	"clear":         "☀️",
	"partly cloudy": "⛅",
	"cloudy":        "☁️",
	"overcast":      "☁️",
	"mist":          "🌫️",
	"rain":          "🌧️",
	"snow":          "❄️",
	"storm":         "⛈️",
	"fog":           "🌫️",
	"thunderstorm":  "🌩️",
	"drizzle":       "🌦️",
	"ice":           "🧊",
	"sleet":         "🌨️",
	"hail":          "🌨️",
	"windy":         "💨",
	"tornado":       "🌪️",
	"hurricane":     "🌀",
}

// GetConditionIcon returns the appropriate weather icon for a condition
func GetConditionIcon(condition string) string {
	conditionLower := strings.ToLower(condition)

	for k, v := range weatherIcons {
		if strings.Contains(conditionLower, k) {
			return v
		}
	}

	// Default icon if no match
	return "🌡️"
}
