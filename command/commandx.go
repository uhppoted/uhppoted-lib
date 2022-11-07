package uhppoted

import (
	"flag"
	"fmt"
	"os"
)

/** EXPERIMENTAL **/

type CommandX interface {
	Name() string
	Configure(map[string]any) 
	FlagSet() *flag.FlagSet
	Execute(...interface{}) error
	Description() string
	Usage() string
	Help()
}

func ParseX(cli []CommandX, run CommandX, help CommandX) (CommandX, error) {
	var cmd CommandX = run
	var args []string

	if flag.Parsed() {
		args = flag.Args()
	} else {
		args = os.Args[1:]
	}

	if len(args) > 0 {
		if alt(help.Name(), args[0]) {
			cmd = help
			args = args[1:]
		} else {
			for _, c := range cli {
				if alt(c.Name(), args[0]) {
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
			panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", name(cmd.Name()), cmd))
		}

		return cmd, flagset.Parse(args)
	}

	return nil, nil
}
