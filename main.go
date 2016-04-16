package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bwalex/rkt-inspect/pods"
	"github.com/spf13/cobra"
)

const rktDataDirDefault = "/var/lib/rkt"

var globalFlags = struct {
	DataDir      string
	UUIDFile     string
	UUID         string
	WaitRunning  bool
	WaitUUIDFile bool
}{
	DataDir:      rktDataDirDefault,
	UUIDFile:     "/tmp/pod.uuid",
	UUID:         "",
	WaitRunning:  false,
	WaitUUIDFile: false,
}

var RootCmd = &cobra.Command{
	Use:   "rkt-inspect [command]",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		println("missing command")
		cmd.HelpFunc()(cmd, args)
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&globalFlags.DataDir, "data-dir", rktDataDirDefault, "rkt data directory")
	RootCmd.PersistentFlags().StringVar(&globalFlags.UUID, "uuid", "", "pod's uuid")
	RootCmd.PersistentFlags().StringVar(&globalFlags.UUIDFile, "uuid-file", "", "file containing pod's uuid")
	RootCmd.PersistentFlags().BoolVar(&globalFlags.WaitRunning, "wait-run", false, "wait for pod to run if it isn't running yet")
	RootCmd.PersistentFlags().BoolVar(&globalFlags.WaitUUIDFile, "wait-uuid", false, "wait for uuid file to appear")

	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if globalFlags.UUIDFile != "" {
			var err error

			if globalFlags.WaitUUIDFile {
				WaitUUIDFile(globalFlags.UUIDFile, 300*time.Second)
			}

			globalFlags.UUID, err = ReadUUID(globalFlags.UUIDFile)
			if err != nil {
				panic(err)
			}
		}

		if globalFlags.WaitRunning {
			err := pods.WaitPodRun(globalFlags.DataDir, globalFlags.UUID, 300*time.Second)
			if err != nil {
				panic(err)
			}
		}

	}
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
