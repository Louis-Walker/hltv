package hltv

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type FullMatch struct {
	SimpleMatch
	Stars int
	Live  bool
}

func (c *Client) GetMatches() (matches []FullMatch, err error) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		el.ForEach(".liveMatch", func(i int, el *colly.HTMLElement) {
			matches = append(matches, getFullMatch(el))
		})
		el.ForEach(".upcomingMatch", func(i int, el *colly.HTMLElement) {
			matches = append(matches, getFullMatch(el))
		})
	})

	collectorError(co, &err)
	co.Visit(c.baseURL + "matches")
	return matches, err
}

func getFullMatch(el *colly.HTMLElement) (m FullMatch) {
	href, _ := el.DOM.Find("a[href]").Attr("href")

	m.MatchID = idFromURL(href, 2)
	m.Maps = el.DOM.Find(".matchMeta").Text()
	m.Time = getMatchTime(el)
	m.Teams = getMatchTeams(el, ".matchTeamName", ".matchTeamLogo")
	m.Event.Name, m.Event.Logo = getMatchEvent(el, ".matchEventName", ".matchEventLogo")

	// Star Rating
	ratingEls := el.DOM.Find(".matchRating")
	m.Stars = 5
	ratingEls.Children().Each(func(i int, s *goquery.Selection) {
		if s.HasClass("faded") {
			m.Stars--
		}
	})

	// Live Indicator
	if el.DOM.HasClass("liveMatch") {
		m.Live = true
	}

	return
}
