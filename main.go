package main

import (
	"auth-app/internal/platform"
	"fmt"
	"net/http"
)


func main() {
	cfg:=	platform.LoadConfig()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(cfg)
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe(":8000", nil)
}
