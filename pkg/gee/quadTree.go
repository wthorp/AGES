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

//BingQuadKeyToTileXY converts from quadkey to x y z
func BingQuadKeyToTileXY(quadkey string) (x, y, z int) {
	z = len(quadkey)
	for i := z; i > 0; i-- {
		mask := 1 << (uint(i) - 1)
		switch quadkey[z-i] {
		case '0':
		case '1':
			x |= mask
		case '2':
			y |= mask
		case '3':
			x |= mask
			y |= mask
		default:
			panic("bad quadkey")
		}
	}
	return
}

//QuadKeyToTileXY converts from quadkey to x y z
func QuadKeyToTileXY(quadkey string) (x, y, z int) {
	z = len(quadkey) - 1
	//for i := z; i >= 0; --i {
	for i := z; i >= 1; {
		i = i - 1
		bitmask := 1 << uint(i)
		digit := quadkey[z-i]

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
	return
}

//OldQuadKeyToTileXY converts from quadkey to x y z
func OldQuadKeyToTileXY(quadkey string) (x, y, z int) {
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
