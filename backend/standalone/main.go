package main

// Standalone version without any server/cloud dependencies

// Run as:
// GOOGLE_APPLICATION_CREDENTIALS=[PATH_TO_JSON] go run main.go

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

const (
	Seq = iota
	Timestamp
	Text
	BlankLine
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

	currentType := BlankLine
	var subTexts []string
	for scanner.Scan() {
		line := scanner.Text()

		switch currentType {
		case Seq:
			currentType = Timestamp
			transSrtFile.WriteString(line + "\n")
		case BlankLine:
			currentType = Seq
			transSrtFile.WriteString(line + "\n")
		case Timestamp:
			currentType = Text
			subTexts = append(subTexts, line)
		case Text:
			if line == "" {
				currentType = BlankLine

				trans, err := translateText(targetLanguage, subTexts)
				if err != nil {
					panic(err)
				}

				for i := range subTexts {
					// fmt.Printf("%s -->  %s \n", subTexts[i], trans[i].Text)
					transSrtFile.WriteString(trans[i].Text + "\n")
				}
				transSrtFile.WriteString(line + "\n")

				subTexts = nil
			} else {
				currentType = Text
				subTexts = append(subTexts, line)
			}
		}
	}

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
