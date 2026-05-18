package cdr

import (
	"bytes"
	"os"
	"testing"
)

func TestToFileHeaderInfo(t *testing.T) {
	data := buildTestFileData(1, []byte{0xAA})
	info := ToFileHeaderInfo(data)
	if info.NumberOfCDRsInFile != 1 {
		t.Errorf("NumberOfCDRsInFile = %d, want 1", info.NumberOfCDRsInFile)
	}
	if info.FileSequenceNumber != 42 {
		t.Errorf("FileSequenceNumber = %d, want 42", info.FileSequenceNumber)
	}
}

func TestToCdrHeaderInfo(t *testing.T) {
	data := buildTestFileData(2, []byte{0xAA, 0xBB})
	info := ToCdrHeaderInfo(data, 1)
	if info.CdrLength != 2 {
		t.Errorf("CdrLength = %d, want 2", info.CdrLength)
	}
	if info.DataRecorderFormat != "BER" {
		t.Errorf("DataRecorderFormat = %q, want BER", info.DataRecorderFormat)
	}
	info2 := ToCdrHeaderInfo(data, 2)
	if info2.TsNumber != "TS32.251" {
		t.Errorf("TsNumber = %q, want TS32.251", info2.TsNumber)
	}
}

func TestCountCdrs(t *testing.T) {
	data := buildTestFileData(5, []byte{0x01})
	if cnt := CountCdrs(data); cnt != 5 {
		t.Errorf("CountCdrs() = %d, want 5", cnt)
	}
}

func TestDumpCdr(t *testing.T) {
	content := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	data := buildTestFileData(1, content)

	f, err := os.CreateTemp("", "tttns-dump-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	DumpCdr(data, 1, f)
	_ = f.Close()

	got, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("DumpCdr wrote %x, want %x", got, content)
	}
}

func TestToCdrInfo(t *testing.T) {
	data := buildTestFileData(3, []byte{0x01, 0x02})
	info := ToCdrInfo(data)
	if info.NumberOfCDRs != 3 {
		t.Errorf("NumberOfCDRs = %d, want 3", info.NumberOfCDRs)
	}
	if len(info.CdrHeaders) != 3 {
		t.Errorf("len(CdrHeaders) = %d, want 3", len(info.CdrHeaders))
	}
}

func TestPrintOutput(t *testing.T) {
	data := struct {
		Name string `json:"name"`
	}{Name: "test"}
	PrintOutput(true, data)
	PrintOutput(false, data)
}

func TestGetContent(t *testing.T) {
	expected := []byte{0x01, 0x02, 0x03}
	f, err := os.CreateTemp("", "tttns-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()
	_, _ = f.Write(expected)
	_ = f.Close()

	result := GetContent(f.Name())
	if len(result) != len(expected) {
		t.Errorf("GetContent() returned %d bytes, want %d", len(result), len(expected))
	}
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("GetContent()[%d] = %d, want %d", i, b, expected[i])
		}
	}
}

func TestIsOutputToTerminal(t *testing.T) {
	if isOutputToTerminal() {
		t.Error("isOutputToTerminal() should be false during tests")
	}
}
