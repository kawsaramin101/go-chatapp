package main

import (
	"fmt"
	"log"
	"net/http"

	auth_views "chatapp/auth/views"
	"chatapp/db"
)

func main() {
	err := db.InitializeDB("your_database.db")

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", auth_views.Index)
	mux.HandleFunc("/login", auth_views.Login)
	mux.HandleFunc("/signup", auth_views.Signup)

	loggedMux := LoggingMiddleware(mux)

	fmt.Println("Listening to port 8000")

	err = http.ListenAndServe(":8000", loggedMux)

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

// TODO: Check if user exist with same username
