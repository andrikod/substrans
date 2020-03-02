package subtranslate

// Function Handler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

const (
	maxUploadSize     = 1024 * 500 // 500K file limit
	defaultTargetLang = "en"
)

// HandleTranslate endpoint
func HandleTranslate(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetLanguage := r.FormValue("language")
	if targetLanguage == "" {
		targetLanguage = defaultTargetLang
	}
	//TODO: language.Parse()

	file, handler, err := r.FormFile("filename")
	if err != nil {
		http.Error(w, "filename is missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	writer, err := parseAndTranslate(scanner, targetLanguage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inputSrtFile := handler.Filename
	ext := path.Ext(inputSrtFile)
	outputSrtFile := fmt.Sprintf("%s-%s%s", inputSrtFile[0:len(inputSrtFile)-len(ext)], targetLanguage, ext)

	w.Header().Set("Content-Disposition", "attachment; filename="+outputSrtFile)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, bytes.NewReader(writer.Bytes()))
}
