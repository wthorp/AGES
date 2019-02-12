package gee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"AGES/pkg/net"
)

//CachingProxy proxies files
type CachingProxy struct {
	URL        string
	ImgHandler func(int, int, int) ([]byte, error)
}

func (p *CachingProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
	quadkey := parts[1]
	filePath := filepath.Join("config", r.URL.RawQuery)
	url := fmt.Sprintf("%s/flatfile?%s-%s-%s.%s", p.URL, parts[0], parts[1], parts[2], parts[3])
	switch parts[0] {
	case "q2": //-q
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			net.DownloadFile(filePath, url)
		}
		metadataHandler2(w, r, quadkey)
	case "f1": //-i
		imageryHandler(w, r, quadkey, p.ImgHandler)
	case "f1c": //-t
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			net.DownloadFile(filePath, url)
		}
		terrainHandler(w, r, quadkey)
	default:
		//Other examples:
		//flatfile?lf-0-icons/shield1_l.png&h=32
		//flatfile?db=tm&qp-0-q.5
		fmt.Printf("unhandled flatfile type %s\n", parts[0])
	}
}

//terrainHandler returns terrain
func terrainHandler(w http.ResponseWriter, r *http.Request, quadkey string) {
	filePath := filepath.Join("config", r.URL.RawQuery)
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}

func unMarshalJSONFile(filePath string, jsonObject interface{}) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, jsonObject)
}
