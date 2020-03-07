package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/julienrbrt/yeego/light/yeelight"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var (
	// Lights are the yeelight present on LAN
	lights []yeelight.Yeelight

	// error messages
	errNotFoundLight    = errors.New("Light not found")
	errYeelightNotFound = errors.New("No Yeelight found. Run `yeego discover` to find lights on your network")

	// configuration file name
	confName = ".yeego"

	// timeout used for discover and effects
	timeout time.Duration
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yeego",
	Short: "Control your Yeelight bulb with Yeego",
	Long: `Yeego is a tool written in Go which permits to control
your Yeelight bulbs in your LAN directly from your terminal.`,
	Example: `yeego discover
yeego on bedroom`,
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// if no light do not write anything
			return nil
		}

		light, err := argToYeelight(args[0])
		if err != nil {
			// if error do not write anything
			return nil
		}

		// Get light properties
		err = light.GetProp()
		if err != nil {
			return nil
		}

		for i := range lights {
			if lights[i].Location == light.Location {
				lights[i] = *light
			}
		}

		err = writeConfig()
		return err
	},
}

// argToYeelight searches a yeelight in the preloaded lights or build a new light if an IP is provided
func argToYeelight(addr string) (*yeelight.Yeelight, error) {
	for _, light := range lights {
		if light.Name == strings.ToLower(addr) || strings.Split(light.Location, ":")[0] == addr {
			return &light, nil
		}
	}

	// parse the value as IP, permits to verify if the user enters an IP
	ip := net.ParseIP(addr)
	if ip != nil {
		return &yeelight.Yeelight{Location: addr + ":" + yeelight.Port}, nil
	}

	return &yeelight.Yeelight{}, errYeelightNotFound
}

// Write the yeego config file
func writeConfig() error {
	// no light found, do not write any config file
	if len(lights) == 0 {
		return nil
	}

	lightsJSON, err := json.Marshal(lights)
	if err != nil {
		return err
	}

	//get program path
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}

	// write config file
	err = ioutil.WriteFile(path.Join(folderPath, confName), lightsJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//get program path
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	// read configuration file - no error check
	file, _ := ioutil.ReadFile(path.Join(folderPath, confName))
	json.Unmarshal(file, &lights)
}
