package hltv

import (
	"time"
)

type Stats struct {
	Name     string
	Nickname string
	ID       int
	Kills    int
	Deaths   int
	ADR      int
	KAST     int
	Rating   int
}

type Team struct {
	Name    string
	Logo    string
	Result  int
	Players []Stats
}

type MapTeam struct {
	Name   string
	Logo   string
	Result struct {
		First  string
		Second string
		Ext    string
	}
}

type Map struct {
	Name  string
	Pick  string
	Teams []MapTeam
}

type Match struct {
	ID    int
	Time  time.Time
	Event struct {
		name string
		logo string
	}
	Teams []Team
	Maps  []Map
}

func (c *Client) GetMatches() (matches []Match, err error) {
	return matches, err
}
