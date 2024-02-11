package main

import (
	"fmt"
	"net/url"
	"regexp"
	"sure/simple-google-images-scraper/pkg/utils"
	"time"
)

func main() {
	fileName := "../assets/input.csv"
	maxFileNameLength := 100

	validExtensionsMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".webp": true,
	}

	for k, v := range validExtensionsMap {
		if !v {
			delete(validExtensionsMap, k)
		}
	}

	lines, err := utils.ReadCSV(fileName)
	if err != nil {
		fmt.Printf("error reading csv: %v\n", err)
		return
	}

	uniqueUrlMap := map[string]bool{}

	emptyUrlCount := 0
	duplicateUrls := []string{}

	downloadCount := 0
	errorCount := 0
	skipCount := 0

	for _, line := range lines {
		imageUrl, err := url.Parse(line)
		if err != nil {
			fmt.Printf("error parsing url [%s]: %v\n", line, err)
			continue
		}

		targetUrl := imageUrl.Query().Get("imgurl")

		if targetUrl == "" {
			emptyUrlCount++
			continue
		}

		if uniqueUrlMap[targetUrl] {
			duplicateUrls = append(duplicateUrls, targetUrl)
		} else {
			uniqueUrlMap[targetUrl] = true
		}
	}

	fmt.Printf("empty url count: %d\n", emptyUrlCount)
	fmt.Printf("duplicate url count: %d:\n", len(duplicateUrls))

	for _, duplicateUrl := range duplicateUrls {
		fmt.Printf("  %s\n", duplicateUrl)
	}

	uniqueUrlSlice := []string{}
	for key := range uniqueUrlMap {
		uniqueUrlSlice = append(uniqueUrlSlice, key)
	}

	for index, uniqueUrl := range uniqueUrlSlice {
		escapedPath := regexp.MustCompile(`(?m)["<>|:*?\\/]`).ReplaceAllString(uniqueUrl, "-")

		skipped, err := downloadImage(
			uniqueUrl,
			escapedPath,
			validExtensionsMap,
			maxFileNameLength,
			index,
			len(uniqueUrlSlice),
		)
		if err != nil {
			errorCount++
			fmt.Printf("[%d/%d] errored [%s]\n", index+1, len(uniqueUrlSlice), uniqueUrl)
			continue
		}

		if skipped {
			skipCount++
		} else {
			downloadCount++
		}
	}

	fmt.Printf("skip count: %d\n", skipCount)
	fmt.Printf("download count: %d\n", downloadCount)
	fmt.Printf("error count: %d\n", errorCount)
}

func downloadImage(
	imageUrl string,
	imageName string,
	validExtensionsMap map[string]bool,
	maxFileNameLength int,
	index int,
	totalCount int,
) (bool, error) {
	if len(imageName) >= maxFileNameLength {
		imageName = imageName[:maxFileNameLength]
	}

	if utils.OutputExists(imageName, validExtensionsMap) {
		fmt.Printf("[%d/%d] skipped [%s]\n", index+1, totalCount, imageUrl)
		return true, nil
	}

	then := time.Now().UnixMilli()

	imageBytes, err := utils.GetBytesFromUrl(imageUrl)
	if err != nil {
		return false, err
	}

	err = utils.SaveFile(imageBytes, imageName, validExtensionsMap)
	if err != nil {
		return false, err
	}

	now := time.Now().UnixMilli()

	fmt.Printf("[%d/%d] [%d] ms [%d] KB downloaded [%s]\n", index+1, totalCount, now-then, len(imageBytes)/1000, imageUrl)

	return false, nil
}
