package commands

import (
	"embed"
	"github.com/urfave/cli/v2"
	"gomarkdownblog/internal/models"
	"gomarkdownblog/internal/parsing"
	"gomarkdownblog/internal/util"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var (
	IndexTemplateHTML string

	NotFoundTemplateHTML string

	PostTemplateHTML string

	ContentDir embed.FS
)

var ParseCommand = &cli.Command{
	Name:        "parse",
	Aliases:     []string{"p"},
	Usage:       "Parse markdown posts to static html files",
	Description: "Parse the markdown files to html in <output folder>",
	ArgsUsage:   "<markdown files> <output folder>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "Where to generate index.html",
		},
	},
	Action: func(c *cli.Context) error {

		mdFiles := c.Args().Slice()[:c.Args().Len()-1]
		outDir := c.Args().Slice()[c.Args().Len()-1]
		outDir = util.ReformatPath(outDir)

		var posts []*models.Post

		for _, fileGlobs := range mdFiles {
			files, _ := filepath.Glob(fileGlobs)

			for _, f := range files {
				p, err := parsing.GetPost(f, nil, false)
				if err != nil {
					continue
				}
				posts = append(posts, p)
			}
		}

		// create folders if not exist
		os.MkdirAll(outDir, os.ModePerm)

		log.Println("parsing posts..")
		for i := 0; i < len(posts); i++ {

			log.Printf("post: %s..", posts[i].File)
			f, err := os.OpenFile(outDir+posts[i].File+".html", os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer f.Close()

			t := template.New("post.html")
			t, _ = t.Parse(PostTemplateHTML)
			t.Execute(f, posts[i])
			f.Close()
		}
		log.Println("Parsed all posts successfully")

		// if also generating index.html..
		if c.String("index") != "" {

			indexPath := c.String("index")

			log.Println("Creating Index.html...")
			f, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer f.Close()
			t := template.New("index.html")
			t, err = t.Parse(IndexTemplateHTML)
			if err != nil {
				return err
			}
			if err := t.Execute(f, posts); err != nil {
				return err
			}
		}

		//
		//log.Println("Creating 404.html...")
		//fileGlobs, err = os.OpenFile(outDir+"404.html", os.O_WRONLY|os.O_CREATE, 0600)
		//if err != nil {
		//	return err
		//}
		//defer fileGlobs.Close()
		//t = template.New("404.html")
		//t, err = t.Parse(NotFoundTemplateHTML)
		//if err != nil {
		//	return err
		//}
		//if err := t.Execute(fileGlobs, posts); err != nil {
		//	return err
		//}
		//log.Println("Copying static files...")
		//err = fs.WalkDir(ContentDir, ".", func(path string, d fs.DirEntry, err error) error {
		//	if d.IsDir() {
		//		return nil
		//	}
		//	log.Printf("%s...", path)
		//	fileGlobs, err := os.OpenFile(outDir+path, os.O_WRONLY|os.O_CREATE, 0600)
		//	if err != nil {
		//		return err
		//	}
		//	defer fileGlobs.Close()
		//	return nil
		//})
		//if err != nil {
		//	return err
		//}

		log.Println("rendered sucessfully")
		return nil
	},
}
