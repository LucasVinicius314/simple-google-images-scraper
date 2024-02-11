package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func GetBytesFromUrl(imageUrl string) ([]byte, error) {
	response, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func ReadCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		lines = append(lines, record...)
	}

	return lines, nil
}

func SaveFile(fileBytes []byte, fileName string) error {
	newFilePath := fmt.Sprintf("../assets/output/%s%s", fileName, mimetype.Detect(fileBytes).Extension())

	file, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(fileBytes))
	return err
}
