package gee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/protobuf/proto"

	"AGES/pkg/core"
	"AGES/pkg/gee/keyhole"
	"AGES/pkg/net"
)

//ImageryProxy proxies imagery
type ImageryProxy struct {
	URL *url.URL
}

//ServeHTTP returns a imagery
func (p *ImageryProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawPath := core.ApplicationDir(r.URL.RawQuery + ".raw")
	jsonPath := core.ApplicationDir(r.URL.RawQuery + ".json")

	if _, err := os.Stat(rawPath); os.IsNotExist(err) {
		err = net.DownloadFile(rawPath, net.RemapURL(p.URL, r.URL))
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	// var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
	// quadkey := parts[1]

	//url := path.Join(proxiedURL, "flatfile?"+r.URL.RawQuery)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		//load raw
		file, e := ioutil.ReadFile(rawPath)
		if e != nil {
			fmt.Printf("File error: %v\n", e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//decode raw
		XOR(file, []byte(defaultKey), true)
		eip := keyhole.EarthImageryPacket{}
		unProto(file, &eip)
		//write image
		imgPath := core.ApplicationDir(r.URL.RawQuery + "." + eip.ImageType.String())
		writeFile(imgPath, eip.ImageData)
		//write JSON
		eip.ImageData = eip.ImageData[0:0]
		b, err := json.MarshalIndent(eip, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile(jsonPath, b)
	}

	//get EarthImageryPacket json data
	eip := &keyhole.EarthImageryPacket{}
	err := unMarshalJSONFile(jsonPath, eip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", err)
		return
	}
	//embed eip image payload in
	imgPath := core.ApplicationDir(r.URL.RawQuery + "." + eip.ImageType.String())
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
