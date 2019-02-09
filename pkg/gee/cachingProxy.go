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
	quadkey, _ /*version */ := parts[1], parts[3]
	filePath := filepath.Join("config", r.URL.RawQuery)
	url := fmt.Sprintf("%s/flatfile?%s-%s-%s.%s", p.URL, parts[0], parts[1], parts[2], parts[3])
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		net.DownloadFile(filePath, url)
	}
	switch parts[0] {
	case "q2": //-q
		q2Handler(w, r, quadkey)
	case "f1": //-i
		f1Handler(w, r, quadkey, p.ImgHandler)
	case "f1c": //-t
		f1cHandler(w, r, quadkey)
	default:
		//Other examples:
		//flatfile?lf-0-icons/shield1_l.png&h=32
		//flatfile?db=tm&qp-0-q.5
		fmt.Printf("unhandled flatfile type %s\n", parts[0])
	}
}

//f1cHandler returns a dbRoot object
func f1cHandler(w http.ResponseWriter, r *http.Request, quadkey string) {
	filePath := filepath.Join("config", r.URL.RawQuery)
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("File error: %v\n", e)
		return
	}
	w.Write(file)
}

func unMarshalJSONFile(filePath string, jsonObject interface{}) {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		jsonObject = nil
		return
	}
	json.Unmarshal(file, jsonObject)
}
