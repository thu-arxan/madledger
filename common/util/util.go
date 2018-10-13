package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

const (
	// MaxUint64 is the max of uint64, used for check overflow
	MaxUint64 = 1<<64 - 1
)

// Hex is the wrapper of fmt.Sprintf("%x", data)
func Hex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

// Now returns the time now
func Now() int64 {
	return time.Now().Unix()
}

// Contain return if the target which is a map or slice contains the obj
func Contain(target interface{}, obj interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// MakeFileAbs makes 'file' absolute relative to 'dir' if not already absolute
func MakeFileAbs(file, dir string) (string, error) {
	if file == "" {
		return "", nil
	}
	if filepath.IsAbs(file) {
		return file, nil
	}
	path, err := filepath.Abs(filepath.Join(dir, file))
	if err != nil {
		return "", fmt.Errorf("Failed making '%s' absolute based on '%s'", file, dir)
	}
	return path, nil
}

// FileExists checks to see if a file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Int64ToBytes turn int64 to []byte
func Int64ToBytes(i int64) []byte {
	return Uint64ToBytes(uint64(i))
}

// Uint64ToBytes turn int64 to []byte
func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// Int32ToBytes turn int64 to []byte
func Int32ToBytes(i int32) []byte {
	return Uint32ToBytes(uint32(i))
}

// Uint32ToBytes turn int64 to []byte
func Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

// BytesToUint32 turn b[0:4] to uint32, if less than 4, panic error
func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// BytesToUint64 turn b[0:4] to uint64
func BytesToUint64(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, errors.New("Wrong lengtg")
	}
	return binary.BigEndian.Uint64(b), nil
}

// BoolToBytes turn bool to []byte
func BoolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// BytesCombine combines some bytes array
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// RandNum return int in [0, num)
func RandNum(num int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(num)
	return randNum
}

// RandUint64 return uint64
func RandUint64() uint64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint64()
}
