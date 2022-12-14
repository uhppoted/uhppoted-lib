package uhppoted

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/uhppoted/uhppoted-lib/config"
)

type Config struct {
	Application string
	Config      string
	debug       bool
}

func (cmd *Config) Name() string {
	return "config"
}

func (cmd *Config) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("config", flag.ExitOnError)

	flagset.BoolVar(&cmd.debug, "debug", cmd.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (cmd *Config) Description() string {
	return fmt.Sprintf("Displays all the configuration information for %s", cmd.Application)
}

func (cmd *Config) Usage() string {
	return ""
}

func (cmd *Config) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s config\n", cmd.Application)
	fmt.Println()
	fmt.Printf("    Displays all the configuration information for %s\n", cmd.Application)
	fmt.Println()
}

func (cmd *Config) Execute(args ...interface{}) error {
	if err := dump(cmd.Config, cmd.debug); err != nil {
		return err
	}

	return nil
}

func dump(path string, debug bool) error {
	if debug {
		fmt.Println()
		fmt.Printf("   ... displaying configuration information from '%s'\n", path)
		fmt.Println()
	}

	cfg := config.NewConfig()
	if f, err := os.Open(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err := cfg.Read(f)
		f.Close()
		if err != nil {
			return err
		}
	}

	var s strings.Builder

	if err := cfg.Write(&s); err != nil {
		return err
	}

	fmt.Printf("%s\n", s.String())
	fmt.Println()

	return nil
}
