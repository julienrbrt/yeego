package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var temperatureCmd = &cobra.Command{
	Use:   "set-temp [name/IP] [color temperature in k]",
	Short: "Change the color temperature of a given light",
	Long: `Change the color temperature of a given light
The range is from 1700 to 6500 (k)`,
	Example: "yeego set-temp bedroom 3500",
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		color, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Color temperature is mandatory")
		}

		_, err = light.SetCtAbx(color, int(timeout.Milliseconds()))
		if err != nil {
			return err
		}

		fmt.Printf("%s color temperature updated\n", args[0])
		return nil
	},
}

var colorCmd = &cobra.Command{
	Use:     "set-color [name/IP] [color in hex]",
	Short:   "Change the color of a given light",
	Example: "yeego set-color bedroom ffffff",
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		value, err := strconv.ParseInt(args[1], 16, 32)
		if err != nil {
			return errors.New("Color is mandatory")
		}

		_, err = light.SetRGBhex(int(value), int(timeout.Milliseconds()))
		if err != nil {
			return err
		}

		fmt.Printf("%s color updated\n", args[0])
		return nil
	},
}

var brightnessCmd = &cobra.Command{
	Use:     "set-bright [name/IP] [level]",
	Short:   "Change the brightness of a given light",
	Example: "yeego set-bright bedroom 75",
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		brightness, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Brightness is mandatory")
		}

		_, err = light.SetBright(brightness, int(timeout.Milliseconds()))
		if err != nil {
			return err
		}

		fmt.Printf("%s brightness updated\n", args[0])
		return nil
	},
}

var adjustCmd = &cobra.Command{
	Use:   "adjust [name/IP] [action] [property]",
	Short: "Adjust the status of a given light",
	Long: `Adjust the status of a given light without knowing its status
by giving an action to perform and the property to change`,
	Example: `yeego adjust bedroom increase bright
yeego adjust bedroom decrease bright
yeego adjust bedroom increase ct
yeego adjust bedroom decrease ct
yeego adjust bedroom color cirle`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		if args[1] == "" || args[2] == "" {
			return errors.New("An action and a property is mandatory")
		}

		_, err = light.SetAdjust(args[1], args[2])
		if err != nil {
			return err
		}

		fmt.Printf("%s adjusted\n", args[0])
		return nil
	},
}

func init() {
	temperatureCmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Millisecond, "Timeout temperature change effect")
	colorCmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Millisecond, "Timeout color change effect")
	brightnessCmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Millisecond, "Timeout brightness change effect")

	rootCmd.AddCommand(temperatureCmd)
	rootCmd.AddCommand(colorCmd)
	rootCmd.AddCommand(brightnessCmd)
	rootCmd.AddCommand(adjustCmd)
}
