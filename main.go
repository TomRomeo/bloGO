package main

import (
	"github.com/urfave/cli/v2"
	"gomarkdownblog/internal/models/commands"
	"log"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:        "blogo",
		HelpName:    "blogo",
		Usage:       "A lightweight markdown renderer and live server",
		Description: "Create a blog from markdown in seconds 🚀",
		Commands: []*cli.Command{
			commands.ServeCommand,
			commands.RenderCommand,
		},
		CommandNotFound: nil,
		OnUsageError:    nil,
		Compiled:        time.Time{},
		Authors: []*cli.Author{
			{
				Name:  "Tom Doil",
				Email: "Tom.Romeo.Doil@gmail.com",
			},
		},
		Copyright: "",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
