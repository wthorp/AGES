package gee

import (
	"AGES/pkg/gee/keyhole"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"sync"

	"github.com/golang/protobuf/proto"
)

var pipeBytes []byte

func init() {
	var e error
	pipeBytes, e = ioutil.ReadFile("pipe.jpg")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}
}

//PipeHandler returns Magritte's pipe image
func PipeHandler(w http.ResponseWriter, r *http.Request, quadkey string, version string) {
	//get EncryptedDbRoot json data
	jpgType := keyhole.EarthImageryPacket_JPEG
	eip := &keyhole.EarthImageryPacket{ImageType: &jpgType, ImageData: pipeBytes}
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

//OSMHandler returns OpenStreetMap imagery
func OSMHandler(w http.ResponseWriter, r *http.Request, quadkey string, version string) {
	//get EncryptedDbRoot json data

	x, y, z := quadKeyToTileXY(quadkey)
	fmt.Printf("%s => %d/%d/%d\n", quadkey, x, y, z)

	url := fmt.Sprintf("https://a.tile.openstreetmap.org/%d/%d/%d.png", z, x, y-1)

	var client http.Client
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("bad osm get")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad osm response status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//osmBytes, _ := ioutil.ReadAll(resp.Body)
	img, err := png.Decode(resp.Body)
	if err != nil {
		fmt.Printf("bad osm png\n%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJpeg(w, img)
}

const tileSize = 256

//FractalHandler returns fractal imagery
func FractalHandler(w http.ResponseWriter, r *http.Request, quadkey string, version string) {
	x, y, z := quadKeyToTileXY(quadkey)
	fmt.Printf("%s => %d/%d/%d\n", quadkey, x, y, z)

	// splits out the URL to get the x,y,z coordinates
	tileZ, tileX, tileY := float64(z), float64(x)-1, float64(y)-1

	myimage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{tileSize, tileSize}})

	i := complex128(complex(0, 1))
	zoom := float64(math.Pow(2, float64(tileZ-2)))
	tileRange := 1 / zoom
	tileStartX := 1/zoom + (tileRange * tileX)
	tileStartY := 1/zoom + (tileRange * tileY)

	// This loop just fills the image tile with fractal data
	var wg sync.WaitGroup
	wg.Add(tileSize)
	for cx := 0; cx < tileSize; cx++ {
		go func(cx int) {
			for cy := 0; cy < tileSize; cy++ {
				x := -2 + tileStartX + (float64(cx)/tileSize)*tileRange
				y := -2 + tileStartY + (float64(cy)/tileSize)*tileRange

				// x and y are now in the range ~-2 -> +2
				z := complex128(complex(x, 0)) + complex128(complex(y, 0))*complex128(i)

				c := complex(0.274, 0.008)
				for n := 0; n < 100; n++ {
					z = z*z + complex128(c)
				}

				z = z * 10
				ratio := float64(2 * (real(z) / 2))
				r := math.Max(0, float64(255*(ratio-1)))
				b := math.Max(0, float64(255*(1-ratio)))
				g := float64(255 - b - r)
				col := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
				myimage.Set(cx, cy, col)
				wg.Done()
			}
		}(cx)
	}
	wg.Wait()
	writeJpeg(w, myimage)
}

func writeJpegRaw(w http.ResponseWriter, img image.Image) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		fmt.Println("bad jpg")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func writeJpeg(w http.ResponseWriter, img image.Image) {
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
