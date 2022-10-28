package hltv

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type SimplePlayer struct {
	ID         int
	Team       string
	Nickname   string
	Slug       string
	MapsPlayed int
	KD         float64
	Rating     float64
}

func (c *Client) GetPlayers() (players []SimplePlayer, err error) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		el.DOM.Find(".stats-table").Find("tr").Each(func(i int, el *goquery.Selection) {
			if i == 0 {
				return
			}

			var sp SimplePlayer

			href, _ := el.Find("a").First().Attr("href")
			sp.ID = idFromURL(href, 3)

			sp.Team, _ = el.Find(".teamCol").Find("img").Attr("title")
			sp.Nickname = el.Find(".playerCol").Find("a").Text()
			sp.Slug = pathFromURL(href, 4)

			mp, _ := strconv.Atoi(el.Find(".statsDetail").Eq(0).Text())
			sp.MapsPlayed = mp

			kd, _ := strconv.ParseFloat(el.Find(".statsDetail").Eq(2).Text(), 64)
			sp.KD = kd

			rating, _ := strconv.ParseFloat(el.Find(".ratingCol").Text(), 64)
			sp.Rating = rating

			players = append(players, sp)
		})
	})

	collectorError(co, &err)
	co.Visit(fmt.Sprintf("%vstats/players", c.baseURL))
	return
}
