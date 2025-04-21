package main


import (
	"github.com/spf13/cobra/doc"
	"github.com/ilia-medvedev-codefresh/s3-aggregated-metrics-collector/cmd"
)

func main() {
	// Generate markdown documentation for the root command
	err := doc.GenMarkdownTree(cmd.GetCmd().Root(), "./docs")
	if err != nil {
		panic(err)
	}
}
