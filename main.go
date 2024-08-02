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

	http.HandleFunc("/", auth_views.Index)
	http.HandleFunc("/login", auth_views.Login)
	http.HandleFunc("/signup", auth_views.Signup)

	fmt.Println("Listening to port 8000")

	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println(err)
	}

}
