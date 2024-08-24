package categories

import (
	"fmt"
)

type Category struct {
	CategoryName string
}

// Implement list.Item for Bubbletea TUI

// Title the project title to display in a list
func (c Category) Title() string { return c.CategoryName }

// Description the project description to display in a list
func (c Category) Description() string { return "" }

// FilterValue choose what field to use for filtering in a Bubbletea list component
func (c Category) FilterValue() string { return c.CategoryName }

// Repository CRUD operations for Categories
type Repository interface {
	PrintCategories() error
	GetAllCategories() ([]Category, error)
}

// InMemoryRepository stores categories in memory
type InMemoryRepository struct {
	categories []Category
}

// NewInMemoryRepository creates a new InMemoryRepository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		categories: []Category{{CategoryName: "Grit"}, {CategoryName: "Perseverance"}, {CategoryName: "Gratitude"}},
	}
}

// PrintCategories print all categories to the console
func (r *InMemoryRepository) PrintCategories() error {
	categories, err := r.GetAllCategories()
	if err != nil {
		return err
	}
	for _, category := range categories {
		fmt.Println(category.CategoryName)
	}
	return nil
}

// GetAllCategories retrieve all categories
func (r *InMemoryRepository) GetAllCategories() ([]Category, error) {
	return r.categories, nil // No error in this simple implementation
}
