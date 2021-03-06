package main

import (
	"fmt"
	"github.com/bwalex/rkt-inspect/pods"
	"github.com/spf13/cobra"
	"os"
)

var (
	flagNetName string = "default"
)

var CmdIp = &cobra.Command{
	Use:   "ip",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		netInfo, err := pods.GetPodNetInfo(globalFlags.DataDir, globalFlags.UUID, flagNetName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print(netInfo.Ip)
	},
}

var CmdIfname = &cobra.Command{
	Use:   "ifname",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		netInfo, err := pods.GetPodNetInfo(globalFlags.DataDir, globalFlags.UUID, flagNetName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}

		fmt.Print(netInfo.IfName)
	},
}

func init() {
	RootCmd.AddCommand(CmdIp)
	RootCmd.AddCommand(CmdIfname)
	CmdIp.Flags().StringVar(&flagNetName, "net", "default", "network name")
	CmdIfname.Flags().StringVar(&flagNetName, "net", "default", "network name")
}
