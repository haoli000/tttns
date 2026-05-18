package cdr

import (
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
		input    uint8
		expected string
	}{
		{0, "TS32.005"},
		{1, "TS32.015"},
		{2, "TS32.205"},
		{3, "TS32.215"},
		{4, "TS32.225"},
		{5, "TS32.235"},
		{6, "TS32.250"},
		{7, "TS32.251"},
		{9, "TS32.260"},
		{10, "TS32.270"},
		{11, "TS32.271"},
		{12, "TS32.272"},
		{13, "TS32.273"},
		{14, "TS32.275"},
		{15, "TS32.274"},
		{16, "TS32.277"},
		{17, "TS32.296"},
		{18, "TS32.278"},
		{19, "TS32.253"},
		{20, "TS32.255"},
		{21, "TS32.254"},
		{22, "TS32.256"},
		{23, "TS28.201"},
		{24, "TS28.202"},
		{25, "TS32.257"},
		{26, "TS32.282"},
		{8, "Unknown 8"},
		{100, "Unknown 100"},
	}
	for _, tt := range tests {
		result := toTsNumber(tt.input)
		if result != tt.expected {
			t.Errorf("toTsNumber(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name     string
		rel      ReleaseVersion
		ext      *uint8
		expected string
	}{
		{
			"normal release",
			ReleaseVersion{ReleaseIdentifier: 5, VersionIdentifier: 3},
			nil,
			"5.3",
		},
		{
			"beyond 9",
			ReleaseVersion{ReleaseIdentifier: 7, VersionIdentifier: 2},
			uint8Ptr(5),
			"15.2",
		},
		{
			"release 7 without extension",
			ReleaseVersion{ReleaseIdentifier: 7, VersionIdentifier: 1},
			nil,
			"7.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatVersion(tt.rel, tt.ext)
			if result != tt.expected {
				t.Errorf("formatVersion() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func uint8Ptr(v uint8) *uint8 { return &v }

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		ts       Timestamp
		expected string
	}{
		{
			"positive offset",
			Timestamp{Day: 15, Month: 3, Hour: 14, Minute: 30, Sign: true, HourDeviation: 2, MinuteDeviation: 0},
			"15/3 14:30:00+0200",
		},
		{
			"negative offset",
			Timestamp{Day: 1, Month: 12, Hour: 8, Minute: 5, Sign: false, HourDeviation: 5, MinuteDeviation: 30},
			"1/12 08:05:00-0530",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestamp(tt.ts)
			if result != tt.expected {
				t.Errorf("formatTimestamp() = %q, want %q", result, tt.expected)
			}
		})
	}
}
