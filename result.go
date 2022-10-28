package hltv

import (
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type SimpleTeam struct {
	Name  string
	Logo  string
	Score int
}

type SimpleMatch struct {
	MatchID int
	Maps    string
	Time    time.Time
	Teams   []SimpleTeam
	Event   struct {
		Name string
		Logo string
	}
}

func (c *Client) GetResults() (results []SimpleMatch, err error) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		el.ForEach(".result-con", func(i int, el *colly.HTMLElement) {
			var r SimpleMatch
			href, _ := el.DOM.Find("a[href]").Attr("href")

			r.MatchID = idFromURL(href, 2)
			r.Maps = el.DOM.Find(".map-text").Text()
			r.Time = getMatchTime(el)
			r.Teams = getMatchTeams(el, ".team", ".team-logo")
			r.Event.Name, r.Event.Logo = getMatchEvent(el, ".event-name", ".event-logo")

			results = append(results, r)
		})
	})

	collectorError(co, &err)
	co.Visit(c.baseURL + "results")
	return results, err
}

func getMatchTime(el *colly.HTMLElement) (matchTime time.Time) {
	var unixInt int64
	unixString := el.Attr("data-zonedgrouping-entry-unix")
	unixInt, _ = strconv.ParseInt(unixString, 10, 64)
	matchTime = time.UnixMilli(unixInt).UTC()
	return
}

func getMatchEvent(el *colly.HTMLElement, nameNode string, logoNode string) (name string, logo string) {
	name = el.DOM.Find(nameNode).Text()
	logo, _ = el.DOM.Find(logoNode).Attr("src")
	return
}

func getMatchTeams(el *colly.HTMLElement, nameNode string, logoNode string) (teams []SimpleTeam) {
	var team1 SimpleTeam
	team1El := el.DOM.Find(".team1")
	team1.Name = team1El.Find(nameNode).Text()
	team1.Logo, _ = team1El.Find(logoNode).Attr("src")

	var team2 SimpleTeam
	team2El := el.DOM.Find(".team2")
	team2.Name = team2El.Find(nameNode).Text()
	team2.Logo, _ = team2El.Find(logoNode).Attr("src")

	scoreEls := el.DOM.Find(".result-score").Children()
	t1Score := scoreEls.First().Text()
	team1.Score, _ = strconv.Atoi(t1Score)
	t2Score := scoreEls.Last().Text()
	team2.Score, _ = strconv.Atoi(t2Score)

	teams = append(teams, team1, team2)
	return
}
