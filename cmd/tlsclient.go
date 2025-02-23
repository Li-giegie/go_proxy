/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Li-giegie/go_proxy/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// proxyclientCmd represents the proxyclient command
var tlsclientCmd = &cobra.Command{
	Use:   "tlsclient",
	Short: "http/s proxy tls tunnel client",
	Long:  `http/s proxy tls tunnel client`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		proxy, _ := cmd.Flags().GetString("proxy")
		pem, _ := cmd.Flags().GetString("pem")
		key, _ := cmd.Flags().GetString("key")
		err := internal.StartForward(addr, &internal.TLSForward{
			ProxyAddr: proxy,
			Pem:       pem,
			Key:       key,
		})
		if err != nil {
			logrus.Errorf("tls client proxy err: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tlsclientCmd)
	tlsclientCmd.Flags().StringP("addr", "a", ":1080", "local listen address")
	tlsclientCmd.Flags().StringP("proxy", "p", "", "proxy server address")
	tlsclientCmd.Flags().String("pem", "", "pem file path")
	tlsclientCmd.Flags().String("key", "", "key file path")
}
