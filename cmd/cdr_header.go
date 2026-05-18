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
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/haoli000/tttns/cdr"
	"github.com/spf13/cobra"
)

// headerCmd represents the header command
var headerCmd = &cobra.Command{
	Use:   "header [file|-] [index|1]",
	Short: "Print CDR header info",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := "-"
		var index uint64 = 1
		if len(args) == 1 {
			if n, err := strconv.ParseUint(args[0], 10, 32); err == nil {
				index = n
			} else {
				fileName = args[0]
			}
		} else if len(args) > 1 {
			fileName = args[0]
			n, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(3)
			}
			index = n
		}
		if index < 1 {
			fmt.Println("Error: Index must be an integer starting from 1")
			os.Exit(3)
		}

		content := cdr.GetContent(fileName)
		info := cdr.ToCdrHeaderInfo(content, uint32(index))
		jsonOutput, _ := cmd.Flags().GetBool("json")
		cdr.PrintOutput(jsonOutput, info)
	},
}

func init() {
	headerCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	cdrCmd.AddCommand(headerCmd)
}
