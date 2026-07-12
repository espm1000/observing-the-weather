package ncei

const (
	nceiBaseUrl      = "https://www.ncei.noaa.gov/cdo-web/api/v2"
	dataEndpoint     = "/data"
	dataSetEndpoint  = "/datasets"
	dataTypeEndpoint = "/datatypes"
	locationEndpoint = "/locations"
	stationEndpoint  = "/stations"
)

type Config struct {
	BaseURL    string
	DataEP     string
	DatasetEP  string
	DatatypeEP string
	LocationEP string
	StationEP  string
	Token      string
}

type WeatherResponse struct {
	Results []WeatherResults `json:"results"`
}

type WeatherResults struct {
	Date       string  `json:"date"`
	Datatype   string  `json:"datatype"`
	Station    string  `json:"station"`
	Attributes string  `json:"attributes"`
	Value      float64 `json:"value"`
}

type WeatherData struct {
	Date     string `json:"date"`
	Location string `json:"location"`
	MaxTemp  string `json:"maxTemp"`
	MinTemp  string `json:"minTemp"`
	Precip   bool   `json:"precip"`
}

func New(token string) *Config {
	return &Config{
		BaseURL:    nceiBaseUrl,
		DataEP:     dataEndpoint,
		DatasetEP:  dataSetEndpoint,
		DatatypeEP: dataTypeEndpoint,
		LocationEP: locationEndpoint,
		StationEP:  stationEndpoint,
		Token:      token,
	}
}
