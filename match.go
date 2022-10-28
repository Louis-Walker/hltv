package hltv

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type stats struct {
	Name     string
	Nickname string
	ID       int
	Kills    int
	Deaths   int
	ADR      float64
	KAST     float64
	Rating   float64
}

type fullTeam struct {
	SimpleTeam
	Players []stats
}

type mapTeam struct {
	Name   string
	Result struct {
		First  int
		Second int
		Ext    int
	}
}

type matchMap struct {
	Name  string
	Pick  string
	Teams []mapTeam
}

type MatchInfo struct {
	ID    int
	Time  time.Time
	Event struct {
		Name string
		Logo string
	}
	Teams []fullTeam
	Maps  []matchMap
}

func (c *Client) GetMatch(ID int) (match MatchInfo, err error) {
	co := c.collector
	matchURL, err := c.getMatchURL(ID)
	if err != nil {
		log.Printf("[HLTV] %v", err)
	}

	co.OnHTML("body", func(el *colly.HTMLElement) {
		match.ID = ID
		match.Event.Name = el.DOM.Find(".event").Children().First().Text()

		// Time
		var unixInt int64
		unixString, _ := el.DOM.Find(".time[data-unix]").Attr("data-unix")
		unixInt, _ = strconv.ParseInt(unixString, 10, 64)
		match.Time = time.UnixMilli(unixInt).UTC()

		// Teams
		teamBoxEls := el.DOM.Find(".standard-box.teamsBox").Children()
		playerBoxEls := el.DOM.Find(".totalstats")
		team1 := getFullTeam(teamBoxEls.First(), playerBoxEls.First())
		team2 := getFullTeam(teamBoxEls.Last(), playerBoxEls.Last())
		match.Teams = append(match.Teams, team1, team2)

		// Maps
		el.DOM.Find(".mapholder").Each(func(i int, el *goquery.Selection) {
			match.Maps = append(match.Maps, getMatchMap(el))
		})
	})

	collectorError(co, &err)
	co.Visit(fmt.Sprintf("%vmatches/%v/%v", c.baseURL, ID, matchURL))
	return
}

func (c *Client) getMatchURL(ID int) (matchURL string, err error) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		el.ForEach(".liveMatch", func(i int, el *colly.HTMLElement) {
			href, _ := el.DOM.Find("a").Attr("href")
			elID := idFromURL(href, 2)
			if ID == elID {
				matchURL = pathFromURL(href, 3)
			}
		})

		el.ForEach(".upcomingMatch", func(i int, el *colly.HTMLElement) {
			href, _ := el.DOM.Find("a").Attr("href")
			elID := idFromURL(href, 2)
			if ID == elID {
				matchURL = pathFromURL(href, 3)
			}
		})
	})

	if len(matchURL) == 0 {
		err = fmt.Errorf("no upcoming match was found with ID: %v", ID)
	}

	co.Visit(fmt.Sprintf("%vmatches", c.baseURL))
	return
}

func getFullTeam(teamEl *goquery.Selection, playerEl *goquery.Selection) (team fullTeam) {
	team.Name = teamEl.Find(".teamName").Text()
	team.Logo, _ = teamEl.Find(".logo").Attr("src")
	scoreStr := teamEl.Find("div[class]").Last().Text()
	team.Score, _ = strconv.Atoi(scoreStr)

	playerEl.Find(".players").Parent().Each(func(i int, pEl *goquery.Selection) {
		if i == 0 {
			return
		}
		team.Players = append(team.Players, getStats(pEl))
	})
	return
}

func getStats(pEl *goquery.Selection) (s stats) {
	n := strings.Split(pEl.Find(".statsPlayerName").First().Text(), "'")
	s.Name = strings.Replace(n[0], " ", "", -1) + n[2]
	s.Nickname = pEl.Find(".player-nick").Text()

	href, _ := pEl.Find("a").First().Attr("href")
	s.ID = idFromURL(href, 2)

	kd := strings.Split(pEl.Find(".kd").Text(), "-")
	k, _ := strconv.Atoi(kd[0])
	s.Kills = k
	d, _ := strconv.Atoi(kd[1])
	s.Deaths = d
	a, _ := strconv.ParseFloat(pEl.Find(".adr").Text(), 64)
	s.ADR = a
	kastTrim := strings.Replace(pEl.Find(".kast").Text(), "%", "", -1)
	kast, _ := strconv.ParseFloat(kastTrim, 64)
	s.KAST = kast
	r, _ := strconv.ParseFloat(pEl.Find(".rating").Text(), 64)
	s.Rating = r
	return
}

func getMatchMap(el *goquery.Selection) (mm matchMap) {
	mm.Name = el.Find(".mapname").Text()
	mm.Pick = el.Find(".pick").Find(".results-teamname").Text()

	names := el.Find(".results-teamname")
	logos := el.Find(".logo")

	var t1, t2 *mapTeam
	t1Name := names.First().Text()
	t1Logo, _ := logos.First().Attr("src")

	t2Name := names.Last().Text()
	t2Logo, _ := logos.Last().Attr("src")

	scores := el.Find(".results-center-half-score").Children()
	t1 = newMapTeam(t1Name, t1Logo, scores.Eq(1).Text(), scores.Eq(5).Text(), scores.Eq(9).Text())
	t2 = newMapTeam(t2Name, t2Logo, scores.Eq(3).Text(), scores.Eq(7).Text(), scores.Eq(11).Text())

	mm.Teams = append(mm.Teams, *t1, *t2)
	return
}

func newMapTeam(name, logo string, scores ...string) *mapTeam {
	f, _ := strconv.Atoi(scores[0])
	s, _ := strconv.Atoi(scores[1])
	e, _ := strconv.Atoi(scores[2])

	mt := &mapTeam{
		Name: name,
		Result: struct {
			First  int
			Second int
			Ext    int
		}{
			First:  f,
			Second: s,
			Ext:    e,
		},
	}
	return mt
}
