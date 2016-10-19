package cmd

import (
	"fmt"
	"log"
	"os"

	accord "github.com/datascienceinc/accord/pkg"
	"github.com/spf13/cobra"
)

var cfgFile string

var acc *accord.Accord

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "accord",
	Short: "A consumer contract peace keeping mission.",
	Long:  ``,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $CWD/.accord)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" { // enable ability to specify config file via flag
		cfgFile = ".accord"
	}

	var err error
	acc, err = accord.Load(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
}
