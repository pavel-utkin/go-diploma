package accrual

import (
	"net/http"
	"time"
)

type Client struct {
	http.Client

	Addr string
}

func NewClient(addr string) *Client {
	client := http.Client{}
	client.Timeout = 1 * time.Second

	return &Client{
		Client: client,
		Addr:   addr,
	}
}
