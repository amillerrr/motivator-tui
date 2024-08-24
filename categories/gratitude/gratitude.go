package gratitude

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Quote struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

func Gratitude() (string, string, error) {
	// Seed the random number generator for different results each time
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Open the JSON file
	file, err := os.Open("./categories/gratitude/gratitude.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", "", err
	}
	defer file.Close()

	// Decode the JSON data into a slice of Quote structs
	var quotes []Quote
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&quotes)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", "", err
	}

	// Check if there are any quotes
	if len(quotes) == 0 {
		fmt.Println("No quotes found in the file.")
		return "", "", err
	}

	// Randomly select a quote
	randomIndex := r.Intn(len(quotes))
	selectedQuote := quotes[randomIndex]

	// Print the selected quote and author
	quote := fmt.Sprintf("%s\n\n", selectedQuote.Quote)
	author := fmt.Sprintf("%s\n", selectedQuote.Author)
	return quote, author, nil
}
