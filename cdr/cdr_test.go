package cdr

import (
	"encoding/binary"
	"testing"
)

func buildTestFileData(numCdrs uint32, cdrContent []byte) []byte {
	// Build a minimal TS 32.297 file with the given number of CDRs.
	// File header: fileLen(4) + headerLen(4) + highRel(1) + lowRel(1) +
	//   openTS(4) + lastTS(4) + numCdrs(4) + seqNum(4) + closureReason(1) +
	//   nodeIP(20) + lostCdr(1) + routeFilterLen(2) + privateExtLen(2) = 52 bytes minimum
	headerLen := uint32(52)

	// Each CDR: cdrLength(2) + version(1) + format+ts(1) + content
	cdrRecordLen := 4 + len(cdrContent)
	totalLen := int(headerLen) + int(numCdrs)*cdrRecordLen

	data := make([]byte, totalLen)
	off := 0

	// File length
	binary.BigEndian.PutUint32(data[off:], uint32(totalLen))
	off += 4
	// Header length
	binary.BigEndian.PutUint32(data[off:], headerLen)
	off += 4
	// High release version: release=5, version=3 => (5<<5)|3 = 0xA3
	data[off] = 0xA3
	off++
	// Low release version: release=4, version=1 => (4<<5)|1 = 0x81
	data[off] = 0x81
	off++
	// File opening timestamp: month=3, day=15, hour=14, min=30, sign=+, hourDev=2, minDev=0
	// bits: 0011 10111 01110 011110 1 00010 000000
	//       month=3(4) day=15(5) hour=14(5) min=30(6) sign=1(1) hDev=2(5) mDev=0(6)
	ts := uint32(3)<<28 | uint32(15)<<23 | uint32(14)<<18 | uint32(30)<<12 | uint32(1)<<11 | uint32(2)<<6 | uint32(0)
	binary.BigEndian.PutUint32(data[off:], ts)
	off += 4
	// Last CDR append timestamp (same for simplicity)
	binary.BigEndian.PutUint32(data[off:], ts)
	off += 4
	// Number of CDRs
	binary.BigEndian.PutUint32(data[off:], numCdrs)
	off += 4
	// File sequence number
	binary.BigEndian.PutUint32(data[off:], 42)
	off += 4
	// File closure trigger reason
	data[off] = 1
	off++
	// Node IP address (20 bytes, put 10.20.30.40 at start)
	data[off] = 10
	data[off+1] = 20
	data[off+2] = 30
	data[off+3] = 40
	off += 20
	// Lost CDR indicator
	data[off] = 0
	off++
	// CDR routeing filter length = 0
	binary.BigEndian.PutUint16(data[off:], 0)
	off += 2
	// Private extension length = 0
	binary.BigEndian.PutUint16(data[off:], 0)
	off += 2

	// CDR records
	for i := uint32(0); i < numCdrs; i++ {
		// CDR length (length of content only, per TS 32.297)
		binary.BigEndian.PutUint16(data[off:], uint16(len(cdrContent)))
		off += 2
		// Version: release=5, version=2 => (5<<5)|2 = 0xA2
		data[off] = 0xA2
		off++
		// Data record format (BER=1) in bits[7:5], TS number (7=TS32.251) in bits[4:0]
		data[off] = (1 << 5) | 7
		off++
		// CDR content
		copy(data[off:], cdrContent)
		off += len(cdrContent)
	}

	return data
}

func TestParseFileHeader(t *testing.T) {
	cdrContent := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	data := buildTestFileData(2, cdrContent)

	file, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	h := file.Header
	if h.FileLength != uint32(len(data)) {
		t.Errorf("FileLength = %d, want %d", h.FileLength, len(data))
	}
	if h.HeaderLength != 52 {
		t.Errorf("HeaderLength = %d, want 52", h.HeaderLength)
	}
	if h.HighReleaseVersionIdentifier.ReleaseIdentifier != 5 {
		t.Errorf("HighRelease = %d, want 5", h.HighReleaseVersionIdentifier.ReleaseIdentifier)
	}
	if h.HighReleaseVersionIdentifier.VersionIdentifier != 3 {
		t.Errorf("HighVersion = %d, want 3", h.HighReleaseVersionIdentifier.VersionIdentifier)
	}
	if h.LowReleaseVersionIdentifier.ReleaseIdentifier != 4 {
		t.Errorf("LowRelease = %d, want 4", h.LowReleaseVersionIdentifier.ReleaseIdentifier)
	}
	if h.NumberOfCdrsInFile != 2 {
		t.Errorf("NumberOfCdrs = %d, want 2", h.NumberOfCdrsInFile)
	}
	if h.FileSequenceNumber != 42 {
		t.Errorf("FileSequenceNumber = %d, want 42", h.FileSequenceNumber)
	}
	if h.FileClosureTriggerReason != 1 {
		t.Errorf("FileClosureTriggerReason = %d, want 1", h.FileClosureTriggerReason)
	}
	if h.LostCdrIndicator != 0 {
		t.Errorf("LostCdrIndicator = %d, want 0", h.LostCdrIndicator)
	}
}

func TestParseTimestamp(t *testing.T) {
	// month=3, day=15, hour=14, min=30, sign=+, hDev=2, mDev=0
	v := uint32(3)<<28 | uint32(15)<<23 | uint32(14)<<18 | uint32(30)<<12 | uint32(1)<<11 | uint32(2)<<6 | uint32(0)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)

	ts := parseTimestamp(buf)
	if ts.Month != 3 {
		t.Errorf("Month = %d, want 3", ts.Month)
	}
	if ts.Day != 15 {
		t.Errorf("Day = %d, want 15", ts.Day)
	}
	if ts.Hour != 14 {
		t.Errorf("Hour = %d, want 14", ts.Hour)
	}
	if ts.Minute != 30 {
		t.Errorf("Minute = %d, want 30", ts.Minute)
	}
	if !ts.Sign {
		t.Error("Sign should be true (positive)")
	}
	if ts.HourDeviation != 2 {
		t.Errorf("HourDeviation = %d, want 2", ts.HourDeviation)
	}
	if ts.MinuteDeviation != 0 {
		t.Errorf("MinuteDeviation = %d, want 0", ts.MinuteDeviation)
	}
}

func TestParseCdrRecords(t *testing.T) {
	cdrContent := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	data := buildTestFileData(3, cdrContent)

	file, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(file.Cdrs) != 3 {
		t.Fatalf("got %d CDRs, want 3", len(file.Cdrs))
	}

	for i, cdr := range file.Cdrs {
		if cdr.CdrLength != 4 {
			t.Errorf("CDR[%d].CdrLength = %d, want 4", i, cdr.CdrLength)
		}
		if cdr.Version.ReleaseIdentifier != 5 {
			t.Errorf("CDR[%d].Version.Release = %d, want 5", i, cdr.Version.ReleaseIdentifier)
		}
		if cdr.DataRecordFormat != 1 {
			t.Errorf("CDR[%d].DataRecordFormat = %d, want 1 (BER)", i, cdr.DataRecordFormat)
		}
		if cdr.TsNumber != 7 {
			t.Errorf("CDR[%d].TsNumber = %d, want 7", i, cdr.TsNumber)
		}
		if len(cdr.CdrContent) != 4 || cdr.CdrContent[0] != 0xDE {
			t.Errorf("CDR[%d].CdrContent mismatch", i)
		}
	}
}

func TestParseReleaseExtension(t *testing.T) {
	// Build a file with release identifier = 7 (beyond R9) to trigger extension parsing
	headerLen := uint32(54) // 52 + 2 extension bytes
	cdrContent := []byte{0x01, 0x02}
	// CDR with release=7 has extra extension byte: cdrLen(2) + ver(1) + fmt(1) + ext(1) + content
	cdrRecordLen := 5 + len(cdrContent)
	totalLen := int(headerLen) + cdrRecordLen

	data := make([]byte, totalLen)
	off := 0
	binary.BigEndian.PutUint32(data[off:], uint32(totalLen))
	off += 4
	binary.BigEndian.PutUint32(data[off:], headerLen)
	off += 4
	// High release = 7 (beyond9), version = 2 => (7<<5)|2 = 0xE2
	data[off] = 0xE2
	off++
	// Low release = 7, version = 1 => (7<<5)|1 = 0xE1
	data[off] = 0xE1
	off++
	// Timestamps (zeros for simplicity)
	off += 8
	// numCdrs = 1
	binary.BigEndian.PutUint32(data[off:], 1)
	off += 4
	// seqNum
	binary.BigEndian.PutUint32(data[off:], 1)
	off += 4
	// closure reason
	data[off] = 0
	off++
	// node IP
	off += 20
	// lost CDR
	off++
	// route filter len = 0
	off += 2
	// private ext len = 0
	off += 2
	// High release extension = 5 (meaning R15)
	data[off] = 5
	off++
	// Low release extension = 3 (meaning R13)
	data[off] = 3
	off++

	// CDR record with release=7
	binary.BigEndian.PutUint16(data[off:], uint16(len(cdrContent)))
	off += 2
	data[off] = 0xE2 // release=7, version=2
	off++
	data[off] = (1 << 5) | 7 // BER, TS32.251
	off++
	data[off] = 5 // extension = R15
	off++
	copy(data[off:], cdrContent)

	file, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if file.Header.HighReleaseIdentifierExtension == nil {
		t.Fatal("expected HighReleaseIdentifierExtension")
	}
	if *file.Header.HighReleaseIdentifierExtension != 5 {
		t.Errorf("HighReleaseIdentifierExtension = %d, want 5", *file.Header.HighReleaseIdentifierExtension)
	}
	if *file.Header.LowReleaseIdentifierExtension != 3 {
		t.Errorf("LowReleaseIdentifierExtension = %d, want 3", *file.Header.LowReleaseIdentifierExtension)
	}

	if file.Cdrs[0].ReleaseIdentifierExtension == nil {
		t.Fatal("expected CDR ReleaseIdentifierExtension")
	}
	if *file.Cdrs[0].ReleaseIdentifierExtension != 5 {
		t.Errorf("CDR ReleaseIdentifierExtension = %d, want 5", *file.Cdrs[0].ReleaseIdentifierExtension)
	}

	ver := formatVersion(file.Header.HighReleaseVersionIdentifier, file.Header.HighReleaseIdentifierExtension)
	if ver != "15.2" {
		t.Errorf("formatVersion() = %q, want \"15.2\"", ver)
	}
}

func TestParseTooShort(t *testing.T) {
	_, err := Parse([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error for short data")
	}
}

func TestParseTruncatedCdr(t *testing.T) {
	data := buildTestFileData(1, []byte{0x01, 0x02, 0x03})
	// Truncate the data mid-CDR
	truncated := data[:len(data)-2]
	_, err := Parse(truncated)
	if err == nil {
		t.Error("expected error for truncated CDR content")
	}
}

func TestParseTruncatedHeader(t *testing.T) {
	data := buildTestFileData(1, []byte{0x01})
	// Truncate mid-header (keep only first 40 bytes of a 52-byte header)
	_, err := Parse(data[:40])
	if err == nil {
		t.Error("expected error for truncated header")
	}
}

func TestParseZeroCdrs(t *testing.T) {
	data := buildTestFileData(0, nil)
	file, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(file.Cdrs) != 0 {
		t.Errorf("expected 0 CDRs, got %d", len(file.Cdrs))
	}
}

func TestParseNodeIPAddress(t *testing.T) {
	cdrContent := []byte{0x01}
	data := buildTestFileData(1, cdrContent)

	file, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	ip := byteArrayToIPv4(file.Header.NodeIpAddress[:])
	if ip != "10.20.30.40" {
		t.Errorf("NodeIPAddress = %q, want \"10.20.30.40\"", ip)
	}
}

func TestEndToEndFileInfo(t *testing.T) {
	cdrContent := []byte{0xAA, 0xBB}
	data := buildTestFileData(2, cdrContent)

	info := ToFileInfo(data)
	if info.HeaderInfo.NumberOfCDRsInFile != 2 {
		t.Errorf("NumberOfCDRsInFile = %d, want 2", info.HeaderInfo.NumberOfCDRsInFile)
	}
	if info.CdrInfo.NumberOfCDRs != 2 {
		t.Errorf("NumberOfCDRs = %d, want 2", info.CdrInfo.NumberOfCDRs)
	}
	if info.HeaderInfo.FileSequenceNumber != 42 {
		t.Errorf("FileSequenceNumber = %d, want 42", info.HeaderInfo.FileSequenceNumber)
	}
	if info.CdrInfo.CdrHeaders[0].DataRecorderFormat != "BER" {
		t.Errorf("DataRecorderFormat = %q, want BER", info.CdrInfo.CdrHeaders[0].DataRecorderFormat)
	}
	if info.CdrInfo.CdrHeaders[0].TsNumber != "TS32.251" {
		t.Errorf("TsNumber = %q, want TS32.251", info.CdrInfo.CdrHeaders[0].TsNumber)
	}
}
