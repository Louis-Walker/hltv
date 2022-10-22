package hltv

import (
	"errors"
	"fmt"

	"github.com/gocolly/colly"
)

const Version = "0.0.1"

type Client struct {
	baseURL   string
	collector *colly.Collector
}

func New() *Client {
	co := colly.NewCollector(
		colly.AllowedDomains("hltv.org", "www.hltv.org"),
	)
	c := &Client{
		baseURL:   "https://www.hltv.org/",
		collector: co,
	}
	return c
}

func collectorError(co *colly.Collector, err *error) {
	co.OnError(func(cr *colly.Response, ce error) {
		if cr.StatusCode != 0 {
			*err = fmt.Errorf("%v : %v", cr.StatusCode, cr.Body)
			return
		}
		*err = errors.New(ce.Error())
	})
}
