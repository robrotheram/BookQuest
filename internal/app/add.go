package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func decodeData(encodedData string) (string, error) {
	// Step 1: URL-decode the input data to handle any percent-encoded characters
	// decodedURL, err := url.QueryUnescape(encodedData)
	// if err != nil {
	// 	return "", err
	// }
	fmt.Println(encodedData)
	// Step 2: Decode the base64 string
	decodedBytes, err := base64.RawStdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	// Step 3: Convert the decoded bytes to a string
	decodedString := string(decodedBytes)

	// Step 4: Replace '+' with space (unescape step)
	decodedString = strings.ReplaceAll(decodedString, "+", " ")

	// Step 5: Percent-decode the string again to handle any additional encoding
	decodedString, err = url.QueryUnescape(decodedString)
	if err != nil {
		return "", err
	}

	return decodedString, nil
}

func (app *App) HandleAdd(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Url  string `json:"url"`
		Html string `json:"html"`
	}
	var request req
	json.NewDecoder(r.Body).Decode(&request)

	fmt.Println(request)
}
