package main

import (
	"fmt"
	"log"
	"net/http"

	auth_views "chatapp/auth/views"
	chat_views "chatapp/chat"

	// "chatapp/chat/hub"
	"chatapp/db"
)

// Start app with `reflex -c reflex.conf`
// TODO: Check if user exist with same username
//

func main() {
	err := db.InitializeDB("your_database.db")

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	hub := chat_views.NewHub()
	go hub.Run()

	mux.HandleFunc("/", chat_views.Index)
	mux.HandleFunc("/login", auth_views.Login)
	mux.HandleFunc("/signup", auth_views.Signup)
	mux.HandleFunc("/request-connection", chat_views.RequestConnection)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		session, _ := auth_views.Store.Get(r, "auth-session")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// username := session.Values["username"].(string)
		// userID := session.Values["userID"].(int)
		// userSecondaryId := session.Values["userSecondaryId"].(string)

		// fmt.Printf("%v %v %v", username, userID, userSecondaryId)
		chat_views.ServeWs(hub, w, r)
	})

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
