package gee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

//q2Handler returns a q2 metadata object
func q2Handler(w http.ResponseWriter, r *http.Request, quadkey string) {
	rawPath := filepath.Join("config", r.URL.RawQuery)
	jsonPath := filepath.Join("config", r.URL.RawQuery+".json")

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
			b, err := json.Marshal(ti)
			if err != nil {
				fmt.Println("error:", err)
			}
			writeFile(jsonPath, b)
		}
	}

	//get  TileInformation map json data
	ti := []TileInformation{}
	unMarshalJSONFile(jsonPath, &ti)
	if ti == nil {
		fmt.Fprintln(w, "ti json")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//the old thing
	filePath := filepath.Join("config", r.URL.RawQuery)
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}
