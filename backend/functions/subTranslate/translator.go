package subtranslate

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

var (
	projectID = os.Getenv("GCP_PROJECT")
	client    *translate.Client
)

func init() {
	var err error
	client, err = translate.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed tranlsate.NewClient: %v", err)
	}
}

const (
	Seq = iota
	Timestamp
	Text
	BlankLine
)

func parseAndTranslate(scanner *bufio.Scanner, targetLanguage string) (bytes.Buffer, error) {
	var writer bytes.Buffer

	currentType := BlankLine
	var subTexts []string
	for scanner.Scan() {
		line := scanner.Text()

		switch currentType {
		case Seq:
			currentType = Timestamp
			writer.WriteString(line + "\n")
		case BlankLine:
			currentType = Seq
			writer.WriteString(line + "\n")
		case Timestamp:
			currentType = Text
			subTexts = append(subTexts, line)
		case Text:
			if line == "" {
				err := appendTranslations(&writer, subTexts, targetLanguage)
				if err != nil {
					log.Fatalln(err)
				}

				currentType = BlankLine
				subTexts = nil
			} else {
				currentType = Text
				subTexts = append(subTexts, line)
			}
		}
	}

	// missing ending new line
	if currentType == Text && len(subTexts) > 0 {
		err := appendTranslations(&writer, subTexts, targetLanguage)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return writer, nil
}

func appendTranslations(writer *bytes.Buffer, subTexts []string, targetLanguage string) error {
	trans, err := translateText(targetLanguage, subTexts)
	if err != nil {
		return err
	}

	for i := range subTexts {
		writer.WriteString(trans[i].Text + "\n")
	}
	writer.WriteString("\n")

	return nil
}

func translateText(targetLanguage string, text []string) ([]translate.Translation, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return nil, fmt.Errorf("language.Parse: %v", err)
	}

	resp, err := client.Translate(context.Background(), text, lang, nil)
	if err != nil {
		return nil, fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("Translate returned empty response to text: %s", text)
	}

	return resp, nil

}
