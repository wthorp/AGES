package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"AGES/pkg/gee"
	//"AGES/pkg/sources"
	"AGES/pkg/sources/proxy"

	"github.com/gorilla/mux"
)

func main() {

	//const wmsURL = "https://svs.gsfc.nasa.gov/cgi-bin/wms?SERVICE=WMS&VERSION=1.3.0&REQUEST=GetMap&LAYERS=2915_21223&FORMAT=image/png&WIDTH=256&HEIGHT=256&CRS=CRS:84&STYLES="
	//source, err := proxy.NewWMS(wmsURL, "PNG", time.Minute)

	const wmsURL = "https://webgate.ec.europa.eu/estat/inspireec/gis/arcgis/services/Basemaps/Blue_marble_4326/MapServer/WMSServer?request=GetMap&service=WMS&VERSION=1.3&LAYERS=0&FORMAT=image/jpeg&WIDTH=256&HEIGHT=256&CRS=CRS:84&STYLES="
	source, err := proxy.NewWMS(wmsURL, "JPEG", time.Minute)

	//const token = `pk.eyJ1IjoiZGlnaXRhbGdsb2JlIiwiYSI6ImNpdHZ6ZDNpazAwNncyc282MHcwZzVsZWEifQ.CjgIsR3Z4V4pUxtcTGCf0g`
	//dgMapBox := `https://api.mapbox.com/styles/v1/digitalglobe/cinvzwtia001hb4nplxp8htn3/tiles/256/%d/%d/%d?access_token=`
	//source, err := sources.NewTMS(dgMapBox+token, "PNG", time.Minute)

	//source, err := sources.NewSingleImage(`C:\Users\Bill\Desktop\go\AGES\pipe.jpg`)
	if err != nil {
		log.Fatal("Imagery source:", err)
	}
	geeProxy := &gee.CachingProxy{
		URL:        "http://www.earthenterprise.org/3d/",
		ImgHandler: source.GetTile,
	}
	//create a url router to handle different endpoints
	r := mux.NewRouter()
	r.HandleFunc("/dbRoot.v5", gee.DBRootHandler)
	r.Handle("/flatfile", geeProxy)
	// Anything we don't yet handle, use a simple reverse proxy
	u, _ := url.Parse(geeProxy.URL)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Using simple proxying for %s\n", r.URL)
		httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
	})
	// Start the server
	if err := http.ListenAndServe(":8085", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
