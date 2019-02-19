package gee

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"AGES/pkg/core"
	"AGES/pkg/net"
)

//TerrainProxy proxies terrain
type TerrainProxy struct {
	URL string
}

//TerrainTile does what it sounds like
type TerrainTile struct {
	Quads [4][]byte
}

//HandleFunc returns terrain
func (p *TerrainProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawFilePath := core.ApplicationDir("AGES", r.URL.RawQuery) + ".raw"
	jsFilePath := core.ApplicationDir("AGES", r.URL.RawQuery) + ".js"

	if _, err := os.Stat(rawFilePath); os.IsNotExist(err) {
		err = net.DownloadFile(rawFilePath, r.URL.RawQuery)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	if _, err := os.Stat(jsFilePath); os.IsNotExist(err) {
		b := readFile(rawFilePath)       //read the protocol buffer
		XOR(b, []byte(defaultKey), true) //unobfuscate

		//https://raw.githubusercontent.com/AnalyticalGraphicsInc/cesium/master/Source/Workers/decodeGoogleEarthEnterprisePacket.js
		offset := 0
		terrainTiles := make([]TerrainTile, 0)
		for offset < len(b) {
			// Each tile is split into 4 parts
			tt := TerrainTile{}
			for quad := 0; quad < 4; quad++ {
				size := int(binary.LittleEndian.Uint32(b[offset : offset+4]))
				offset += 4
				tt.Quads[quad] = b[offset : offset+size]
				offset += size
			}
			terrainTiles = append(terrainTiles, tt)
		}
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseTerrainProvider.js
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseTerrainData.js

		b, err := json.MarshalIndent(terrainTiles, "", "  ") //convert to json
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile(jsFilePath, b) //write to disk
	}

	//todo:  read/write JSON file instead
	file, e := ioutil.ReadFile(rawFilePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}
