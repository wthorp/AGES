package gee

import (
	"fmt"
	"net/http"
)

const maxDepth = 15

//metadataHandler2 returns a q2 metadata object
func metadataHandler2(w http.ResponseWriter, r *http.Request, quadkey string) {
	// level := len(quadkey)
	// numLevels := 4
	// if maxDepth-level < numLevels {
	// 	numLevels = maxDepth % numLevels
	// }
	// numInstances := ([]int{1, 5, 21, 85, 341, 1365})[4]
	//numInstances := ((1 << uint(numLevels*2)) + 1) / 3 //4^n+1/3
	tiles := populateTiles(quadkey, true, make([]TileInformation, 0, 0))
	fmt.Printf("FOUND %d @ %s\n", len(tiles), quadkey)
	qp := &QtPacket{
		Header:     NewQtHeader(len(tiles)),
		Tiles:      tiles,
		DataBuffer: nil,
		MetaBuffer: nil,
	}

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

func populateTiles(quadkey string, isRoot bool, tileInfos []TileInformation) []TileInformation {
	level := len(quadkey)
	isLeaf := !isRoot && level%4 == 1
	ti := TileInformation{}
	ti.SetDefaults(quadkey, !isLeaf)
	tileInfos = append(tileInfos, ti)
	if isLeaf {
		return tileInfos
	}
	for i := uint(0); i < 4; i++ { //depth first packing
		if ti.Bits&(1<<i) != 0 {
			tileInfos = populateTiles(fmt.Sprintf("%s%d", quadkey, i), false, tileInfos)
		}
	}
	return tileInfos
}
