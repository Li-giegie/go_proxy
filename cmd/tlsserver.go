/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Li-giegie/go_proxy/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tlsserverCmd = &cobra.Command{
	Use:   "tlsserver",
	Short: "http/s proxy tls tunnel server",
	Long:  `http/s proxy tls tunnel server`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		pem, _ := cmd.Flags().GetString("pem")
		cpem, _ := cmd.Flags().GetString("cpem")
		key, _ := cmd.Flags().GetString("key")
		l, err := internal.NewTLSListen(addr, pem, key, cpem)
		if err != nil {
			logrus.Errorf("new tls listen err: %s", err)
			return
		}
		if err = internal.StartProxy(l); err != nil {
			logrus.Errorf("err: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tlsserverCmd)
	tlsserverCmd.Flags().StringP("addr", "a", ":1080", "listen address")
	tlsserverCmd.Flags().String("pem", "", "pem file path")
	tlsserverCmd.Flags().String("cpem", "", "client pem file path")
	tlsserverCmd.Flags().String("key", "", "key file path")
}
