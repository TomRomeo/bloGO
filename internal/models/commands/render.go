package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var RenderCommand = &cli.Command{
	Name:        "render",
	Aliases:     []string{"r"},
	Usage:       "parse markdown to static html files",
	Description: "Render the markdown from <input folder> to html in <output folder>",
	ArgsUsage:   "<input folder> <output folder>",
	Action: func(c *cli.Context) error {
		fmt.Println("render command")
		return nil
	},
}
