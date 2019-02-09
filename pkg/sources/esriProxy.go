package sources

import (
	"fmt"
	"net/url"
	"time"
)

//ESRIProxy gets images from a WGS84 ESRI REST endpoint
type ESRIProxy struct {
	URL     *url.URL
	Timeout time.Duration
}

//NewESRIProxy return an ESRI REST proxy
func NewESRIProxy(baseURL string, timeout time.Duration) (*ESRIProxy, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &ESRIProxy{URL: base, Timeout: timeout}, nil
}

//GetTile returns an image tile
func (p *ESRIProxy) GetTile(x, y, z int) ([]byte, error) {
	//todo
	return nil, fmt.Errorf("not implemented")
}
