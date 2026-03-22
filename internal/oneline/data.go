package oneline

// WeatherData represents the top-level structure of the weather API response.
type WeatherData struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	NearestArea      []NearestArea      `json:"nearest_area"`
	Request          []Request          `json:"request"`
	Weather          []WeatherDay       `json:"weather"`
}

// CurrentCondition represents the current observed weather conditions.
type CurrentCondition struct {
	FeelsLikeC       string      `json:"FeelsLikeC"`
	FeelsLikeF       string      `json:"FeelsLikeF"`
	Cloudcover       string      `json:"cloudcover"`
	Humidity         string      `json:"humidity"`
	LocalObsDateTime string      `json:"localObsDateTime"`
	ObservationTime  string      `json:"observation_time"`
	PrecipInches     string      `json:"precipInches"`
	PrecipMM         string      `json:"precipMM"`
	Pressure         string      `json:"pressure"`
	PressureInches   string      `json:"pressureInches"`
	TempC            string      `json:"temp_C"`
	TempF            string      `json:"temp_F"`
	UVIndex          string      `json:"uvIndex"`
	Visibility       string      `json:"visibility"`
	VisibilityMiles  string      `json:"visibilityMiles"`
	WeatherCode      string      `json:"weatherCode"`
	WeatherDesc      []ValueItem `json:"weatherDesc"`
	WeatherIconURL   []ValueItem `json:"weatherIconUrl"`
	Winddir16Point   string      `json:"winddir16Point"`
	WinddirDegree    string      `json:"winddirDegree"`
	WindspeedKmph    string      `json:"windspeedKmph"`
	WindspeedMiles   string      `json:"windspeedMiles"`
}

// NearestArea represents information about the nearest location/area.
type NearestArea struct {
	AreaName   []ValueItem `json:"areaName"`
	Country    []ValueItem `json:"country"`
	Latitude   string      `json:"latitude"`
	Longitude  string      `json:"longitude"`
	Population string      `json:"population"`
	Region     []ValueItem `json:"region"`
	WeatherURL []ValueItem `json:"weatherUrl"`
}

// Request represents the parameters used for the weather query.
type Request struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

// WeatherDay represents daily weather forecast data (one entry per date).
type WeatherDay struct {
	Astronomy   []Astronomy `json:"astronomy"`
	AvgTempC    string      `json:"avgtempC"`
	AvgTempF    string      `json:"avgtempF"`
	Date        string      `json:"date"`
	Hourly      []Hourly    `json:"hourly"`
	MaxTempC    string      `json:"maxtempC"`
	MaxTempF    string      `json:"maxtempF"`
	MinTempC    string      `json:"mintempC"`
	MinTempF    string      `json:"mintempF"`
	SunHour     string      `json:"sunHour"`
	TotalSnowCM string      `json:"totalSnow_cm"`
	UVIndex     string      `json:"uvIndex"`
}

// Astronomy contains sunrise, sunset, moon phase and related information for a day.
type Astronomy struct {
	MoonIllumination string `json:"moon_illumination"`
	MoonPhase        string `json:"moon_phase"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
}

// Hourly represents hourly weather forecast data.
type Hourly struct {
	DewPointC        string      `json:"DewPointC"`
	DewPointF        string      `json:"DewPointF"`
	FeelsLikeC       string      `json:"FeelsLikeC"`
	FeelsLikeF       string      `json:"FeelsLikeF"`
	HeatIndexC       string      `json:"HeatIndexC"`
	HeatIndexF       string      `json:"HeatIndexF"`
	WindChillC       string      `json:"WindChillC"`
	WindChillF       string      `json:"WindChillF"`
	WindGustKmph     string      `json:"WindGustKmph"`
	WindGustMiles    string      `json:"WindGustMiles"`
	ChanceOfFog      string      `json:"chanceoffog"`
	ChanceOfFrost    string      `json:"chanceoffrost"`
	ChanceOfHighTemp string      `json:"chanceofhightemp"`
	ChanceOfOvercast string      `json:"chanceofovercast"`
	ChanceOfRain     string      `json:"chanceofrain"`
	ChanceOfRemDry   string      `json:"chanceofremdry"`
	ChanceOfSnow     string      `json:"chanceofsnow"`
	ChanceOfSunshine string      `json:"chanceofsunshine"`
	ChanceOfThunder  string      `json:"chanceofthunder"`
	ChanceOfWindy    string      `json:"chanceofwindy"`
	Cloudcover       string      `json:"cloudcover"`
	DiffRad          string      `json:"diffRad"`
	Humidity         string      `json:"humidity"`
	PrecipInches     string      `json:"precipInches"`
	PrecipMM         string      `json:"precipMM"`
	Pressure         string      `json:"pressure"`
	PressureInches   string      `json:"pressureInches"`
	ShortRad         string      `json:"shortRad"`
	TempC            string      `json:"tempC"`
	TempF            string      `json:"tempF"`
	Time             string      `json:"time"` // usually "0", "300", "600", ..., "2100"
	UVIndex          string      `json:"uvIndex"`
	Visibility       string      `json:"visibility"`
	VisibilityMiles  string      `json:"visibilityMiles"`
	WeatherCode      string      `json:"weatherCode"`
	WeatherDesc      []ValueItem `json:"weatherDesc"`
	WeatherIconURL   []ValueItem `json:"weatherIconUrl"`
	Winddir16Point   string      `json:"winddir16Point"`
	WinddirDegree    string      `json:"winddirDegree"`
	WindspeedKmph    string      `json:"windspeedKmph"`
	WindspeedMiles   string      `json:"windspeedMiles"`
}

// ValueItem is a common wrapper used for most description/icon fields
// (single object with "value" key).
type ValueItem struct {
	Value string `json:"value"`
}
