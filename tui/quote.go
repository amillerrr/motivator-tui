package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"motivator-tui/quote"
	"motivator-tui/tui/constants"
)

type (
	errMsg struct{ error }
	// NewUpdatedQuotesMsg holds the new quotes
	NewUpdatedQuotesMsg struct {
		quoteMsg quote.Quote
	}
)

type styles struct {
	quoteBox lipgloss.Style
	category lipgloss.Style
	quote    lipgloss.Style
}

var quoteStyles = styles{
	quoteBox: lipgloss.NewStyle().
		Border(customBorder).
		BorderForeground(lipgloss.Color("#C772F1")).
		Padding(0, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true),
	category: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7FE098")).
		Height(1).
		Bold(true).
		Underline(true),
	quote: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F9EECE")),
	// Foreground(lipgloss.Color("#FFF7DB")),
}

// custom border for quote result
var customBorder = lipgloss.Border{
	Top:         "-",
	Bottom:      "-",
	Left:        "|",
	Right:       "|",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

// Quote implements tea.Model
type QuoteModel struct {
	width          int
	height         int
	activeCategory string
	wrappedQuote   string
	error          string
	quotes         quote.Quote
	quitting       bool
}

// Init run any intial IO on program start
func (m QuoteModel) Init() tea.Cmd {
	return nil
}

// InitQuote initialize the quoteui model for your program
func InitQuote(qr *quote.QuoteRepository, activeCategory string, p *tea.Program) *QuoteModel {
	m := &QuoteModel{activeCategory: activeCategory}
	m.updateDimensions()

	updatedQuoteMsg := m.setupQuote().(NewUpdatedQuotesMsg)
	m.quotes = updatedQuoteMsg.quoteMsg
	m.wrapQuote()
	m.setContent()

	return m
}

// setupQuote fetches the quote for the active category
func (m *QuoteModel) setupQuote() tea.Msg {
	quotes, err := constants.Qr.GetQuoteByCategoryName(m.activeCategory)
	if err != nil {
		return errMsg{fmt.Errorf("Cannot find project: %v", err)}
	}
	return NewUpdatedQuotesMsg{quotes}
}

// wrapQuote formats the quote and author into a wrapped string
func (m *QuoteModel) wrapQuote() string {
	m.updateDimensions()

	selectedQuote := quote.FormatQuote(m.quotes)
	selectedAuthor := quote.FormatAuthor(m.quotes)
	var wrappedLines []string

	for len(selectedQuote) > 0 {
		if utf8.RuneCountInString(selectedQuote) > m.width-8 {
			splitPos := strings.LastIndex(selectedQuote[:m.width-8], " ")
			if splitPos == -1 {
				splitPos = m.width
			}

			wrappedLines = append(wrappedLines, selectedQuote[:splitPos])
			selectedQuote = selectedQuote[splitPos:]
			selectedQuote = strings.TrimSpace(selectedQuote)
		} else {
			wrappedLines = append(wrappedLines, selectedQuote)
			break
		}
	}

	indent := strings.Repeat(" ", 8)

	wrappedLines = append(wrappedLines, indent+selectedAuthor)

	return strings.Join(wrappedLines, "\n")
}

// setContent updates the formatted content
func (m *QuoteModel) setContent() string {
	m.updateDimensions()

	selectedCategory := quoteStyles.category.Align(lipgloss.Center, lipgloss.Center).Render(quote.FormatCategoryName(m.quotes))
	formattedQuote := quoteStyles.quote.Render(m.wrappedQuote)
	stackedContent := lipgloss.JoinVertical(lipgloss.Center, selectedCategory, formattedQuote)

	borderedContent := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		quoteStyles.quoteBox.Render(stackedContent),
	)

	return borderedContent
}

// updateDimensions recalculates dimensions for layout changes
func (m *QuoteModel) updateDimensions() {
	top, right, bottom, left := constants.DocStyle.GetMargin()
	m.width = constants.WindowSize.Width - left - right
	m.height = constants.WindowSize.Height - top - bottom - 1
}

// Update handle IO and commands
func (m QuoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateDimensions()
	case errMsg:
		m.error = msg.Error()
	case NewUpdatedQuotesMsg:
		m.quotes = msg.quoteMsg
		m.wrappedQuote = m.wrapQuote()
		m.setContent()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Back):
			return InitCategory()
		case key.Matches(msg, constants.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.wrappedQuote = m.wrapQuote()
	m.setContent() // Refresh the content on every Update call

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// helpView returns the help text
func (m QuoteModel) helpView() string {
	return constants.HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

// View return the text UI to be output to the terminal
func (m QuoteModel) View() string {
	if m.quitting {
		return ""
	}

	var errorView string
	if m.error != "" {
		errorView = constants.ErrStyle(m.error)
	}

	formatted := lipgloss.JoinVertical(
		lipgloss.Top,
		m.setContent(),
		m.helpView(),
		errorView,
	)
	return constants.DocStyle.Render(formatted)
}
