// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
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
		return 0, errors.New("Wrong length")
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

// RandomString return random string includes upper and lowwer chars
func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// IsLegalChannelName return if a channel id is a legal channel name
func IsLegalChannelName(channelID string) bool {
	if m, err := regexp.MatchString("^[a-z0-9]{1,32}$", channelID); err != nil || !m {
		return false
	}
	return true
}

// GetAllFiles return all files of a dir
func GetAllFiles(dirPth string, abs bool) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth+PthSep+fi.Name(), abs)
		} else {
			if abs {
				files = append(files, dirPth+PthSep+fi.Name())
			} else {
				files = append(files, PthSep+fi.Name())
			}

		}
	}

	for _, table := range dirs {
		temp, _ := GetAllFiles(table, abs)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

// IsDirSame return is two dir contain same files
func IsDirSame(a, b string) bool {
	aFiles, err := GetAllFiles(a, false)
	if err != nil {
		return false
	}
	bFiles, err := GetAllFiles(b, false)
	if err != nil {
		return false
	}
	if len(aFiles) != len(bFiles) {
		return false
	}
	sort.Strings(aFiles)
	sort.Strings(bFiles)
	for i := range aFiles {
		if aFiles[i] != bFiles[i] {
			fmt.Println(aFiles[i])
			return false
		}
	}
	return true
}

// CopyBytes copy bytes
func CopyBytes(origin []byte) []byte {
	if origin == nil {
		return nil
	}
	var res = make([]byte, len(origin))
	for i := range origin {
		res[i] = origin[i]
	}
	return res
}

// HexToBytes will remove 0x or 0X of begin ,and then call hex.DecodeString
func HexToBytes(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = strings.Replace(s, "0x", "", 1)
	} else if strings.HasPrefix(s, "0X") {
		s = strings.Replace(s, "0X", "", 1)
	}
	return hex.DecodeString(s)
}
