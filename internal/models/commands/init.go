package commands

import (
	"github.com/urfave/cli/v2"
	"gomarkdownblog/internal/models"
	"gomarkdownblog/internal/util"
	"html/template"
	"io/fs"
	"log"
	"os"
)

var InitCommand = &cli.Command{
	Name:        "init",
	Aliases:     []string{"i"},
	Usage:       "Create a markdown blog scaffold!",
	Description: "Create a blog scaffold in <output folder>",
	ArgsUsage:   "<output folder>",
	Action: func(c *cli.Context) error {
		outDir := c.Args().Get(0)
		outDir = util.ReformatPath(outDir)

		// create folders if not exist
		os.MkdirAll(outDir+"posts/", os.ModePerm)
		os.MkdirAll(outDir+"content/css", os.ModePerm)
		os.MkdirAll(outDir+"content/js", os.ModePerm)

		log.Println("Creating Index.html...")
		f, err := os.OpenFile(outDir+"index.html", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		t := template.New("index.html")
		t, err = t.Parse(IndexTemplateHTML)
		if err != nil {
			return err
		}
		if err := t.Execute(f, []*models.Post{}); err != nil {
			return err
		}

		log.Println("Creating 404.html...")
		f, err = os.OpenFile(outDir+"404.html", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		t = template.New("404.html")
		t, err = t.Parse(NotFoundTemplateHTML)
		if err != nil {
			return err
		}
		if err := t.Execute(f, []*models.Post{}); err != nil {
			return err
		}
		log.Println("Creating static files...")
		err = fs.WalkDir(ContentDir, ".", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			log.Printf("%s...", path)
			f, err := os.OpenFile(outDir+path, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer f.Close()
			return nil
		})
		if err != nil {
			return err
		}

		log.Println("initialized sucessfully")
		return nil
	},
}
