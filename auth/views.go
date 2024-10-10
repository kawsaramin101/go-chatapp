package auth

import (
	db "chatapp/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("your-256-bit-secret")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequestData struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

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
	var creds Credentials

	// Decode the JSON request body into the Credentials struct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Could not decode request body."})
		return
	}
	username := creds.Username
	password := creds.Password

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

	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString, "username": user.Username})

}

func Signup(w http.ResponseWriter, r *http.Request) {
	if db.DB == nil {
		http.Error(w, "Database connection not initialized", http.StatusInternalServerError)
		return
	}

	var requestData SignupRequestData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Could not decode request body."})
		return
	}

	if requestData.Password != requestData.ConfirmPassword {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Password and confirm password do not match."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	secondary_id := uuid.New()

	user := db.User{Username: requestData.Username, Password: string(hashedPassword), SecondaryID: secondary_id.String()}

	err = db.DB.Create(&user).Error
	if err != nil {
		http.Error(w, "Internel Server Error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created"})

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
