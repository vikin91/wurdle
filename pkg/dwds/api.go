package dwds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURLV1 = "https://www.dwds.de/api"
)

type FreqReply struct {
	Hits      int
	Frequency int
	Total     string
	Q         string
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseURLV1,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) Frequency(ctx context.Context, word string) (int, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/frequency/?q=%s", c.BaseURL, word), nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)

	res := &FreqReply{}
	if err := c.sendRequest(req, res); err != nil {
		return 0, err
	}
	return res.Frequency, nil
}

func (c *Client) sendRequest(req *http.Request, obj *FreqReply) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&obj); err != nil {
		return err
	}

	return nil
}
