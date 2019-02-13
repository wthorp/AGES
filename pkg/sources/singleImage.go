package sources

import (
	"io/ioutil"
)

//SingleImage is a single image 'cache'
type SingleImage struct {
	imageBytes []byte
}

//NewSingleImage return a single image for all tile requests
func NewSingleImage(imgPath string) (*SingleImage, error) {
	imageBytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}
	return &SingleImage{imageBytes: imageBytes}, nil
}

//GetTile returns the cached image
func (s *SingleImage) GetTile(x, y, z int) ([]byte, error) {
	//fmt.Printf("x %d y %d z %d\n", x, y, z)
	return s.imageBytes, nil
}
