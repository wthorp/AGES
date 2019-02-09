package gee

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	khdb "AGES/pkg/gee/keyhole_dbroot"
)

func TestQ2(t *testing.T) {
	rq := "q2-03130003-q.3"
	//	func q2Handler(w http.ResponseWriter, r *http.Request, quadkey string, version string) {
	rawPath := filepath.Join("config", rq)
	jsonPath := filepath.Join("config", rq+".json")

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
