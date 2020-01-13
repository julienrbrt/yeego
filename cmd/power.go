package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var turnOnCmd = &cobra.Command{
	Use:   "on [name/IP]",
	Short: "Turn on the given light",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		_, err = light.On()
		if err != nil {
			return err
		}

		fmt.Printf("%s turned on\n", args[0])
		return nil
	},
}

var turnOffCmd = &cobra.Command{
	Use:   "off [name/IP]",
	Short: "Turn off the given light",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		_, err = light.Off()
		if err != nil {
			return err
		}

		fmt.Printf("%s turned off\n", args[0])
		return nil

	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle [name/IP]",
	Short: "Toggle the given light",
	Long:  `Toggle inverts the status off a light (on -> off and off -> on).`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		_, err = light.Toggle()
		if err != nil {
			return err
		}

		fmt.Printf("%s toggled\n", args[0])
		return nil

	},
}

func init() {
	rootCmd.AddCommand(turnOnCmd)
	rootCmd.AddCommand(turnOffCmd)
	rootCmd.AddCommand(toggleCmd)
}
