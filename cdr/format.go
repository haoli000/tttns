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

import "fmt"

func formatVersion(rel ReleaseVersion, ext *uint8) string {
	if rel.ReleaseIdentifier == 7 && ext != nil {
		return fmt.Sprintf("1%d.%d", *ext, rel.VersionIdentifier)
	}
	return fmt.Sprintf("%d.%d", rel.ReleaseIdentifier, rel.VersionIdentifier)
}

func formatTimestamp(ts Timestamp) string {
	sign := "+"
	if !ts.Sign {
		sign = "-"
	}
	return fmt.Sprintf("%d/%d %02d:%02d:00%s%02d%02d",
		ts.Day, ts.Month, ts.Hour, ts.Minute, sign, ts.HourDeviation, ts.MinuteDeviation)
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
	startIndex := 0
	for startIndex < len(byteArray) && byteArray[startIndex] == 0xFF {
		startIndex++
	}
	var ipBytes []byte
	if startIndex+4 <= len(byteArray) {
		ipBytes = byteArray[startIndex : startIndex+4]
	} else {
		ipBytes = byteArray[startIndex:]
		for len(ipBytes) < 4 {
			ipBytes = append(ipBytes, 0)
		}
	}
	return fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

func decodeLostCdrIndicator(value uint8) string {
	msb := value >> 7
	lowerBits := value & 0x7F

	if msb == 0 {
		if lowerBits == 0 {
			return "No CDRs have been lost"
		} else if lowerBits <= 126 {
			return fmt.Sprintf("CGF has identified that %d CDR(s) were lost, while it is unknown whether more CDRs were lost", lowerBits)
		}
		return "CGF has identified that 127 or more CDRs were lost, while it is unknown whether more CDRs were lost"
	}
	if lowerBits == 0 {
		return "CDRs have been lost but CGF cannot determine the number of lost CDRs"
	} else if lowerBits <= 126 {
		return fmt.Sprintf("CGF has calculated the number of lost CDRs as %d", lowerBits)
	}
	return "CGF has calculated the number of lost CDRs to be 127 or more"
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

func toTsNumber(num uint8) string {
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

func fileToHeaderInfo(file *ThreegppCdrFile) FileHeaderInfo {
	h := file.Header
	return FileHeaderInfo{
		FileLength:               int(h.FileLength),
		HeaderLength:             int(h.HeaderLength),
		HighReleaseVersion:       formatVersion(h.HighReleaseVersionIdentifier, h.HighReleaseIdentifierExtension),
		LowReleaseVersion:        formatVersion(h.LowReleaseVersionIdentifier, h.LowReleaseIdentifierExtension),
		FileOpeningTimestamp:     formatTimestamp(h.FileOpeningTimestamp),
		LastCDRAppendTimestamp:   formatTimestamp(h.LastCdrAppendTimestamp),
		NumberOfCDRsInFile:       int(h.NumberOfCdrsInFile),
		FileSequenceNumber:       int(h.FileSequenceNumber),
		FileClosureTriggerReason: fmt.Sprintf("%d - %s", h.FileClosureTriggerReason, toFileClosureTriggerReason(h.FileClosureTriggerReason)),
		NodeIPAddress:            byteArrayToIPv4(h.NodeIpAddress[:]),
		LostCDRIndicator:         decodeLostCdrIndicator(h.LostCdrIndicator),
	}
}

func cdrToHeaderInfo(row CdrRecord) CdrHeaderInfo {
	return CdrHeaderInfo{
		CdrLength:          int(row.CdrLength),
		ReleaseVersion:     formatVersion(row.Version, row.ReleaseIdentifierExtension),
		DataRecorderFormat: toCdrEncoding(uint64(row.DataRecordFormat)),
		TsNumber:           toTsNumber(row.TsNumber),
	}
}
