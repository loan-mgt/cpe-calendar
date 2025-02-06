package main

import (
	"cpe/calendar/handlers"
	"cpe/calendar/metrics"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var tpl *template.Template

func init() {
	// Load environment variables from .env file
	godotenv.Load()
	// Parse templates
	tpl = template.Must(template.ParseFiles(filepath.Join("static", "index.html")))

	prometheus.Register(metrics.TotalRequests)
	prometheus.Register(metrics.ResponseStatus)
	prometheus.Register(metrics.HttpDuration)
}

func main() {
	r := mux.NewRouter()
	r.Use(metrics.PrometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	// Serve dynamic index page
	r.HandleFunc("/", serveIndex).Methods("GET")

	// Serve static files like JavaScript, CSS, images, etc.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Serve calendar.ics route
	r.HandleFunc("/your-cpe-calendar.ics", handlers.GenerateICSHandler).Methods("GET")

	//validate route
	r.HandleFunc("/validate", handlers.ValidateHandler).Methods("GET")

	// check app health
	r.HandleFunc("/health", handlers.Health).Methods("GET")

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
