package main

import (
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/canonicalizecmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdutil"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/exportdotcmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/exportgoiricmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/pipecmd"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "rdfkit",
	}

	app := &cmdutil.App{
		Registry: rdfio.Registry,
	}

	cmd.AddCommand(
		pipecmd.New(app),
		canonicalizecmd.New(app),
		exportdotcmd.New(app),
		exportgoiricmd.New(app),
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
