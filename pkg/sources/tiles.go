package main

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
