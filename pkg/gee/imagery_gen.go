package gee

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"AGES/pkg/core"
	"AGES/pkg/gee/keyhole"
)

//ImageryProvider is things with GetTile
type ImageryProvider interface {
	GetTile(int, int, int) ([]byte, error)
}

//ImageryGen returns imagery
type ImageryGen struct {
	Provider ImageryProvider
}

//ImageryHandler returns an image
func (p *ImageryGen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
	quadkey := parts[1]

	jpgType := keyhole.EarthImageryPacket_JPEG
	imageBytes, err := p.Provider.GetTile(QuadKeyToTileXY(quadkey))
	if err != nil {
		fmt.Printf("bad image source\n%v\n", err)
		//fmt.Fprintf(w, "bad image source\n%v", err)
		//w.WriteHeader(http.StatusInternalServerError)
		//return
		imageBytes, err = createTextImage(err.Error())
		if err != nil {
			fmt.Printf("createTextImage err\n%v\n", err)
			fmt.Fprintf(w, "createTextImage err\n%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	imageBytes, err = demo(imageBytes)
	if err != nil {
		fmt.Printf("bad demo source\n%v\n", err)
	}

	eip := &keyhole.EarthImageryPacket{ImageType: &jpgType, ImageData: imageBytes}
	eipBytes, err := proto.Marshal(eip) //convert to protobuf
	if err != nil {
		fmt.Printf("eip proto\n%+v\n%v", eip, err)
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

func demo(imgBytes []byte) ([]byte, error) {
	orig, err := jpeg.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(), orig, image.ZP, draw.Src)
	col := color.RGBA{200, 100, 0, 255}
	x, y := 8, 8
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{Dst: img, Src: image.NewUniform(col), Face: basicfont.Face7x13, Dot: point}
	d.DrawString("FOR DEMO ONLY")
	return core.JPEGBytes(img)
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
