package net

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

//RemapURL appends the path and query of one URL to another
func RemapURL(newBase, oldURL *url.URL) string {
	newURL := *newBase
	newURL.RawPath = path.Join(newBase.RawPath, oldURL.RawPath)
	newURL.RawQuery = oldURL.RawQuery
	return newURL.String()
}

//DownloadFile persist HTTP content to disk
func DownloadFile(path string, url string) error {
	//ensure directory
	dirPath := filepath.Dir(path)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
