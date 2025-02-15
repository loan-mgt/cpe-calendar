package main

import (
	"cpe/calendar/handlers"
	"cpe/calendar/logger"
	"cpe/calendar/metrics"
	"html/template"
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
	err := godotenv.Load()
	if err != nil {
		// Log error and exit if environment variables can't be loaded
		logger.Log.Warn().Err(err).Msg("Error loading .env file")
	}

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

	// Start HTTP server and log any errors that occur
	logger.Log.Info().Msg("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		// Log any errors that occur while starting the server
		logger.Log.Fatal().Err(err).Msg("Error starting server")
	}
}

// serveIndex renders the index.html Go template with environment variables
func serveIndex(w http.ResponseWriter, r *http.Request) {
	publicKey, err := os.ReadFile(filepath.Join("static", "public.pem"))
	if err != nil {
		// Log error if public key can't be read
		logger.Log.Error().Err(err).Msg("Error reading public.pem")
		http.Error(w, "Error reading public key", http.StatusInternalServerError)
		return
	}

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
		// Log error if template rendering fails
		logger.Log.Error().Err(err).Msg("Error rendering template")
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
