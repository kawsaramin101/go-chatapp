package chat

import (
	auth_views "chatapp/auth/views"
	"fmt"
	"html/template"
	"net/http"
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

	fmt.Println("run")
	username := r.FormValue("username")

	fmt.Println(username)

}
