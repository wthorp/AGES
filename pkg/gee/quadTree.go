package gee

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
