package views

import (
	db "chatapp/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginGet(w, r)
	} else if r.Method == http.MethodPost {
		loginPost(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		signupGet(w, r)
	} else if r.Method == http.MethodPost {
		signupPost(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getHtmlFilePath(relativefilePath string) string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	return filepath.Join(workingDir, relativefilePath)
}

func loginGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(getHtmlFilePath("/auth/templates/login.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(username, password)

	var hashedPassword string
	query := "SELECT password FROM user WHERE username = ?"
	err := db.DB.QueryRow(query, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Username or password didn't match", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		http.Error(w, "Username or password didn't match", http.StatusUnauthorized)
		return
	}
	session, _ := store.Get(r, "session-name")

	session.Values["authenticated"] = true
    session.Values["username"] = username
    session.Values["userID"] = 12345 // Example user ID
    session.Save(r, w)
    http.Redirect(w, r, "/", http.StatusFound)

}

func signupGet(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(getHtmlFilePath("/auth/templates/signup.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func signupPost(w http.ResponseWriter, r *http.Request) {

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

	secondary_id := uuid.New()

	stmt, err := db.DB.Prepare("INSERT INTO user (secondary_id, username, password) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Database preparation error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(secondary_id, data.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Database execution error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
