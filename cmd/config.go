package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var getPropsCmd = &cobra.Command{
	Use:   "props [name/IP]",
	Short: "Get properties of a given light",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		err = light.GetProp()
		if err != nil {
			return err
		}

		lightJSON, err := json.Marshal(light)
		if err != nil {
			return err
		}

		fmt.Printf("%s properties\n %s\n", args[0], lightJSON)
		return nil

	},
}

var setDefaultCmd = &cobra.Command{
	Use:   "set-default [name/IP]",
	Short: "Set state of given light as default",
	Long: `Save state of given light as default.
If the yeelight is turned off from power, the saved status is used when powered on`,
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		_, err = light.SetDefault()
		if err != nil {
			return err
		}

		fmt.Printf("%s settings saved as default\n", args[0])
		return nil

	},
}

var setNameCmd = &cobra.Command{
	Use:     "set-name [name/IP] [new name]",
	Short:   "Gives a name to a given light",
	Example: "yeego set-name bedroom",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		if len(args[1]) == 0 {
			return errors.New("Please enter the new name of the light")
		}

		_, err = light.SetName(args[1])
		if err != nil {
			return err
		}

		// TODO write name in config

		fmt.Printf("%s name saved\n", args[0])
		return nil

	},
}

func init() {
	rootCmd.AddCommand(getPropsCmd)
	rootCmd.AddCommand(setDefaultCmd)
	rootCmd.AddCommand(setNameCmd)
}
