package core

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
)

type Tile struct {
	Row, Level, Column, EpsgCode int
}

type TileBundle struct {
	Level, MinRow, MaxRow, MinCol, MaxCol, EpsgCode int
}

type TileCache struct {
	HasTransparency                                       bool
	TileColumnSize, TileRowSize, ColsPerFile, RowsPerFile int
	EpsgCode, MinLevel, MaxLevel                          int
}

type BBox struct {
	Left, Bottom, Right, Top float64
}

//PNGBytes encodes an Image as a PNG and returns its bytes
func PNGBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//JPEGBytes encodes an Image as a JPEG and returns its bytes
func JPEGBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
