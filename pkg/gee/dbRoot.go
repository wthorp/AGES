package gee

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func check(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
	}
}

func decompress(data []byte) []byte {
	memReader := bytes.NewReader(data)
	zData, err := zlib.NewReader(memReader)
	check(err, "failed to create zlib reader 1")
	defer zData.Close()
	bytes, err := ioutil.ReadAll(zData)
	check(err, "failed to create zlib reader 2")
	return bytes
}

func writeFile(path string, data []byte) {
	ioutil.WriteFile(path, data, 0644)
}

func readFile(path string) []byte {
	in, err := ioutil.ReadFile(path)
	check(err, "failed to read file from disk")
	return in
}

func unProto(in []byte, pb proto.Message) {
	err := proto.Unmarshal(in, pb)
	check(err, fmt.Sprintf("failed to unmarshal protocol buffer %v\n", reflect.TypeOf(pb)))
}

func isMagic(data []byte) bool {
	return (data[0] == 116 && data[1] == 104 && data[2] == 222 && data[3] == 173) ||
		(data[3] == 116 && data[2] == 104 && data[1] == 222 && data[0] == 173)
}

//XOR encrypts / obfuscates data (in place)
func XOR(data []byte, key []byte, isDecode bool) {
	if isDecode && isMagic(data[0:4]) {
		return //this data isn't XOR encoded
	}

	var dp = 0
	var dpend = len(data)
	var dpend64 = dpend - (dpend % 8)
	var kpend = len(key)
	var kp = 0
	var off = 8

	// while we have a full uint64 (8 bytes) left to do
	// assumes buffer is 64bit aligned (or processor doesn't care)
	for dp < dpend64 {
		// rotate the key each time through by using the offets 16,0,8,16,0,8,...
		off = (off + 8) % 24
		kp = off

		// run through one key length xor'ing one uint64 at a time
		// then drop out to rotate the key for the next bit
		for (dp < dpend64) && (kp < kpend) {
			data[dp] = data[dp] ^ key[kp]
			data[dp+1] = data[dp+1] ^ key[kp+1]
			data[dp+2] = data[dp+2] ^ key[kp+2]
			data[dp+3] = data[dp+3] ^ key[kp+3]
			data[dp+4] = data[dp+4] ^ key[kp+4]
			data[dp+5] = data[dp+5] ^ key[kp+5]
			data[dp+6] = data[dp+6] ^ key[kp+6]
			data[dp+7] = data[dp+7] ^ key[kp+7]
			dp += 8
			kp += 24
		}
	}

	// now the remaining 1 to 7 bytes
	if dp < dpend {
		if kp >= kpend {
			// rotate the key one last time (if necessary)
			off = (off + 8) % 24
			kp = off
		}

		for dp < dpend {
			data[dp] = data[dp] ^ key[kp]
			dp++
			kp++
		}
	}
}

var compressedMagic uint32 = 0x7468dead
var compressedMagicSwap uint32 = 0xadde6874

func uncompressPacket(data []byte) ([]byte, error) {
	// The layout of this decoded data is
	// Magic Uint32 / Size Uint32 / [GZipped chunk of Size bytes]

	// Pullout magic and verify we have the correct data
	buf := bytes.NewBuffer(data)

	var size uint32
	var magic uint32
	var ndn binary.ByteOrder

	binary.Read(buf, binary.LittleEndian, &magic)
	if magic == compressedMagic {
		ndn = binary.LittleEndian
	} else if magic == compressedMagicSwap {
		fmt.Println("SWAP!")
		ndn = binary.BigEndian
	} else {
		fmt.Printf("Bad magic %d %d %d\n", magic, compressedMagic, compressedMagicSwap)
		return []byte{}, fmt.Errorf("bad magic")
	}
	binary.Read(buf, ndn, &size)
	memReader := bytes.NewReader(data[8:])
	zData, err := zlib.NewReader(memReader)
	defer zData.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to create zlib reader")
	}
	bytes, err := ioutil.ReadAll(zData)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read from zlib stream")
	}

	if size != uint32(len(bytes)) {
		return []byte{}, fmt.Errorf("zlib uncompress returned unexpected length")
	}
	return bytes, nil
}

func compressPacket(data []byte) ([]byte, error) {
	// pack as Magic Uint32 / Size Uint32 / [GZipped chunk of Size bytes]
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, compressedMagic)
	binary.Write(&b, binary.LittleEndian, uint32(len(data)))
	wz := zlib.NewWriter(&b)
	wz.Write(data)
	wz.Flush()
	wz.Close() //not sure why this can't be deferred
	return b.Bytes(), nil
}
