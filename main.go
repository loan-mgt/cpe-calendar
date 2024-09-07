package main

import (
	"cpe/calendar/handlers"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var tpl *template.Template

func init() {
	// Load environment variables from .env file
	godotenv.Load()
	// Parse templates
	tpl = template.Must(template.ParseFiles(filepath.Join("static", "index.html")))
}

func main() {
	r := mux.NewRouter()

	// Serve dynamic index page
	r.HandleFunc("/", serveIndex).Methods("GET")

	// Serve static files like JavaScript, CSS, images, etc.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Serve calendar.ics route
	r.HandleFunc("/your-cpe-calendar.ics", generate3IRCHandler).Methods("GET")

	// Use the router in the http server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// serveIndex renders the index.html Go template with environment variables
func serveIndex(w http.ResponseWriter, r *http.Request) {
	publicKey, _ := os.ReadFile(filepath.Join("static", "public.pem"))
	publicKey = []byte(strings.ReplaceAll(string(publicKey), "\n", ""))
	separator := os.Getenv("SEPARATOR")

	data := struct {
		PublicKey string
		Separator string
	}{
		PublicKey: string(publicKey),
		Separator: separator,
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// generate3IRCHandler is a wrapper around generateICSHandler that uses a specific filename
func generate3IRCHandler(w http.ResponseWriter, r *http.Request) {
	// Call generateICSHandler with the specific filename and calendar name
	handlers.GenerateICSHandler(w, r, "3irc_calendar.ics", "3IRC Calendar")
}
