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

var VisolFeatureMap = map[string]string{
	"1": "VISOL_TAG",
	"2": "GAS",
	"3": "SUPER_MARKET",
	"4": "CHANGE_OIL",
	"5": "CAR_WASH",
	"6": "UKNOWN",
}

type WissolStationItem struct {
	StationId string `json:"stationid"`
	Lang      string `json:"lang"`
	Poster    string `json:"poster"`
	Address   string `json:"address"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
	City      string `json:"city"`

	Fuels    []string `json:"fuels"`
	Features []string `json:"features"`

	CenterType string `json:"center_type"`
}

func GetWissolStationsEntities() (stations []models.AutoStationEntity, err error) {
	stationsKa, err := GetWissolStations("geo")
	if err != nil {
		return
	}

	stationsEn, err := GetWissolStations("eng")
	if err != nil {
		return
	}

	if len(stationsKa) != len(stationsEn) {
		return nil, fmt.Errorf("socar stetions are not equal KA: %d vs EN: %d", len(stationsKa), len(stationsEn))
	}

	for index, stationKa := range stationsKa {
		stationEn := stationsEn[index]
		var products []models.ServiceType

		for _, fuel := range stationKa.Fuels {
			products = append(products, models.ServiceType{
				Code: fuel,
			})
		}
		productsJson, err := json.Marshal(products)
		if err != nil {
			return nil, err
		}

		var objects []models.ServiceType

		for _, feature := range stationKa.Features {
			objects = append(objects, models.ServiceType{
				Code: VisolFeatureMap[feature],
			})
		}
		objectsJson, err := json.Marshal(objects)
		if err != nil {
			return nil, err
		}
		stations = append(stations, models.AutoStationEntity{
			IdByProvider: stationKa.StationId,
			ProviderCode: "WISSOL",
			StationType:  "AUTO_STATION",
			Name:         stationKa.Address,
			NameEn:       stationEn.Address,

			Address:   stationKa.Address,
			AddressEn: stationEn.Address,

			City:   stationKa.City,
			CityEn: stationEn.City,

			Active:    true,
			Latitude:  stationKa.Latitude,
			Longitude: stationKa.Longitude,
			Picture:   stationKa.Poster,

			ProductTypes: productsJson,
			ObjectTypes:  objectsJson,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}
	return
}

func GetWissolStations(lang string) ([]WissolStationItem, error) {
	databytes, err := makeResponseForWissol(lang)
	if err != nil {
		return nil, err
	}

	var stations []WissolStationItem

	unmarshallErr := json.Unmarshal(databytes, &stations)
	if unmarshallErr != nil {
		log.Println("Error parsing JSON:", unmarshallErr)
		return nil, unmarshallErr
	}
	return stations, nil
}

func makeResponseForWissol(lang string) ([]byte, error) {
	source := fmt.Sprintf("http://wissol.ge/adminarea/api/ajaxapi/map?search=&location_id=&location_id=&location_id=&location_id=&location_id=&lang=%s", lang)
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
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(res.Body)
}
