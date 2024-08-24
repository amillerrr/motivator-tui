package quote

import (
	"fmt"
)

// FormattedOutputFromQuotes format all quotes as a single string in reverse chronological order
func FormattedOutputFromQuote(quote Quote) []byte {
	output := fmt.Sprintf("Category: %s\n\n %s\n", quote.CategoryName, quote.Message)
	return []byte(output)
}

// FormatCategoryName return the entry details as a formatted string
func FormatCategoryName(quote Quote) string {
	return fmt.Sprintf("Category: %s\n\n", quote.CategoryName)
}

// FormatQuote return the entry details as a formatted string
func FormatQuote(quote Quote) string {
	return fmt.Sprintf("%s\n", quote.Message)
}

// FormatAuthor return the entry details as a formatted string
func FormatAuthor(quote Quote) string {
	return fmt.Sprintf("-- %s\n", quote.Author)
}
