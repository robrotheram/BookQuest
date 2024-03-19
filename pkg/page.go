package BookQuest

import (
	"BookQuest/pkg/stopwords"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Page struct {
	ID               string `json:"id"`
	Url              string `json:"url"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Text             string
	Tags             []string `json:"tags"`
	Meta             Meta     `json:"Meta"`
	Favourite        []string
	FavouritedByUser bool
}

type Meta struct {
	CreatedAt      time.Time `json:"created"`
	CreatedBy      string    `json:"created_by"`
	LastModified   time.Time `json:"modified"`
	LastModifiedBy string    `json:"modified_by"`
}

func NewPage(url, user string) Page {
	return Page{
		Url: url,
		ID:  urlToID(url),
		Meta: Meta{
			CreatedAt: time.Now(),
			CreatedBy: user,
		},
	}
}

func urlToID(url string) string {
	// Calculate SHA-256 hash of the URL
	hash := sha256.New()
	hash.Write([]byte(url))
	hashBytes := hash.Sum(nil)
	// Encode the hash bytes to a hexadecimal string
	id := hex.EncodeToString(hashBytes)
	return id
}

func (p *Page) ToggleFavoirte(profile Profile) {
	for i, v := range p.Favourite {
		if v == profile.Id {
			p.Favourite = append(p.Favourite[:i], p.Favourite[i+1:]...)
			return
		}
	}
	p.Favourite = append(p.Favourite, profile.Id)
}

func (p *Page) isFavoirte(profile Profile) bool {
	for _, v := range p.Favourite {
		if v == profile.Id {
			p.FavouritedByUser = true
			return true
		}
	}
	return false
}

func (p *Page) parseHTML(htmlCode string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader((htmlCode)))
	if err != nil {
		log.Fatal(err)
	}
	p.Name = doc.Find("title").Text()
	if description, exists := doc.Find("meta[name='description']").Attr("content"); exists {
		p.Description = description
	}
	visibleText := findVisibleText(doc)
	p.Text = visibleText
	p.Tags = getTags(visibleText, 6)
}

// Recursively find and return visible text
func findVisibleText(doc *goquery.Document) string {
	var visibleText strings.Builder
	// Iterate over top-level elements
	doc.Contents().Each(func(i int, s *goquery.Selection) {
		visibleText.WriteString(processElement(s))
	})
	return visibleText.String()
}

// Recursively process elements and return visible text
func processElement(s *goquery.Selection) string {
	var text string

	// Process current element
	if s.Nodes[0].Type == html.ElementNode {
		switch s.Nodes[0].Data {
		case "script", "style":
			// Ignore script and style elements
			return ""
		case "br":
			// Append newline for <br> tags
			return "\n"
		case "p", "div", "span", "h1", "h2", "h3", "h4", "h5", "h6":
			// Append a space for block-level elements
			text += " "
		}
	}

	// Recursively process child elements
	s.Contents().Each(func(i int, child *goquery.Selection) {
		text += processElement(child)
	})

	// Append text if it's visible
	if s.Nodes[0].Type == html.TextNode {
		trimmedText := strings.TrimSpace(s.Text())
		if trimmedText != "" {
			text += trimmedText + " "
		}
	}

	return text
}

func getTags(largeString string, numTags int) []string {
	//words := strings.Fields(largeString)

	words := strings.FieldsFunc(largeString, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	// Create a map to count word occurrences
	wordCount := make(map[string]int)
	for _, word := range words {
		if !stopwords.StopWordsMap[strings.ToLower(word)] {
			wordCount[word]++
		}
	}

	// Convert the map to a slice for sorting
	type kv struct {
		Key   string
		Value int
	}
	var sortedWords []kv
	for k, v := range wordCount {
		sortedWords = append(sortedWords, kv{k, v})
	}

	// Sort the slice by word count in descending order
	sort.Slice(sortedWords, func(i, j int) bool {
		return sortedWords[i].Value > sortedWords[j].Value
	})
	results := []string{}
	for i := 0; i < numTags && i < len(sortedWords); i++ {
		results = append(results, sortedWords[i].Key)
	}
	return results
}
