package uhppoted

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

type HelpV struct {
	service string
	cli     []Command
	run     Command
}

func NewHelp(service string, cli []Command, run Command) *HelpV {
	return &HelpV{
		service: service,
		cli:     cli,
		run:     run,
	}
}

func (h *HelpV) Name() string {
	return "help"
}

func (h *HelpV) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("help", flag.ExitOnError)
}

func (h *HelpV) Description() string {
	return "Displays the help information"
}

func (h *HelpV) Usage() string {
	return ""
}

func (h *HelpV) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s help <command>\n", h.service)
	fmt.Println()
	fmt.Println("  Commands:")

	for _, c := range h.cli {
		fmt.Printf("    %-13s %s\n", c.FlagSet().Name(), c.Description())
	}

	fmt.Println("    help          Displays this message")
}

func (h *HelpV) Execute(args ...interface{}) error {
	if len(os.Args) > 2 {
		if os.Args[2] == "commands" {
			h.helpCommands()
			return nil
		}

		if os.Args[2] == h.Name() {
			h.Help()
			return nil
		}

		for _, c := range h.cli {
			if os.Args[2] == c.Name() {
				c.Help()
				return nil
			}
		}

		fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
	} else {
		h.usage()
	}

	return nil
}

func (h *HelpV) usage() {
	fmt.Println()
	fmt.Printf("  Usage: %s [options] <command> [parameters]\n", h.service)
	fmt.Println()

	fmt.Println("  Commands:")
	for _, c := range h.cli {
		fmt.Printf("    %-13s %s\n", c.FlagSet().Name(), c.Description())
	}
	fmt.Printf("    help          Displays this message. For help on a specific command use '%s help <command>'\n", h.service)

	var options bytes.Buffer
	var count = 0

	fmt.Fprintln(&options)
	fmt.Fprintln(&options, "  Options:")
	flag.VisitAll(func(f *flag.Flag) {
		count++
		fmt.Fprintf(&options, "    --%-7s %s\n", f.Name, f.Usage)
	})

	if count > 0 {
		fmt.Printf(string(options.Bytes()))
	}

	fmt.Println()

	if h.run != nil {
		fmt.Printf("  Defaults to '%s'.\n", h.run.Name())
		fmt.Println()
		fmt.Printf("   '%s' options:\n", h.run.Name())

		h.run.FlagSet().VisitAll(func(f *flag.Flag) {
			fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
		})

		fmt.Println()
	}
}

func (h *HelpV) helpCommands() {
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range h.cli {
		fmt.Printf("     %-16s %s\n", c.FlagSet().Name(), c.Description())
	}

	fmt.Println()

	if h.run != nil {
		fmt.Printf("     Defaults to '%s'.\n", h.run.Name())
		fmt.Println()
	}
}
