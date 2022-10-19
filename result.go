package hltv

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type team struct {
	name   string
	logo   string
	result int
}

type Result struct {
	MatchID int
	Maps    string
	Time    time.Time
	Teams   []team
	Event   struct {
		name string
		logo string
	}
}

func (c *Client) GetResults() (results []Result, err error) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		el.ForEach(".result-con", func(i int, el *colly.HTMLElement) {
			var sr Result

			sr.Maps = el.DOM.Find(".map-text").Text()
			sr.Event.name = el.DOM.Find(".event-name").Text()
			sr.Event.logo, _ = el.DOM.Find(".event-logo").Attr("src")

			// Match ID
			var href string
			href, _ = el.DOM.Find("a[href]").Attr("href")
			id := strings.Split(href, "/")[2]
			sr.MatchID, _ = strconv.Atoi(id)

			// Time
			var unixInt int64
			unixString := el.Attr("data-zonedgrouping-entry-unix")
			unixInt, _ = strconv.ParseInt(unixString, 10, 64)
			sr.Time = time.UnixMilli(unixInt).UTC()

			// Teams
			var team1 team
			team1El := el.DOM.Find(".team1")
			team1.name = team1El.Find("div").Text()
			team1.logo, _ = team1El.Find(".team-logo").Attr("src")

			var team2 team
			team2El := el.DOM.Find(".team2")
			team2.name = team2El.Find("div").Text()
			team2.logo, _ = team2El.Find(".team-logo").Attr("src")

			scoreEls := el.DOM.Find(".result-score").Children()
			t1Score := scoreEls.First().Text()
			team1.result, _ = strconv.Atoi(t1Score)
			t2Score := scoreEls.Last().Text()
			team1.result, _ = strconv.Atoi(t1Score)
			team2.result, _ = strconv.Atoi(t2Score)

			sr.Teams = append(sr.Teams, team1, team2)

			results = append(results, sr)
		})
	})

	co.OnError(func(cr *colly.Response, ce error) {
		collectorError(cr, ce, &err)
	})

	co.Visit(c.baseURL + "results")
	return results, err
}
