package gee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"AGES/pkg/core"
)

//metadataHandler returns a q2 metadata object
func metadataHandler(w http.ResponseWriter, r *http.Request, quadkey string) {
	rawPath := core.ApplicationDir("config", r.URL.RawQuery)
	jsonPath := core.ApplicationDir("config", r.URL.RawQuery+".json")

	//url := path.Join(proxiedURL, "flatfile?"+r.URL.RawQuery)
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
		mdBytes, err := uncompressPacket(file)
		ti, err := processMetadata(mdBytes, len(mdBytes), quadkey)
		if err != nil {
			fmt.Printf("Err in q2 metdata:\n%v\n", err)
		} else {
			//write JSON
			b, err := json.MarshalIndent(ti, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			writeFile(jsonPath, b)
		}
	}

	//get  TileInformation map json data
	qp := &QtPacket{}
	err := unMarshalJSONFile(jsonPath, qp)
	if err != nil {
		fmt.Println("ti json")
		fmt.Fprintln(w, "ti json")
		w.WriteHeader(http.StatusNotImplemented)
		return
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
