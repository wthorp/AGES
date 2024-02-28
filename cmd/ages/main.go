package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"AGES/pkg/gee"
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
	//https://basemap.nationalmap.gov/arcgis/rest/services/USGSTopo/MapServer/WMSServer
	//https://neo.sci.gsfc.nasa.gov/wms/wms | BlueMarbleNG-TB

	var wms, layers string
	//get command line options
	flag.StringVar(&wms, "wms", "https://basemap.nationalmap.gov/arcgis/services/USGSTopo/MapServer/WMSServer", "WMS URL")
	flag.StringVar(&layers, "layers", "0", "WMS layers parameter")
	flag.Parse()
	if wms == "" || layers == "" {
		flag.PrintDefaults()
		return
	}
	//configure providers
	wmsURL := wms + "?request=GetMap&service=WMS&VERSION=1.3&LAYERS=" + layers + "&FORMAT=image/jpeg&WIDTH=256&HEIGHT=256&CRS=CRS:84&STYLES="
	imgProvider, err := proxy.NewWMS(wmsURL, "JPEG", time.Minute)
	if err != nil {
		log.Fatal("Imagery source:", err)
	}

	rootHandler := &gee.DBRootGen{}
	metadataHandler := &gee.MetadataGen{MaxDepth: 15, HasTerrain: false}
	imageryHandler := &gee.ImageryGen{Provider: imgProvider}
	terrainHandler := &gee.TerrainGen{Provider: imgProvider}

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
		}
	})
	// Start the server
	if err := http.ListenAndServe(":8085", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
