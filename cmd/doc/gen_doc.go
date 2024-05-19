package main

import (
	"log"

	"github.com/haoli000/tttns/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "../../docs")
	if err != nil {
		log.Fatal(err)
	}
}
