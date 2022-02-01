package main

import (
	"database/sql"
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
	Status   string
	Title    string
	Date     string
	Summary  string
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
	uniquepost := strings.Replace(r.URL.Path, "/posts/", "", -1)

	// declare an array to keep all comments
	var comments []Comment

	if conf.AllowComments {

		rows, err := db.Query("select id, name, comment from comments where uniquepost = ?", uniquepost)
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

	f := "posts/" + uniquepost + ".md"
	fileread, err := ioutil.ReadFile(f)
	if err != nil {
		handle404(w, r)
		return
	}

	lines := strings.Split(string(fileread), "\n")
	status := string(lines[0])
	title := string(lines[1])
	date := string(lines[2])
	summary := string(lines[3])
	body := strings.Join(lines[4:len(lines)], "\n")
	htmlBody := template.HTML(blackfriday.MarkdownCommon([]byte(body)))

	post := Post{status, title, date, summary, htmlBody, uniquepost, comments}
	t := template.New("post.html")
	t, _ = t.ParseFiles("post.html")
	t.Execute(w, post)

}
func handle404(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "404.html")
}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)

		// replace backticks for windows use
		file = strings.Replace(file, "posts\\", "", -1)
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		status := string(lines[0])
		title := string(lines[1])
		date := string(lines[2])
		summary := string(lines[3])
		body := strings.Join(lines[4:len(lines)], "\n")
		htmlBody := template.HTML(blackfriday.MarkdownCommon([]byte(body)))

		a = append(a, Post{status, title, date, summary, htmlBody, file, nil})
	}
	return a
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
