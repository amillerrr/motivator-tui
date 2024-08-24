package main

import (
	"log"

	"motivator-tui/categories"
	"motivator-tui/quote"
	"motivator-tui/tui"
)

func main() {
	cr := categories.NewInMemoryRepository()
	qr := &quote.QuoteRepository{}
	categories, err := cr.GetAllCategories()
	if err != nil {
		log.Fatal(err)
	}
	if len(categories) < 1 {
		log.Fatal("Error loading categories")
	} else {
		tui.StartTea(*cr, *qr)
	}
}
