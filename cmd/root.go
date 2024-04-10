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
	"encoding/json"
	"fmt"
	"os"

	"github.com/haoli000/tttns/cdr"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tttns [file|-]",
	Short: "3GPP TS 32.297 cdr decoder",
	Long: `A cli tool to inspect 3GPP TS 32.297 CDR files.
The name is plainly from the first letters of 32297.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := cdr.GetContent(args[0])
		fileInfo := cdr.ToFileInfo(content)
		jsonBytes, err := json.MarshalIndent(fileInfo, "", "    ")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			cdr.PrettyPrintJSON(jsonBytes)
		} else {
			cdr.PrettyPrintYAML(jsonBytes)
		}
	}}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
