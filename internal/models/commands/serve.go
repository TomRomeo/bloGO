package commands

import (
	"github.com/urfave/cli/v2"
	"gomarkdownblog/internal/server"
)

var ServeCommand = &cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Usage:   "Serve markdown files live",
	Description: `This starts a local dev server where visiting a link parses the markdown on the fly.
Great for development, not very great for performance.`,
	ArgsUsage: "<input folder>",
	Action: func(c *cli.Context) error {
		dir := c.Args().Get(0)

		// add trailing slash if missing
		if dir == "" {
			dir = "./"
		} else if dir[len(dir)-1:] != "/" {
			dir += "/"
		}
		server.ServeBlogoServer(dir)
		return nil
	},
}
