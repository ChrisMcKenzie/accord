package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/datascienceinc/accord/pkg/module"
	getter "github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

// DefaultDataDir is the default directory for storing local data
const DefaultDataDir = ".accord"

var cfgFile string

var ctx *Context

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

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $CWD/.accord)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" { // enable ability to specify config file via flag
		cfgFile = "accord.hcl"
	}

	tree, err := module.NewTreeModule("", cfgFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = tree.Load(&getter.FolderStorage{
		StorageDir: filepath.Join(DefaultDataDir, "modules"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx = NewContext(tree)
}
