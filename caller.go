package sgql

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Doer         http.Client
	EndPoint     string
	Password     string
	RestoreCount int
}

const (
	PAUSELONG = 1000
)

func (c *Client) GraphQuery(query string) (resp *http.Response, errResp error) {
	qlData := strings.NewReader(query)
	req, errReq := http.NewRequest("POST", c.EndPoint, qlData)
	if errReq != nil {
		return nil, errReq
	}
	req.Header.Add("Content-Type", "application/graphql")
	req.Header.Add("X-Shopify-Access-Token", c.Password)
	// try 10 times get normal response - hi CloudFlare
	count := c.RestoreCount
	for {
		count--
		if count == 0 {
			return
		}
		resp, errResp = c.Doer.Do(req)
		if resp == nil {
			errResp = fmt.Errorf("uncaught error")
			continue
		}
		if resp.StatusCode == 200 {
			break // all ok - we can return
		}
		// response status != 200
		if errResp == nil {
			errResp = fmt.Errorf("response status: %s", resp.Status)
		}
		time.Sleep(time.Duration(c.RestoreCount-count) * PAUSELONG * time.Millisecond)
	}
	return
}
