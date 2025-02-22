/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Li-giegie/go_proxy/interval"

	"github.com/spf13/cobra"
)

// proxyclientCmd represents the proxyclient command
var proxyclientCmd = &cobra.Command{
	Use:   "proxyclient",
	Short: "http/s tls tunnel proxy client",
	Long:  `http/s tls tunnel proxy client`,
	Run: func(cmd *cobra.Command, args []string) {
		err := interval.StartProxyClient(*listenAddr, *proxyAddr, *pemPath, *keyPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var (
	proxyAddr  *string
	listenAddr *string
	pemPath    *string
	keyPath    *string
)

func init() {
	rootCmd.AddCommand(proxyclientCmd)
	listenAddr = proxyclientCmd.Flags().StringP("listen", "l", ":1080", "local listen address")
	proxyAddr = proxyclientCmd.Flags().StringP("proxy", "p", ":1080", "proxy server address")
	pemPath = proxyclientCmd.Flags().String("pem", "", "pem file path")
	keyPath = proxyclientCmd.Flags().String("key", "", "key file path")
}
