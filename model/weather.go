package model

type WeatherData struct {
}

type CurrentWeather struct {
	TempC      float64   `json:"temp_c"`
	TempF      float64   `json:"temp_f"`
	IsDay      int       `json:"is_day"`
	Condition  Condition `json:"condition"`
	WindKph    float64   `json:"wind_kph"`
	WindMph    float64   `json:"wind_mph"`
	WindDir    string    `json:"wind_dir"`
	PressureMb float64   `json:"pressure_mb"`
	PrecipMm   float64   `json:"precip_mm"`
	FeelslikeC float64   `json:"feelslike_c"`
	FeelslikeF float64   `json:"feelslike_f"`
	VisKm      float64   `json:"vis_km"`
	UV         float64   `json:"uv"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Localtime      string  `json:"localtime"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
}

type Forecast struct {
	Forecastday []Forecastday `json:"forecastday"`
}

type Forecastday struct {
	Date      string `json:"date"`
	DateEpoch int    `json:"date_epoch"`
	Day       Day    `json:"day"`
	Astro     Astro  `json:"astro"`
	Hour      []Hour `json:"hour"`
}

type Day struct {
	MaxTempC          float64   `json:"maxtemp_c"`
	MinTempC          float64   `json:"mintemp_c"`
	AvgTempC          float64   `json:"avgtemp_c"`
	MaxWindKph        float64   `json:"maxwind_kph"`
	TotalPrecipMm     float64   `json:"totalprecip_mm"`
	DailyChanceOfRain int       `json:"daily_chance_of_rain"`
	Condition         Condition `json:"condition"`
}

type Astro struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

type Hour struct {
	TimeEpoch    int64     `json:"time_epoch"`
	Time         string    `json:"time"`
	TempC        float64   `json:"temp_c"`
	Condition    Condition `json:"condition"`
	ChanceOfRain int       `json:"chance_of_rain"`
}

// HistoricalData for comparison
type HistoricalData struct {
	Location Location `json:"location"`
	Forecast Forecast `json:"forecast"`
}
