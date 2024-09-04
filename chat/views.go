package chat

import (
	auth_views "chatapp/auth/views"
	db "chatapp/db"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := auth_views.Store.Get(r, "auth-session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	username := session.Values["username"].(string)
	userID := session.Values["userID"].(uint)
	// userSecondaryId := session.Values["userSecondaryId"].(string)

	var userOne db.User
	db.DB.First(&userOne, userID)

	// Add users to the chat

	db.DB.Preload("Chats").First(&userOne, userOne.ID)

	data := struct {
		Username string
		Chats    []db.Chat
	}{
		Username: username,
		Chats:    userOne.Chats,
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

func ChatBox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatID"]

	data := struct {
		ChatID string
	}{
		ChatID: chatID,
	}


	tmpl, err := template.ParseFiles("chat/templates/chat.html")
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
	// session, _ := auth_views.Store.Get(r, "auth-session")

	// if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }

	// requestTo := r.FormValue("username")

	// var requestToUser db.User
	// // Find the user by username
	// err := db.DB.Where("username = ?", requestTo).First(&requestToUser).Error
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		// No user found
	// 		http.Error(w, `{"message": "User not found"}`, http.StatusNotFound)
	// 		return
	// 	}
	// 	// Other errors
	// 	http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
	// 	return
	// }

	// // requestFrom := session.Values["username"].(string)
	// requestFromID := session.Values["userID"].(uint)
	// // userSecondaryId := session.Values["userSecondaryId"].(string)

	// var requestFromUser db.User
	// err = db.DB.First(&requestFromUser, requestFromID).Error
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		// No user found
	// 		http.Error(w, `{"message": "User not found"}`, http.StatusNotFound)
	// 		return
	// 	}
	// 	// Other errors
	// 	http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
	// 	return
	// }
	// var connection db.ConnectionRequest

	// // Query the entire row and scan it into the struct
	// result := db.DB.Where(
	// 	"(send_by = ? AND send_to = ?) OR (send_by = ? AND send_to = ?)",
	// 	requestFromUser, requestToUser, requestToUser, requestFromUser,
	// ).Find(&connection)

	// // Check if any record was found
	// if result.Error != nil {
	// 	// Handle other types of errors (e.g., connection errors, SQL syntax errors)
	// 	http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
	// 	return
	// }

	// if result.RowsAffected == 0 {
	// 	// No matching record found
	// 	// Proceed with creating a new connection
	// 	secondary_id := uuid.New()

	// 	newChat := db.Chat{
	// 		SecondaryID:   secondary_id.String(),
	// 		IsPrivateChat: true,
	// 	}

	// 	err = db.DB.Create(&newChat).Error
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	secondary_id = uuid.New()
	// 	newConnection := db.ConnectionRequest{SendBy: requestFromUser, SendTo: requestToUser, SecondaryID: secondary_id.String()}

	// 	err := db.DB.Create(&newConnection).Error
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusCreated)
	// 	fmt.Fprintf(w, `{"secondary_id": "%s"}`, newConnection.SecondaryID)
	// } else {
	// 	// If no error occurred and a record was found, use the connection instance as needed
	// 	if !connection.IsAccepted.Valid {
	// 		// is_accepted is null
	// 		http.Error(w, `{"message": "Already request sent"}`, http.StatusConflict)
	// 		return
	// 	} else if connection.IsAccepted.Bool {
	// 		// is_accepted is true
	// 		http.Error(w, `{"message": "Already connection created"}`, http.StatusConflict)
	// 		return
	// 	} else {
	// 		// is_accepted is false
	// 		http.Error(w, `{"message": "Connection request declined"}`, http.StatusConflict)
	// 		return
	// 	}
	// }

}
