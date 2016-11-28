package rtmp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAMF3EncodeFalse(t *testing.T) {
	testAMF3Encode(t, false, []byte{0x02})
}

func TestAMF3EncodeTrue(t *testing.T) {
	testAMF3Encode(t, true, []byte{0x03})
}

func TestAMF3EncodeInteger(t *testing.T) {
	expected := []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF}
	testAMF3Encode(t, uint32(0x3FFFFFFF), expected)
}

func testAMF3Encode(t *testing.T, data interface{}, expected []byte) {
	w := new(bytes.Buffer)
	e := NewAMF3Encoder(w)
	l, err := e.Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result:%v\n", w.Bytes())
	if l != len(expected) {
		t.Fatalf("Expecting result length %d, got %d", len(expected), l)
	}
	if !bytes.Equal(w.Bytes(), expected) {
		t.Fatalf("Expecting buf %v, got %v", expected, w.Bytes())
	}
}

func TestAMF3DecodeInteger(t *testing.T) {
	b := []byte{0x04, 0xFF, 0xFF, 0xFF, 0xFF}
	r := bytes.NewBuffer(b)
	d := NewAMF3Decoder(r)
	data, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("result:%x\n", data.(uint32))
}
