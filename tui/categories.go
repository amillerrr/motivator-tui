package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"motivator-tui/categories"
	"motivator-tui/tui/constants"
)

type (
	updateCategoryListMsg struct{}
)

// SelectMsg the message to change the view to the selected entry
type SelectMsg struct {
	ActiveCategory string
}

// Model the categories model definition
type Model struct {
	list     list.Model
	quitting bool
}

// InitCategory initialize the projectui model for your program
func InitCategory() (tea.Model, tea.Cmd) {
	items, err := displayCategoryList()
	if err != nil {
		return nil, func() tea.Msg { return errMsg{err} }
	}

	m := Model{list: list.New(items, list.NewDefaultDelegate(), 8, 8)}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}

	m.list.Title = "Categories"
	m.list.SetShowStatusBar(false)
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Back,
		}
	}
	return m, nil
}

func displayCategoryList() ([]list.Item, error) {
	cat := categories.NewInMemoryRepository()
	categories, err := cat.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("cannot get all projects: %w", err)
	}
	return categoriesToItems(categories), err
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)
	case updateCategoryListMsg:
		m.updateCategoryList()
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return constants.DocStyle.Render(m.list.View() + "\n")
}

// handleWindowSize updates the list size when the window size changes
func (m *Model) handleWindowSize(msg tea.WindowSizeMsg) {
	constants.WindowSize = msg
	top, right, bottom, left := constants.DocStyle.GetMargin()
	m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
}

// updateCategoryList refreshes the list of categories
func (m *Model) updateCategoryList() {
	cat := categories.NewInMemoryRepository()
	categoriesList, _ := cat.GetAllCategories()
	m.list.SetItems(categoriesToItems(categoriesList))
}

// handleKeyMsg processes key inputs
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, constants.Keymap.Quit):
		m.quitting = true
		return m, tea.Quit
	case key.Matches(msg, constants.Keymap.Enter):
		activeCategory := m.list.SelectedItem().(categories.Category)
		quote := InitQuote(&constants.Qr, activeCategory.CategoryName, constants.P)
		return quote.Update(constants.WindowSize)
	default:
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
}

// categoriesToItems converts []categories.Category to []list.Item
func categoriesToItems(categories []categories.Category) []list.Item {
	items := make([]list.Item, len(categories))
	for i := range categories {
		items[i] = categories[i]
	}
	return items
}

// func (m Model) getActiveCategoryName() string {
// 	items := m.list.Items()
// 	activeCategory := items[m.list.Index()]
// 	return activeCategory.(categories.Category).CategoryName
// }
