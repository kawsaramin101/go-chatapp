package chat

import (
	"fmt"
	"net/http"

	auth_views "chatapp/auth/views"
)

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := auth_views.Store.Get(r, "auth-session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	username := session.Values["username"].(string)
	userID := session.Values["userID"].(int)
	userSecondaryId := session.Values["userSecondaryId"].(string)

	fmt.Fprintf(w, "%v %v %v", username, userID, userSecondaryId)
}
