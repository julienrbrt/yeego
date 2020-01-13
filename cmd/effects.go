package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var temperatureCmd = &cobra.Command{
	Use:   "temp [name or IP] [color temperature in k]",
	Short: "Change the color temperature of a given light",
	Long: `Change the color temperature of a given light
The range is from 1700 to 6500 (k)`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		color, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Color temperature is mandatory and integer")
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
	Use:   "color [name or IP] [color in hex]",
	Short: "Change the color of a given light",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(lights, args[0])
		if err != nil {
			return err
		}

		value, err := strconv.ParseInt(args[1], 16, 32)
		if err != nil {
			return errors.New("Color is mandatory and hexademical")
		}

		_, err = light.SetRGBhex(int(value), int(timeout.Milliseconds()))
		if err != nil {
			return err
		}

		fmt.Printf("%s color updated\n", args[0])
		return nil
	},
}

func init() {
	temperatureCmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Millisecond, "Timeout temperature change effect")
	colorCmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Millisecond, "Timeout color change effect")

	rootCmd.AddCommand(temperatureCmd)
	rootCmd.AddCommand(colorCmd)
}
