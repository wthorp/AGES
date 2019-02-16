package net

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

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
