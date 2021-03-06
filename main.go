package main

import (
	"blogo/internal/models/commands"
	"blogo/internal/parsing"
	"blogo/internal/server"
	"embed"
	_ "embed"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	//go:embed index.html
	IndexTemplateHTML string

	//go:embed 404.html
	NotFoundTemplateHTML string

	//go:embed post.html
	PostTemplateHTML string

	//go:embed content
	ContentDir embed.FS

	//go:embed config.yml
	ConfigFile string
)

func main() {

	// inject embedded files into subpackage
	server.IndexTemplateHTML = IndexTemplateHTML
	server.NotFoundTemplateHTML = NotFoundTemplateHTML
	server.PostTemplateHTML = PostTemplateHTML
	server.ContentDir = ContentDir

	commands.IndexTemplateHTML = IndexTemplateHTML
	commands.NotFoundTemplateHTML = NotFoundTemplateHTML
	commands.PostTemplateHTML = PostTemplateHTML
	commands.ContentDir = ContentDir
	commands.ConfigFile = ConfigFile

	parsing.IndexTemplateHTML = IndexTemplateHTML
	parsing.NotFoundTemplateHTML = NotFoundTemplateHTML

	app := &cli.App{
		Name:        "blogo",
		HelpName:    "blogo",
		Usage:       "A lightweight markdown renderer and live server",
		Description: "Create a blog from markdown in seconds 🚀",
		Commands: []*cli.Command{
			commands.ServeCommand,
			commands.ParseCommand,
			commands.InitCommand,
			commands.UpdateCommand,
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
