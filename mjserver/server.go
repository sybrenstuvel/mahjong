package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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
