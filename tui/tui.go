package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"motivator-tui/categories"
	"motivator-tui/quote"
	"motivator-tui/tui/constants"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(cr categories.InMemoryRepository, qr quote.QuoteRepository) error {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	constants.Cr = cr
	constants.Qr = qr
	m, _ := InitCategory() // TODO: can we acknowledge this error
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
