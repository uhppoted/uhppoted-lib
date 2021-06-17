package uhppoted

import (
	"flag"
	"fmt"
)

type Version struct {
	Application string
	Version     string
}

func (cmd *Version) Name() string {
	return "version"
}

func (cmd *Version) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("version", flag.ExitOnError)
}

func (cmd *Version) Execute(args ...interface{}) error {
	fmt.Printf("%v\n", cmd.Version)

	return nil
}

func (cmd *Version) Description() string {
	return "Displays the current version"
}

func (cmd *Version) Usage() string {
	return ""
}

func (cmd *Version) Help() {
	fmt.Println()
	fmt.Printf("  Displays the %s version in the format v<major>.<minor>.<build> e.g. v1.00.10\n", cmd.Application)
	fmt.Println()
}
