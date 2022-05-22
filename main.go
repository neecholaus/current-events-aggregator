package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	pulls := map[string]string{
		"reuters":        "https://reuters.com",
		"associatepress": "https://apnews.com",
	}

	for sourceName, sourceHost := range pulls {
		// Pull html
		res, err := http.Get(sourceHost)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		newsPull := Extraction{SourceName: sourceName}

		newsPull.ExtractHeadlinesFromHTML(string(body))

		if newsPull.LastError != "" {
			fmt.Println("Failed to extract headlines.")
			fmt.Println(newsPull.LastError)
			return
		}

		for i := 0; i < len(newsPull.Headlines); i++ {
			fmt.Println(fmt.Sprintf("%s\t%s", newsPull.SourceName, newsPull.Headlines[i]))
		}
	}
}
