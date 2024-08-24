package quote

import (
	"fmt"

	"motivator-tui/categories/gratitude"
	"motivator-tui/categories/grit"
	"motivator-tui/categories/perseverance"
)

// Quote the entry model
type Quote struct {
	CategoryName string
	Message      string
	Author       string
}

// QuoteRepository implements the Repository interface
type QuoteRepository struct{}

// GetQuoteByCategoryName implements the Repository interface
func (qr *QuoteRepository) GetQuoteByCategoryName(categoryName string) (Quote, error) {
	quoteFuncs := map[string]func() (string, string, error){
		"Grit":         grit.Grit,
		"Perseverance": perseverance.Perseverance,
		"Gratitude":    gratitude.Gratitude,
		// Add more categories and their functions here as needed
	}

	quoteFunc, ok := quoteFuncs[categoryName]
	if !ok {
		return Quote{}, fmt.Errorf("invalid category: %s", categoryName)
	}

	message, author, err := quoteFunc()
	if err != nil {
		return Quote{}, err
	}

	return Quote{CategoryName: categoryName, Message: message, Author: author}, nil
}
