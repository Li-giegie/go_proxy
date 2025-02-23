/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Li-giegie/go_proxy/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
)

var xorserverCmd = &cobra.Command{
	Use:   "xorserver",
	Short: "http/s proxy XOR encryption tunnel server",
	Long:  `http/s proxy XOR encryption tunnel server`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		key, _ := cmd.Flags().GetString("key")
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logrus.Errorf("listen err: %v", err)
			return
		}
		err = internal.StartProxy(&internal.XORProxy{
			Key:      []byte(key),
			Listener: l,
		})
		if err != nil {
			logrus.Errorf("err: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(xorserverCmd)
	xorserverCmd.Flags().StringP("addr", "a", ":1080", "listen address")
	xorserverCmd.Flags().StringP("key", "k", "!vjhkm&#45^78(", "proxy tunnel key")
}
