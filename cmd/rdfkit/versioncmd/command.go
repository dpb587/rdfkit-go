package versioncmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

type Properties struct {
	Name    string
	Version string

	BuildTag    string
	BuildCommit string
	BuildClean  string
	BuildTime   string
}

func New(v Properties) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout := cmd.OutOrStdout()

			fmt.Fprintf(stdout, "%s version=%s\n", v.Name, v.Version)
			fmt.Fprintf(stdout, "build tag=%s commit=%s clean=%s time=%s\n", v.BuildTag, v.BuildCommit, v.BuildClean, v.BuildTime)
			fmt.Fprintf(stdout, "runtime os=%s arch=%s\n", runtime.GOOS, runtime.GOARCH)

			return nil
		},
	}

	return cmd
}
