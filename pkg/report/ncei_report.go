package report

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/espm1000/observing-the-weather/pkg/ncei"
)

type Params struct {
	StartDate       string
	EndDate         string
	StationId       string
	Dataset         string
	Units           string
	Limit           string
	IncludeMetadata string
}

type Report struct {
	header   []string
	datarows ncei.WeatherData
}

func NCEIReport(p Params, token string) (*Report, error) {
	log.Println("generating ncei weather report")
	rd, err := getRawData(p, token)
	if err != nil {
		return nil, err
	}

	return rd, nil
}

func BuildNCEI(token string) *ncei.Config {
	return ncei.New(token)
}

func getRawData(p Params, nceiToken string) (*Report, error) {
	cfg := BuildNCEI(nceiToken)

	params := url.Values{}
	params.Add("startdate", p.StartDate)
	params.Add("enddate", p.EndDate)
	params.Add("datasetid", p.Dataset)
	params.Add("stationid", p.Dataset+":"+p.StationId)
	params.Add("units", p.Units)
	params.Add("limit", p.Limit)
	params.Add("includemetadata", p.IncludeMetadata)

	url := cfg.BaseURL + cfg.DataEP + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("error generating http request: ", err)
	}

	req.Header.Add("token", nceiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error fetching ncei data: ", err)
	}

	defer resp.Body.Close()
	// bytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("error decoding body: ", err)
	// }
	// fmt.Println(string(bytes))
	dataResponse := ncei.WeatherResponse{}
	err = json.NewDecoder(resp.Body).Decode(&dataResponse)
	if err != nil {
		log.Println("error decoding json: ", err)
		return nil, err
	}
	fmt.Printf("length of response: %v\n", len(dataResponse.Results))

	return nil, nil
}
