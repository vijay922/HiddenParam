package main

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "strings"

        "github.com/PuerkitoBio/goquery"
)

// Function to append hidden input field values to the URL
func appendHiddenInputValues(url string, doc *goquery.Document) string {
        doc.Find("input[name]").Each(func(i int, s *goquery.Selection) {
                name, _ := s.Attr("name")
                value, _ := s.Attr("value")
                url += fmt.Sprintf("?%s=%s&", name, value)
        })
        return url
}

func main() {
        // Read the URL list from urls.txt
        urlsData, err := ioutil.ReadFile("urls.txt")
        if err != nil {
                fmt.Println("Error reading urls.txt:", err)
                os.Exit(1)
        }

        urlList := strings.Split(string(urlsData), "\n")

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

                        // Append hidden input field values to the URL
                        modifiedURL := appendHiddenInputValues(url, doc)

                        // Print the modified URL
                        fmt.Println("Modified URL:", modifiedURL)

                        // Write the modified URL to a text file
                        outputFile, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                        if err != nil {
                                fmt.Println("Error opening output.txt:", err)
                                continue
                        }
                        defer outputFile.Close()

                        outputFile.WriteString(modifiedURL + "\n")
                } else {
                        fmt.Printf("Failed to retrieve the page. Status code: %d\n", response.StatusCode)
                }
        }
}
