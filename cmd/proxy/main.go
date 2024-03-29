package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"AGES/pkg/gee"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println(`
	Copyright (C) 2018 William Patrick Thorp - All Rights Reserved.

	This software is for demo purposes only.   Do not distribute.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
	INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A 
	PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT 
	HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF 
	CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE 
	OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

	
	Please point Google Earth Enterprise to http://localhost:8085/
	Press Ctrl-C to exit
	`)

	var source string
	//get command line options
	flag.StringVar(&source, "source", `http://www.earthenterprise.org/3d/`, "GEE URL to proxy")
	flag.Parse()
	if source == "" {
		flag.PrintDefaults()
		return
	}
	sourceURL, err := url.Parse(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	rootHandler := &gee.DBRootProxy{URL: sourceURL}
	metadataHandler := &gee.MetadataProxy{URL: sourceURL}
	imageryHandler := &gee.ImageryProxy{URL: sourceURL}
	terrainHandler := &gee.TerrainProxy{URL: sourceURL}
	otherHandler := httputil.NewSingleHostReverseProxy(sourceURL)

	//create a url router to handle different endpoints
	r := mux.NewRouter()
	r.Handle("/dbRoot.v5", rootHandler)
	r.HandleFunc("/flatfile", func(w http.ResponseWriter, r *http.Request) {
		var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
		switch parts[0] {
		case "q2": //-q
			metadataHandler.ServeHTTP(w, r)
		case "f1": //-i
			imageryHandler.ServeHTTP(w, r)
		case "f1c": //-t
			terrainHandler.ServeHTTP(w, r)
		default:
			//Other examples:
			//flatfile?lf-0-icons/shield1_l.png&h=32
			//flatfile?db=tm&qp-0-q.5
			fmt.Printf("unhandled URL %s\n", r.URL)
			otherHandler.ServeHTTP(w, r)
		}
	})
	// Anything we don't yet handle, use a simple reverse proxy
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("unhandled URL %s\n", r.URL)
		otherHandler.ServeHTTP(w, r)
	})
	// Start the server
	if err := http.ListenAndServe(":8085", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
