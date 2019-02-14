package gee

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"AGES/pkg/core"
)

//TerrainHandler returns terrain
func TerrainHandler(w http.ResponseWriter, r *http.Request, quadkey string) {
	filePath := core.ApplicationDir("config", r.URL.RawQuery)
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}
