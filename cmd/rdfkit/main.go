package main

import (
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/canonicalizecmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdutil"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/exportdotcmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/exportgoiricmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/pipecmd"
	"github.com/dpb587/rdfkit-go/cmd/rdfkit/versioncmd"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/spf13/cobra"
)

var (
	Version     = "dev"
	BuildTag    = "unknown"
	BuildCommit = "0000000000"
	BuildClean  = "unknown"
	BuildTime   = "0000-00-00T00:00:00Z"
)

func main() {
	cmd := &cobra.Command{
		Use:           "rdfkit",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	app := &cmdutil.App{
		Registry: rdfio.Registry,
	}

	cmd.AddCommand(
		pipecmd.New(app),
		canonicalizecmd.New(app),
		exportdotcmd.New(app),
		exportgoiricmd.New(app),
		versioncmd.New(versioncmd.Properties{
			Name:        "rdfkit",
			Version:     Version,
			BuildTag:    BuildTag,
			BuildCommit: BuildCommit,
			BuildClean:  BuildClean,
			BuildTime:   BuildTime,
		}),
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
