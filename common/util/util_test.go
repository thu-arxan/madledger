package util

import (
	"bytes"
	"testing"
)

func TestBytesCombine(t *testing.T) {
	if !bytes.Equal(BytesCombine([]byte("Hello"), []byte(" World")), []byte("Hello World")) {
		t.Fatal()
	}
	if !bytes.Equal(BytesCombine([]byte("Hello"), []byte(" "), []byte("World")), []byte("Hello World")) {
		t.Fatal()
	}
	if !bytes.Equal(BytesCombine([]byte("Hello")), []byte("Hello")) {
		t.Fatal()
	}
}
