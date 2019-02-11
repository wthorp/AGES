package gee

import (
	"fmt"
	"net/http"
)

//metadataHandler2 returns a q2 metadata object
func metadataHandler2(w http.ResponseWriter, r *http.Request, quadkey string) {
	numLevels := 4
	numInstances := ((1 << uint(numLevels*2)) - 1) / 3 //4^n-1/3
	qp := &QtPacket{
		Header:     NewQtHeader(numInstances),
		Tiles:      make([]TileInformation, numInstances, numInstances),
		DataBuffer: nil,
		MetaBuffer: nil,
	}
	level = len(quadkey)
	populateTiles(quadkey, len(quadkey))

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

func populateTiles(parentKey string, parent TileInformation, level int, tileInfo map[string]TileInformation, index *int) {
	isLeaf := false
	if level == 4 {
		if parent.HasSubtree() {
			return // We have a subtree, so just return
		}
		isLeaf = true // No subtree, so set all children to null
	}
	for i := uint(0); i <= 4; i++ {
		var childKey = fmt.Sprintf("%s%d", parentKey, i)
		if isLeaf {
			// No subtree so set all children to null
			// tileInfo[childKey] = nil
		} else if level < 4 {
			// We are still in the middle of the subtree, so add child
			//  only if their bits are set, otherwise set child to null.
			if !parent.HasChild(i) {
				//tileInfo[childKey] = nil
			} else {
				if index == numInstances {
					console.log("Incorrect number of instances")
					return nil
				}

				var instance = instances[index]
				index++
				tileInfo[childKey] = instance
				populateTiles(childKey, instance, level+1)
			}
		}
	}
}
