/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/kilianp07/simu/core"
	"github.com/kilianp07/simu/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simu",
	Short: "A simulator for Modbus interfaces of PV meters and batteries",
	Long: `Simu is a command-line application that simulates Modbus interfaces for PV (Photovoltaic) meters and batteries. 
	This application helps you test and emulate the behavior of these devices, facilitating development and testing.`,
	Run: launch,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	var configFile string
	rootCmd.PersistentFlags().StringVar(&configFile, "conf", "", "path to config file")
	_ = rootCmd.MarkFlagRequired("conf")
}

func launch(cmd *cobra.Command, args []string) {
	var (
		confpath string
		err      error
	)
	if confpath, err = cmd.Flags().GetString("conf"); err != nil || confpath == "" {
		log := logger.Get()
		log.Fatal().Msg("No config file provided")
		return
	}
	core.Launch(confpath)

}
