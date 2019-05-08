package main

import (
	"encoding/json"
	"fmt"
	"go_spider/core/common/page"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/YonnaR/go-asciibot"
)

/*
TripAdvisorProcessor structure data reception and navigation between the in crawlling pages
==================================================================
	@startURI is for the CLI, define start crawling uri
	@nextURI is used for navigation between pages
	@output is the desired output for the data
	@data is the data collected during process
*/
type TripAdvisorProcessor struct {
	startURI string
	nextURI  string
	output   string
	data     []string
}

/*NewTripAdvisorProcessor return new instance of TA processor  */
func NewTripAdvisorProcessor(url string, output string) *TripAdvisorProcessor {
	return &TripAdvisorProcessor{
		startURI: url,
		output:   output,
	}
}

/*Process is life cycle spider
==================================================================
	=>Get first page
	=>Parse html dom
	=>Check pagination for get other pages
	=>Collect data and store it to processor array
Package goquery is used to parse html.
*/
func (t *TripAdvisorProcessor) Process(p *page.Page) {

	//if page is crawled successfully
	if p.IsSucc() {

		//Parse HTML current page with goquery
		query := p.GetHtmlParser()
		page := query.Find(".ui_pagination")
		cPage := t.nextURI

		//Set next page
		getNextPage(page, t)

		//if next page is a valid url
		if cPage != t.nextURI && t.nextURI != "" {
			p.AddTargetRequest("https://www.tripadvisor.fr"+t.nextURI, "html")
		}

		/*
			Check gived uri and call proper parser.
		*/
		if strings.Contains(t.startURI, "Restaurant") {
			review := query.Find(".review-container")
			review.Each(func(i int, s *goquery.Selection) {
				c := restaurantParser(s)
				d, _ := json.Marshal(c)
				t.addToPipeLine(string(d))
			})

		} else if strings.Contains(t.startURI, "Hotel") {
			getNextPage(page, t)
			review := query.Find(".hotels-hotel-review-community-content-Card__ui_card--3kTH_")
			review.Each(func(i int, s *goquery.Selection) {
				c := hotelParser(s)
				d, _ := json.Marshal(c)
				t.addToPipeLine(string(d))
			})

		} else { // can't recognize uri
			log.Fatal("can't handle requested uri")
		}

	} else {
		log.Fatal("can't fetch website, check your internet connection")
		t.Finish()
	}

}

/*
Finish is last task of spider function
==================================================================
	Write data on file when process is end
*/
func (t *TripAdvisorProcessor) Finish() {

	fmt.Printf("\n\r Writing Data ... \n")

	//Serialise data
	file, err := json.Marshal(t.data)
	if err != nil {
		fmt.Println("err : ", err.Error())
	}

	//Write data to output file
	werr := ioutil.WriteFile(t.output, file, 0644)
	if werr != nil {
		fmt.Println("err : ", err.Error())
	}

	//Last Message
	fmt.Printf(asciibot.Random() + "\n\r Spidering ended \n")
}

/* addToPipeLine add a serialized commentary to data[]*/
func (t *TripAdvisorProcessor) addToPipeLine(item string) {
	t.data = append(t.data, item)
}
