package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

type vulnerability struct {
	Line         string `json:"line"`
	Column       string `json:"column"`
	File         string `json:"file"`
	Code         string `json:"code"`
	Details      string `json:"details"`
	SecurityTool string `json:"securityTool"`
	Language     string `json:"language"`
	Severity     string `json:"severity"`
	Type         string `json:"type"`
	RuleId       string `json:"rule_id"`
	VulnHash     string `json:"vulnHash"`
}

func (v *vulnerability) IsMultiLine() bool {
	return v.Column == ""
}

func (v *vulnerability) StartLine() int {
	number, err := strconv.Atoi(v.Line)
	if err != nil {
		return -1
	}
	return number
}

func (v *vulnerability) EndLine() int {
	if v.IsMultiLine() {
		re := regexp.MustCompile("[0-9]+")
		numbersStr := re.FindAllString(v.Code, -1)
		if len(numbersStr) == 2 {
			number, err := strconv.Atoi(numbersStr[1])
			if err != nil {
				fmt.Println("Error convert string to integer")
				return -1
			}
			return number
		} else {
			fmt.Println("Error Parse EndLine")
			return -1
		}
	}

	number, err := strconv.Atoi(v.Line)
	if err != nil {
		return -1
	}
	return number
}

type result struct {
	vulnerability `json:"vulnerabilities"`
}

const resultsFile = "results.json"

func loadResultsFile() ([]result, error) {
	results := struct {
		AnalysisVulnerabilities []result `json:"analysisVulnerabilities"`
	}{}

	file, err := ioutil.ReadFile(resultsFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &results)
	if err != nil {
		return nil, err
	}
	return results.AnalysisVulnerabilities, nil
}
