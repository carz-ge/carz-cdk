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

type ServiceType struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Code  string `json:"code"`
}

type SocarStationItem struct {
	Id        string `json:"id"`
	Lang      string `json:"lang"`
	Title     string `json:"title"`
	Code      string `json:"code"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	RegionId  string `json:"region_id"`
	Text      string `json:"text"`
	Publish   string `json:"publish"`

	ProductType interface{} `json:"product_type"`
	ObjectType  interface{} `json:"object_type"`

	PaymentType  interface{} `json:"payment_type"`
	ServiceTypes interface{} `json:"service_type"`
}

func GetSocarStationsEntities() (stations []models.AutoStationEntity, err error) {
	stationsKa, err := GetSocarStations("ge")
	if err != nil {
		return
	}

	stationsEn, err := GetSocarStations("en")
	if err != nil {
		return
	}

	if len(stationsKa) != len(stationsEn) {
		return nil, fmt.Errorf("socar stetions are not equal KA: %d vs EN: %d", len(stationsKa), len(stationsEn))
	}

	for index, stationKa := range stationsKa {
		stationEn := stationsEn[index]
		objects, err := convertToStationType(stationKa.ObjectType)
		if err != nil {
			return nil, err
		}
		products, err := convertToStationType(stationKa.ProductType)
		if err != nil {
			return nil, err
		}
		services, err := convertToStationType(stationKa.ServiceTypes)
		if err != nil {
			return nil, err
		}
		payments, err := convertToStationType(stationKa.PaymentType)
		if err != nil {
			return nil, err
		}

		stations = append(stations, models.AutoStationEntity{
			IdByProvider: stationKa.Id,
			ProviderCode: "SOCAR",
			StationType:  "AUTO_STATION",
			Name:         stationKa.Title,
			NameEn:       stationEn.Title,
			TextHtml:     stationKa.Text,
			TextHtmlEn:   stationEn.Text,
			Active:       stationKa.Publish == "1",
			Latitude:     stationKa.Latitude,
			Longitude:    stationKa.Longitude,

			ObjectTypes:  objects,
			ProductTypes: products,
			ServiceTypes: services,
			PaymentTypes: payments,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}
	return
}

func convertToStationType(v interface{}) (services []byte, err error) {
	switch v := v.(type) {
	case map[string]interface{}:
		var services []models.ServiceType

		// it's an object
		for _, val := range v {
			switch val := val.(type) {
			case map[string]interface{}:

				serviceType := models.ServiceType{}

				if val[`title`] != nil {
					serviceType.Title = val[`title`].(string)
				}
				if val[`code`] != nil {
					serviceType.Code = val[`code`].(string)
				}
				if val[`image`] != nil {
					serviceType.Id = val[`image`].(string)
				}
				services = append(services, serviceType)
				break
			default:
				continue
			}
		}
		return json.Marshal(services)

	case []interface{}:
		return
	default:
		err = fmt.Errorf("objectType is not map %s", v)
		return
	}
}

func GetSocarStations(lang string) ([]SocarStationItem, error) {
	databytes, err := makeResponseForSocar(lang)
	if err != nil {
		return nil, err
	}

	var stations []SocarStationItem

	unmarshallErr := json.Unmarshal(databytes, &stations)
	if unmarshallErr != nil {
		log.Println("Error parsing JSON:", unmarshallErr)
		return nil, unmarshallErr
	}
	return stations, nil
}

func makeResponseForSocar(lang string) ([]byte, error) {
	source := fmt.Sprintf("https://www.sgp.ge/%s/map/getResult", lang)
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
