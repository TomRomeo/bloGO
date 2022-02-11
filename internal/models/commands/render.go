package commands

import (
	"embed"
	"github.com/urfave/cli/v2"
	"gomarkdownblog/internal/parsing"
	"gomarkdownblog/internal/util"
	"html/template"
	"io/fs"
	"log"
	"os"
)

var (
	IndexTemplateHTML string

	NotFoundTemplateHTML string

	PostTemplateHTML string

	ContentDir embed.FS
)

var RenderCommand = &cli.Command{
	Name:        "render",
	Aliases:     []string{"r"},
	Usage:       "parse markdown to static html files",
	Description: "Render the markdown from <input folder> to html in <output folder>",
	ArgsUsage:   "<input folder> <output folder>",
	Action: func(c *cli.Context) error {
		inDir := c.Args().Get(0)
		outDir := c.Args().Get(1)

		inDir = util.ReformatPath(inDir)
		outDir = util.ReformatPath(outDir)

		posts := parsing.GetPosts(inDir + "posts/")

		// create folders if not exist
		os.MkdirAll(outDir+"posts/", os.ModePerm)
		os.MkdirAll(outDir+"content/css", os.ModePerm)
		os.MkdirAll(outDir+"content/js", os.ModePerm)

		log.Println("parsing posts..")
		for i := 0; i < len(posts); i++ {

			log.Printf("post: %s..", posts[i].File)
			f, err := os.OpenFile(outDir+"posts/"+posts[i].File+".html", os.O_WRONLY|os.O_CREATE, 0600)
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
		if err := t.Execute(f, posts); err != nil {
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
		if err := t.Execute(f, posts); err != nil {
			return err
		}
		log.Println("Copying static files...")
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

		log.Println("rendered sucessfully")
		return nil
	},
}
