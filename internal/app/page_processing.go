package app

import (
	"BookQuest/internal/stopwords"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// TFIDF represents a term frequency-inverse document frequency score for a term
type TFIDF struct {
	Term  string
	Score float64
}

func Process(paragraph string) []string {
	paragraph, _ = url.QueryUnescape(paragraph)
	// Tokenize the paragraph
	tokens := tokenize(paragraph)

	// Remove stopwords and perform stemming
	tokens = preprocess(tokens)

	// Calculate TF-IDF scores
	tfidfScores := calculateTFIDF(tokens)

	// Sort terms by TF-IDF score
	sortedTerms := sortByScore(tfidfScores)
	terms := []string{}

	// Print top significant terms
	topTerms := 5
	// fmt.Printf("Top %d significant terms:\n", topTerms)
	for i, term := range sortedTerms {
		if i >= topTerms {
			break
		}
		terms = append(terms, term.Term)
	}
	return terms
}

// Tokenize breaks a string into tokens based on whitespace and punctuation
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})
}

// Preprocess removes stopwords and performs stemming
func preprocess(tokens []string) []string {
	var result []string
	for _, token := range tokens {
		// Check if the token is a stopword

		if _, ok := stopwords.StopWordsMap[strings.ToLower(token)]; !ok {
			//remove numbers
			if _, err := strconv.Atoi(token); err != nil {
				//not a number
				result = append(result, token)
			}
			// Stem the token
			// stemmedToken := english.Stem(token, true)

		}
	}
	return result
}

// calculateTFIDF calculates the TF-IDF scores for terms
func calculateTFIDF(tokens []string) map[string]float64 {
	// Dummy TF-IDF calculation
	tfidfScores := make(map[string]float64)
	for _, token := range tokens {
		tfidfScores[token]++
	}
	//Normalize TF-IDF scores
	totalTokens := float64(len(tokens))
	for term, count := range tfidfScores {
		tfidfScores[term] = count / totalTokens
	}
	return tfidfScores
}

// sortByScore sorts terms by TF-IDF score
func sortByScore(tfidfScores map[string]float64) []TFIDF {
	var sortedTerms []TFIDF
	for term, score := range tfidfScores {
		sortedTerms = append(sortedTerms, TFIDF{Term: term, Score: score})
	}
	// Sort terms by TF-IDF score (descending order)
	sort.Slice(sortedTerms, func(i, j int) bool {
		return sortedTerms[i].Score > sortedTerms[j].Score
	})
	return sortedTerms
}
