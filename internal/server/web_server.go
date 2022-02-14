package server

import (
	"blogo/internal/models"
	"blogo/internal/models/errors"
	"blogo/internal/parsing"
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db         *sql.DB
	conf       models.Conf
	rootFolder string

	IndexTemplateHTML string

	NotFoundTemplateHTML string

	PostTemplateHTML string

	ContentDir embed.FS
)

func parseConfig() {
	yConfig, err := ioutil.ReadFile(rootFolder + "config.yml")
	if err != nil {
		log.Fatalf("Could not open config.yml: %s", err)
	}
	if err := yaml.Unmarshal(yConfig, &conf); err != nil {
		log.Fatalf("Could not parse config.yml: %s", err)
	}
}

func InitServer() {
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

	posts := parsing.GetPosts(rootFolder + "/posts/")
	t := template.New("index.html")
	t, err := t.Parse(IndexTemplateHTML)
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
	var comments []models.Comment

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
			comments = append(comments, models.Comment{name, comment})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	postMarkdownPath := rootFolder + "posts/" + postName + ".md"

	// check if file exist -> else return 404
	_, err := ioutil.ReadFile(postMarkdownPath)
	if err != nil {
		handle404(w, r)
		return
	}

	post, err := parsing.GetPost(postMarkdownPath, comments, conf.AllowComments)
	if err != nil {
		if _, ok := err.(*errors.FrontMatterMissingError); ok {
			handle404(w, r)
			log.Println(fmt.Sprintf("%s seems to be missing yml attributes, the request has been dropped", postMarkdownPath))
		} else {
			log.Fatal(err)
		}
	}

	t := template.New("post.html")
	t, _ = t.Parse(PostTemplateHTML)
	t.Execute(w, post)

}
func handle404(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "404.html")
	w.Write([]byte(NotFoundTemplateHTML))
}

func ServeBlogoServer(folderpath string) {
	rootFolder = folderpath

	http.HandleFunc("/", handlerequest)
	http.HandleFunc("/posts/", handlePosts)
	fss := fs.FS(ContentDir)
	cssDir, err := fs.Sub(fss, "content/css")
	if err != nil {
		log.Fatal(err)
	}
	jsDir, err := fs.Sub(fss, "content/js")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.FS(cssDir))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.FS(jsDir))))

	srv := &http.Server{Addr: ":8000", Handler: nil}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func(srv *http.Server) {
		log.Println("Blog Server deployed on :8000")
		if err := srv.ListenAndServe(); err != nil {
			c <- syscall.SIGINT
		}

	}(srv)

	<-c

	log.Println("Shutting down webserver")
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Server could not shut down: %s", err)
	}
}
