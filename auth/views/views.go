package views

import (
	db "chatapp/db"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func getHtmlFilePath(relativefilePath string) string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	return filepath.Join(workingDir, relativefilePath)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func Login(w http.ResponseWriter, r *http.Request) {

}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles(getHtmlFilePath("/auth/templates/signup.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else if r.Method == http.MethodPost {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		type RequestData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var data RequestData
		err = json.Unmarshal(body, &data)

		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		stmt, err := db.DB.Prepare("INSERT INTO user (secondary_id, username, password) VALUES (?, ?, ?)")
		if err != nil {
			http.Error(w, "Database preparation error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(data.Username, data.Username, string(hashedPassword))
		if err != nil {
			http.Error(w, "Database execution error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
