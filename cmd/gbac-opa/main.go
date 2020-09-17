package main

import (
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/cmd"
	grego "github.com/richicoder1/gbac/pkg/rego"
)

func main() {
	grego.RegisterBuiltins()

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
