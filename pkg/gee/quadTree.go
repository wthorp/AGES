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
// https://github.com/google/earthenterprise/blob/950c525ce56aca2ee6624199339feafe47af113e/earth_enterprise/src/common/quadtreepath.cpp#L122
// https://github.com/AnalyticalGraphicsInc/cesium/blob/5d00c8ea29d18748dd8871b77c10b184986774bc/Source/Core/GoogleEarthEnterpriseMetadata.js#L247
func QuadKeyToTileXY(quadkey string) (x, y, z int) {
	z = len(quadkey) - 1
	for i := 0; i < z; i++ {
		bitmask := 1 << uint(i)
		switch quadkey[z-i] {
		case '0':
			y |= bitmask
		case '1':
			x |= bitmask
			y |= bitmask
		case '2':
			x |= bitmask
		}
	}
	return x, y, z
}
