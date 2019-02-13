package tilecache_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"AGES/pkg/sources/tilecache"
)

const layersDir1 = `C:\Users\Bill\Desktop\EsriTileCache\Layers`
const layersDir2 = `C:\Users\Bill\AppData\Local\ESRI\Local Caches\MapCacheV1\L9C4C\L9C4D`
const layersDir3 = `C:\Users\Bill\AppData\Local\ESRI\Local Caches\MapCacheV1\LCC8D\LCC8E`

//TestNewEsri tests the Esri constructor
func TestNewEsri(t *testing.T) {
	testProperties(t, layersDir1, 512, 4326, 0, 9)
	testProperties(t, layersDir2, 256, 102100, 1, 18)
	testProperties(t, layersDir3, 256, 102100, 2, 2)
}

func testProperties(t *testing.T, layersDir string, pixels, epsg, min, max int) {
	tc, err := tilecache.NewEsri(filepath.Join(layersDir, "conf.xml"))
	require.NoError(t, err)
	//Esri properties
	require.Equal(t, "esriMapCacheStorageModeCompact", tc.CacheFormat)
	require.Equal(t, layersDir, tc.BaseDirectory)
	require.Equal(t, "JPEG", tc.FileFormat)
	//TileCache properties
	require.Equal(t, false, tc.HasTransparency)
	require.Equal(t, pixels, tc.TileColumnSize)
	require.Equal(t, pixels, tc.TileRowSize)
	require.Equal(t, 128, tc.ColsPerFile)
	require.Equal(t, 128, tc.RowsPerFile)
	require.Equal(t, epsg, tc.EpsgCode)
	require.Equal(t, min, tc.MinLevel)
	require.Equal(t, max, tc.MaxLevel)
}
