package uhppoted

import (
	"flag"
	"fmt"
	"os"
)

type Command interface {
	Name() string
	FlagSet() *flag.FlagSet
	Execute(...interface{}) error
	Description() string
	Usage() string
	Help()
}

func Parse(cli []Command, run Command, help Command) (Command, error) {
	var cmd Command = run
	var args []string

	if flag.Parsed() {
		args = flag.Args()
	} else {
		args = os.Args[1:]
	}

	if len(args) > 0 {
		if args[0] == help.Name() {
			cmd = help
			args = args[1:]
		} else {
			for _, c := range cli {
				if args[0] == c.Name() {
					cmd = c
					args = args[1:]
					break
				}
			}
		}
	}

	if cmd != nil {
		flagset := cmd.FlagSet()
		if flagset == nil {
			panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", cmd.Name(), cmd))
		}

		return cmd, flagset.Parse(args)
	}

	return nil, nil
}
