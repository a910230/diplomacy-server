package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Set router
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("res"))))
	mux.HandleFunc("/solver", solver)
	mux.Handle("/game", fileHandler("game/game.html"))
	mux.HandleFunc("/403", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "403 forbidden", 403) })
	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "404 not found", 404) })

	// Start a web server.
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/diplomacy.guru/cert.pem", "/etc/letsencrypt/live/diplomacy.guru/privkey.pem", mux)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

func index(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) == 1 {
		http.ServeFile(w, r, "index.html")
	} else {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
	}
}

func solver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var message [][]string
	json.NewDecoder(r.Body).Decode(&message)
	info := message[0]
	var orders []Order
	for i := 1; i < len(message); i++ {
		var objs [3]string
		copy(objs[:], message[i][1:])
		order := Order{unit: message[i][0], objs: objs}
		orders = append(orders, order)
	}
	_ = info

}

func fileHandler(filename string) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
	return http.HandlerFunc(f)
}

// var info = [infoObj.getAttribute("user"), infoObj.getAttribute("role"), infoObj.getAttribute("gameid"), infoObj.getAttribute("turn")];
