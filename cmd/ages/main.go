package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"AGES/pkg/gee"
	//"AGES/pkg/sources"
	"AGES/pkg/sources/proxy"

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

	//"https://webgate.ec.europa.eu/estat/inspireec/gis/arcgis/services/Basemaps/Blue_marble_4326/MapServer/WMSServer"
	//0
	//"https: //gis.apfo.usda.gov/arcgis/rest/services/NAIP/USDA_CONUS_PRIME/ImageServer"
	//"http://nsdig2gapps.ncsi.gov.om/arcgis/rest/services/Base_Map_EN/MapServer"
	//const token = `pk.eyJ1IjoiZGlnaXRhbGdsb2JlIiwiYSI6ImNpdHZ6ZDNpazAwNncyc282MHcwZzVsZWEifQ.CjgIsR3Z4V4pUxtcTGCf0g`
	//dgMapBox := `https://api.mapbox.com/styles/v1/digitalglobe/cinvzwtia001hb4nplxp8htn3/tiles/256/%d/%d/%d?access_token=`
	//source, err := sources.NewTMS(dgMapBox+token, "PNG", time.Minute)
	//source, err := sources.NewSingleImage(`C:\Users\Bill\Desktop\go\AGES\pipe.jpg`)

	var wms, layers string
	proxyURL := "http://www.earthenterprise.org/3d/"

	//get command line options
	flag.StringVar(&wms, "wms", "https://neo.sci.gsfc.nasa.gov/wms/wms", "WMS URL")
	flag.StringVar(&layers, "layers", "BlueMarbleNG-TB", "WMS layers parameter")
	flag.Parse()
	if wms == "" || layers == "" {
		flag.PrintDefaults()
		return
	}
	//configure providers
	wmsURL := wms + "?request=GetMap&service=WMS&VERSION=1.3&LAYERS=" + layers + "&FORMAT=image/jpeg&WIDTH=256&HEIGHT=256&CRS=CRS:84&STYLES="
	imgHandler, err := proxy.NewWMS(wmsURL, "JPEG", time.Minute)
	if err != nil {
		log.Fatal("Imagery source:", err)
	}

	//create a url router to handle different endpoints
	r := mux.NewRouter()
	r.HandleFunc("/dbRoot.v5", gee.DBRootHandler2)
	r.HandleFunc("/flatfile", func(w http.ResponseWriter, r *http.Request) {
		var parts = strings.FieldsFunc(r.URL.RawQuery, func(c rune) bool { return c == '-' || c == '.' })
		quadkey := parts[1]
		switch parts[0] {
		case "q2": //-q
			gee.MetadataHandler2(w, r, quadkey)
		case "f1": //-i
			gee.ImageryHandler(w, r, quadkey, imgHandler.GetTile)
		case "f1c": //-t
			//note:  this is functionally disabled by MetadataHandler2
			// filePath := core.ApplicationDir("config", r.URL.RawQuery)
			// url := fmt.Sprintf("%s/flatfile?%s-%s-%s.%s", p.URL, parts[0], parts[1], parts[2], parts[3])
			// if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// 	net.DownloadFile(filePath, url)
			// }
			//gee.TerrainHandler(w, r, quadkey)
		default:
			//Other examples:
			//flatfile?lf-0-icons/shield1_l.png&h=32
			//flatfile?db=tm&qp-0-q.5
			fmt.Printf("unhandled flatfile type %s\n", parts[0])
		}
	})
	// Anything we don't yet handle, use a simple reverse proxy
	u, _ := url.Parse(proxyURL)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Using simple proxying for %s\n", r.URL)
		httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
	})
	// Start the server
	if err := http.ListenAndServe(":8085", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
