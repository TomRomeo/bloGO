package commands

import (
	"blogo/internal/server"
	"blogo/internal/util"
	"github.com/urfave/cli/v2"
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

		dir = util.ReformatPath(dir)

		server.InitServer()
		server.ServeBlogoServer(dir)
		return nil
	},
}
