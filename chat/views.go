package chat

import (
	auth_views "chatapp/auth/views"
	db "chatapp/db"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

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
	// Find the user by username using GORM
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
	requestFromID := session.Values["userID"].(int)
	// userSecondaryId := session.Values["userSecondaryId"].(string)

	var isAccepted sql.NullBool
	// Check the status of the existing connection using GORM
	err = db.DB.Table("connections").Select("is_accepted").
		Where("(send_by = ? AND send_to = ?) OR (send_by = ? AND send_to = ?)", requestFromID, requestToUser.ID, requestToUser.ID, requestFromID).
		Scan(&isAccepted).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		// Handle unexpected errors
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check the status of the existing connection
	if err == nil {
		if !isAccepted.Valid {
			// is_accepted is null
			http.Error(w, `{"message": "Already request sent"}`, http.StatusConflict)
			return
		} else if isAccepted.Bool {
			// is_accepted is true
			http.Error(w, `{"message": "Already connection created"}`, http.StatusConflict)
			return
		} else {
			// is_accepted is false
			http.Error(w, `{"message": "Connection request declined"}`, http.StatusConflict)
			return
		}
	}

	// At this point, no existing connection was found, proceed with creating a new connection
	// Add logic here to create the new connection

}
