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
	"strconv"
	"time"
)

type ChargerItem struct {
	Name      string `json:"mtitle"`
	Latitude  string `json:"mlocation"`
	Longitude string `json:"mlocatioy"`
}

func GetChargerEntities() (stations []models.AutoStationEntity, err error) {
	chargers, err := GetChargers()

	for index, charger := range chargers {
		lat, err := strconv.ParseFloat(charger.Latitude, 64)
		if err != nil {
			return nil, err
		}
		lng, err := strconv.ParseFloat(charger.Longitude, 64)
		if err != nil {
			return nil, err
		}
		stations = append(stations, models.AutoStationEntity{
			IdByProvider: fmt.Sprintf("%d", index),
			Name:         charger.Name,
			NameEn:       charger.Name,
			ProviderCode: "ESPACE",
			StationType:  "EV_CHARGER",
			Active:       true,
			Latitude:     lat,
			Longitude:    lng,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}
	return
}

func GetChargers() ([]ChargerItem, error) {
	substr, err := getChargersJsonData()
	if err != nil {
		return nil, err
	}

	var stations []ChargerItem

	unmarshallErr := json.Unmarshal([]byte(substr), &stations)
	if unmarshallErr != nil {
		log.Println("Error parsing JSON:", unmarshallErr)
		return nil, unmarshallErr
	}
	return stations, nil
}

func getChargersJsonData() (string, error) {
	resBody, err := makeResponseForChargers()
	if err != nil {
		return "", err
	}

	line, lineErr := utils.GetDataLine(resBody, 465)

	if lineErr != nil {
		return "", lineErr
	}

	substr, err := utils.SubStringTheDataLine(line, "[", "]")
	if err != nil {
		return "", err
	}
	return substr, nil
}

func makeResponseForChargers() (io.ReadCloser, error) {
	source := "http://e-space.ge/%e1%83%a0%e1%83%a3%e1%83%99%e1%83%90/"
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
		log.Fatal(getErr)
	}
	return res.Body, nil
}
