package hltv

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func (c *Client) GetPlayer(ID int) (player SimplePlayer, err error) {
	co := c.collector
	name := c.getNameFromID(ID)

	co.OnHTML("body", func(el *colly.HTMLElement) {
		player.ID = ID
		player.Team = el.DOM.Find(".SummaryTeamname").Find("a").Text()
		player.Nickname = el.DOM.Find(".summaryNickname").Text()
		player.Slug = name

		sr := el.DOM.Find(".stats-row")
		mp, _ := strconv.Atoi(sr.Eq(6).Children().Last().Text())
		player.MapsPlayed = mp

		kd, _ := strconv.ParseFloat(sr.Eq(3).Children().Last().Text(), 64)
		player.KD = kd

		r, _ := strconv.ParseFloat(sr.Eq(13).Children().Last().Text(), 64)
		player.Rating = r
	})

	collectorError(co, &err)
	co.Visit(fmt.Sprintf("%vstats/players/%v/%v", c.baseURL, ID, name))
	return
}

func (c *Client) getNameFromID(ID int) (name string) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		el.DOM.Find(".stats-table").Find("tr").Each(func(i int, el *goquery.Selection) {
			if i == 0 {
				return
			}

			href, _ := el.Find("a[href]").Attr("href")
			if idFromURL(href, 3) == ID {
				name = pathFromURL(href, 4)
			}
		})
	})

	co.Visit(fmt.Sprintf("%vstats/players", c.baseURL))
	return
}
