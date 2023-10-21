package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Function to append hidden input field values to the URL in URL format
func appendHiddenInputValues(url string, doc *goquery.Document) string {
	paramsAppended := make(map[string]bool) // Keep track of appended parameters

	// Initialize the modified URL with the base URL
	modifiedURL := url

	doc.Find("input[name]").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")

		// Check if the parameter has already been appended
		if _, exists := paramsAppended[name]; !exists {
			// Append the parameter to the modified URL
			if strings.Contains(modifiedURL, "?") {
				modifiedURL += "&"
			} else {
				modifiedURL += "?"
			}
			modifiedURL += fmt.Sprintf("%s=%s", name, value)
			paramsAppended[name] = true // Mark as appended
		}
	})

	return modifiedURL
}

func main() {
	// Read the URL list from urls.txt
	urlsData, err := ioutil.ReadFile("urls.txt")
	if err != nil {
		fmt.Println("Error reading urls.txt:", err)
		os.Exit(1)
	}

	urlList := strings.Split(string(urlsData), "\n")

	// Create or open the output file for writing
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating output.txt:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	for _, url := range urlList {
		// Send an HTTP GET request
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error requesting URL:", err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			// Parse the HTML content of the response
			doc, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				fmt.Println("Error parsing HTML:", err)
				continue
			}

			// Append hidden input field values to the URL in URL format
			modifiedURL := appendHiddenInputValues(url, doc)

			// Print the modified URL
			fmt.Println("Modified URL:", modifiedURL)

			// Write the modified URL to the output file
			_, err = outputFile.WriteString(modifiedURL + "\n")
			if err != nil {
				fmt.Println("Error writing to output.txt:", err)
			}
		} else {
			fmt.Printf("Failed to retrieve the page. Status code: %d\n", response.StatusCode)
		}
	}
}
