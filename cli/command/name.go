package command

import (
	"errors"
	"strings"

	"github.com/taubyte/dreamland/cli/flags"
	"github.com/urfave/cli/v2"
)

func NameWithDefault(c *cli.Command, def string) {
	flag := flags.Name
	flag.DefaultText = def
	flag.Value = def

	attachName(c, &flag)
}

func Name(c *cli.Command) {
	attachName(c, &flags.Name)
}

// adds the name flag to a command and prepends name parsing to the command's action
// also adds 'name' to the args usage
func attachName(c *cli.Command, flag cli.Flag) {
	c.Flags = append(c.Flags, flag)

	if len(c.ArgsUsage) == 0 {
		c.ArgsUsage = "name"
	} else {
		c.ArgsUsage = "name," + c.ArgsUsage
	}

	originalAction := c.Action
	c.Action = func(ctx *cli.Context) error {
		name, err := getName(ctx)
		if err != nil {
			return err
		}
		ctx.Set("name", name)

		// execute the original action at the end
		return originalAction(ctx)
	}

}

// when name is args0 or flag -n this method will get
// or return an error
func getName(c *cli.Context) (name string, err error) {
	name = c.Args().First()
	if len(name) == 0 {
		name = c.String("name")
		if len(name) == 0 {
			err = errors.New("Please provide a name")
		}
	} else if strings.HasPrefix(c.Args().Get(1), "-") {
		// If the args0 is present and there is an unparsed flag following it, return an error
		// (urfave cli does not support flags after arguments: https://github.com/urfave/cli/issues/427#issue-156378589)
		err = errors.New("Parse arguments failed: write [arguments] after -flags")
	}

	return
}
