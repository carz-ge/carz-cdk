package stations

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/src/models"
	"net/http"
	"time"
)

type RompetrolStationItem struct {
	Id        int16   `json:"id"`
	Lang      string  `json:"lang"`
	Name      string  `json:"name"`
	NameEn    string  `json:"name_en"`
	County    string  `json:"county"`
	CountyEn  string  `json:"county_en"`
	City      string  `json:"city"`
	CityEn    string  `json:"city_en"`
	Address   string  `json:"address"`
	AddressEn string  `json:"address_en"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`

	Type       string   `json:"type"`
	Services   []string `json:"services"`
	Infowindow string   `json:"infowindow"`
}

func GetRompetrolStationsEntities() (stations []models.AutoStationEntity, err error) {
	stationsKa, err := GetRompetrolStations()
	if err != nil {
		return
	}

	for _, stationKa := range stationsKa {
		//var products []models.ServiceType
		//
		//for _, fuel := range stationKa.Services {
		//	products = append(products, models.ServiceType{
		//		Code: fuel,
		//	})
		//}

		stations = append(stations, models.AutoStationEntity{
			IdByProvider: fmt.Sprintf("%d", stationKa.Id),
			ProviderCode: "ROMPETROLL",
			StationType:  "AUTO_STATION",
			Name:         stationKa.Name,
			NameEn:       stationKa.NameEn,

			Address:   stationKa.Address,
			AddressEn: stationKa.AddressEn,

			City:   stationKa.City,
			CityEn: stationKa.CityEn,

			Region:   stationKa.County,
			RegionEn: stationKa.CountyEn,

			TextHtml:  []byte(stationKa.Infowindow),
			Active:    true,
			Latitude:  stationKa.Latitude,
			Longitude: stationKa.Longitude,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return
}

func GetRompetrolStations() ([]RompetrolStationItem, error) {
	databytes, err := makeResponseForRompetrol()
	if err != nil {
		return nil, err
	}

	var stations []RompetrolStationItem

	unmarshallErr := json.Unmarshal(databytes, &stations)
	if unmarshallErr != nil {
		log.Println("Error parsing JSON:", unmarshallErr)
		return nil, unmarshallErr
	}
	return stations, nil
}

func makeResponseForRompetrol() ([]byte, error) {
	source := "https://www.rompetrol.ge/routeplanner/stations?language_id=1"
	// Create a transport with insecure skip verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	httpClient := &http.Client{Transport: transport, Timeout: time.Second * 5}

	req, err := http.NewRequest(http.MethodPost, source, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("sec-fetch-site", "none")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(res.Body)
}
