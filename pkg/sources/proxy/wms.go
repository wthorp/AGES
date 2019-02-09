package proxy

import (
	"fmt"
	"time"

	"AGES/pkg/core"
)

//WMS gets images from a WGS84 WMS
type WMS struct {
	TMS
}

//NewWMS return a WMS proxy
func NewWMS(url, imageType string, timeout time.Duration) (*WMS, error) {
	//todo: validate URL, derive imageType
	return &WMS{TMS{URL: url, ImageType: imageType, Timeout: timeout}}, nil
}

//GetTile returns WMS imagery tiles
func (p *WMS) GetTile(x, y, z int) ([]byte, error) {
	bbox := core.TileXYToBBox(x, y, z)
	//q := p.URL.Query()
	//q.Set("BBOX", )
	//p.URL.RawQuery = q.Encode()
	return p.FetchAsJPEGBytes(p.URL + "?" + fmt.Sprintf("BBOX=%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top))
}
