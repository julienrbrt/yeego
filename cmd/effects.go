package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		light, err := argToYeelight(args[0])
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
		light, err := argToYeelight(args[0])
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
		light, err := argToYeelight(args[0])
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
yeego adjust bedroom cirle color`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(args[0])
		if err != nil {
			return err
		}

		if args[1] == "" || args[2] == "" {
			return errors.New("Action and property are mandatory")
		}

		_, err = light.SetAdjust(args[1], args[2])
		if err != nil {
			return err
		}

		fmt.Printf("%s adjusted\n", args[0])
		return nil
	},
}

var colorFlowCmd = &cobra.Command{
	Use:   "start-cf [name/IP] [count] [action] [expression]",
	Short: "Start running a color flow (cf)",
	Long: `Start running a color flow (cf)
"count" is the total number of visible state changing before color flow
	stopped. 0 means infinite loop on the state changing.
"action" is the action taken after the flow is stopped.
	"recover-state" means smart LED recover to the state before the color flow started.
	"keep-state" means smart LED stay at the state when the flow is stopped.
	"turn-off" means turn off the smart LED after the flow is stopped.
"flow_expression" is the expression of the state changing series.
	Each visible state changing is defined to be a flow tuple that contains 4
	elements: [duration, mode, value, brightness]. A flow expression is a series of flow tuples.
	So for above request example, it means: change CT to 2700K & maximum brightness
	gradually in 1000ms, then change color to red & 10% \brightness gradually in 500ms, then
	stay at this state for 5 seconds, then change CT to 5000K & minimum brightness gradually in
	500ms. After 4 changes reached, stopped the flow and power off the smart LED`,
	Example: `yeego start-cf bedroom 4 turn-off 1000,2,2700,100
yeego start-cf bedroom 4 recover-state 100,2,2700,100,50,1,255,10,500,7,0,0,500,2,5000,1`,
	Args: cobra.MinimumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(args[0])
		if err != nil {
			return err
		}

		count, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("The number of time to repeat the flow is mandatory")
		}

		var action int
		switch args[2] {
		case "recover-state":
			action = 0
		case "keep-state":
			action = 1
		case "turn-off":
			action = 2
		default:
			return errors.New("Action invalid. Please check help")
		}

		exp := strings.Split(args[3], ",")
		if len(exp) < 4 {
			return errors.New("Action invalid. Please check help")
		}

		for i := range exp {
			tmp, err := strconv.Atoi(exp[i])
			if err != nil {
				return errors.New("All the numbers of the flow should be integer: [duration, mode, value, brightness]")
			}

			// convert seconds to ms
			if (i % 4) == 0 {
				exp[i] = strconv.FormatInt(int64(tmp*1000), 10)
			}
		}

		_, err = light.StartCf(count, action, strings.Join(exp, ","))
		if err != nil {
			return err
		}

		fmt.Printf("%s color flow started\n", args[0])
		return nil

	},
}

var stopColorFlowCmd = &cobra.Command{
	Use:   "stop-cf [name/IP]",
	Short: "Stop a running color flow (cf)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		light, err := argToYeelight(args[0])
		if err != nil {
			return err
		}

		_, err = light.StopCf()
		if err != nil {
			return err
		}

		fmt.Printf("%s color flow stopped\n", args[0])
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
	rootCmd.AddCommand(colorFlowCmd)
	rootCmd.AddCommand(stopColorFlowCmd)
}
