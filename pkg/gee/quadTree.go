package gee

import (
	"AGES/pkg/core"
	"math"
)

//TileXYToQuadKey converts from x y z to quadkey
func TileXYToQuadKey(x, y, z int) (quadkey string) {
	for i := z; i >= 0; i-- {
		bitmask := 1 << uint(i)
		digit := 0
		if y&bitmask == 0 {
			digit |= 2
			if x&bitmask == 0 {
				digit |= 1
			}
		} else if x&bitmask != 0 {
			digit |= 1
		}
		quadkey += string(digit + '0')
	}
	return quadkey
}

//QuadKeyToTileXY converts from quadkey to x y z
func QuadKeyToTileXY(quadkey string) (x, y, z int) {
	z = len(quadkey) - 1
	for i := z; i >= 0; i-- {
		bitmask := 1 << uint(i)
		digit := '0' + quadkey[z-i]

		if digit&2 != 0 {
			if digit&1 == 0 {
				x |= bitmask
			}
		} else {
			y |= bitmask
			if digit&1 != 0 {
				x |= bitmask
			}
		}
	}
	return x, y, z
}

//TileXYToBBox converts from x y z to bounding box
func TileXYToBBox(x, y, z int) (bbox core.BBox) {
	scale := 360.0 / (math.Pow(2.0, float64(z)))
	return core.BBox{
		Left:   -180.0 + (scale * float64(x)),
		Bottom: -90.0 + (scale * float64(y)),
		Right:  -180.0 + (scale * (float64(x) + 1)),
		Top:    -90.0 + (scale * (float64(y) + 1)),
	}
}
