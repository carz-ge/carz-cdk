package utils

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strings"
)

func SubStringTheDataLine(line *string, startChar string, endChar string) (string, error) {
	lineValue := *line

	startIndex := strings.Index(lineValue, startChar) - 1 // Find the index of the first occurrence of '['
	endIndex := strings.LastIndex(lineValue, endChar) + 1 // Find the index of the last occurrence of ']'

	if startIndex == -2 || endIndex == 0 || endIndex <= startIndex {
		log.Println("No valid substring found")
		return "", errors.New("no valid substring found")
	}

	substr := lineValue[startIndex+1 : endIndex] // Extract the substring
	return substr, nil
}

func GetDataLine(data io.ReadCloser, lineNumber int) (*string, error) {
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
