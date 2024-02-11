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

func OutputExists(fileName string, validExtensionMap map[string]bool) bool {
	for extension := range validExtensionMap {
		newFilePath := fmt.Sprintf("../assets/output/%s%s", fileName, extension)
		if _, err := os.Stat(newFilePath); err == nil {
			return true
		}
	}

	return false
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

func SaveFile(fileBytes []byte, fileName string, validExtensionsMap map[string]bool) error {
	extension := mimetype.Detect(fileBytes).Extension()
	if !validExtensionsMap[extension] {
		return fmt.Errorf("invalid extension [%s] for file [%s]", extension, fileName)
	}

	file, err := os.Create(fmt.Sprintf("../assets/output/%s%s", fileName, extension))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(fileBytes))
	return err
}
