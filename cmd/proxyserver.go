/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Li-giegie/go_proxy/interval"
	"github.com/spf13/cobra"
)

// proxyserverCmd represents the proxyserver command
var proxyserverCmd = &cobra.Command{
	Use:   "proxyserver",
	Short: "http/s tls tunnel proxy server",
	Long:  `http/s tls tunnel proxy server`,
	Run: func(cmd *cobra.Command, args []string) {
		err := interval.StartProxyServer(*listenAddr, *pemPathSrv, *keyPathSrv, *pemPathCli)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
var (
	listenAddrSrv *string
	pemPathSrv    *string
	pemPathCli    *string
	keyPathSrv    *string
)

func init() {
	rootCmd.AddCommand(proxyserverCmd)
	listenAddrSrv = proxyserverCmd.Flags().StringP("listen", "l", ":1080", "local listen address")
	pemPathSrv = proxyserverCmd.Flags().String("pem", "", "pem file path")
	pemPathCli = proxyserverCmd.Flags().String("clientpem", "", "client pem file path")
	keyPathSrv = proxyserverCmd.Flags().String("key", "", "key file path")
}
