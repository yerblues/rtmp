package rtmp

import (
	"bytes"
	"testing"
)

func TestAMF3EncodeBool(t *testing.T) {
	w := new(bytes.Buffer)
	e := NewAMF3Encoder(w)
	l, err := e.Encode(true)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x03}
	if l != len(expected) {
		t.Fatalf("Expecting result length %d, got %d", len(expected), l)
	}
	if !bytes.Equal(w.Bytes(), expected) {
		t.Fatalf("Expecting buf %v, got %v", expected, w.Bytes())
	}
}
