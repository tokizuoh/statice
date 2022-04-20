package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templatelogin.html")
	if err != nil {
		fmt.Fprintf(w, err.Error(), nil)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		fmt.Fprintf(w, err.Error(), nil)
		return
	}
}

func validLogin(w http.ResponseWriter, r *http.Request) {
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

	// TODO: セッションID発行してCookieに入れて送る
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	// TODO: 直接 /home にリクエストできてしまうので解決する
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
	http.HandleFunc("/validLogin", validLogin)
	http.HandleFunc("/home", home)
	http.ListenAndServe(":8081", nil)
}
