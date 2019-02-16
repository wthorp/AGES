package gee

import (
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

//HandleFunc returns terrain
func (p *TerrainProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
	// quadkey := parts[1]
	filePath := core.ApplicationDir("AGES", r.URL.RawQuery)
	url := fmt.Sprintf("%s/flatfile?%s", p.URL, r.URL.RawQuery)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		net.DownloadFile(filePath, url)
	}
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}
