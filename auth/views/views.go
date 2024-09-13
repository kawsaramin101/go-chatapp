package views

import (
	db "chatapp/db"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

var Store = sessions.NewCookieStore([]byte("your-secret-key"))

var mySigningKey = []byte("your-256-bit-secret")

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func generateJWT(username string, userSecondadyId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":        username,
		"userSecondaryId": userSecondadyId,
		"exp":             time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Decode the JSON request body into the Credentials struct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Could not decode request body."})
		return
	}
	username := creds.Username
	password := creds.Password

	fmt.Println(username, password)

	var user db.User
	// Find the user by username
	err = db.DB.Where("username = ?", username).First(&user).Error
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

	tokenString, err := generateJWT(user.Username, user.SecondaryID)

	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})

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

type SignupRequestData struct {
	username        string
	password        string
	confirmPassword string
}

func signupPost(w http.ResponseWriter, r *http.Request) {

	if db.DB == nil {
		http.Error(w, "Database connection not initialized", http.StatusInternalServerError)
		return
	}

	var requestData SignupRequestData

	// Decode the JSON request body into the Credentials struct
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Could not decode request body."})
		return
	}

	if requestData.password != requestData.confirmPassword {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Password and confirm password do not match."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	secondary_id := uuid.New()

	user := db.User{Username: requestData.username, Password: string(hashedPassword), SecondaryID: secondary_id.String()}

	err = db.DB.Create(&user).Error
	if err != nil {
		http.Error(w, "Internel Server Error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User created"})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		tokenString := r.Header.Get("Authorization")

		// Validate the token and extract the username and secondary ID
		username, userSecondaryId, err := ValidateToken(tokenString)
		if err != nil {
			respondWithJSON(w, http.StatusForbidden, err.Error())
			return
		}

		// Store the extracted values in the request context
		ctx := context.WithValue(r.Context(), "username", username)
		ctx = context.WithValue(ctx, "userSecondaryId", userSecondaryId)

		// Pass the request to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateToken(tokenString string) (string, string, error) {
	if tokenString == "" {
		return "", "", fmt.Errorf("no token provided")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil // Replace with your actual key
	})

	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims (username and userSecondaryId)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, usernameOk := claims["username"].(string)
		userSecondaryId, secondaryOk := claims["userSecondaryId"].(string)
		if !usernameOk || !secondaryOk {
			return "", "", fmt.Errorf("invalid token claims")
		}
		return username, userSecondaryId, nil
	}

	return "", "", fmt.Errorf("invalid token claims")
}
