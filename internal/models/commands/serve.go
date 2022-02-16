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
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "live",
			Aliases: []string{"l"},
			Usage:   "live parses the markdown files on the fly",
		},
	},
	Action: func(c *cli.Context) error {
		dir := c.Args().Get(0)
		live := c.Bool("live")

		dir = util.ReformatPath(dir)

		server.InitServer()
		server.ServeBlogoServer(dir, live)
		return nil
	},
}
