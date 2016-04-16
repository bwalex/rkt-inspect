package main

import (
	"fmt"
	"github.com/bwalex/rkt-inspect/pods"
	"github.com/spf13/cobra"
	"os"
)

var CmdPpid = &cobra.Command{
	Use:   "ppid",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		ppid, err := pods.GetPodPpid(globalFlags.DataDir, globalFlags.UUID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		fmt.Print(ppid)
	},
}

func init() {
	RootCmd.AddCommand(CmdPpid)
}
