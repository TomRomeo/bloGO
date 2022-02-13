package commands

import (
	"github.com/urfave/cli/v2"
	"gomarkdownblog/internal/parsing"
	"html/template"
	"log"
	"os"
)

var UpdateCommand = &cli.Command{
	Name:        "update",
	Aliases:     []string{"u"},
	Usage:       "Update markdown posts to static html files",
	Description: "Parse the markdown posts to html",
	ArgsUsage:   "",
	Action: func(c *cli.Context) error {
		posts := parsing.GetPosts("posts/")

		if len(posts) == 0 {
			log.Println("No posts found in subfolder posts/")
			log.Println("Make sure to call this command from the bloGO rootfolder")
			return nil
		}

		outDir := "postsHTML/"

		// create folders if not exist
		os.MkdirAll(outDir+"postsHTML/", os.ModePerm)

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

		log.Println("Rebuilding Index.html...")
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

		log.Println("rendered sucessfully")
		return nil
	},
}
