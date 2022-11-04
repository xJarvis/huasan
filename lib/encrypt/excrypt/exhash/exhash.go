// Copyright 2019 gocrypt Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exhash

import (
	"crypto/hmac"
	"github.com/xjarvis/huashan/lib/encrypt/excrypt"
)

type hash struct {
	hashType excrypt.Hash
}

func NewHash(hashType excrypt.Hash) *hash {
	return &hash{hashType}
}

//Get gets hashed bytes with defined hashType
func (h *hash) Get(src []byte) (dst []byte, err error) {
	_, dst, err = excrypt.GetHash(src, h.hashType)
	return
}

//EncodeToString gets hashed bytes with defined hashType and then encode to string
func (h *hash) EncodeToString(src []byte, encodeType excrypt.Encode) (dst string, err error) {
	data, err := GetHash(src, h.hashType)
	return excrypt.EncodeToString(data, encodeType)
}

type hmacHash struct {
	hashType excrypt.Hash
	key      []byte
}

func NewHMAC(hashType excrypt.Hash, key []byte) *hmacHash {
	return &hmacHash{hashType, key}
}

//Get gets hmac hashed bytes with defined hashType & key
func (hh *hmacHash) Get(src []byte) (dst []byte, err error) {
	h, _ := excrypt.GetHashFunc(hh.hashType)
	hm := hmac.New(h, hh.key)
	hm.Write(src)
	dst = hm.Sum(nil)
	return
}

//EncodeToString gets hmac hashed bytes with defined hashType & key then encode to string
func (hh *hmacHash) EncodeToString(src []byte, encodeType excrypt.Encode) (dst string, err error) {
	data, err := GetHMACHash(src, hh.hashType, hh.key)
	return excrypt.EncodeToString(data, encodeType)
}

//GetHash gets hashed bytes with defined hashType
func GetHash(src []byte, hashType excrypt.Hash) (dst []byte, err error) {
	_, dst, err = excrypt.GetHash(src, hashType)
	return
}

//GetHashEncodeToString gets hashed bytes with defined hashType and then encode to string
func GetHashEncodeToString(encodeType excrypt.Encode, src []byte, hashType excrypt.Hash) (dst string, err error) {
	data, err := GetHash(src, hashType)
	return excrypt.EncodeToString(data, encodeType)
}

//GetHMACHash gets hmac hashed bytes with defined hashType & key
func GetHMACHash(src []byte, hashType excrypt.Hash, key []byte) (dst []byte, err error) {
	h, _ := excrypt.GetHashFunc(hashType)
	hm := hmac.New(h, key)
	hm.Write(src)
	dst = hm.Sum(nil)
	return
}

//GetHMACHashEncodeToString gets hmac hashed bytes with defined hashType & key then encode to string
func GetHMACHashEncodeToString(encodeType excrypt.Encode, src []byte, hashType excrypt.Hash, key []byte) (dst string, err error) {
	data, err := GetHMACHash(src, hashType, key)
	return excrypt.EncodeToString(data, encodeType)
}
