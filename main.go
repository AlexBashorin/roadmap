package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	mysql "github.com/go-sql-driver/mysql"
)

type Page struct {
	Title string
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DB query
func dbConn() (db *sql.DB) {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return db
}

// PageVariables is used to pass data to templates
type PageVariables struct {
	Title string
}

func main() {
	// Define your routes
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/projects", ProjPage)

	// Parse the navigation template
	navTmpl, err := template.ParseFiles("templates/nav.html")
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register the navigation template
	http.HandleFunc("/nav", func(w http.ResponseWriter, r *http.Request) {
		err := navTmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the web server
	http.ListenAndServe(":90", nil)
}

// HomePage handles requests to the home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	// Create a PageVariables instance with the title for the template
	pageVariables := PageVariables{
		Title: "Home Page",
	}

	// Parse the template file
	tmpl, err := template.ParseFiles("templates/index.html", "templates/nav.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the PageVariables
	err = tmpl.ExecuteTemplate(w, "index.html", pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AboutPage handles requests to the about page
func ProjPage(w http.ResponseWriter, r *http.Request) {
	// Create a PageVariables instance with the title for the template
	pageVariables := PageVariables{
		Title: "About Page",
	}

	// Parse the template file
	tmpl, err := template.ParseFiles("templates/projects.html", "templates/nav.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the PageVariables
	err = tmpl.ExecuteTemplate(w, "projects.html", pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
