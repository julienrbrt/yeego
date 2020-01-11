package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/julienrbrt/yeego/light/yeelight"
	"github.com/spf13/cobra"
)

// error messages
var errDuration = errors.New("Duration must be in ms")

var temperatureCmd = &cobra.Command{
	Use:   "temp [name or IP] [color temperature in k] [duration in ms]",
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
			return errors.New("Color temperature is mandatory and an integer")
		}

		var duration int

		if len(args) > 2 {
			duration, err = strconv.Atoi(args[2])
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
	Use:   "color [name or IP] [color in rgb] [duration in ms]",
	Short: "Change the color of a given light",
	Args:  cobra.MinimumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(lights) == 0 {
			return errYeelightNotFound
		}

		red, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Red is mandatory and an integer")
		}

		green, err := strconv.Atoi(args[2])
		if err != nil {
			return errors.New("Green is mandatory and an integer")
		}

		blue, err := strconv.Atoi(args[3])
		if err != nil {
			return errors.New("Blue is mandatory and an integer")
		}

		var duration int

		if len(args) > 4 {
			duration, err = strconv.Atoi(args[4])
			if err != nil {
				return errDuration
			}
		}

		for i := range lights {
			if yeelight.Matching(lights[i], args[0]) {
				_, err := lights[i].SetRGB(red, green, blue, duration)
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
