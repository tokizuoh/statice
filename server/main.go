package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			statusCode := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(statusCode), statusCode)
		}
		fmt.Fprintln(w, "Hello, World")
	})
	http.ListenAndServe(":8080", nil)
}
