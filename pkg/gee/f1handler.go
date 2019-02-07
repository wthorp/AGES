package gee

import (
	"AGES/pkg/gee/keyhole"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
)

//f1Handler returns an image
func f1Handler(w http.ResponseWriter, r *http.Request, quadkey string, imgSource func(int, int, int) ([]byte, error)) {
	jpgType := keyhole.EarthImageryPacket_JPEG
	imageBytes, err := imgSource(QuadKeyToTileXY(quadkey))
	eip := &keyhole.EarthImageryPacket{ImageType: &jpgType, ImageData: imageBytes}
	eipBytes, err := proto.Marshal(eip) //convert to protobuf
	if err != nil {
		fmt.Fprintf(w, "eip proto\n%+v\n%v", eip, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	XOR(eipBytes, []byte(defaultKey), false) //encode raw
	w.Write(eipBytes)                        //send bytes
}

//oldF1Handler returns a dbRoot object
func oldF1Handler(w http.ResponseWriter, r *http.Request, quadkey string) {
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

func WriteJpegRaw(w http.ResponseWriter, img image.Image) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		fmt.Println("bad jpg")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func WriteJpeg(w http.ResponseWriter, img image.Image) {
	jpgType := keyhole.EarthImageryPacket_JPEG
	o := &jpeg.Options{Quality: 85}
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, o)
	if err != nil {
		fmt.Println("bad jpg")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	imgBytes := buf.Bytes()
	eip := &keyhole.EarthImageryPacket{ImageType: &jpgType, ImageData: imgBytes}
	//convert to protobuf
	eipBytes, err := proto.Marshal(eip)
	if err != nil {
		fmt.Printf("bad jpg proto\n%+v\n%v", eip, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//encode raw
	XOR(eipBytes, []byte(defaultKey), false)
	//send bytes
	w.Write(eipBytes)
}
