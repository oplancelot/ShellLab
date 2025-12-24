package models

// Category represents a loot category (instance, boss, set, etc.)
type Category struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	ParentID  *int   `json:"parentId,omitempty"`
	Type      string `json:"type"`
	SortOrder int    `json:"sortOrder"`
}

// CategoryItem represents an item in a category
type CategoryItem struct {
	CategoryID int    `json:"categoryId"`
	ItemID     int    `json:"itemId"`
	DropRate   string `json:"dropRate,omitempty"`
	SortOrder  int    `json:"sortOrder"`
}
