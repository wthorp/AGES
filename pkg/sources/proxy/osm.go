package proxy

import (
	"fmt"
	"time"
)

//OSM is an Open Street Map proxy
type OSM struct {
	TMS
}

//NewOSM return a Open Street Map proxy
func NewOSM(url string, timeout time.Duration) (*OSM, error) {
	//todo: validate URL, derive imageType? // https://a.tile.openstreetmap.org
	return &OSM{TMS{URL: url, ImageType: "PNG", Timeout: timeout}}, nil
}

//GetTile returns OpenStreetMap imagery
func (p *OSM) GetTile(x, y, z int) ([]byte, error) {
	url := fmt.Sprintf("%s/%d/%d/%d.png", p.URL, z, x, y-1)
	return p.FetchAsJPEGBytes(url)
}
