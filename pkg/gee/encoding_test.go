package gee

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//test Compress circular loops
func TestCompress(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := make([]byte, rand.Intn(10000))
		rand.Read(input)
		compressed, err := compressPacket(input)
		if err != nil {
			t.Errorf("failed during compress : %+v", err)
			break
		}
		output, err := uncompressPacket(compressed)
		if err != nil {
			t.Errorf("failed during uncompress : %+v", err)
			break
		}
		if !bytes.Equal(input, output) {
			t.Error("input != ouptut")
			break
		}
	}
	fmt.Printf("TestCompress passed\n")
}

//test Xor circular loops
func TestXor(t *testing.T) {
	key := []byte(defaultKey)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		size := rand.Intn(10000)
		input := make([]byte, size)
		output := make([]byte, size)
		rand.Read(input)
		copy(output, input)
		XOR(output, key, false)
		XOR(output, []byte(defaultKey), true)
		if !bytes.Equal(input, output) {
			t.Error("XOR input != ouptut")
			break
		}
	}
	fmt.Printf("TestXor passed\n")
}
