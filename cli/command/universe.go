package command

import "github.com/urfave/cli/v2"

func Universe(c *cli.Command) {
	universeAtPos(c, 1)
}
