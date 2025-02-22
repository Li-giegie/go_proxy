/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Li-giegie/go_proxy/interval"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// xorclientCmd represents the xorclient command
var xorclientCmd = &cobra.Command{
	Use:   "xorclient",
	Short: "http/s proxy XOR encryption tunnel client",
	Long:  `http/s proxy XOR encryption tunnel client`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		proxy, _ := cmd.Flags().GetString("proxy")
		key, _ := cmd.Flags().GetString("key")
		err := interval.StartForward(addr, &interval.XORForward{
			ProxyAddr: proxy,
			Key:       []byte(key),
		})
		if err != nil {
			logrus.Errorf("forward err: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(xorclientCmd)
	xorclientCmd.Flags().StringP("addr", "a", ":1080", "listen address")
	xorclientCmd.Flags().StringP("proxy", "p", ":1080", "proxy server address")
	xorclientCmd.Flags().StringP("key", "k", "!vjhkm&#45^78(", "proxy tunnel key")
}
