package sources

import (
	"fmt"
	"image/png"
	"net/http"
	"net/url"
	"time"

	"AGES/pkg/core"
)

//OSMProxy is an Open Street Map proxy
type OSMProxy struct {
	URL     *url.URL
	Timeout time.Duration
}

//NewOSMProxy return a Open Street Map proxy
func NewOSMProxy(baseURL string, timeout time.Duration) (*OSMProxy, error) {
	base, err := url.Parse(baseURL) // https://a.tile.openstreetmap.org
	if err != nil {
		return nil, err
	}
	return &OSMProxy{URL: base, Timeout: timeout}, nil
}

//GetTile returns OpenStreetMap imagery
func (p *OSMProxy) GetTile(x, y, z int) ([]byte, error) {
	url := fmt.Sprintf("%s/%d/%d/%d.png", p.URL, z, x, y-1)
	var client http.Client
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad osm response status")
	}
	//osmBytes, _ := ioutil.ReadAll(resp.Body)
	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return core.JPEGBytes(img)
}
