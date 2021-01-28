package main

import (
	"net/http"

	"github.com/acharyab15/test-with-go/app"
)

func main() {
	http.ListenAndServe(":3000", &app.Server{})
}
