package hltv

import (
	"fmt"
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
	Logo   string
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
	MatchID int
	Time    time.Time
	Event   struct {
		Name string
		Logo string
	}
	Teams []fullTeam
	Maps  []matchMap
}

func (c *Client) GetMatch(matchID int, matchURL string) (match MatchInfo, err error) {
	co := c.collector

	co.OnHTML("body", func(el *colly.HTMLElement) {
		match.MatchID = matchID
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
		el.DOM.Find(".mapholder").Each(func(i int, mEl *goquery.Selection) {
			match.Maps = append(match.Maps, getMatchMap(mEl))
		})
	})

	collectorError(co, &err)

	co.Visit(fmt.Sprintf("%vmatches/%v/%v", c.baseURL, matchID, matchURL))
	return match, err
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
	id, _ := strconv.Atoi(strings.Split(href, "/")[2])
	s.ID = id

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

func getMatchMap(mEl *goquery.Selection) (mm matchMap) {
	mm.Name = mEl.Find(".mapname").Text()
	mm.Pick = mEl.Find(".pick").Find(".results-teamname").Text()

	names := mEl.Find(".results-teamname")
	logos := mEl.Find(".logo")

	scoreHTML := mEl.Find(".results-center-half-score").Children().Text()
	replacer := strings.NewReplacer("(", "", ")", ",", ":", ",", ";", ",", " ", "")
	scores := strings.Split(replacer.Replace(scoreHTML), ",")

	t1Name := names.First().Text()
	t1Logo, _ := logos.First().Attr("src")
	t1 := newMapTeam(t1Name, t1Logo, scores[0], scores[2], scores[4])

	t2Name := names.Last().Text()
	t2Logo, _ := logos.Last().Attr("src")
	t2 := newMapTeam(t2Name, t2Logo, scores[1], scores[3], scores[5])

	mm.Teams = append(mm.Teams, *t1, *t2)
	return
}

func newMapTeam(name, logo, first, second, ext string) *mapTeam {
	f, _ := strconv.Atoi(first)
	s, _ := strconv.Atoi(second)
	e, _ := strconv.Atoi(ext)

	mt := &mapTeam{
		Name: name,
		Logo: logo,
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
