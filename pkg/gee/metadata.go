package gee

// https://github.com/google/earthenterprise/blob/master/earth_enterprise/src/common/qtpacket/quadtreepacket.h
// https://github.com/AnalyticalGraphicsInc/cesium/blob/master/Source/Core/GoogleEarthEnterpriseMetadata.js

import (
	"encoding/binary"
	"fmt"
)

const (
	qtMagic              = 32301
	anyChildBitmask byte = 0x0F //0x01 & 0x02 & 0x04 & 0x08
	leafBitmask     byte = 0x10
	vectorBitmask   byte = 0x20
	imageBitmask    byte = 0x40
	terrainBitmask  byte = 0x80
)

//QtPacket is a quadtree packet
type QtPacket struct {
	Header     QtHeader
	Tiles      []TileInformation
	DataBuffer []byte
	MetaBuffer []byte
}

//QtHeader is a quadtree header packet
type QtHeader struct {
	MagicID          uint32
	DataTypeID       uint32
	Version          uint32 //Version of the request for subtree metadata.
	NumInstances     int32
	DataInstanceSize int32
	DataBufferOffset int32
	DataBufferSize   int32
	MetaBufferSize   int32
}

//TileInformation describes a tile
type TileInformation struct {
	Bits            byte //junk uint8
	CnodeVersion    uint16
	ImageryVersion  uint16
	TerrainVersion  uint16
	NumChannels     uint16 //junk uint16
	TypeOffset      int32
	VersionOffset   int32
	ImageNeighbors  [8]byte
	ImageryProvider uint8
	TerrainProvider uint8 //junk uint16
}

func processMetadata(buffer []byte, totalSize int, quadKey string) (*QtPacket, error) {
	qp := &QtPacket{}
	err := qp.Header.UnmarshalBinary(buffer[0:32])
	if err != nil {
		return nil, err
	}
	// verify the packets is all there header + instances + dataBuffer + metaBuffer
	if qp.Header.DataBufferOffset+qp.Header.DataBufferSize+qp.Header.MetaBufferSize != int32(totalSize) {
		return nil, fmt.Errorf("invalid packet offsets")
	}
	// read all the instances
	qp.Tiles = make([]TileInformation, qp.Header.NumInstances, qp.Header.NumInstances)
	for i := int32(0); i < qp.Header.NumInstances; i++ {
		// i+1 because dataInstanceSize == sizeof(QtHeader) == 32
		if err = qp.Tiles[i].UnmarshalBinary(buffer[32*(i+1) : 32*(i+2)]); err != nil {
			return nil, err
		}
	}
	// read data and meta buffers
	dbStart := 32 * (qp.Header.NumInstances + 1)
	mbStart := dbStart + qp.Header.DataBufferSize
	qp.DataBuffer = buffer[dbStart:mbStart]
	qp.MetaBuffer = buffer[mbStart : mbStart+qp.Header.MetaBufferSize]
	return qp, nil
}

func unprocessMetadata(quadKey string, qp *QtPacket) ([]byte, error) {
	bufferSize := (len(qp.Tiles) + 1) * 32 // +1 because dataInstanceSize == sizeof(QtHeader) == 32
	buffer := make([]byte, bufferSize)
	err := qp.Header.MarshalBinary(buffer[0:32])
	if err != nil {
		return nil, err
	}
	// Read all the instances
	for i := int32(0); i < qp.Header.NumInstances; i++ {
		// i+1 because dataInstanceSize == sizeof(QtHeader) == 32
		instanceBuffer := buffer[32*(i+1) : 32*(i+2)]
		if err = qp.Tiles[i].MarshalBinary(instanceBuffer); err != nil {
			return nil, err
		}
		//qp.Tiles[i].ImageryProvider = 1 // hack the system
	}
	//qp.Header.DataBufferSize = int32(len(qp.DataBuffer))
	buffer = append(buffer, qp.DataBuffer...)
	//qp.Header.MetaBufferSize = int32(len(qp.MetaBuffer))
	buffer = append(buffer, qp.MetaBuffer...)

	totalSize := len(buffer)
	if qp.Header.DataBufferOffset+qp.Header.DataBufferSize+qp.Header.MetaBufferSize != int32(totalSize) {
		return nil, fmt.Errorf("invalid packet offsets")
	}
	return buffer, nil
}

//Validate checks a quadtree header for correctness
func (qt *QtHeader) Validate() error {
	if qt.MagicID != qtMagic {
		return fmt.Errorf("invalid quadtree header magic")
	}
	if qt.DataTypeID != 1 {
		return fmt.Errorf("invalid quadtree header data type; must be 1 for QuadTreePacket")
	}
	// Tile format version
	if qt.Version != 2 {
		return fmt.Errorf("invalid quadtree header version; only version 2 is supported")
	}
	if qt.DataInstanceSize != 32 {
		return fmt.Errorf("invalid quadtree header instance size")
	}
	// Offset from beginning of packet (instances + current offset)
	if qt.DataBufferOffset != (qt.NumInstances*qt.DataInstanceSize + 32) {
		return fmt.Errorf("invalid quadtree header dataBufferOffset")
	}
	return nil
}

//NewQtHeader returns a new quadtree header packet
func NewQtHeader(numInstances int) QtHeader {
	return QtHeader{
		MagicID:          32301,
		DataTypeID:       1,
		Version:          2,
		NumInstances:     int32(numInstances),
		DataInstanceSize: 32,
		DataBufferOffset: int32((numInstances + 1) * 32),
		DataBufferSize:   0,
		MetaBufferSize:   0,
	}
}

//SetDefaults set a quadtree tile node to default values
func (ti *TileInformation) SetDefaults(quadkey string, hasSubTree bool) {
	childBitMasks := anyChildBitmask
	if len(quadkey) == 2 {
		if quadkey[1] < '2' {
			childBitMasks = 0x04 | 0x08 // hide < -180
		} else {
			childBitMasks = 0x01 | 0x02 // hide > 180
		}
	}
	if hasSubTree {
		ti.Bits = imageBitmask | childBitMasks
		ti.CnodeVersion = 0
	} else {
		ti.Bits = imageBitmask | leafBitmask
		ti.CnodeVersion = 3
	}
	ti.ImageryProvider = 1
}

//MarshalBinary returns QtHeader to a binary form
func (qt *QtHeader) MarshalBinary(data []byte) error {
	if len(data) != 32 {
		return fmt.Errorf("Bad QtHeader byte length in MarshalBinary")
	}
	dv := binary.LittleEndian
	dv.PutUint32(data[0:4], qtMagic)
	dv.PutUint32(data[4:8], qt.DataTypeID)
	dv.PutUint32(data[8:12], qt.Version)
	dv.PutUint32(data[12:16], uint32(qt.NumInstances))
	dv.PutUint32(data[16:20], uint32(qt.DataInstanceSize))
	dv.PutUint32(data[20:24], uint32(qt.DataBufferOffset))
	dv.PutUint32(data[24:28], uint32(qt.DataBufferSize))
	dv.PutUint32(data[28:32], uint32(qt.MetaBufferSize))
	return qt.Validate()
}

//UnmarshalBinary returns QtHeader from a binary form
func (qt *QtHeader) UnmarshalBinary(data []byte) error {
	dv := binary.LittleEndian
	qt.MagicID = dv.Uint32(data[0:4])
	qt.DataTypeID = dv.Uint32(data[4:8])
	qt.Version = dv.Uint32(data[8:12])
	qt.NumInstances = int32(dv.Uint32(data[12:16]))
	qt.DataInstanceSize = int32(dv.Uint32(data[16:20]))
	qt.DataBufferOffset = int32(dv.Uint32(data[20:24]))
	qt.DataBufferSize = int32(dv.Uint32(data[24:28]))
	qt.MetaBufferSize = int32(dv.Uint32(data[28:32]))
	return qt.Validate()
}

//UnmarshalBinary returns TileInformation from a binary form
func (ti *TileInformation) UnmarshalBinary(data []byte) error {
	dv := binary.LittleEndian
	ti.Bits = data[0]
	ti.CnodeVersion = dv.Uint16(data[2:4])
	ti.ImageryVersion = dv.Uint16(data[4:6])
	ti.TerrainVersion = dv.Uint16(data[6:8])
	ti.NumChannels = dv.Uint16(data[8:10])
	ti.TypeOffset = int32(dv.Uint32(data[12:16]))
	ti.VersionOffset = int32(dv.Uint32(data[16:20]))
	for x := 0; x < 8; x++ {
		ti.ImageNeighbors[x] = data[20+x]
	}
	ti.ImageryProvider = data[28]
	ti.TerrainProvider = data[29]
	return nil
}

//MarshalBinary returns TileInformation in a binary form
func (ti *TileInformation) MarshalBinary(data []byte) error {
	if len(data) != 32 {
		return fmt.Errorf("Bad TileInformation byte length in MarshalBinary")
	}
	dv := binary.LittleEndian
	data[0] = ti.Bits
	dv.PutUint16(data[2:4], ti.CnodeVersion)
	dv.PutUint16(data[4:6], ti.ImageryVersion)
	dv.PutUint16(data[6:8], ti.TerrainVersion)
	dv.PutUint16(data[8:10], ti.NumChannels)
	dv.PutUint32(data[12:16], uint32(ti.TypeOffset))
	dv.PutUint32(data[16:20], uint32(ti.VersionOffset))
	for x := 0; x < 8; x++ {
		data[20+x] = ti.ImageNeighbors[x]
	}
	data[28] = ti.ImageryProvider
	data[29] = ti.TerrainProvider
	return nil
}
