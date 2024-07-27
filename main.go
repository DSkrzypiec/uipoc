package main

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/ppacer/core/timeutils"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
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

	// Serve static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

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
		fmt.Fprintf(w, "%d", rand.Intn(100)+8)
	})

	return mux
}

/*
func main() {
	port := ":8181"
	templates := newTemplates()

	// Serve static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Handler for the root URL
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.Render(w, "index", dagRuns); err != nil {
			msg := fmt.Sprintf("Failed to render template: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/dagruns", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.Render(w, "dagruns", dagRuns); err != nil {
			msg := fmt.Sprintf("Failed to render template: %s", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/random-status", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(randomStatus()))
	})

	http.HandleFunc("/rand", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "%d", rand.Intn(100)+8)
	})

	log.Println("Starting server on ", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
*/
