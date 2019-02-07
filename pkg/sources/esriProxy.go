package sources

import (
	"fmt"
	"net/url"
	"time"
)

//ESRIProxy gets images from a WGS84 ESRI REST endpoint
type ESRIProxy struct {
	URL     url.URL
	Timeout time.Duration
}

//NewESRIProxy return an ESRI REST proxy
func NewESRIProxy(url string, timeout time.Duration) (EsriProxy, error) {
	url, err := url.Parse(BaseURL)
	if err != nil {
		return err
	}
	return EsriProxy{URL: url, Timeout: timeout}
}

func (p *EsriProxy) GetTile(x, y, z int) ([]byte, error) {
	//todo
	return nil, fmt.Errorf("not implemented")
}
