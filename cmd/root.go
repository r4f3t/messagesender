/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	cfgFile      string
	cobraCommand *cobra.Command
}

// rootCmd represents the base command when called without any subcommands
var RootCommand = RootCmd{
	cobraCommand: &cobra.Command{
		Use:   "webapi",
		Short: "Servis Solution",
		Long:  "Service Solution Smple",
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCommand.cobraCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//defines flags that willbe used for all commands

}

func (c *RootCmd) AddCommand(cmd *cobra.Command) {
	c.cobraCommand.AddCommand(cmd)
}
