package chargers

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ChargerItem struct {
	Name      string `json:"mtitle"`
	Latitude  string `json:"mlocation"`
	Longitude string `json:"mlocatioy"`
}

func GetChargers() ([]ChargerItem, bool) {
	substr, done := getChargersJsonData()
	if done {
		return nil, true
	}

	var chargers []ChargerItem

	unmarshallErr := json.Unmarshal([]byte(substr), &chargers)
	if unmarshallErr != nil {
		log.Println("Error parsing JSON:", unmarshallErr)
		return nil, true
	}
	return chargers, false
}

func getChargersJsonData() (string, bool) {
	resBody := makeResponseForChargers()

	line, lineErr := getDataLine(resBody, 465)

	if lineErr != nil {
		log.Fatal(lineErr)
	}

	substr, failed := subStringTheDataLine(line)
	if failed {
		return "", true
	}
	return substr, false
}

func makeResponseForChargers() io.ReadCloser {
	source := "http://e-space.ge/%e1%83%a0%e1%83%a3%e1%83%99%e1%83%90/"
	// Create a transport with insecure skip verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	httpClient := &http.Client{Transport: transport, Timeout: time.Second * 5}

	req, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		log.Fatal(err)
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
	return res.Body
}

func subStringTheDataLine(line *string) (string, bool) {
	lineValue := *line

	startIndex := strings.Index(lineValue, "[") - 1   // Find the index of the first occurrence of '['
	endIndex := strings.LastIndex(lineValue, "]") + 1 // Find the index of the last occurrence of ']'

	if startIndex == -2 || endIndex == 0 || endIndex <= startIndex {
		log.Println("No valid substring found")
		return "", true
	}

	substr := lineValue[startIndex+1 : endIndex] // Extract the substring
	return substr, false
}

func getDataLine(data io.ReadCloser, lineNumber int) (*string, error) {
	reader := bufio.NewReader(data)

	currentLine := 1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading file:", err)
			return nil, err
		}

		if currentLine == lineNumber {
			return &line, nil
		}

		currentLine++
	}
}
