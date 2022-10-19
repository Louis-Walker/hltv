package hltv

import (
	"context"
	"time"
)

type SimpleNews struct {
	Title       string
	Description string
	Link        string
	Time        time.Time
}

func (c *Client) GetNews(ctx context.Context) (news *[]SimpleNews, err error) {
	return news, err
}
