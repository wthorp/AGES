package gee

import (
	"net/http"
)

//TerrainGen returns terrain DEMs
type TerrainGen struct {
	Provider ImageryProvider
}

//ServeHTTP returns a terrain DEMs
func (p *TerrainGen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	return
}
