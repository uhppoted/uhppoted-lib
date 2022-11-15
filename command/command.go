package uhppoted

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Command interface {
	Name() string
	FlagSet() *flag.FlagSet
	Execute(...any) error
	Description() string
	Usage() string
	Help()
}

type CommandExt interface {
	ParseCmd(args ...string) error
}

func name(name string) string {
	return strings.Split(name, "|")[0]
}

func alt(name string, arg string) bool {
	tokens := strings.Split(name, "|")

	for _, t := range tokens {
		if t == arg {
			return true
		}
	}

	return false
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
		if ext, ok := cmd.(CommandExt); ok {
			return cmd, ext.ParseCmd(args...)
		} else if flagset := cmd.FlagSet(); flagset == nil {
			panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", name(cmd.Name()), cmd))
		} else {
			return cmd, flagset.Parse(args)
		}
		
		// flagset := cmd.FlagSet()
		// if flagset == nil {
		// 	panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", name(cmd.Name()), cmd))
		// }

		// return cmd, flagset.Parse(args)
	}

	return nil, nil
}
