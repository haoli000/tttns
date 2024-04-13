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
package cdr

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

func GetContent(filename string) []byte {
	var inputReader io.Reader
	if filename == "-" {
		// Read from stdin
		inputReader = os.Stdin
	} else {
		// Read from file
		file, err := os.Open(filename)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		inputReader = file
	}
	content, err := io.ReadAll(inputReader)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return content
}

func ParseCdr(content []byte) *ThreegppCdr {
	file := NewThreegppCdr()
	err := file.Read(kaitai.NewStream(bytes.NewReader(content)), nil, file)
	if err != nil {
		fmt.Println("Failed to parse CDR file.")
		os.Exit(2)
	}
	return file
}

func ToFileHeaderInfo(content []byte) FileHeaderInfo {
	file := ParseCdr(content)
	return FileHeaderInfo{
		FileLength:               int(file.Header.FileLength),
		HeaderLength:             int(file.Header.HeaderLength),
		HighReleaseVersion:       toVersion(file.Header.HighReleaseVersionIdentifier, file.Header.HighReleaseIdentifierExtension),
		LowReleaseVersion:        toVersion(file.Header.LowReleaseVersionIdentifier, file.Header.LowReleaseIdentifierExtension),
		FileOpeningTimestamp:     toTimeStamp(file.Header.FileOpeningTimestamp),
		LastCDRAppendTimestamp:   toTimeStamp(file.Header.LastCdrAppendTimestamp),
		NumberOfCDRsInFile:       int(file.Header.NumberOfCdrsInFile),
		FileSequenceNumber:       int(file.Header.FileSequenceNumber),
		FileClosureTriggerReason: fmt.Sprintf("%d - %s", file.Header.FileClosureTriggerReason, toFileClosureTriggerReason(file.Header.FileClosureTriggerReason)),
		NodeIPAddress:            byteArrayToIPv4([]byte(file.Header.NodeIpAddress.IpAddress)),
		LostCDRIndicator:         decodeLostCdrIndicator(file.Header.LostCdrIndicator),
	}
}

func ToCdrHeaderInfo(content []byte, index uint32) CdrHeaderInfo {
	row := getCdrContent(content, index)
	return CdrHeaderInfo{
		CdrLength:          int(row.CdrLength),
		ReleaseVersion:     toVersion(row.Version, row.ReleaseIdentifierExtension),
		DataRecorderFormat: toCdrEncoding(row.DataRecordFormat),
		TsNumber:           toTsNumber(row.TsNumber),
	}
}

func CountCdrs(content []byte) uint32 {
	file := ParseCdr(content)
	return file.Header.NumberOfCdrsInFile
}

func DumpCdr(content []byte, index uint32, file *os.File) {
	row := getCdrContent(content, index)
	_, err := file.Write(row.CdrContent)
	if err != nil {
		fmt.Println("Error dumping CDR:", err)
		os.Exit(4)
	}
}

func ToCdrInfo(content []byte) CdrInfo {
	cnt := CountCdrs(content)
	var cdrHeaderInfos []CdrHeaderInfo
	for i := uint32(1); i <= cnt; i++ {
		info := ToCdrHeaderInfo(content, i)
		cdrHeaderInfos = append(cdrHeaderInfos, info)
	}
	return CdrInfo{
		NumberOfCDRs: int(cnt),
		CdrHeaders:   cdrHeaderInfos,
	}
}

func ToFileInfo(content []byte) FileInfo {
	return FileInfo{
		HeaderInfo: ToFileHeaderInfo(content),
		CdrInfo:    ToCdrInfo(content),
	}
}

func isOutputToTerminal() bool {
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
		//Display info to the terminal
		return true
	}
	return false
}

func PrettyPrintYAML(json []byte) {
	backend := logging.SetBackend(logging.NewLogBackend(os.Stderr, "", log.LstdFlags))
	backend.SetLevel(logging.CRITICAL, "")
	yqlib.GetLogger().SetBackend(backend)
	decoder := yqlib.NewJSONDecoder()
	decoder.Init(bytes.NewReader(json))
	node, err := decoder.Decode()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	prefs := yqlib.NewDefaultYamlPreferences()
	prefs.ColorsEnabled = isOutputToTerminal()
	enc := yqlib.NewYamlEncoder(prefs)
	err = enc.Encode(os.Stdout, node)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

func PrettyPrintJSON(json []byte) {
	backend := logging.SetBackend(logging.NewLogBackend(os.Stderr, "", log.LstdFlags))
	backend.SetLevel(logging.CRITICAL, "")
	yqlib.GetLogger().SetBackend(backend)
	decoder := yqlib.NewJSONDecoder()
	decoder.Init(bytes.NewReader(json))
	node, err := decoder.Decode()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	prefs := yqlib.NewDefaultJsonPreferences()
	prefs.ColorsEnabled = isOutputToTerminal()
	enc := yqlib.NewJSONEncoder(prefs)
	err = enc.Encode(os.Stdout, node)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

func toVersion(rel *ThreegppCdr_ReleaseVersionIdentifier, ext *ThreegppCdr_ReleaseIdentifierExtension) string {
	if ThreegppCdr_ReleaseVersionIdentifier_Rel(rel.ReleaseIdentifier) == ThreegppCdr_ReleaseVersionIdentifier_Rel__Beyond9 {
		return fmt.Sprintf("1%d.%d", ext.ThreegppRelease, rel.VersionIdentifier)
	}
	return fmt.Sprintf("%d.%d", rel.threegppRelease, rel.VersionIdentifier)
}

func toTimeStamp(ts *ThreegppCdr_Timestamp) string {
	s := ""
	if ts.Sign {
		s = "+"
	} else {
		s = "-"
	}
	dev := fn(ts.HourDeviation) + fn(ts.MinuteDeviation)
	time := fn(ts.Hour) + ":" + fn(ts.Minute) + ":00"
	return strconv.FormatInt(int64(int(ts.Date)), 10) + "/" + strconv.FormatInt(int64(int(ts.Month)), 10) + " " + time + s + dev
}

func fn(num uint64) string {
	return fmt.Sprintf("%02d", num)
}

func toFileClosureTriggerReason(num uint8) string {
	switch num {
	case 0:
		return "Normal closure (Undefined normal closure reason)"
	case 1:
		return "File size limit reached (OAM&P configured)"
	case 2:
		return "File open-time limit reached (OAM&P configured)"
	case 3:
		return "Maximum number of CDRs in file reached (OAM&P configured)"
	case 4:
		return "File closed by manual intervention"
	case 5:
		return "CDR release, version or encoding change"
	case 128:
		return "Abnormal file closure (Undefined error closure reason)"
	case 129:
		return "File system error"
	case 130:
		return "File system storage exhausted"
	case 131:
		return "File integrity error"
	default:
		return "Reserved for future use"
	}
}

func byteArrayToIPv4(byteArray []byte) string {
	// Remove leading bytes with value 0xFF
	startIndex := 0
	for startIndex < len(byteArray) && byteArray[startIndex] == 0xFF {
		startIndex++
	}
	// Extract the first four bytes
	var bytes []byte
	if startIndex+4 <= len(byteArray) {
		bytes = byteArray[startIndex : startIndex+4]
	} else {
		// Handle case where there are fewer than 4 bytes after removing leading 0xFF bytes
		bytes = byteArray[startIndex:]
		for len(bytes) < 4 {
			bytes = append(bytes, 0)
		}
	}
	// Convert bytes to IPv4 format
	return fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3])
}

func decodeLostCdrIndicator(value uint8) string {
	msb := value >> 7         // Get the value of the most significant bit
	lowerBits := value & 0x7F // Get the value of the lower 7 bits

	switch msb {
	case 0:
		if lowerBits == 0 {
			return "No CDRs have been lost"
		} else if lowerBits <= 126 {
			return fmt.Sprintf("CGF has identified that %d CDR(s) were lost, while it is unknown whether more CDRs were lost", lowerBits)
		} else {
			return "CGF has identified that 127 or more CDRs were lost, while it is unknown whether more CDRs were lost"
		}
	case 1:
		if lowerBits == 0 {
			return "CDRs have been lost but CGF cannot determine the number of lost CDRs"
		} else if lowerBits <= 126 {
			return fmt.Sprintf("CGF has calculated the number of lost CDRs as %d", lowerBits)
		} else {
			return "CGF has calculated the number of lost CDRs to be 127 or more"
		}
	default:
		return "Invalid input"
	}
}

func toCdrEncoding(num uint64) string {
	switch num {
	case 1:
		return "BER"
	case 2:
		return "Unaligned PER"
	case 3:
		return "Aligned PER"
	case 4:
		return "XML"
	default:
		return fmt.Sprintf("Unknown %d", num)
	}
}

func toTsNumber(num ThreegppCdr_Cdr_Ts) string {
	switch num {
	case 0:
		return "TS32.005"
	case 1:
		return "TS32.015"
	case 2:
		return "TS32.205"
	case 3:
		return "TS32.215"
	case 4:
		return "TS32.225"
	case 5:
		return "TS32.235"
	case 6:
		return "TS32.250"
	case 7:
		return "TS32.251"
	case 9:
		return "TS32.260"
	case 10:
		return "TS32.270"
	case 11:
		return "TS32.271"
	case 12:
		return "TS32.272"
	case 13:
		return "TS32.273"
	case 14:
		return "TS32.275"
	case 15:
		return "TS32.274"
	case 16:
		return "TS32.277"
	case 17:
		return "TS32.296"
	case 18:
		return "TS32.278"
	case 19:
		return "TS32.253"
	case 20:
		return "TS32.255"
	case 21:
		return "TS32.254"
	case 22:
		return "TS32.256"
	case 23:
		return "TS28.201"
	case 24:
		return "TS28.202"
	case 25:
		return "TS32.257"
	case 26:
		return "TS32.282"
	default:
		return fmt.Sprintf("Unknown %d", num)
	}
}

func getCdrContent(content []byte, index uint32) *ThreegppCdr_Cdr {
	file := ParseCdr(content)
	if index > file.Header.NumberOfCdrsInFile {
		fmt.Printf("Error: Number of CDRS in file: %d\n", file.Header.NumberOfCdrsInFile)
		os.Exit(2)
	}
	return file.Cdrs[index-1]
}
