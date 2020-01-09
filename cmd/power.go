package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var turnOnCmd = &cobra.Command{
	Use:   "on [name or IP]",
	Short: "Turn on the given light",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(lights) == 0 {
			return errYeelightNotFound
		}

		for i := range lights {
			if lights[i].Name == strings.ToLower(args[0]) || strings.Split(lights[i].Location, ":")[0] == args[0] {
				_, err := lights[i].On()
				if err != nil {
					return err
				}
				fmt.Printf("%s turned on\n", args[0])
				return nil
			}
		}

		return errNotFoundLight
	},
}

var turnOffCmd = &cobra.Command{
	Use:   "off [name or IP]",
	Short: "Turn off the given light",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(lights) == 0 {
			return errYeelightNotFound
		}

		for i := range lights {
			if lights[i].Name == strings.ToLower(args[0]) || strings.Split(lights[i].Location, ":")[0] == args[0] {
				_, err := lights[i].Off()
				if err != nil {
					return err
				}
				fmt.Printf("%s turned off\n", args[0])
				return nil
			}
		}

		return errNotFoundLight
	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle [name or IP]",
	Short: "Toggle the given light",
	Long:  `Toggle inverts the status off a light (on -> off and off -> on).`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(lights) == 0 {
			return errYeelightNotFound
		}

		for i := range lights {
			if lights[i].Name == strings.ToLower(args[0]) || strings.Split(lights[i].Location, ":")[0] == args[0] {
				_, err := lights[i].Toggle()
				if err != nil {
					return err
				}

				fmt.Printf("%s toggled\n", args[0])
				return nil
			}
		}

		return errNotFoundLight
	},
}

func init() {
	rootCmd.AddCommand(turnOnCmd)
	rootCmd.AddCommand(turnOffCmd)
	rootCmd.AddCommand(toggleCmd)
}
