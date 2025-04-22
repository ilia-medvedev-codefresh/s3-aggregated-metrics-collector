package main

import (
	"github.com/ilia-medvedev-codefresh/s3-aggregated-metrics-collector/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	// Generate markdown documentation for the root command
	err := doc.GenMarkdownTree(cmd.GetCmd().Root(), "./docs")
	if err != nil {
		panic(err)
	}
}
