package Main

import (
    "encoding/xml"
    "fmt"
    "log"
    "os"
	"launchpad.net/xmlpath"
)
type EsriTileCache struct {
	string CacheFormat, BaseDirectory, FileFormat
}

type EsriTileBundle struct {
	TileCache    TileCache
	tilesPerFile int
	isCompact    bool
}

func (bundle *EsriTileBundle) Init(basePath string){

}