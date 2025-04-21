package main

import (
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/inspectdecodercmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/irigencmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/pipecmd"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "rdfkit",
	}
	cmd.AddCommand(
		pipecmd.New(),
		irigencmd.New(),
		inspectdecodercmd.New(),
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
