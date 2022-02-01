package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"gopkg.in/yaml.v3"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/russross/blackfriday"
)

type Post struct {
	Status   string `yaml:"status"`
	Title    string `yaml:"title"`
	Date     string `yaml:"date"`
	Summary  string `yaml:"summary"`
	Body     template.HTML
	File     string
	Comments []Comment
}

type Comment struct {
	Name, Comment string
}

var (
	db   *sql.DB
	conf Conf
)

type Conf struct {
	AllowComments bool `yaml:"allow_comments"`
}

type FrontMatterMissingError struct {
	fileName string
}

func (e *FrontMatterMissingError) Error() string {
	return fmt.Sprintf("Missing YML Frontmatter: %s", e.fileName)
}

func parseConfig() {
	yConfig, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Could not open config.yml: %s", err)
	}
	if err := yaml.Unmarshal(yConfig, &conf); err != nil {
		log.Fatalf("Could not parse config.yml: %s", err)
	}
}

func init() {
	parseConfig()

	//only connect to sqlite if we want comment functionality
	if conf.AllowComments {

		// you do not have to open the db connection on every request
		// it can be done once at the start of the app
		db, err := sql.Open("mysql", "username:password(localhost:3306)/databasename")
		if err != nil {
			log.Fatal(err)
		}

		// Open doesn't open a connection. Validate DSN data:
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handlerequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case http.MethodGet:

		if r.URL.Path == "/" {
			handleIndex(w, r)
			return
		}

		// 404 handler
		handle404(w, r)

	case http.MethodPost:
		handlePostComment(w, r)
	}

}

func handleIndex(w http.ResponseWriter, r *http.Request) {

	posts := getPosts()
	t := template.New("index.html")
	t, err := t.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Execute(w, posts); err != nil {
		log.Fatal(err)
	}
}

func handlePostComment(w http.ResponseWriter, r *http.Request) {

	if conf.AllowComments {
		uniquepost := r.FormValue("uniquepost")
		namein := r.FormValue("name")
		commentin := r.FormValue("comment")

		_, err := db.Exec(
			"INSERT INTO comments (uniquepost, name, comment) VALUES (?, ?, ?)",
			uniquepost,
			namein,
			commentin,
		)
		if err != nil {
			log.Fatal(err)
		}
		//when done inserting comment redirect back to this page
		http.Redirect(w, r, r.URL.Path, 301)
		return
	}
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
	postName := strings.Replace(r.URL.Path, "/posts/", "", -1)

	// declare an array to keep all comments
	var comments []Comment

	if conf.AllowComments {

		rows, err := db.Query("select id, name, comment from comments where postName = ?", postName)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var name, comment string
			err := rows.Scan(&id, &name, &comment)
			if err != nil {
				log.Fatal(err)
			}
			//append the comment into the array when done
			comments = append(comments, Comment{name, comment})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	postMarkdownPath := "posts/" + postName + ".md"

	// check if file exist -> else return 404
	_, err := ioutil.ReadFile(postMarkdownPath)
	if err != nil {
		handle404(w, r)
		return
	}

	post, err := getPost(postMarkdownPath, comments)
	if err != nil {
		if _, ok := err.(*FrontMatterMissingError); ok {
			handle404(w, r)
			log.Println(fmt.Sprintf("%s seems to be missing yml attributes, the request has been dropped", postMarkdownPath))
		} else {
			log.Fatal(err)
		}
	}

	t := template.New("post.html")
	t, _ = t.ParseFiles("post.html")
	t.Execute(w, post)

}
func handle404(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "404.html")
}

func getPosts() []Post {
	posts := []Post{}
	files, _ := filepath.Glob("posts/*")

	for _, filePath := range files {
		post, err := getPost(filePath, nil)
		if err != nil {
			if _, ok := err.(*FrontMatterMissingError); ok {
				return posts
			} else {
				log.Fatal(err)
			}
		}
		posts = append(posts, *post)
	}
	return posts
}

// function that returns a go struct post for a path
func getPost(path string, comments []Comment) (*Post, error) {

	fileName := strings.Replace(path, "posts/", "", -1)
	fileName = strings.Replace(fileName, ".md", "", -1)
	// also replace backticks for servers on Windows
	fileName = strings.Replace(fileName, "posts\\", "", -1)

	var post Post

	readFile, _ := ioutil.ReadFile(path)
	// first, parse the frontmatter with yaml
	split := bytes.SplitN(readFile, []byte("---"), 2)

	// throw error if frontmatter was not found
	if len(split) < 2 {
		return nil, &FrontMatterMissingError{fileName: fileName}
	}
	if err := yaml.Unmarshal(split[0], &post); err != nil {
		return nil, err
	}

	post.Body = template.HTML(blackfriday.MarkdownCommon(split[1]))
	post.File = fileName
	post.Comments = comments

	return &post, nil
}

func main() {
	http.HandleFunc("/", handlerequest)
	http.HandleFunc("/posts/", handlePosts)
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./content/css"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./content/js"))))

	log.Println("Blog Server deployed on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}

}
