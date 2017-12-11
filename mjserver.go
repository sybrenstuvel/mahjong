package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	stdlog "log"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const serverVersion = "0.1-dev"

var cliArgs struct {
	version bool
	verbose bool
	debug   bool
}

func parseCliArgs() {
	flag.BoolVar(&cliArgs.version, "version", false, "Shows the application version, then exits.")
	flag.BoolVar(&cliArgs.verbose, "verbose", false, "Enable info-level logging.")
	flag.BoolVar(&cliArgs.debug, "debug", false, "Enable debug-level logging.")
	flag.Parse()
}

func configLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Only log the warning severity or above by default.
	level := log.WarnLevel
	if cliArgs.debug {
		level = log.DebugLevel
	} else if cliArgs.verbose {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	stdlog.SetOutput(log.StandardLogger().Writer())
}

func logStartup() {
	level := log.GetLevel()
	defer log.SetLevel(level)

	log.SetLevel(log.InfoLevel)
	log.WithFields(log.Fields{
		"version": serverVersion,
	}).Info("Starting Mahjong Server")
}

func main() {
	parseCliArgs()
	if cliArgs.version {
		fmt.Println(serverVersion)
		return
	}

	configLogging()
	logStartup()

	// Set some more or less sensible limits & timeouts.
	http.DefaultTransport = &http.Transport{
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   3 * time.Second,
		IdleConnTimeout:       15 * time.Minute,
		ResponseHeaderTimeout: 15 * time.Second,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/score", scoreHand)

	listen := ":8080"
	log.Println("Listening on", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Welcome!<br>")
	fmt.Fprintln(w, "<a href='/score'>Score your hand</a>")
}

func scoreHand(w http.ResponseWriter, r *http.Request) {
	//hand := r.URL.RawQuery
	score := 1 // score.Score(hand)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, score)
}

func todoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoID)
}
