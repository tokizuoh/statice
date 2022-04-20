package main

import (
	"log"
	"net/http"
)

func validLogin(id string, password string) bool {
	return id == "test-id" && password == "test-password"
}

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: HEADリクエストは受け付ける必要ある？
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

		ok := validLogin(id, password)
		if ok != true {
			statusCode := http.StatusUnauthorized
			http.Error(w, http.StatusText(statusCode), statusCode)
			return
		}

		log.Println(ok)

		// TODO: リダイレクト
	})
	http.ListenAndServe(":8080", nil)
}
