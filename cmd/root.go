package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/julienrbrt/yeego/lib/yeelight"
	"github.com/spf13/cobra"
)

var (
	// Lights are the yeelight present on LAN
	lights []yeelight.Yeelight

	// error messages
	errNotFoundLight    = errors.New("Light not found")
	errYeelightNotFound = errors.New("No Yeelight found. Run `yeego discover` to find lights under your network")

	// configuration
	filename = ".yeego"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yeego",
	Short: "Control your Yeelight bulb with Yeego",
	Long: `Yeego is a tool written in Go which permits to control
your Yeelight bulbs in your LAN directly from your terminal.`,
	Example: `yeego discover
yeego on bedroom`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Welcome in Yeego!")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// read configuration file - no error check
	file, _ := ioutil.ReadFile(filename)
	json.Unmarshal(file, &lights)
}
