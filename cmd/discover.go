package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/julienrbrt/yeego/light/yeelight"
	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover Yeelight bulbs on your network",
	RunE: func(cmd *cobra.Command, args []string) error {
		lights, err := yeelight.Discover(time.Duration(timeout))
		if err != nil {
			return err
		}

		fmt.Printf("%v Yeelight found on your network.\n", len(lights))

		//write configuration file
		err = writeConfig()
		return err
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List the saved Yeelight",
	Long:    "List the saved Yeelight from the .yeego configuration file",
	Example: "yeego list",
	RunE: func(cmd *cobra.Command, args []string) error {
		// no light found
		if len(lights) == 0 {
			fmt.Println("No Yeelight saved in configuration")
			return nil
		}

		fmt.Printf("%v Yeelight saved in configuration:\n", len(lights))
		for i, light := range lights {
			if light.Name == "" {
				light.Name = "Unknown [no name]"
			}
			fmt.Printf("- %d: %s on %v\n", i+1, light.Name, strings.Split(light.Location, ":")[0])
		}

		return nil
	},
}

func init() {
	discoverCmd.Flags().DurationVarP(&timeout, "timeout", "t", time.Second, "Timeout for discover")
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(listCmd)
}
