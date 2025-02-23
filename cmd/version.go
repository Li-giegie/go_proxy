/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

const (
	Version = "0.3"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go_proxy")
		fmt.Println("version:", Version)
		fmt.Println("Go     :", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
