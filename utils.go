package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

/*CModel is a commentary model structure */
type CModel struct {
	Author     string   `json:"author" bson:"author"`
	Title      string   `json:"title" bson:"title"`
	Commentary string   `json:"commentary" bson:"commentay"`
	Created    string   `json:"created" bson:"created"`
	Images     []string `json:"images" bson:"images"`
}

/*CHModel is a commentary model structure for hotel */
type CHModel struct {
	Author     string   `json:"author" bson:"author"`
	Title      string   `json:"title" bson:"title"`
	Commentary string   `json:"commentary" bson:"commentay"`
	Created    string   `json:"created" bson:"created"`
	Book       string   `json:"book" bson:"book"`
	Images     []string `json:"images" bson:"images"`
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
