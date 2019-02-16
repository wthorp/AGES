package gee

import (
	"fmt"
	"net/http"
	"strings"
)

//MetadataGen proxies terrain
type MetadataGen struct {
	MaxDepth   int
	HasTerrain bool
}

//MetadataHandler returns a q2 metadata object
func (p *MetadataGen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
	quadkey := parts[1]

	level := len(quadkey)
	numLevels := 4
	if p.MaxDepth-level < numLevels {
		numLevels = p.MaxDepth % numLevels
	}
	numInstances := 0
	skipBadLatitudesModifier := 0 // -1 to skip >|+/-180|
	for l := 0; l <= numLevels; l++ {
		if level+l == 3 { // second level starts the issue
			skipBadLatitudesModifier = -1
		}
		numInstances += 1 << uint((l*2)+skipBadLatitudesModifier)
	}
	//numInstances := ((1 << uint(numLevels*2)) + 1) / 3 //4^n+1/3
	//tiles := populateTiles(quadkey, true, make([]TileInformation, numInstances, numInstances))
	//fmt.Printf("FOUND %d @ %s\n", numInstances, quadkey)
	qp := &QtPacket{
		Header:     NewQtHeader(numInstances),
		Tiles:      make([]TileInformation, numInstances, numInstances),
		DataBuffer: nil,
		MetaBuffer: nil,
	}
	index := 0
	populateTiles(&index, quadkey, true, p.HasTerrain, qp.Tiles)

	mdBytes, err := unprocessMetadata(quadkey, qp)
	if err != nil {
		fmt.Println("unprocessMetadata", err.Error())
		fmt.Fprintln(w, "unprocessMetadata", err.Error())
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	compressedBytes, err := compressPacket(mdBytes)
	if err != nil {
		fmt.Println("compressPacket", err.Error())
		fmt.Fprintln(w, "compressPacket", err.Error())
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	XOR(compressedBytes, []byte(defaultKey), false)
	w.Write(compressedBytes)
}

//populateTiles creates a populated TileInformation array, four levels deep
func populateTiles(index *int, quadkey string, isRoot, hasTerrain bool, tileInfos []TileInformation) {
	level := len(quadkey)
	isLeaf := !isRoot && level%4 == 1
	tileInfos[*index].SetDefaults(quadkey, !isLeaf, hasTerrain)
	if isLeaf {
		return
	}
	bits := tileInfos[*index].Bits
	for i := uint(0); i < 4; i++ { //depth first packing
		if bits&(1<<i) != 0 {
			*index = *index + 1
			populateTiles(index, fmt.Sprintf("%s%d", quadkey, i), false, hasTerrain, tileInfos)
		}
	}
	return
}
