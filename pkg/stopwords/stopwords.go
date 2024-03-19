package stopwords

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed stop_words_english.json
var stopWordsEnglishJSON []byte

//go:embed stop_words_html.json
var stopWordsHTMLJSON []byte

var StopWordsMap map[string]bool

func init() {
	stopWordsEnglish := []string{}
	if err := json.Unmarshal(stopWordsEnglishJSON, &stopWordsEnglish); err != nil {
		fmt.Println("Error loading stop words:", err)
		return
	}
	stopWordsHTML := []string{}
	if err := json.Unmarshal(stopWordsHTMLJSON, &stopWordsHTML); err != nil {
		fmt.Println("Error loading stop words:", err)
		return
	}

	stopWords := append(stopWordsEnglish, stopWordsHTML...)
	// Create a map for fast lookup
	StopWordsMap = make(map[string]bool)
	for _, word := range stopWords {
		StopWordsMap[strings.ToLower(word)] = true
	}
}
