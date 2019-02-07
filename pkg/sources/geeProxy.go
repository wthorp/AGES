package sources

import (
	"AGES/pkg/gee"
	"AGES/pkg/gee/keyhole"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
)

type GEEProxy struct {
	URL     *url.URL
	Timeout time.Duration
}

//NewGEEProxy return an GEE proxy
func NewGEEProxy(baseURL string, timeout time.Duration) (*GEEProxy, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &GEEProxy{URL: base, Timeout: timeout}, nil
}

//f1Handler returns a dbRoot object
func (g *GEEProxy) GetTile(x, y, z int) ([]byte, error) {
	quadkey := gee.TileXYToQuadKey(x, y, z)
	rawPath := filepath.Join("config", g.URL.RawQuery)
	jsonPath := filepath.Join("config", g.URL.RawQuery+".json")

	//url := path.Join(proxiedURL, "flatfile?"+r.URL.RawQuery)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		//load raw
		file, e := ioutil.ReadFile(rawPath)
		if e != nil {
			return nil, e
		}
		//decode raw
		gee.XOR(file, []byte(gee.defaultKey), true)
		eip := keyhole.EarthImageryPacket{}
		gee.unProto(file, &eip)
		//write image
		imgPath := filepath.Join("config", g.URL.RawQuery+"."+eip.ImageType.String())
		writeFile(imgPath, eip.ImageData)
		//write JSON
		eip.ImageData = eip.ImageData[0:0]
		b, err := json.Marshal(eip)
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile(jsonPath, b)
	}

	//get EarthImageryPacket json data
	eip := &keyhole.EarthImageryPacket{}
	gee.unMarshalJSONFile(jsonPath, eip)
	//embed eip image payload in
	imgPath := filepath.Join("config", g.URL.RawQuery+"."+eip.ImageType.String())
	imgBytes, e := ioutil.ReadFile(imgPath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	eip.ImageData = imgBytes
	//convert to protobuf
	eipBytes, err := proto.Marshal(eip)
	if err != nil {
		fmt.Fprintf(w, "eip proto\n%+v\n%v", eip, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//encode raw
	XOR(eipBytes, []byte(defaultKey), false)
	//send bytes
	w.Write(eipBytes)
}
