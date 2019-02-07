package sources

import (
	"AGES/pkg/core"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//EsriTileCache implements TileCache for ESRI local files
type EsriTileCache struct {
	CacheFormat   string
	BaseDirectory string
	FileFormat    string
	core.TileCache
}

//CacheInfo corresponds to an ESRI conf.xml document
type CacheInfo struct {
	TileCacheInfo struct {
		LODInfos struct {
			LODInfo []struct {
				LevelID int
			}
		}
		SpatialReference struct {
			WKID int
		}
		TileCols int
		TileRows int
	}
	TileImageInfo struct {
		CacheTileFormat string
	}
	CacheStorageInfo struct {
		StorageFormat string
		PacketSize    *int
	}
}

//NewEsriTileCache returns a new EsriTileCache, based on a conf.xml path
func NewEsriTileCache(confPath string) (*EsriTileCache, error) {
	tc := &EsriTileCache{}
	confXML, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var cache CacheInfo
	err = xml.Unmarshal(confXML, &cache)
	if err != nil {
		return nil, err
	}
	tc.BaseDirectory = filepath.Dir(confPath)
	tc.MinLevel, tc.MaxLevel = calcMinMaxLevels(&cache, tc.BaseDirectory)
	tc.FileFormat = cache.TileImageInfo.CacheTileFormat
	tc.CacheFormat = cache.CacheStorageInfo.StorageFormat
	packetSize := cache.CacheStorageInfo.PacketSize
	tc.HasTransparency = (tc.FileFormat == "PNG" || tc.FileFormat == "PNG32" || tc.FileFormat == "MIXED")
	tc.EpsgCode = cache.TileCacheInfo.SpatialReference.WKID
	tc.TileColumnSize = cache.TileCacheInfo.TileCols
	tc.TileRowSize = cache.TileCacheInfo.TileRows
	if packetSize != nil {
		tc.ColsPerFile, tc.RowsPerFile = *packetSize, *packetSize
	} else {
		tc.ColsPerFile, tc.RowsPerFile = 1, 1
	}
	return tc, nil
}

//calcMinMaxLevels is called by NewEsriTileCache to return min and max levels
func calcMinMaxLevels(cache *CacheInfo, baseDir string) (int, int) {
	minLevel := int(^uint(0) >> 1)
	maxLevel := 0
	for _, li := range cache.TileCacheInfo.LODInfos.LODInfo {
		levelPath := path.Join(baseDir, "_alllayers", fmt.Sprintf("L%02d", li.LevelID))
		if _, err := os.Stat(levelPath); err != nil {
			continue
		}
		if li.LevelID > maxLevel {
			maxLevel = li.LevelID
		}
		if li.LevelID < minLevel {
			minLevel = li.LevelID
		}
	}
	if minLevel > maxLevel {
		minLevel = maxLevel
	}
	return minLevel, maxLevel
}

//ReadTile returns a 256x256 tile
func (tc *EsriTileCache) ReadTile(tile core.Tile) ([]byte, error) {
	if tc.CacheFormat == "esriMapCacheStorageModeCompact" {
		return tc.ReadCompactTile(tile)
	}
	return tc.ReadExplodedTile(tile)
}

//WriteTile writes a 256x256 tile
func (tc *EsriTileCache) WriteTile(tile core.Tile, tileData []byte) error {
	if tc.CacheFormat == "esriMapCacheStorageModeCompact" {
		return tc.WriteCompactTile(tile, tileData)
	} else {
		return tc.WriteExplodedTile(tile, tileData)
	}
}

//ReadCompactTile returns a bundled 256x256 tile
func (tc *EsriTileCache) ReadCompactTile(tile core.Tile) ([]byte, error) {
	bundlxPath, bundlePath, imgDataIndex := tc.GetFileInfo(tile)
	bundlx, err := os.Open(bundlxPath)
	if err != nil {
		return nil, err
	}
	defer bundlx.Close()
	bundlx.Seek((16 + (5 * imgDataIndex)), io.SeekStart)
	bOffset := make([]byte, 5, 5)
	bundlx.Read(bOffset)
	offset := int64(binary.LittleEndian.Uint64(bOffset))
	bundle, err := os.Open(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundle.Close()
	bundle.Seek(offset, io.SeekStart)
	bLength := make([]byte, 4, 4)
	bundle.Read(bLength)
	length := binary.LittleEndian.Uint64(bLength)
	imgBytes := make([]byte, length, length)
	bundle.Read(imgBytes)
	return imgBytes, nil
}

//WriteCompactTile writes a bundled 256x256 tile
func (tc *EsriTileCache) WriteCompactTile(tile core.Tile, tileData []byte) error {
	return fmt.Errorf("not implemented")
}

//GetFileInfo returns file paths and indexes into those files
func (tc *EsriTileCache) GetFileInfo(tile core.Tile) (bundlxPath, bundlePath string, imgDataIndex int64) {
	internalRow := tile.Row % tc.RowsPerFile
	internalCol := tile.Column % tc.ColsPerFile
	bundleRow := tile.Row - internalRow
	bundleCol := tile.Column - internalCol
	bundleBasePath := path.Join(tc.BaseDirectory, "_alllayers", fmt.Sprintf("L%02d", tile.Level), fmt.Sprintf("R%04xC%04x", bundleRow, bundleCol))
	bundlxPath = bundleBasePath + ".bundlx"
	bundlePath = bundleBasePath + ".bundle"
	imgDataIndex = int64((tc.ColsPerFile * internalCol) + internalRow)
	return bundlxPath, bundlePath, imgDataIndex
}

//ReadExplodedTile returns a standalone 256x256 tile
func (tc *EsriTileCache) ReadExplodedTile(tile core.Tile) ([]byte, error) {
	return ioutil.ReadFile(tc.GetFilePath(tile))
}

//WriteExplodedTile writes a standalone 256x256 tile
func (tc *EsriTileCache) WriteExplodedTile(tile core.Tile, tileData []byte) error {
	return ioutil.WriteFile(tc.GetFilePath(tile), tileData, 0644)
}

//GetFilePath return the primary file path, sans extension
func (tc *EsriTileCache) GetFilePath(tile core.Tile) string {
	level := fmt.Sprintf("L%02d", tile.Level)
	row := fmt.Sprintf("R%08x", tile.Row)
	column := fmt.Sprintf("C%08x", tile.Column)
	filePath := path.Join(tc.BaseDirectory, level, row, column)
	if tc.FileFormat == "JPEG" {
		return filePath + ".jpg" //JPEG
	}
	if tc.FileFormat != "MIXED" {
		return filePath + ".png" //PNG, PNG8, PNG24, PNG32
	}
	if _, err := os.Stat(filePath + ".jpg"); err == nil {
		return filePath + ".jpg" //MIXED...
	}
	if _, err := os.Stat(filePath + ".png"); err == nil {
		return filePath + ".png"
	}
	return filePath
}
