package uipoc

//package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/ppacer/core/timeutils"
)

//go:embed views/*.html
var viewsFS embed.FS

//go:embed assets/* css/*
var staticFS embed.FS

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseFS(viewsFS, "views/*.html")),
	}
}

type DagRunRow struct {
	Id       int
	DagId    string
	Schedule string
	Status   string
}

type DagRunRows []DagRunRow

func newDagRunRows() DagRunRows {
	dagId := "sample_dag"
	ts := time.Now()
	return DagRunRows{
		DagRunRow{1, dagId, timeutils.ToString(ts), "RUNNING"},
		DagRunRow{2, dagId, timeutils.ToString(ts.Add(-100 * time.Minute)), "RUNNING"},
		DagRunRow{3, dagId, timeutils.ToString(ts.Add(-300 * time.Minute)), "RUNNING"},
		DagRunRow{4, dagId, timeutils.ToString(ts.Add(-900 * time.Minute)), "RUNNING"},
		DagRunRow{5, dagId, timeutils.ToString(ts.Add(-1300 * time.Minute)), "SUCCESS"},
	}
}

func randomStatus() string {
	statuses := []string{"RUNNING", "SUCCESS", "FAILED"}
	return statuses[rand.Intn(len(statuses))]
}

func UIServer() http.Handler {
	mux := http.NewServeMux()
	dagRuns := newDagRunRows()
	templates := newTemplates()

	// Serve static files from embedded filesystem
	mux.Handle("/assets/", http.FileServer(http.FS(staticFS)))
	mux.Handle("/css/", http.FileServer(http.FS(staticFS)))

	// Handler for the root URL
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.Render(w, "index", dagRuns); err != nil {
			msg := fmt.Sprintf("Failed to render template: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/dagruns", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.Render(w, "dagruns", dagRuns); err != nil {
			msg := fmt.Sprintf("Failed to render template: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/random-status", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(randomStatus()))
	})

	mux.HandleFunc("/rand", func(w http.ResponseWriter, _ *http.Request) {
		rnd := time.Duration(rand.Intn(800) + 200)
		time.Sleep(rnd * time.Millisecond)
		fmt.Fprintf(w, "%d", rand.Intn(100)+8)
	})

	return mux
}

func main() {
	port := ":8181"
	log.Println("Starting server on ", port)
	if err := http.ListenAndServe(port, UIServer()); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
