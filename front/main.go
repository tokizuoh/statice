package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var ctx context.Context
var rdb *redis.Client

func login(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// クッキーにセッションIDが含まれてないと判断
		t, err := template.ParseFiles("template/login.html")
		if err != nil {
			fmt.Fprintf(w, err.Error(), nil)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			fmt.Fprintf(w, err.Error(), nil)
			return
		}
	} else {
		_, err := rdb.Get(ctx, cookie.Value).Result()
		if err != nil {
			log.Println(err.Error())
			t, err := template.ParseFiles("template/login.html")
			if err != nil {
				fmt.Fprintf(w, err.Error(), nil)
				return
			}

			if err := t.Execute(w, nil); err != nil {
				fmt.Fprintf(w, err.Error(), nil)
				return
			}
			return
		}
		http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
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
	if err := rdb.Set(ctx, sessionID, id, 0).Err(); err != nil {
		log.Println(err.Error())
		statusCode := http.StatusFound
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
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

func init() {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // TODO: .env で管理する
		DB:       0,
	})
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/validateLogin", validateLogin)
	http.HandleFunc("/home", home)
	http.ListenAndServe(":8081", nil)
}
