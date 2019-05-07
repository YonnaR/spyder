package main

import (
	"github.com/PuerkitoBio/goquery"
)

func restaurantParser(s *goquery.Selection) CModel {

	/* Parse website informations */
	created, _ := s.Find(".ratingDate").Attr("title")
	auth := s.Find(".info_text").Text()
	com := s.Find(".partial_entry").Text()
	title := s.Find(".noQuotes").Text()
	image := s.Find(".centeredImg")

	nCom := CModel{
		Author:     auth,
		Created:    created,
		Title:      title,
		Commentary: com,
	}

	/* Get commentaries images urls */
	image.Each(func(i int, s *goquery.Selection) {
		imgURI, _ := s.Attr("data-lazyurl")
		nCom.Images = append(nCom.Images, imgURI)
	})

	return nCom
}

func hotelParser(s *goquery.Selection) CHModel {

	/* Parse website informations */
	d := s.Find(".social-member-event-MemberEventOnObjectBlock__event_type--3njyv span")
	auth := d.Find("a").Remove().Text()
	created := s.Find(".social-member-event-MemberEventOnObjectBlock__event_type--3njyv span a").Text()
	com := s.Find("q").Text()
	title := s.Find(".hotels-review-list-parts-ReviewTitle__reviewTitle--2Fauz").Text()
	image := s.Find("img")
	book := s.Find(".hotels-review-list-parts-EventDate__event_date--CRXs4").Text()

	//fmt.Println("bookDate :", book)
	nCom := CHModel{
		Author:     auth,
		Created:    created,
		Title:      title,
		Commentary: com,
		Book:       book,
	}

	/* Get commentaries images urls */
	image.Each(func(i int, s *goquery.Selection) {
		imgURI, _ := s.Attr("data-lazyurl")
		nCom.Images = append(nCom.Images, imgURI)
	})

	return nCom
}
