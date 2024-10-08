package main

import (
	"fmt"
	"log"
	"net/http"

	auth_views "chatapp/auth"
	chat_views "chatapp/chat"

	// "chatapp/chat/hub"
	"chatapp/db"
	// "github.com/gorilla/mux"
)

// Start app with `reflex -c reflex.conf`
// TODO: Check if user exist with same username
// Number one rule in go
// If a channel is not working make sure you initialized it

func main() {
	err := db.InitializeDB("main.db")

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to database")
	}

	router := http.NewServeMux()

	hub := chat_views.NewHub()
	go hub.Run()

	router.HandleFunc("POST /", chat_views.Index)
	router.HandleFunc("POST /login", auth_views.Login)
	router.HandleFunc("POST /signup", auth_views.Signup)
	// router.HandleFunc("/request-connection", chat_views.RequestConnection)
	// router.HandleFunc("/chat/{chatID}", chat_views.ChatBox)

	router.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		chat_views.ServeWs(hub, w, r)
	})

	handlerWithMiddleware := LoggingMiddleware(CorsMiddleware(router))

	fmt.Println("Listening to port 8000")
	err = http.ListenAndServe(":8000", handlerWithMiddleware)

	if err != nil {
		fmt.Println(err)
	}

}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, URL: %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
