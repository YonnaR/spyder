package main

import (
	"encoding/json"
	"fmt"
	"go_spider/core/common/page"
	"io/ioutil"
	"log"
	"net/http"
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
==================================================================
*/
type TripAdvisorProcessor struct {
	startURI string
	nextURI  string
	output   string
	data     []string
}

/*CModel is a commentary model structure */
type CModel struct {
	Author     string   `json:"author" bson:"author"`
	Title      string   `json:"title" bson:"title"`
	Commentary string   `json:"commentary" bson:"commentay"`
	Created    string   `"json:"created" bson:"created"`
	Images     []string `json:"images" bson:"images"`
}

/*CHModel is a commentary model structure for hotel */
type CHModel struct {
	Author     string   `json:"author" bson:"author"`
	Title      string   `json:"title" bson:"title"`
	Commentary string   `json:"commentary" bson:"commentay"`
	Created    string   `"json:"created" bson:"created"`
	Book       string   `"json:"book" bson:"book"`
	Images     []string `json:"images" bson:"images"`
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
==================================================================
*/
func (t *TripAdvisorProcessor) Process(p *page.Page) {

	if p.IsSucc() {
		query := p.GetHtmlParser()
		page := query.Find(".ui_pagination")
		cPage := t.nextURI

		getNextPage(page, t)
		if cPage != t.nextURI && t.nextURI != "" {
			p.AddTargetRequest("https://www.tripadvisor.fr"+t.nextURI, "html")
		} else {
			result := p.GetPageItems().GetAll()
			for k, i := range result {
				fmt.Println(k, i)

			}
		}

		if strings.Contains(t.startURI, "Restaurant") {
			review := query.Find(".review-container")
			review.Each(func(i int, s *goquery.Selection) {
				c := restaurantParser(s)
				d, _ := json.MarshalIndent(c, "", "")
				t.AddToPipeLine(string(d))
				//p.AddField(bson.NewObjectId().String(), string(d))
			})

		}

		if strings.Contains(t.startURI, "Hotel") {
			getNextPage(page, t)
			review := query.Find(".hotels-hotel-review-community-content-Card__ui_card--3kTH_")
			review.Each(func(i int, s *goquery.Selection) {
				c := hotelParser(s)
				d, _ := json.MarshalIndent(c, "", "")
				//p.AddField(bson.NewObjectId().String(), string(d))
				t.AddToPipeLine(string(d))
				//t.data = append(t.data, string(d))

			})

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
	send data?
	recup and show log?

==================================================================
*/
func (t *TripAdvisorProcessor) Finish() {
	//fmt.Println(t.data[1])
	file, err := json.MarshalIndent(t.data, "", " ")
	if err != nil {
		fmt.Println("err : ", err.Error())
	}
	fmt.Println(len(t.data))
	werr := ioutil.WriteFile(t.output, file, 0644)
	if werr != nil {
		fmt.Println("err : ", err.Error())
	}

	fmt.Printf(asciibot.Random() + "\n\r Spidering ended \n")
}

func (t *TripAdvisorProcessor) AddToPipeLine(item string) {
	//fmt.Println("current item : ", item)
	t.data = append(t.data, item)
	fmt.Println("Item created")
}

func checkConnection() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func getNextPage(p *goquery.Selection, t *TripAdvisorProcessor) {

	p.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Text() == "Suivant" || s.Text() == "Next" {
			u, _ := s.Attr("href")
			if u != "" {
				t.nextURI = u
				return false

			}
			if u != t.startURI {
				t.nextURI = u
				return false
			}
		}
		return true
	})
}
