package gee

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"AGES/pkg/core"
	"AGES/pkg/gee/keyhole"
)

//f1Handler returns an image
func f1Handler(w http.ResponseWriter, r *http.Request, quadkey string, imgSource func(int, int, int) ([]byte, error)) {
	jpgType := keyhole.EarthImageryPacket_JPEG
	imageBytes, err := imgSource(QuadKeyToTileXY(quadkey))
	if err != nil {
		fmt.Fprintf(w, "bad image source\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

//f1RawHandler returns an image
func f1RawHandler(w http.ResponseWriter, quadkey string, imgSource func(int, int, int) ([]byte, error)) {
	imgBytes, err := imgSource(QuadKeyToTileXY(quadkey))
	if err != nil {
		fmt.Fprintf(w, "bad image source\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(imgBytes)
}

//createTextImage is intended for use with error messges
func createTextImage(label string) ([]byte, error) {
	//maybe see https://github.com/golang/freetype/blob/master/example/drawer/main.go
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.ZP, draw.Src)
	col := color.RGBA{200, 100, 0, 255}
	x, y := 8, 8
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{Dst: img, Src: image.NewUniform(col), Face: basicfont.Face7x13, Dot: point}
	d.DrawString(label)
	return core.JPEGBytes(img)
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
