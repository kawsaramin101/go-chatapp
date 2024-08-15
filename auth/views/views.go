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

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Store = sessions.NewCookieStore([]byte("your-secret-key"))

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginGet(w, r)
	} else if r.Method == http.MethodPost {
		loginPost(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session, err := Store.Get(r, "auth-session")

		if err != nil {
			http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = -1

		// Save the session to delete it
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Failed to delete session", http.StatusInternalServerError)
			return
		}

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged out successfully"))

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

	var user db.User
	// Find the user by username
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"message": "Username or password didn't match"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, `{"message": "Username or password didn't match"}`, http.StatusUnauthorized)
		return
	}

	session, err := Store.Get(r, "auth-session")

	if err != nil {
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	session.Values["userID"] = user.ID
	session.Values["userSecondaryId"] = user.SecondaryID
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
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

	if db.DB == nil {
		http.Error(w, "Database connection not initialized", http.StatusInternalServerError)
		return
	}

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

	user := db.User{Username: data.Username, Password: string(hashedPassword), SecondaryID: secondary_id.String()}

	err = db.DB.Create(&user).Error
	if err != nil {
		http.Error(w, "Internel Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
