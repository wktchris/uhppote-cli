package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/uhppoted/uhppoted-api/config"
)

var SetPCControlCmd = SetPCControl{}

// Command implementation for set-pc-control to enable or disable
// access control from a remote application.
type SetPCControl struct {
}

// Gets the device ID and enable/disable value from the command line
// and sends a set-pc-control to the designated controller.
func (c *SetPCControl) Execute(ctx Context) error {
	deviceID, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	enable := true
	if len(flag.Args()) > 2 {
		v := strings.ToLower(flag.Arg(2))
		if matches, _ := regexp.MatchString("true|false", v); !matches {
			return fmt.Errorf("Invalid command - expected 'true' or 'false', got '%v'", flag.Arg(2))
		}

		if v == "false" {
			enable = false
		}
	}

	succeeded, err := ctx.uhppote.SetPCControl(deviceID, enable)
	if err != nil {
		return err
	}

	if !succeeded {
		if enable {
			return fmt.Errorf("Failed to enable 'set pc control' for %v", deviceID)
		} else {
			return fmt.Errorf("Failed to disable 'set pc control' for %v", deviceID)
		}
	}

	fmt.Printf("%v %v\n", deviceID, enable)

	return nil
}

// Returns the 'set-pc-control' command string for the CLI interface.
func (c *SetPCControl) CLI() string {
	return "set-pc-control"
}

// Returns the 'set-pc-control' command summary for the CLI interface.
func (c *SetPCControl) Description() string {
	return "Enables or disables remote access control"
}

// Returns the 'set-pc-control' command parameters for the CLI interface.
func (c *SetPCControl) Usage() string {
	return "<serial number> <enable>"
}

// Outputs the 'set-pc-control' command help for the CLI interface.
func (c *SetPCControl) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-pc-control <serial number> <enable>")
	fmt.Println()
	fmt.Println(" Enables or disables remote access control")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <enable>         (optional) 'true' or 'false'. Defaults to 'true'")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config  File path for the 'conf' file containing the controller configuration")
	fmt.Printf("              (defaults to %s)\n", config.DefaultConfig)
	fmt.Println("    --debug   Displays internal information for diagnosing errors")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --debug --config .config set-pc-control 12345678 true")
	fmt.Println()
}

// Returns false - configuration is useful but optional.
func (c *SetPCControl) RequiresConfig() bool {
	return false
}
