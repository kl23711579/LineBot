package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func GetData() string {
	url := "http://www.bom.gov.au/qld/forecasts/brisbane.shtml"

	result := ""

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Fatal Error:", err.Error())
		os.Exit(0)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".forecast").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// get min temperature
		min := s.Find(".min").Text()
		if min == "" {
			min = "Unknown"
		}

		// get max temperature
		max := s.Find(".max").Text()
		if max == "" {
			max = "Unknown"
		}

		//get chance of rain
		rainpop := s.Find(".rain").Find(".pop").Text()
		ele := regexp.MustCompile("\\d+")
		rainProbability := ele.Find([]byte(rainpop))

		rainNumber, err := strconv.ParseInt(string(rainProbability), 10, 0)
		if err != nil {
			fmt.Println(err)
		}

		rainHint := GetRainHint(rainNumber)
		rain := string(rainProbability) + "%    " + rainHint

		// get short summary
		shortSummary := s.Find(".summary").Text()
		if shortSummary == "" {
			shortSummary = "Unknown";
		}

		longSummary := s.Find("p").Text()

		result = GetWeatherDetail([]string {shortSummary, min, max, rain, longSummary})
		//fmt.Println(result)

		return false
	})
	//doc.Find(".title").Each(func(i int, s *goquery.Selection) {
	//	title := s.Find("a").Text()
	//	href, exist := s.Find("a").Attr("href")
	//	if exist {
	//		// fmt.Println("https://www.ptt.cc/" + href + "   =>   " + title)
	//		time := transTime(href)
	//		titleResult = append(titleResult, time+"  =>  "+title)
	//
	//		hrefResult = append(hrefResult, "https://www.ptt.cc/"+href)
	//	}
	//})

	return result
}

func GetRainHint(rainProbability int64) string {
	var hint string

	if rainProbability <= 30 {
		hint = "No necessary to bring Umbrella."
	} else if rainProbability >30 && rainProbability <= 50 {
		hint = "Maybe bring umbrella."
	} else if rainProbability > 50 && rainProbability <= 80 {
		hint = "Recommendation for bringing umbrella."
	} else if rainProbability > 80 &&rainProbability <= 100 {
		hint = "Please bring umbrella."
	} else {
		hint = "Something wrong."
	}

	return hint
}

func GetWeatherDetail(wIndividual []string) string {
	// short min max rain summary
	detail := "Today's weather is " + wIndividual[0] + "\nTemperature: \n  min: " + wIndividual[1] + " \n  Max: " + wIndividual[2]
	detail += " \nRain: " + wIndividual[3] + " \nSummary: " + wIndividual[4]
	return detail
}