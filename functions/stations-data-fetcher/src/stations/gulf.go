package stations

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/src/models"
	"main/src/utils"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type GulfStationItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      string `json:"is_active"`
	Region      string `json:"region"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`

	FuelTypes []string `json:"fuel_types"`
	PoiTypes  []string `json:"poi_types"`
	FoodTypes []string `json:"food_types"`
	Picture   string   `json:"picture"`
}

func GetGulfStationsEntities() (stations []models.AutoStationEntity, err error) {
	stationsKa, err := GetGulfStations("ka")
	if err != nil {
		return
	}
	sort.Slice(stationsKa, func(i, j int) bool {
		return stationsKa[i].Id < stationsKa[j].Id
	})
	stationsEn, err := GetGulfStations("en")
	if err != nil {
		return
	}
	sort.Slice(stationsEn, func(i, j int) bool {
		return stationsEn[i].Id < stationsEn[j].Id
	})
	if len(stationsKa) != len(stationsEn) {
		return nil, fmt.Errorf("socar stetions are not equal KA: %d vs EN: %d", len(stationsKa), len(stationsEn))
	}

	for index, stationKa := range stationsKa {
		stationEn := stationsEn[index]
		var products []models.ServiceType

		for _, fuel := range stationKa.FuelTypes {
			products = append(products, models.ServiceType{
				Code: fuel,
			})
		}

		productsJson, err := json.Marshal(products)
		if err != nil {
			return nil, err
		}

		var objects []models.ServiceType

		for _, poi := range stationKa.PoiTypes {
			objects = append(objects, models.ServiceType{
				Code: poi,
			})
		}

		objectsJson, err := json.Marshal(objects)
		if err != nil {
			return nil, err
		}
		lat, err := strconv.ParseFloat(stationKa.Latitude, 64)
		if err != nil {
			return nil, err
		}
		lng, err := strconv.ParseFloat(stationKa.Longitude, 64)
		if err != nil {
			return nil, err
		}
		stations = append(stations, models.AutoStationEntity{
			IdByProvider: stationKa.Id,
			ProviderCode: "GULF",
			StationType:  "AUTO_STATION",
			Name:         stationKa.Name,
			NameEn:       stationEn.Name,

			Description:   stationKa.Description,
			DescriptionEn: stationEn.Description,
			Active:        stationKa.Active == "1",
			Latitude:      lat,
			Longitude:     lng,
			Picture:       stationKa.Picture,

			ProductTypes: productsJson,
			ObjectTypes:  objectsJson,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}
	return
}
func GetGulfStations(lang string) ([]GulfStationItem, error) {
	substr, err := getGulfJsonData(lang)
	if err != nil {
		return nil, err
	}

	var stationsMap map[string]GulfStationItem

	unmarshallErr := json.Unmarshal([]byte(substr), &stationsMap)
	if unmarshallErr != nil {
		log.Println("Error parsing JSON:", unmarshallErr)
		return nil, unmarshallErr
	}

	var stations []GulfStationItem
	for _, value := range stationsMap {
		stations = append(stations, value)
	}
	return stations, nil
}

func getGulfJsonData(lang string) (string, error) {
	resBody, err := makeResponseForGulf(lang)

	if err != nil {
		return "", err
	}
	line, lineErr := utils.GetDataLine(resBody, 169)

	if lineErr != nil {
		return "", lineErr
	}

	substr, err := utils.SubStringTheDataLine(line, "{", "}")
	if err != nil {
		return "", err
	}
	return substr, nil
}

func makeResponseForGulf(lang string) (io.ReadCloser, error) {
	source := fmt.Sprintf("https://gulf.ge/%s/map", lang)
	// Create a transport with insecure skip verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	httpClient := &http.Client{Transport: transport, Timeout: time.Second * 5}

	req, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	return res.Body, nil
}
