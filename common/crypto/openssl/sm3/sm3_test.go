package sm3

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/tjfoc/gmsm/sm3"
)

func TestSm3Nil(t *testing.T) {
	h1 := hex.EncodeToString(Sm3Sum(nil))
	h2 := hex.EncodeToString(sm3.Sm3Sum(nil))
	if h1 != h2 {
		t.Error("Not equal", h1, h2)
	}
}

func TestSm3(t *testing.T) {
	for i := 0; i < 10000; i++ {
		data := randomByte()
		d1 := hex.EncodeToString(Sm3Sum(data))
		d2 := hex.EncodeToString(sm3.Sm3Sum(data))
		if d1 != d2 {
			t.Error(data)
		}
	}
}

func randomByte() []byte {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, rand.Intn(30)+1)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return []byte(string(b))
}

func BenchmarkSm3(b *testing.B) {
	data := randomByte()

	for i := 0; i < b.N; i++ {
		Sm3Sum(data)
	}
}

func BenchmarkTjSm3(b *testing.B) {
	data := randomByte()
	for i := 0; i < b.N; i++ {
		sm3.Sm3Sum(data)
	}
}
