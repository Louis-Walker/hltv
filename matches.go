package hltv

import (
	"time"
)

type Match struct {
	ID    int
	Time  time.Time
	Maps  string
	Teams []simpleTeam
	Event struct {
		name string
		logo string
	}
}

func (c *Client) GetMatches() (matches []Match, err error) {
	return matches, err
}
