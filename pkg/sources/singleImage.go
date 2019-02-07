package sources

import (
	"io/ioutil"
)

type SingleImage struct {
	imageBytes []byte
}

//NewSingleImage return a single image for all tile requests
func NewSingleImage(imgPath string) (SingleImage, error) {
	imageBytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return err
	}
	return SingleImage{imageBytes: imageBytes}
}

//Pipe returns Magritte's pipe image
func (s *SingleImage) GetTile(x, y, z int) ([]byte, error) {
	return s.imageBytes
}
