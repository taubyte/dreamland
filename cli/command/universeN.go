package command

import (
	"errors"
	"log"
	"strings"

	"github.com/taubyte/dreamland/cli/common"
	"github.com/taubyte/dreamland/cli/flags"
	"github.com/urfave/cli/v2"
)

// Attaches universe argument parsing to a command.
// `pos` is position of the universe argument, expected to be 0 or 1.
func universeAtPos(c *cli.Command, pos int) {
	c.Flags = append(c.Flags, &flags.Universe)

	// Add arg name to command's ArgUsage based on arg position
	if len(c.ArgsUsage) == 0 {
		if pos == 1 {
			// command has no other arguments (with a usage string) but it should
			log.Fatal("universe expected to be second argument")
		}
		c.ArgsUsage = "universe"
	} else {
		if pos == 0 {
			c.ArgsUsage = "universe," + c.ArgsUsage
		} else {
			c.ArgsUsage += ", universe"
		}
	}

	originalAction := c.Action
	c.Action = func(ctx *cli.Context) error {
		universe, err := getUniverseAtPos(ctx, pos)
		if err != nil {
			return err
		}
		ctx.Set("universe", universe)
		// execute the original action at the end
		return originalAction(ctx)
	}
}

// get the universe from the flag or the argument at position `pos`
func getUniverseAtPos(c *cli.Context, pos int) (universe string, err error) {
	universe = c.String("universe")
	if universe != common.DefaultUniverseName {
		return
	}

	universeArg := c.Args().Get(pos)
	if len(universeArg) == 0 {
		return
	}

	universe = universeArg
	if strings.HasPrefix(universe, "-") {
		err = errors.New("Parse arguments failed: write [arguments] after -flags")
	}
	return
}
