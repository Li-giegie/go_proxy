/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Li-giegie/go_proxy/interval"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "http/s proxy server (no tunnels)",
	Long:  "http/s proxy server (no tunnels)",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logrus.Errorf("listen %s error: %v", addr, err)
			return
		}
		defer l.Close()
		if err = interval.StartProxy(l); err != nil {
			logrus.Errorf("err:%v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("addr", "a", ":1080", "listen address")
}
