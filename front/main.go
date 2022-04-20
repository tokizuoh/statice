package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("login.html")
		if err != nil {
			fmt.Fprintf(w, err.Error(), nil)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			fmt.Fprintf(w, err.Error(), nil)
			return
		}
	})
	http.ListenAndServe(":8081", nil)
}
