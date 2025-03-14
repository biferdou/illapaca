// model/weather.go
package model

// Original models - we'll keep these as our unified internal format
type WeatherData struct {
	Current  CurrentWeather `json:"current"`
	Location Location       `json:"location"`
	Forecast Forecast       `json:"forecast"`
}

type CurrentWeather struct {
	TempC      float64   `json:"temp_c"`
	TempF      float64   `json:"temp_f"`
	IsDay      int       `json:"is_day"`
	Condition  Condition `json:"condition"`
	WindMph    float64   `json:"wind_mph"`
	WindKph    float64   `json:"wind_kph"`
	WindDir    string    `json:"wind_dir"`
	PressureMb float64   `json:"pressure_mb"`
	PrecipMm   float64   `json:"precip_mm"`
	Humidity   int       `json:"humidity"`
	FeelsLikeC float64   `json:"feelslike_c"`
	FeelsLikeF float64   `json:"feelslike_f"`
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
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Forecast struct {
	ForecastDay []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date      string `json:"date"`
	DateEpoch int64  `json:"date_epoch"`
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

// OpenWeatherMap specific models
type Coordinates struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Name    string
	Country string
}

type GeoLocation struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type OpenWeatherCurrent struct {
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Rain struct {
		OneHour float64 `json:"1h"`
	} `json:"rain"`
	Snow struct {
		OneHour float64 `json:"1h"`
	} `json:"snow"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

type OpenWeatherForecast struct {
	List []OpenWeatherForecastItem `json:"list"`
	City struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Coord   struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Sunrise  int `json:"sunrise"`
		Sunset   int `json:"sunset"`
		Timezone int `json:"timezone"`
	} `json:"city"`
}

type OpenWeatherForecastItem struct {
	Dt   int `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
		TempKf    float64 `json:"temp_kf"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float64 `json:"pop"`
	Rain       struct {
		ThreeHour float64 `json:"3h"`
	} `json:"rain"`
	Snow struct {
		ThreeHour float64 `json:"3h"`
	} `json:"snow"`
	Sys struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
}

type OpenWeatherHistorical struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
	Current  struct {
		Dt         int     `json:"dt"`
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Hourly []OpenWeatherHourlyItem `json:"hourly"`
}

type OpenWeatherHourlyItem struct {
	Dt         int     `json:"dt"`
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feels_like"`
	Pressure   int     `json:"pressure"`
	Humidity   int     `json:"humidity"`
	DewPoint   float64 `json:"dew_point"`
	Clouds     int     `json:"clouds"`
	Visibility int     `json:"visibility"`
	WindSpeed  float64 `json:"wind_speed"`
	WindDeg    int     `json:"wind_deg"`
	WindGust   float64 `json:"wind_gust"`
	Weather    []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Pop  float64 `json:"pop"`
	Rain struct {
		OneHour float64 `json:"1h"`
	} `json:"rain"`
	Snow struct {
		OneHour float64 `json:"1h"`
	} `json:"snow"`
}
