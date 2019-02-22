package gee

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"

	"AGES/pkg/core"
	"AGES/pkg/net"
)

//TerrainProxy proxies terrain
type TerrainProxy struct {
	URL *url.URL
}

//TerrainTile does what it sounds like
type TerrainTile struct {
	Quads [4][]Face
}

//Face is one triangle of 3D points
type Face [3]Point3D

//Point3D is a point w/ XY and height
type Point3D struct {
	X float64
	Y float64
	Z float32
}

//HandleFunc returns terrain
func (p *TerrainProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawPath := core.ApplicationDir(r.URL.RawQuery + ".raw")
	jsonPath := core.ApplicationDir(r.URL.RawQuery + ".json")

	if _, err := os.Stat(rawPath); os.IsNotExist(err) {
		err = net.DownloadFile(rawPath, net.RemapURL(p.URL, r.URL))
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		//load raw
		file, e := ioutil.ReadFile(rawPath)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("File error: %v\n", e)
			return
		}
		//decode raw
		XOR(file, []byte(defaultKey), true)
		b, err := uncompressPacket(file)
		if err != nil {
			fmt.Printf("Err in TerrainProxy uncompressPacket:\n%v\n", err)
		}
		//https://raw.githubusercontent.com/AnalyticalGraphicsInc/cesium/master/Source/Workers/decodeGoogleEarthEnterprisePacket.js
		offset := 0
		terrainTiles := make([]TerrainTile, 0)
		for offset < len(b) {
			// Each tile is split into 4 parts
			tt, size := UnpackTile(b[offset:], 1, 0)
			offset += size
			terrainTiles = append(terrainTiles, tt)
		}
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseTerrainProvider.js
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseTerrainData.js
		//https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Workers/createVerticesFromGoogleEarthEnterpriseBuffer.js

		jb, err := json.MarshalIndent(terrainTiles, "", "  ") //convert to json
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile(jsonPath, jb) //write to disk
	}

	//todo:  read/write JSON file instead
	file, e := ioutil.ReadFile(rawPath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}

//UnpackTile converts bytes to 4 quadrants consisting of many faces / 3D triangles
func UnpackTile(buffer []byte, negativeAltitudeExponentBias, negativeElevationThreshold float32) (TerrainTile, int) {
	tt := TerrainTile{}
	dv := binary.LittleEndian
	offset := 0
	// Compute sizes
	for quad := 0; quad < 4; quad++ {
		b := buffer[offset:]
		quadSize := int(int32(dv.Uint32(b[0:4])))
		originX := math.Float64frombits(dv.Uint64(b[4:12])) * 180.0
		originY := math.Float64frombits(dv.Uint64(b[12:20])) * 180.0
		stepX := math.Float64frombits(dv.Uint64(b[20:28])) * 180.0
		stepY := math.Float64frombits(dv.Uint64(b[28:36])) * 180.0
		numPoints := int(int32(dv.Uint32(b[36:40])))
		numFaces := int(int32(dv.Uint32(b[40:44])))
		//level := int32(dv.Uint32(b[44:48]))
		//process points
		points := make([]Point3D, numPoints)
		for i := 0; i < numPoints; i++ {
			pBuf := b[48+(i*6):]
			points[i].X = originX + (float64(uint8(pBuf[0])) * stepX)
			points[i].Y = originY + (float64(uint8(pBuf[1])) * stepY)
			points[i].Z = math.Float32frombits(dv.Uint32(pBuf[2:6])) * 6371010.0
			// negative altitude values are stored as height/-2^32
			if points[i].Z < negativeElevationThreshold {
				points[i].Z *= negativeAltitudeExponentBias
			}

		}
		//process faces
		faces := make([]Face, numFaces)
		for j := 0; j < numFaces; j++ {
			fBuf := b[48+(numPoints*6)+(j*6):]
			faces[j][0] = points[dv.Uint16(fBuf[0:2])]
			faces[j][1] = points[dv.Uint16(fBuf[2:4])]
			faces[j][2] = points[dv.Uint16(fBuf[4:6])]
		}
		tt.Quads[quad] = faces

		if quadSize != 44+((numPoints+numFaces)*6) {
			fmt.Printf("unexpected quadsize %d vs %d\n", quadSize, 44+(numPoints+numFaces)*6)
		}
		offset += quadSize + 4 // + len(quadSize)
	}
	return tt, offset
}
