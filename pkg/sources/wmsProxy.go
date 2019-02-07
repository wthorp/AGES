package sources

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//WMSProxy gets images from a WGS84 WMS
type WMSProxy struct {
	URL     url.URL
	Timeout time.Duration
}

//NewWMSProxy return a WMS proxy
func NewWMSProxy(baseURL string, timeout time.Duration) (*WMSProxy, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &WMSProxy{URL: base, Timeout: timeout}, nil
}

func (w *WMSProxy) GetTile(x, y, z int) ([]byte, error) {
	bbox := TileXYToBBox(x, y, z)
	q := w.URL.Query()
	q.Set("BBOX", fmt.Sprintf("%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top))
	url.RawQuery = q.Encode()
	var client = &http.Client{Timeout: w.Timeout}
	resp, err := client.Get(url.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
