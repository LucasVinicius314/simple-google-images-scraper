package main

import (
	"fmt"
	"net/url"
	"strconv"
	"sure/simple-google-images-scraper/pkg/utils"
)

func main() {
	filename := "../assets/input.csv"

	lines, err := utils.ReadCSV(filename)
	if err != nil {
		fmt.Printf("error reading csv: %v", err)
		return
	}

	for index, line := range lines {
		imageUrl, err := url.Parse(line)
		if err != nil {
			fmt.Printf("error parsing url [%s]: %v", line, err)
			continue
		}

		targetUrl := imageUrl.Query().Get("imgurl")

		if targetUrl == "" {
			continue
		}

		err = downloadImage(targetUrl, strconv.FormatInt(int64(index), 10))
		if err != nil {
			fmt.Printf("error downloading image [%s]: %v", line, err)
			continue
		}
	}
}

func downloadImage(imageUrl string, imageName string) error {
	imageBytes, err := utils.GetBytesFromUrl(imageUrl)
	if err != nil {
		return err
	}

	err = utils.SaveFile(imageBytes, imageName)
	if err != nil {
		return err
	}

	fmt.Printf("downloaded [%s] [%d] KB", imageUrl, len(imageBytes)/1000)

	return nil
}
