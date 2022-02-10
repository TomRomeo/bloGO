package parsing

import (
	"bytes"
	readingtime "github.com/begmaroman/reading-time"
	"github.com/russross/blackfriday"
	"gomarkdownblog/internal/models"
	"gomarkdownblog/internal/models/errors"
	"gopkg.in/yaml.v3"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
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

func GetPosts(rootFolder string) []models.Post {
	posts := []models.Post{}
	files, _ := filepath.Glob(rootFolder + "posts/*.md")

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
			posts = append(posts, *post)
		}
	}
	return posts
}
