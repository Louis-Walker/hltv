package hltv

import (
	"time"
)

type stats struct {
	Name     string
	Nickname string
	ID       int
	Kills    int
	Deaths   int
	ADR      int
	KAST     int
	Rating   int
}

type fullTeam struct {
	SimpleTeam
	Players []stats
}

type mapTeam struct {
	Name   string
	Logo   string
	Result struct {
		First  string
		Second string
		Ext    string
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
		name string
		logo string
	}
	Teams []fullTeam
	Maps  []matchMap
}

func (c *Client) GetMatch() (matches []MatchInfo, err error) {
	return matches, err
}
