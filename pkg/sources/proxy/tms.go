package proxy

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"time"

	"AGES/pkg/core"
)

//TMS is a generic tiled map service proxy
type TMS struct {
	URL       string
	ImageType string
	Timeout   time.Duration
}

//NewTMS return a generic tiled map service proxy
func NewTMS(url, imageType string, timeout time.Duration) (*TMS, error) {
	//todo: validate URL, derive imageType
	return &TMS{URL: url, ImageType: imageType, Timeout: timeout}, nil
}

//GetTile returns OpenStreetMap imagery
func (p *TMS) GetTile(x, y, z int) ([]byte, error) {
	return p.FetchAsJPEGBytes(fmt.Sprintf(p.URL, z, x, y))
}

//FetchAsJPEGBytes fetches a remote image as JPEG bytes
func (p *TMS) FetchAsJPEGBytes(url string) ([]byte, error) {
	var client = &http.Client{Timeout: p.Timeout}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad proxy response status")
	}
	//always return JPEG for now
	if p.ImageType == "JPEG" {
		return ioutil.ReadAll(resp.Body)
	}
	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return core.PNGBytes(img)
}
