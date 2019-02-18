package gee

import (
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
	if _, err := os.Stat(rawFilePath); os.IsNotExist(err) {
		b := readFile("config/dbRoot.raw") //read the protocol buffer
		XOR(b, []byte(defaultKey), true)   //unobfuscate


		//https://raw.githubusercontent.com/AnalyticalGraphicsInc/cesium/master/Source/Workers/decodeGoogleEarthEnterprisePacket.js
        offset := 0
        terrainTiles = make([]TerrainTile, 0)
        for offset < len(b){
            // Each tile is split into 4 parts
			tileStart = offset
			tt := TerrainTile{}
            for quad = 0; quad < 4; quad++ {
				size = dv.getUint32(offset, true)
				tt.Quads[quad] = b[offset:]
                offset += sizeOfUint32
                offset += size
            }
            tile = buffer.slice(tileStart, offset)
            terrainTiles = append(terrainTiles, tile)
		}
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseTerrainProvider.js
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseTerrainData.js



		b, err := json.MarshalIndent(drp, "", "  ") //convert to json
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile("config/dbRoot.js", b) //write to disk
		edrp.DbrootData = nil
		e, err := json.MarshalIndent(edrp, "", "  ") //convert to json
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile("config/encDbRoot.js", e) //write to disk
	}

	url := fmt.Sprintf("%s/flatfile?%s", p.URL, r.URL.RawQuery)
	if _, err := os.Stat(filePath) os.IsNotExist(err) {
		err = net.DownloadFile(filePath, url)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}
