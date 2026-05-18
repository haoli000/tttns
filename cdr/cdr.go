/*
Copyright © 2024 Hao Li <mr.hao.li@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cdr implements a parser for 3GPP TS 32.297 CDR file format.
package cdr

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TS 32.297 File Header structure (fixed portion is 68 bytes minimum)
type FileHeader struct {
	FileLength                     uint32
	HeaderLength                   uint32
	HighReleaseVersionIdentifier   ReleaseVersion
	LowReleaseVersionIdentifier    ReleaseVersion
	FileOpeningTimestamp            Timestamp
	LastCdrAppendTimestamp          Timestamp
	NumberOfCdrsInFile             uint32
	FileSequenceNumber             uint32
	FileClosureTriggerReason       uint8
	NodeIpAddress                  [20]byte
	LostCdrIndicator               uint8
	CdrRouteingFilter              []byte
	PrivateExtension               []byte
	HighReleaseIdentifierExtension *uint8
	LowReleaseIdentifierExtension  *uint8
}

// ReleaseVersion is a single byte: bits[7:5] = release identifier (3 bits), bits[4:0] = version (5 bits)
type ReleaseVersion struct {
	ReleaseIdentifier uint8
	VersionIdentifier uint8
}

// Timestamp is 4 bytes packed per TS 32.297:
// bits[31:28] = month (4), bits[27:23] = day (5), bits[22:18] = hour (5),
// bits[17:12] = minute (6), bit[11] = sign (1=positive), bits[10:6] = hour deviation (5), bits[5:0] = minute deviation (6)
type Timestamp struct {
	Month           uint8
	Day             uint8
	Hour            uint8
	Minute          uint8
	Sign            bool
	HourDeviation   uint8
	MinuteDeviation uint8
}

// CdrRecord represents a single CDR within the file
type CdrRecord struct {
	CdrLength                  uint16
	Version                    ReleaseVersion
	DataRecordFormat           uint8
	TsNumber                   uint8
	ReleaseIdentifierExtension *uint8
	CdrContent                 []byte
}

// ThreegppCdrFile represents a fully parsed TS 32.297 CDR file
type ThreegppCdrFile struct {
	Header FileHeader
	Cdrs   []CdrRecord
}

func parseReleaseVersion(b byte) ReleaseVersion {
	return ReleaseVersion{
		ReleaseIdentifier: (b >> 5) & 0x07,
		VersionIdentifier: b & 0x1F,
	}
}

func parseTimestamp(data []byte) Timestamp {
	// 4 bytes, big-endian bit packing
	v := binary.BigEndian.Uint32(data)
	return Timestamp{
		Month:           uint8((v >> 28) & 0x0F),
		Day:             uint8((v >> 23) & 0x1F),
		Hour:            uint8((v >> 18) & 0x1F),
		Minute:          uint8((v >> 12) & 0x3F),
		Sign:            ((v >> 11) & 0x01) == 1,
		HourDeviation:   uint8((v >> 6) & 0x1F),
		MinuteDeviation: uint8(v & 0x3F),
	}
}

// Parse reads a TS 32.297 CDR file from a byte slice.
func Parse(data []byte) (*ThreegppCdrFile, error) {
	if len(data) < 52 {
		return nil, fmt.Errorf("data too short for file header: %d bytes", len(data))
	}

	r := newReader(data)
	var h FileHeader

	h.FileLength = r.readU32()
	h.HeaderLength = r.readU32()
	h.HighReleaseVersionIdentifier = parseReleaseVersion(r.readU8())
	h.LowReleaseVersionIdentifier = parseReleaseVersion(r.readU8())
	h.FileOpeningTimestamp = parseTimestamp(r.readBytes(4))
	h.LastCdrAppendTimestamp = parseTimestamp(r.readBytes(4))
	h.NumberOfCdrsInFile = r.readU32()
	h.FileSequenceNumber = r.readU32()
	h.FileClosureTriggerReason = r.readU8()
	copy(h.NodeIpAddress[:], r.readBytes(20))
	h.LostCdrIndicator = r.readU8()

	routeFilterLen := r.readU16()
	h.CdrRouteingFilter = r.readBytes(int(routeFilterLen))

	privateExtLen := r.readU16()
	h.PrivateExtension = r.readBytes(int(privateExtLen))

	if h.HighReleaseVersionIdentifier.ReleaseIdentifier == 7 {
		v := r.readU8()
		h.HighReleaseIdentifierExtension = &v
		v2 := r.readU8()
		h.LowReleaseIdentifierExtension = &v2
	}

	if r.err != nil {
		return nil, fmt.Errorf("error parsing file header: %w", r.err)
	}

	file := &ThreegppCdrFile{Header: h}

	for r.remaining() > 0 {
		cdr, err := parseCdrRecord(r)
		if err != nil {
			return nil, fmt.Errorf("error parsing CDR record %d: %w", len(file.Cdrs)+1, err)
		}
		file.Cdrs = append(file.Cdrs, cdr)
	}

	return file, nil
}

func parseCdrRecord(r *reader) (CdrRecord, error) {
	var rec CdrRecord

	rec.CdrLength = r.readU16()

	// CDR header byte 1: release version (1 byte)
	rec.Version = parseReleaseVersion(r.readU8())

	// CDR header byte 2: bits[7:5] = data record format (3 bits), bits[4:0] = TS number (5 bits)
	b := r.readU8()
	rec.DataRecordFormat = (b >> 5) & 0x07
	rec.TsNumber = b & 0x1F

	if rec.Version.ReleaseIdentifier == 7 {
		v := r.readU8()
		rec.ReleaseIdentifierExtension = &v
	}

	rec.CdrContent = r.readBytes(int(rec.CdrLength))

	if r.err != nil {
		return rec, r.err
	}
	return rec, nil
}

// reader is a simple big-endian binary reader
type reader struct {
	data []byte
	pos  int
	err  error
}

func newReader(data []byte) *reader {
	return &reader{data: data}
}

func (r *reader) remaining() int {
	return len(r.data) - r.pos
}

func (r *reader) readU8() uint8 {
	if r.err != nil || r.pos+1 > len(r.data) {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := r.data[r.pos]
	r.pos++
	return v
}

func (r *reader) readU16() uint16 {
	if r.err != nil || r.pos+2 > len(r.data) {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := binary.BigEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return v
}

func (r *reader) readU32() uint32 {
	if r.err != nil || r.pos+4 > len(r.data) {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := binary.BigEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return v
}

func (r *reader) readBytes(n int) []byte {
	if r.err != nil || r.pos+n > len(r.data) {
		r.err = io.ErrUnexpectedEOF
		return make([]byte, n)
	}
	v := make([]byte, n)
	copy(v, r.data[r.pos:r.pos+n])
	r.pos += n
	return v
}
