package main

import (
	"fmt"
	"github.com/bwalex/rkt-inspect/pods"
	"github.com/spf13/cobra"
	"os"
)

var CmdPid = &cobra.Command{
	Use:   "pid",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := pods.GetPodPid(globalFlags.DataDir, globalFlags.UUID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		fmt.Print(pid)
	},
}

func init() {
	RootCmd.AddCommand(CmdPid)
}
