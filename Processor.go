package main

import (
	"encoding/json"
	"fmt"
	"go_spider/core/common/page"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/YonnaR/go-asciibot"
	"gopkg.in/mgo.v2/bson"
)

type TripAdvisorProcessor struct {
	startURI string
	nextURI  string
	output   string
	thread   int
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
func NewTripAdvisorProcessor(url string, output string, thread int) *TripAdvisorProcessor {
	return &TripAdvisorProcessor{
		startURI: url,
		output:   output,
		thread:   thread,
	}
}

/*Process is life cycle spider
Parse html dom here and record the parse result that we want to Page.
Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.*/
func (t *TripAdvisorProcessor) Process(p *page.Page) {

	if p.IsSucc() {
		query := p.GetHtmlParser()
		page := query.Find(".ui_pagination")
		cPage := t.nextURI

		if strings.Contains(t.startURI, "Restaurant") {
			review := query.Find(".review-container")
			review.Each(func(i int, s *goquery.Selection) {
				fmt.Println(s.Text())
				c := restaurantParser(s)
				d, _ := json.Marshal(c)
				p.AddField(bson.NewObjectId().String(), string(d))
			})
			getNextPage(page, t)

		}

		if strings.Contains(t.startURI, "Hotel") {
			review := query.Find(".hotels-hotel-review-community-content-Card__ui_card--3kTH_")
			review.Each(func(i int, s *goquery.Selection) {
				c := hotelParser(s)
				d, _ := json.Marshal(c)
				p.AddField(bson.NewObjectId().String(), string(d))
			})
			getNextPage(page, t)

		}

		if cPage != t.nextURI && t.nextURI != "" {
			p.AddTargetRequest("https://www.tripadvisor.fr"+t.nextURI, "html")
		}

	} else {
		log.Fatal("can't fetch website, check your internet connection")
		t.Finish()
	}

}

/*Finish is last task of spider function */
func (t *TripAdvisorProcessor) Finish() {
	fmt.Printf(asciibot.Random() + "\n\r Spidering ended \n")
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
