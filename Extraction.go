package main

import (
	"errors"
	"fmt"
	"regexp"
)

type Extraction struct {
	SourceName string   `json:"sourceName"`
	Headlines  []string `json:"headlines"`
	LastError  string
}

func (extraction *Extraction) ExtractHeadlinesFromHTML(html string) error {
	regexString, err := extraction.getRegexForSource()
	if err != nil {
		extraction.LastError = err.Error()
		return err
	}

	regExp, err := regexp.Compile(regexString)
	if err != nil {
		extraction.LastError = err.Error()
		return err
	}

	x := regExp.FindAllStringSubmatch(html, -1)

	for i := 0; i < len(x); i++ {
		extraction.Headlines = append(extraction.Headlines, x[i][1])
	}

	return nil
}

func (extraction Extraction) getRegexForSource() (string, error) {
	if extraction.SourceName == "reuters" {
		return `media-story-card__heading__eqhp9"><span>([^<]*)`, nil
	} else if extraction.SourceName == "associatepress" {
		return `-cardHeading">([^<]*)`, nil
	}

	return "", errors.New(
		fmt.Sprintf("Source [%s] does not have an extractor.", extraction.SourceName),
	)
}
