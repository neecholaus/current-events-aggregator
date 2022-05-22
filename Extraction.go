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
	// Get the regex string for the appropriate source.
	regexString, err := extraction.getRegexForSource()
	if err != nil {
		extraction.LastError = err.Error()
		return err
	}

	// Make regex object
	regex, err := regexp.Compile(regexString)
	if err != nil {
		extraction.LastError = err.Error()
		return err
	}

	matches := regex.FindAllStringSubmatch(html, -1)

	// Determine index of headline regex grouping
	indexOfHealineGrouping := 0
	for i, name := range regex.SubexpNames() {
		if name == "headline" {
			indexOfHealineGrouping = i
		}
	}

	// Assign headlines to self
	for i := 0; i < len(matches); i++ {
		if matches[i][indexOfHealineGrouping] != "" {
			extraction.Headlines = append(extraction.Headlines, matches[i][indexOfHealineGrouping])
		}
	}

	return nil
}

func (extraction Extraction) getRegexForSource() (string, error) {
	switch source := extraction.SourceName; source {
	case "reuters":
		return `media-story-card__heading__eqhp9"><span>(?P<headline>[^<]*)`, nil
	case "associatepress":
		return `-cardHeading">(?P<headline>[^<]*)`, nil
	case "theeconomist":
		return `data-analytics="(topical_content|top_stories)(_\d{0,2})?:headline_\d{1,2}">(?P<headline>[^<]*)`, nil
	}

	return "", errors.New(
		fmt.Sprintf("Source [%s] does not have an extractor.", extraction.SourceName),
	)
}
