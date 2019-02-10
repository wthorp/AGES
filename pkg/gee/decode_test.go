package gee

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	khdb "AGES/pkg/gee/keyhole_dbroot"
)

const basePath = `C:\Users\Bill\Desktop\go\AGES\config\`

func TestQ2(t *testing.T) {
	rq := "q2-0313-q.3.json"
	//	func metadataHandler(w http.ResponseWriter, r *http.Request, quadkey string, version string) {
	rawPath := filepath.Join(basePath, rq)
	jsonPath := filepath.Join(basePath, rq+".json")

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		//load raw
		file, e := ioutil.ReadFile(rawPath)
		if e != nil {
			t.Error(e)
		}
		//decode raw
		XOR(file, []byte(defaultKey), true)
		drp := khdb.DbRootProto{}
		unProto(file, &drp)
		fmt.Printf("%+v\n", drp)
	}
}
