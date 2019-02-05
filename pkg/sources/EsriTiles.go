package sources

import (
    "encoding/xml"
    "fmt"
    "log"
	"os"
	"path"
	"strconv"
	"io/ioutil"
	"encoding/binary"
)

//EsriTileCache implements TileCache for ESRI local files
type EsriTileCache struct {
	CacheFormat string
	BaseDirectory string
	FileFormat string
	TileCache
}

//CacheInfo corresponds to an ESRI conf.xml document
type CacheInfo struct {
	TileCacheInfo struct {
		LODInfos []struct {
			LODInfo struct {
				LevelID int
			}
		}
		SpatialReference struct {
			WKID string
		}
		TileCols int
		TileRows int
	}
	TileImageInfo struct {
		CacheTileFormat string
	}
	CacheStorageInfo struct{
		StorageFormat string
		PacketSize *int
	}
}

//New returns a new EsriTileCache
func New(string basePath) (*EsriTileCache, error) {
	e := &EsriTileCache{}
	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil{
		return er
	}
	var cache CacheInfo
	err = xml.Unmarshal(bytes, &cache)
	if err != nil{
		return er
	}
	levelIds := make([]int, len(cache.TileCacheInfo.LODInfos))
	e.MinLevel = cache.TileCacheInfo.LODInfos[0]	
	e.MaxLevel = MinLevel
	for i, li := range(cache.TileCacheInfo.LODInfos){
		levelIds[i] = li.LODInfo.LevelID
		e.MaxLevel = li.LODInfo.LevelID
	}
	e.FileFormat = cache.TileImageInfo.CacheTileFormat
	e.BaseDirectory = Path.GetDirectoryName(basePath)
	e.CacheFormat = cache.CacheStorageInfo.StorageFormat
	packetSize := cache.CacheStorageInfo.PacketSize
	e.HasTransparency = (FileFormat == "PNG" || FileFormat == "PNG32" || FileFormat == "MIXED")
	e.EpsgCode = cache.TileCacheInfo.SpatialReference.WKID
	e.TileColumnSize = cache.TileCacheInfo.TileCols
	e.TileRowSize = cache.TileCacheInfo.TileRows
	e.ColsPerFile = 1
	e.RowsPerFile = 1
	if packetSize != nil{
		e.ColsPerFile, e.RowsPerFile = packetSize, packetSize
	}
}

func (etci *EsriTileCache) ReadTile(tile Tile, bool ignore512) []byte {
	//mandatory file size check
	if (!ignore512 && TileRowSize == 512 && TileColumnSize == 512){
		return ReadTile512(tile)
	}

	if (CacheFormat == "esriMapCacheStorageModeCompact"){
		return ReadCompactTile(tile)
	}
	return ReadExplodedTile(tile)
}

func (etci *EsriTileCache) WriteTile(tile Tile, tileData []byte, ignore512 bool) {
	//mandatory file size check
	if (!ignore512 && TileRowSize == 512 && TileColumnSize == 512) {
		WriteTile512(tile, tileData)
		return
	}
	if (CacheFormat == "esriMapCacheStorageModeCompact"){
		WriteCompactTile(tile, tileData)
	}else{
		WriteExplodedTile(tile, tileData)
	}
}

func (etci *EsriTileCache) ReadCompactTile(tile Tile) []byte {
	bundlxPath, bundlePath, imgDataIndex := GetFileInfo(tile)
	if (!File.Exists(bundlxPath)){
		return nil
	}
	using (FileStream bundlx = new FileStream(bundlxPath, FileMode.Open, FileAccess.Read, FileShare.Read))
	{
		bundlx.Seek((16 + (5 * imgDataIndex)), SeekOrigin.Begin)
		var buffer [8]byte
		bundlx.Read(buffer, 0, 5)
		imageStartIndex = BitConverter.ToInt64(buffer, 0)
		using (FileStream bundle = new FileStream(bundlePath, FileMode.Open, FileAccess.Read, FileShare.Read))
		{
			bundle.Seek(imageStartIndex, SeekOrigin.Begin)
			var imgLength [4]byte
			bundle.Read(imgLength, 0, 4)
			count = BitConverter.ToInt32(imgLength, 0)
			imgBytes := make([]byte, count, count)
			bundle.Read(imgBytes, 0, count)
			return imgBytes
		}
	}
}

func (etci *EsriTileCache) WriteCompactTile(tile Tile, tileData []byte){
	//int imgDataIndex
	//string bundlePath, bundlxPath
	//GetFileInfo(tile, out bundlxPath, out bundlePath, out imgDataIndex)
}

func (etci *EsriTileCache) GetFileInfo(tile Tile) (bundlxPath, bundlePath, imgDataIndex string){
	internalRow = tile.Row % RowsPerFile
	internalCol = tile.Column % ColsPerFile
	bundleRow = tile.Row - internalRow
	bundleCol = tile.Column - internalCol
	bundleBasePath = path.Join(BaseDirectory, "_alllayers", fmt.Sprintf("L{0:d2}", tile.Level), fmt.Sprintf("R{0:x4}C{1:x4}", bundleRow, bundleCol))
	bundlxPath = bundleBasePath + ".bundlx"
	bundlePath = bundleBasePath + ".bundle"
	imgDataIndex = (ColsPerFile * internalCol) + internalRow
	return bundlxPath, bundlePath, imgDataIndex
}
 
func (etci *EsriTileCache) ReadExplodedTile(tile Tile) []byte {
	return binary.Read(GetFilePath(tile))
}

func (etci *EsriTileCache) WriteExplodedTile(tile Tile, tileData []byte){
	binary.Write(GetFilePath(tile), tileData)
}

//GetFilePath return the primary file path, sans extension
func (etci *EsriTileCache) GetFilePath(tile Tile) string {
	level := fmt.Sprintf("L{0:d2}", tile.Level)
	row := fmt.Sprintf("R{0:x8}", tile.Row)
	column := fmt.Sprintf("C{0:x8}", tile.Column)
	filePath := path.Join(BaseDirectory, level, row, column)
	if (FileFormat == "JPEG"){
		return filePath + ".jpg" //JPEG
	}
	if (FileFormat != "MIXED"){
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
