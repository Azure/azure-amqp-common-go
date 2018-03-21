package uuid

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"crypto/rand"
	"encoding/hex"
)

// Size of a UUID in bytes.
const Size = 16

// UUID versions
const (
	_ byte = iota
	_
	_
	_
	V4
	_

	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

var (
	randomReader = rand.Reader

	// Nil is special form of UUID that is specified to have all
	// 128 bits set to zero.
	Nil = UUID{}
)

type (
	// UUID representation compliant with specification
	// described in RFC 4122.
	UUID [Size]byte
)

// NewV4 returns random generated UUID.
func NewV4() (UUID, error) {
	u := UUID{}
	if _, err := randomReader.Read(u[:]); err != nil {
		return Nil, err
	}
	u.setVersion(V4)
	u.setVariant(VariantRFC4122)

	return u, nil
}

func (u *UUID) setVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

func (u *UUID) setVariant(v byte) {
	switch v {
	case VariantNCS:
		u[8] = u[8]&(0xff>>1) | (0x00 << 7)
	case VariantRFC4122:
		u[8] = u[8]&(0xff>>2) | (0x02 << 6)
	case VariantMicrosoft:
		u[8] = u[8]&(0xff>>3) | (0x06 << 5)
	case VariantFuture:
		fallthrough
	default:
		u[8] = u[8]&(0xff>>3) | (0x07 << 5)
	}
}

func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}
