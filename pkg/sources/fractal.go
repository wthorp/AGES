package sources

import (
	"AGES/pkg/gee"
	"image"
	"image/color"
	"math"
	"sync"
)

type Fractal struct{}

//Fractal returns fractal imagery
func (f *Fractal) GetTile(x, y, z int) ([]byte, error) {
	// splits out the URL to get the x,y,z coordinates
	tileZ, tileX, tileY := float64(z), float64(x)-1, float64(y)-1

	image := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{tileSize, tileSize}})

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
				image.Set(cx, cy, col)
				wg.Done()
			}
		}(cx)
	}
	wg.Wait()
	gee.WriteJpeg(w, image)
}
