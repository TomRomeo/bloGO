package parsing

import (
	"blogo/internal/models"
	"blogo/internal/models/errors"
	"bytes"
	readingtime "github.com/begmaroman/reading-time"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

var (
	IndexTemplateHTML    string
	NotFoundTemplateHTML string
)

// function that returns a go struct post for a path
func GetPost(path string, comments []models.Comment, allowComments bool) (*models.Post, error) {

	fileName := strings.Replace(path, "posts/", "", -1)
	fileName = strings.Replace(fileName, ".md", "", -1)
	// also replace backticks for servers on Windows
	fileName = strings.Replace(fileName, "posts\\", "", -1)

	var post models.Post

	// default should be true, not false
	post.EnableComments = true

	readFile, _ := ioutil.ReadFile(path)
	// first, parse the frontmatter with yaml
	split := bytes.SplitN(readFile, []byte("---"), 2)

	// throw error if frontmatter was not found
	if len(split) < 2 {
		return nil, &errors.FrontMatterMissingError{FileName: fileName}
	}
	if err := yaml.Unmarshal(split[0], &post); err != nil {
		return nil, err
	}

	post.EnableComments = post.EnableComments && allowComments

	// estimate reading time
	estimation := readingtime.Estimate(string(split[1]))
	post.Ert = estimation.Text

	post.Body = template.HTML(blackfriday.MarkdownCommon(split[1]))
	post.File = fileName
	post.Comments = comments

	return &post, nil
}

func GetPosts(rootFolder string) []*models.Post {
	posts := []*models.Post{}
	files, _ := filepath.Glob(rootFolder + "*.md")

	for _, filePath := range files {

		post, err := GetPost(filePath, nil, false)
		if err != nil {
			if _, ok := err.(*errors.FrontMatterMissingError); ok {
				return posts
			} else {
				log.Fatal(err)
			}
		}

		// ignore wip files
		if !strings.Contains(filePath, "wip") && !(strings.ToLower(post.Status) == "wip") {
			posts = append(posts, post)
		}
	}
	return posts
}

func ParseIndex(w io.Writer, posts []*models.Post) error {

	t := template.New("index.html")
	t, err := t.Parse(IndexTemplateHTML)
	if err != nil {
		return err
	}
	if err := t.Execute(w, posts); err != nil {
		return err
	}
	return nil
}
func Parse404(w io.Writer) error {
	t := template.New("404.html")
	t, err := t.Parse(NotFoundTemplateHTML)
	if err != nil {
		return err
	}
	var l interface{}
	if err = t.Execute(w, l); err != nil {
		return err
	}
	return nil
}
