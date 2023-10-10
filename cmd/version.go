/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示版本类型`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("当前使用的版本为: 1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
