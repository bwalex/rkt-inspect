package main

import (
    "os"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/bwalex/rkt-inspect/pods"
)

var (
    flagSubsystem string = "memory"
)

var CmdCgroup = &cobra.Command{
    Use: "cgroup",
    Short: "",
    Run: func(cmd *cobra.Command, args []string) {
        cgroup,err := pods.GetPodCgroup(globalFlags.DataDir, globalFlags.UUID, flagSubsystem)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s", err)
            os.Exit(1)
        }

        fmt.Print(cgroup)
    },
}

func init() {
    RootCmd.AddCommand(CmdCgroup)
    CmdCgroup.Flags().StringVar(&flagSubsystem, "subsystem", "memory", "subsystem name")
}
