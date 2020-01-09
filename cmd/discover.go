package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/julienrbrt/yeego/lib/yeelight"
	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover Yeelight bulbs on your network",
	RunE: func(cmd *cobra.Command, args []string) error {
		Lights, err := yeelight.Discover(time.Duration(time.Second))
		if err != nil {
			return err
		}

		if len(Lights) == 0 {
			return errYeelightNotFound
		}

		fmt.Printf("%v Yeelight found on your network\n", len(Lights))

		lightsJSON, err := json.Marshal(Lights)
		if err != nil {
			return err
		}

		// write config file
		err = ioutil.WriteFile(filename, lightsJSON, 0644)
		if err != nil {
			panic(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
