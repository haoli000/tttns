package cdr

import (
	"os"
	"strings"
	"testing"
)

func TestByteArrayToIPv4(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"simple", []byte{192, 168, 1, 1}, "192.168.1.1"},
		{"with leading 0xFF", []byte{0xFF, 0xFF, 10, 0, 0, 1}, "10.0.0.1"},
		{"all zeros", []byte{0, 0, 0, 0}, "0.0.0.0"},
		{"max values", []byte{255, 255, 255, 255}, "0.0.0.0"},
		{"short after strip", []byte{0xFF, 0xFF, 0xFF, 10, 1}, "10.1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := byteArrayToIPv4(tt.input)
			if result != tt.expected {
				t.Errorf("byteArrayToIPv4(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecodeLostCdrIndicator(t *testing.T) {
	tests := []struct {
		name     string
		input    uint8
		expected string
	}{
		{"no loss", 0x00, "No CDRs have been lost"},
		{"msb0 known count", 0x05, "CGF has identified that 5 CDR(s) were lost, while it is unknown whether more CDRs were lost"},
		{"msb0 max", 0x7F, "CGF has identified that 127 or more CDRs were lost, while it is unknown whether more CDRs were lost"},
		{"msb1 unknown count", 0x80, "CDRs have been lost but CGF cannot determine the number of lost CDRs"},
		{"msb1 known count", 0x83, "CGF has calculated the number of lost CDRs as 3"},
		{"msb1 max", 0xFF, "CGF has calculated the number of lost CDRs to be 127 or more"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decodeLostCdrIndicator(tt.input)
			if result != tt.expected {
				t.Errorf("decodeLostCdrIndicator(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToFileClosureTriggerReason(t *testing.T) {
	tests := []struct {
		input    uint8
		contains string
	}{
		{0, "Normal closure"},
		{1, "File size limit"},
		{2, "File open-time limit"},
		{3, "Maximum number of CDRs"},
		{4, "manual intervention"},
		{5, "encoding change"},
		{128, "Abnormal file closure"},
		{129, "File system error"},
		{130, "storage exhausted"},
		{131, "integrity error"},
		{99, "Reserved"},
	}
	for _, tt := range tests {
		result := toFileClosureTriggerReason(tt.input)
		if !strings.Contains(result, tt.contains) {
			t.Errorf("toFileClosureTriggerReason(%d) = %q, want it to contain %q", tt.input, result, tt.contains)
		}
	}
}

func TestToCdrEncoding(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{1, "BER"},
		{2, "Unaligned PER"},
		{3, "Aligned PER"},
		{4, "XML"},
		{99, "Unknown 99"},
	}
	for _, tt := range tests {
		result := toCdrEncoding(tt.input)
		if result != tt.expected {
			t.Errorf("toCdrEncoding(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToTsNumber(t *testing.T) {
	tests := []struct {
		input    ThreegppCdr_Cdr_Ts
		expected string
	}{
		{0, "TS32.005"},
		{1, "TS32.015"},
		{7, "TS32.251"},
		{23, "TS28.201"},
		{100, "Unknown 100"},
	}
	for _, tt := range tests {
		result := toTsNumber(tt.input)
		if result != tt.expected {
			t.Errorf("toTsNumber(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToVersion(t *testing.T) {
	tests := []struct {
		name     string
		rel      *ThreegppCdr_ReleaseVersionIdentifier
		ext      *ThreegppCdr_ReleaseIdentifierExtension
		expected string
	}{
		{
			"normal release",
			&ThreegppCdr_ReleaseVersionIdentifier{ReleaseIdentifier: 5, VersionIdentifier: 3, threegppRelease: 5},
			nil,
			"5.3",
		},
		{
			"beyond 9",
			&ThreegppCdr_ReleaseVersionIdentifier{ReleaseIdentifier: uint64(ThreegppCdr_ReleaseVersionIdentifier_Rel__Beyond9), VersionIdentifier: 2},
			&ThreegppCdr_ReleaseIdentifierExtension{ThreegppRelease: 5},
			"15.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toVersion(tt.rel, tt.ext)
			if result != tt.expected {
				t.Errorf("toVersion() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestToTimeStamp(t *testing.T) {
	tests := []struct {
		name     string
		ts       *ThreegppCdr_Timestamp
		expected string
	}{
		{
			"positive offset",
			&ThreegppCdr_Timestamp{Date: 15, Month: 3, Hour: 14, Minute: 30, Sign: true, HourDeviation: 2, MinuteDeviation: 0},
			"15/3 14:30:00+0200",
		},
		{
			"negative offset",
			&ThreegppCdr_Timestamp{Date: 1, Month: 12, Hour: 8, Minute: 5, Sign: false, HourDeviation: 5, MinuteDeviation: 30},
			"1/12 08:05:00-0530",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTimeStamp(tt.ts)
			if result != tt.expected {
				t.Errorf("toTimeStamp() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPrintOutput(t *testing.T) {
	data := struct {
		Name string `json:"name"`
	}{Name: "test"}
	// Just verify it doesn't panic
	// (actual output goes to stdout, which is fine for a smoke test)
	PrintOutput(true, data)
	PrintOutput(false, data)
}

func TestGetContent(t *testing.T) {
	expected := []byte{0x01, 0x02, 0x03}
	f, err := os.CreateTemp("", "tttns-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.Write(expected)
	f.Close()

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
	// When running under go test, stdout is not a terminal
	if isOutputToTerminal() {
		t.Error("isOutputToTerminal() should be false during tests")
	}
}
