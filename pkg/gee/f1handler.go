package gee

import (
	"AGES/pkg/gee/keyhole"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
)

//f1Handler returns a dbRoot object
func f1Handler(w http.ResponseWriter, r *http.Request, quadkey string, version string) {
	rawPath := filepath.Join("config", r.URL.RawQuery)
	jsonPath := filepath.Join("config", r.URL.RawQuery+".json")

	//url := path.Join(proxiedURL, "flatfile?"+r.URL.RawQuery)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		//load raw
		file, e := ioutil.ReadFile(rawPath)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("File error: %v\n", e)
			return
		}
		//decode raw
		XOR(file, []byte(defaultKey), true)
		eip := keyhole.EarthImageryPacket{}
		unProto(file, &eip)
		//write image
		imgPath := filepath.Join("config", r.URL.RawQuery+"."+eip.ImageType.String())
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
	unMarshalJSONFile(jsonPath, eip)
	//embed eip image payload in
	imgPath := filepath.Join("config", r.URL.RawQuery+"."+eip.ImageType.String())
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
