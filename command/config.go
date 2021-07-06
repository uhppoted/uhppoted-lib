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
}

func (cmd *Config) Name() string {
	return "config"
}

func (cmd *Config) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("config", flag.ExitOnError)
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
	if err := dump(cmd.Config); err != nil {
		return err
	}

	return nil
}

func dump(path string) error {
	fmt.Println()
	fmt.Printf("   ... displaying configuration information from '%s'\n", path)

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

	fmt.Println()
	fmt.Printf("%s\n", s.String())
	fmt.Println()

	return nil
}
