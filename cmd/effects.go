package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/julienrbrt/yeego/light/yeelight"
	"github.com/spf13/cobra"
)

// error messages
var errDuration = errors.New("Duration must be in seconds")

var temperatureCmd = &cobra.Command{
	Use:   "temp [name or IP] [color temperature in k] [duration in sec]",
	Short: "Change the color temperature of a given light",
	Long: `Change the color temperature of a given light
The range is from 1700 to 6500 (k)`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(lights) == 0 {
			return errYeelightNotFound
		}

		color, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Color temperature is mandatory and integer")
		}

		var duration int

		if len(args) > 2 {
			duration, err = strconv.Atoi(args[2])
			duration = duration / 1000 // second to ms
			if err != nil {
				return errDuration
			}
		}

		for i := range lights {
			if yeelight.Matching(lights[i], args[0]) {

				_, err := lights[i].SetCtAbx(color, duration)
				if err != nil {
					return err
				}
				fmt.Printf("%s color temperature updated\n", lights[i].Name)
				return nil
			}
		}

		return nil
	},
}

var colorCmd = &cobra.Command{
	Use:   "color [name or IP] [color in hex] [duration in sec]",
	Short: "Change the color of a given light",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(lights) == 0 {
			return errYeelightNotFound
		}

		value, err := strconv.ParseInt(args[1], 16, 32)
		if err != nil {
			return errors.New("Color is mandatory and hexademical")
		}

		var duration int

		if len(args) > 2 {
			duration, err = strconv.Atoi(args[2])
			duration = duration / 1000 // second to ms
			if err != nil {
				return errDuration
			}
		}

		for i := range lights {
			if yeelight.Matching(lights[i], args[0]) {
				_, err := lights[i].SetRGBhex(int(value), duration)
				if err != nil {
					return err
				}
				fmt.Printf("%s color updated\n", lights[i].Name)
				return nil
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(temperatureCmd)
	rootCmd.AddCommand(colorCmd)
}
