package rtmp

import (
	"bytes"
	"testing"
)

func TestAMF3EncodeFalse(t *testing.T) {
	testAMF3Encode(t, false, 0x02, nil)
}

func TestAMF3EncodeTrue(t *testing.T) {
	testAMF3Encode(t, true, 0x03, nil)
}

func TestAMF3EncodeInteger(t *testing.T) {
	for _, d := range testU29Data {
		testAMF3Encode(t, d.data, 0x04, d.expected)
	}
}

func TestAMF3EncodeIntegerOverRange(t *testing.T) {
	e := NewAMF3Encoder(new(bytes.Buffer))
	_, err := e.Encode(uint32(0x40000000))
	if err != ErrAMF3U29OverRange {
		t.Fatalf("Expecting %v, got %v", ErrAMF3U29OverRange, err)
	}
}

func testAMF3Encode(t *testing.T, data interface{}, marker byte, expected []byte) {
	w := new(bytes.Buffer)
	e := NewAMF3Encoder(w)
	l, err := e.Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	expected = append([]byte{marker}, expected...)
	if l != len(expected) {
		t.Fatalf("Expecting result length %d, got %d", len(expected), l)
	}
	if !bytes.Equal(w.Bytes(), expected) {
		t.Fatalf("Expecting buf %v, got %v", expected, w.Bytes())
	}
}

var testU29Data = []struct {
	data     uint32
	expected []byte
}{
	{0x00, []byte{0x00}},
	{0x01, []byte{0x01}},
	{0x7F, []byte{0x7F}},
	{0x80, []byte{0x81, 0x00}},
	{0x3FFF, []byte{0xFF, 0x7F}},
	{0x4000, []byte{0x81, 0x80, 0x00}},
	{0x1FFFFF, []byte{0xFF, 0xFF, 0x7F}},
	{0x200000, []byte{0x80, 0xC0, 0x80, 0x00}},
	{0x1FFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestAMF3DecodeFalse(t *testing.T) {
	data := testAMF3Decode(t, 0x02, nil)
	if data.(bool) {
		t.Fatal("Expecting false, got true")
	}
}

func TestAMF3DecodeTrue(t *testing.T) {
	data := testAMF3Decode(t, 0x03, nil)
	if !data.(bool) {
		t.Fatal("Expecting true, got false")
	}
}

func TestAMF3DecodeInteger(t *testing.T) {
	for _, d := range testU29Data {
		data := testAMF3Decode(t, 0x04, d.expected)
		if data.(uint32) != d.data {
			t.Fatalf("Expecting %x, got %x", d.data, data.(uint32))
		}
	}
}

func testAMF3Decode(t *testing.T, marker byte, b []byte) interface{} {
	b = append([]byte{marker}, b...)
	r := bytes.NewBuffer(b)
	d := NewAMF3Decoder(r)
	data, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	return data
}
