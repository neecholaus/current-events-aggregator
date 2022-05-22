package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	pulls := map[string]string{
		"reuters":        "https://reuters.com",
		"associatepress": "https://apnews.com",
		"theeconomist":   "https://www.economist.com",
	}

	extractions := []Extraction{}

	for sourceName, sourceHost := range pulls {
		// Pull html
		res, err := http.Get(sourceHost)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		defer res.Body.Close()

		// Isolate response body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		newsPull := Extraction{
			SourceName: sourceName,
			CreatedAt:  time.Now(),
		}

		newsPull.ExtractHeadlinesFromHTML(string(body))
		if newsPull.LastError != "" {
			fmt.Println("Failed to extract headlines.")
			fmt.Println(newsPull.LastError)
			return
		}

		extractions = append(extractions, newsPull)
	}

	storeExtractions(extractions)
}

func storeExtractions(extractions []Extraction) {
	// Create storage for current pull
	filepath := fmt.Sprintf("./headlines/%s.tsv", extractions[0].CreatedAt.Format(time.UnixDate))
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0640))
	if err != nil {
		panic(err)
	}

	// Create slice of formatted lines
	var entries []string
	for _, extraction := range extractions {
		for _, headline := range extraction.Headlines {
			entries = append(entries, fmt.Sprintf("%s\t%s\t%s\n",
				extraction.CreatedAt.Format(time.UnixDate),
				extraction.SourceName,
				headline,
			))
		}
	}

	file.WriteString(strings.Join(entries, ""))
}
