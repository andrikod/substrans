package main

// Local Server for testing the Function
// Run as:
// GOOGLE_APPLICATION_CREDENTIALS=[PATH_TO_JSON] go run *.go

import (
	"log"
	"net/http"

	subtranslate "subtrans.com/subtrans/handler"
)

func main() {
	http.HandleFunc("/HandleTranslate", subtranslate.HandleTranslate)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
