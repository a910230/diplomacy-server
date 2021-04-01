package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/solver", solverHandler)
	mux.Handle("/game", fileHandler("game/game.html"))
	mux.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("res"))))
	mux.Handle("/403", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "403 forbidden", 403) }))

	// Start a web server.
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/diplomacy.guru/cert.pem", "/etc/letsencrypt/live/diplomacy.guru/privkey.pem", mux)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func connectDB() {
	_, err := sql.Open("mysql", "username:password@tcp(diplomacy.c6kpxx0eowrf.us-east-2.rds.amazonaws.com:3306)?charset=utf8")
	check(err)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func solverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var orders [][]string
	json.NewDecoder(r.Body).Decode(&orders)
	if orders != nil {
		fmt.Fprintln(w, orders)
	} else {
		fmt.Fprintln(w, "Hello, World")
	}
}

func fileHandler(filename string) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
	return http.HandlerFunc(f)
}
