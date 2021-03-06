package main

// Standalone version without any server/cloud dependencies

// Run as:
// GOOGLE_APPLICATION_CREDENTIALS=[PATH_TO_JSON] go run main.go

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func main() {
	inputSrtFile := "./subs/test.srt"
	targetLanguage := "en"

	ext := path.Ext(inputSrtFile)
	outputSrtFile := fmt.Sprintf(inputSrtFile[0:len(inputSrtFile)-len(ext)], targetLanguage, ext)

	// input file
	srtFile, err := os.Open(inputSrtFile)
	if err != nil {
		panic(err)
	}
	defer srtFile.Close()

	scanner := bufio.NewScanner(srtFile)
	scanner.Split(bufio.ScanLines)

	// output file
	transSrtFile, err := os.Create(outputSrtFile)
	if err != nil {
		panic(err)
	}

	defer transSrtFile.Close()
	defer transSrtFile.Sync()

	for scanner.Scan() {
		if err := translateBlock(transSrtFile, scanner, targetLanguage); err != nil {
			panic(err)
		}
	}

}

func translateBlock(transSrtFile *os.File, scanner *bufio.Scanner, targetLanguage string) error {
	// seq
	transSrtFile.WriteString(scanner.Text() + "\n")

	// timestamp
	if hasNext := scanner.Scan(); !hasNext {
		return errors.New("Unexpected end of file")
	}
	transSrtFile.WriteString(scanner.Text() + "\n")

	// text
	var subTexts []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line == "\n" {
			break
		}
		subTexts = append(subTexts, line)
	}

	trans, err := translateText(targetLanguage, subTexts)
	if err != nil {
		return err
	}

	for i := range trans {
		transSrtFile.WriteString(trans[i].Text + "\n")
	}

	transSrtFile.WriteString("\n")
	return nil
}

func translateText(targetLanguage string, text []string) ([]translate.Translation, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return nil, fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, text, lang, nil)
	if err != nil {
		return nil, fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp, nil
}
