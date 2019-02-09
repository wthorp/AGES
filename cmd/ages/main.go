package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"AGES/pkg/gee"
	"AGES/pkg/sources"

	"github.com/gorilla/mux"
)

func main() {
	source, err := sources.NewSingleImage("pipe.jpg")
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