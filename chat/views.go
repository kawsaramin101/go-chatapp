package chat

import (
	auth_views "chatapp/auth/views"
	db "chatapp/db"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := auth_views.Store.Get(r, "auth-session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	username := session.Values["username"].(string)
	// userID := session.Values["userID"].(int)
	// userSecondaryId := session.Values["userSecondaryId"].(string)

	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl, err := template.ParseFiles("chat/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RequestConnection(w http.ResponseWriter, r *http.Request) {
	session, _ := auth_views.Store.Get(r, "auth-session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Println("run")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	requestTo := r.FormValue("username")

	var requestToUser db.User
	// Find the user by username
	err := db.DB.Where("username = ?", requestTo).First(&requestToUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No user found
			http.Error(w, `{"message": "User not found"}`, http.StatusNotFound)
			return
		}
		// Other errors
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// requestFrom := session.Values["username"].(string)
	requestFromID := session.Values["userID"].(uint)
	// userSecondaryId := session.Values["userSecondaryId"].(string)
	//
	var connection db.Connection

	// Query the entire row and scan it into the struct
	result := db.DB.Table("connections").
		Where("(send_by = ? AND send_to = ?) OR (send_by = ? AND send_to = ?)", requestFromID, requestToUser.ID, requestToUser.ID, requestFromID).
		Scan(&connection)

	// Check if any record was found
	if result.Error != nil {
		// Handle other types of errors (e.g., connection errors, SQL syntax errors)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		// No matching record found
		// Proceed with creating a new connection
		secondary_id := uuid.New()
		newConnection := db.Connection{SendBy: uint(requestFromID), SendTo: requestToUser.ID, SecondaryID: secondary_id.String()}

		err := db.DB.Create(&newConnection).Error
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"secondary_id": "%s"}`, newConnection.SecondaryID)
	} else {
		// If no error occurred and a record was found, use the connection instance as needed
		fmt.Println(connection)
		if !connection.IsAccepted.Valid {
			// is_accepted is null
			fmt.Println("Run 1")
			http.Error(w, `{"message": "Already request sent"}`, http.StatusConflict)
			return
		} else if connection.IsAccepted.Bool {
			fmt.Println("Run 2")

			// is_accepted is true
			http.Error(w, `{"message": "Already connection created"}`, http.StatusConflict)
			return
		} else {
			fmt.Println("Run 3")

			// is_accepted is false
			http.Error(w, `{"message": "Connection request declined"}`, http.StatusConflict)
			return
		}
	}

}
