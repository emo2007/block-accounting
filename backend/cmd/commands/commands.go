package commands

import "github.com/urfave/cli/v2"

var commandPool []*cli.Command

func Register(c *cli.Command) {
	commandPool = append(commandPool, c)
}

func Commands() []*cli.Command {
	return commandPool
}
