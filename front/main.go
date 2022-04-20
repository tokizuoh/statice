package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

// key: sessionID, value: id (user)
var mapSessionID = map[string]string{}

func login(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		if len(mapSessionID[cookie.Value]) > 0 {
			http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
			return
		}
	}

	t, err := template.ParseFiles("template/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error(), nil)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		fmt.Fprintf(w, err.Error(), nil)
		return
	}
}

func validateLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		statusCode := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	if err := r.ParseForm(); err != nil {
		statusCode := http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	id := r.Form.Get("id")
	password := r.Form.Get("password")

	if !(id == "test-id" && password == "test-password") {
		statusCode := http.StatusUnauthorized
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	sessionID := uuid.NewString()
	mapSessionID[sessionID] = id
	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	if err != nil {
		statusCode := http.StatusUnauthorized
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	t, err := template.ParseFiles("template/home.html")
	if err != nil {
		fmt.Fprintf(w, err.Error(), nil)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		fmt.Fprintf(w, err.Error(), nil)
		return
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/validateLogin", validateLogin)
	http.HandleFunc("/home", home)
	http.ListenAndServe(":8081", nil)
}
