package proxy

import (
	"fmt"
	"time"
)

//ESRI gets images from a WGS84 ESRI REST endpoint
type ESRI struct {
	TMS
}

//NewESRI return an ESRI REST proxy
func NewESRI(url, imageType string, timeout time.Duration) (*ESRI, error) {
	//todo: validate URL, derive imageType
	return &ESRI{TMS{URL: url, ImageType: imageType, Timeout: timeout}}, nil
}

//GetTile returns an image tile
func (p *ESRI) GetTile(x, y, z int) ([]byte, error) {
	return p.FetchAsJPEGBytes(fmt.Sprintf(p.URL, z, x, y))
}
