package sources

import (
	"image"
	"image/color"
	"math"
	"sync"

	"AGES/pkg/core"
)

// Fractal returns fractal imagery
type Fractal struct{}

// NewFractal return a fractal imagery for all tile requests
func NewFractal() (*Fractal, error) {
	return &Fractal{}, nil
}

// GetTile returns fractal imagery
func (f *Fractal) GetTile(x, y, z int) ([]byte, error) {
	// splits out the URL to get the x,y,z coordinates
	tileZ, tileX, tileY := float64(z), float64(x)-1, float64(y)-1

	image := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{256, 256}})

	i := complex128(complex(0, 1))
	zoom := float64(math.Pow(2, float64(tileZ-2)))
	tileRange := 1 / zoom
	tileStartX := 1/zoom + (tileRange * tileX)
	tileStartY := 1/zoom + (tileRange * tileY)

	// This loop just fills the image tile with fractal data
	var wg sync.WaitGroup
	wg.Add(256)
	for cx := 0; cx < 256; cx++ {
		go func(cx int) {
			defer wg.Done() // Ensure that Done is called after the goroutine finishes
			for cy := 0; cy < 256; cy++ {
				x := -2 + tileStartX + (float64(cx)/256)*tileRange
				y := -2 + tileStartY + (float64(cy)/256)*tileRange

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
			}
		}(cx)
	}
	wg.Wait()
	return core.JPEGBytes(image)
}
