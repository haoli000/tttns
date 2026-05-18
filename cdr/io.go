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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

func GetContent(filename string) []byte {
	var inputReader io.Reader
	if filename == "-" {
		fi, err := os.Stdin.Stat()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error checking stdin:", err)
			os.Exit(1)
		}
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			_, _ = fmt.Fprintln(os.Stderr, "Error: No data available on stdin")
			os.Exit(1)
		}
		inputReader = os.Stdin
	} else {
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

func ParseCdr(content []byte) *ThreegppCdrFile {
	file, err := Parse(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse CDR file:", err)
		os.Exit(2)
	}
	return file
}

func ToFileHeaderInfo(content []byte) FileHeaderInfo {
	file := ParseCdr(content)
	return fileToHeaderInfo(file)
}

func ToCdrHeaderInfo(content []byte, index uint32) CdrHeaderInfo {
	row := getCdrRecord(content, index)
	return cdrToHeaderInfo(row)
}

func CountCdrs(content []byte) uint32 {
	file := ParseCdr(content)
	return file.Header.NumberOfCdrsInFile
}

func DumpCdr(content []byte, index uint32, file *os.File) {
	row := getCdrRecord(content, index)
	_, err := file.Write(row.CdrContent)
	if err != nil {
		fmt.Println("Error dumping CDR:", err)
		os.Exit(4)
	}
}

func ToCdrInfo(content []byte) CdrInfo {
	file := ParseCdr(content)
	cnt := file.Header.NumberOfCdrsInFile
	cdrHeaderInfos := make([]CdrHeaderInfo, 0, cnt)
	for i := uint32(0); i < cnt; i++ {
		cdrHeaderInfos = append(cdrHeaderInfos, cdrToHeaderInfo(file.Cdrs[i]))
	}
	return CdrInfo{
		NumberOfCDRs: int(cnt),
		CdrHeaders:   cdrHeaderInfos,
	}
}

func ToFileInfo(content []byte) FileInfo {
	file := ParseCdr(content)
	cnt := file.Header.NumberOfCdrsInFile
	cdrHeaderInfos := make([]CdrHeaderInfo, 0, cnt)
	for i := uint32(0); i < cnt; i++ {
		cdrHeaderInfos = append(cdrHeaderInfos, cdrToHeaderInfo(file.Cdrs[i]))
	}
	return FileInfo{
		HeaderInfo: fileToHeaderInfo(file),
		CdrInfo: CdrInfo{
			NumberOfCDRs: int(cnt),
			CdrHeaders:   cdrHeaderInfos,
		},
	}
}

func isOutputToTerminal() bool {
	o, _ := os.Stdout.Stat()
	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}

func PrintOutput(jsonOutput bool, v any) {
	jsonBytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if jsonOutput {
		prettyPrintJSON(jsonBytes)
	} else {
		prettyPrintYAML(jsonBytes)
	}
}

func initYqLogger() {
	backend := logging.SetBackend(logging.NewLogBackend(os.Stderr, "", log.LstdFlags))
	backend.SetLevel(logging.CRITICAL, "")
	yqlib.GetLogger().SetBackend(backend)
}

func prettyPrintYAML(jsonData []byte) {
	initYqLogger()
	decoder := yqlib.NewJSONDecoder()
	_ = decoder.Init(bytes.NewReader(jsonData))
	node, err := decoder.Decode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	prefs := yqlib.NewDefaultYamlPreferences()
	prefs.ColorsEnabled = isOutputToTerminal()
	enc := yqlib.NewYamlEncoder(prefs)
	err = enc.Encode(os.Stdout, node)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func prettyPrintJSON(jsonData []byte) {
	initYqLogger()
	decoder := yqlib.NewJSONDecoder()
	_ = decoder.Init(bytes.NewReader(jsonData))
	node, err := decoder.Decode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	prefs := yqlib.NewDefaultJsonPreferences()
	prefs.ColorsEnabled = isOutputToTerminal()
	enc := yqlib.NewJSONEncoder(prefs)
	err = enc.Encode(os.Stdout, node)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func getCdrRecord(content []byte, index uint32) CdrRecord {
	file := ParseCdr(content)
	if index > file.Header.NumberOfCdrsInFile || index == 0 {
		fmt.Printf("Error: Number of CDRS in file: %d\n", file.Header.NumberOfCdrsInFile)
		os.Exit(2)
	}
	return file.Cdrs[index-1]
}
