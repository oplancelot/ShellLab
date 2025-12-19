package database

import (
	"fmt"
	"strings"
)

// SearchFilter defines criteria for advanced item search
type SearchFilter struct {
	Query         string `json:"query"`
	Quality       []int  `json:"quality,omitempty"`
	Class         []int  `json:"class,omitempty"`
	SubClass      []int  `json:"subClass,omitempty"`
	InventoryType []int  `json:"inventoryType,omitempty"`
	MinLevel      int    `json:"minLevel,omitempty"`
	MaxLevel      int    `json:"maxLevel,omitempty"`
	MinReqLevel   int    `json:"minReqLevel,omitempty"`
	MaxReqLevel   int    `json:"maxReqLevel,omitempty"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
}

// SearchResult represents the search output
type SearchResult struct {
	Items      []*Item `json:"items"`
	TotalCount int     `json:"totalCount"`
}

// AdvancedSearch performs a multi-dimensional search on items
func (r *ItemRepository) AdvancedSearch(filter SearchFilter) (*SearchResult, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}

	var conditions []string
	var args []interface{}

	// Name filter
	if filter.Query != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+filter.Query+"%")
	}

	// Quality filter
	if len(filter.Quality) > 0 {
		placeholders := make([]string, len(filter.Quality))
		for i, q := range filter.Quality {
			placeholders[i] = "?"
			args = append(args, q)
		}
		conditions = append(conditions, fmt.Sprintf("quality IN (%s)", strings.Join(placeholders, ",")))
	}

	// Class filter
	if len(filter.Class) > 0 {
		placeholders := make([]string, len(filter.Class))
		for i, c := range filter.Class {
			placeholders[i] = "?"
			args = append(args, c)
		}
		conditions = append(conditions, fmt.Sprintf("class IN (%s)", strings.Join(placeholders, ",")))
	}

	// SubClass filter
	if len(filter.SubClass) > 0 {
		placeholders := make([]string, len(filter.SubClass))
		for i, sc := range filter.SubClass {
			placeholders[i] = "?"
			args = append(args, sc)
		}
		conditions = append(conditions, fmt.Sprintf("subclass IN (%s)", strings.Join(placeholders, ",")))
	}

	// InventoryType filter
	if len(filter.InventoryType) > 0 {
		placeholders := make([]string, len(filter.InventoryType))
		for i, it := range filter.InventoryType {
			placeholders[i] = "?"
			args = append(args, it)
		}
		conditions = append(conditions, fmt.Sprintf("inventory_type IN (%s)", strings.Join(placeholders, ",")))
	}

	// Level Range
	if filter.MinLevel > 0 {
		conditions = append(conditions, "item_level >= ?")
		args = append(args, filter.MinLevel)
	}
	if filter.MaxLevel > 0 {
		conditions = append(conditions, "item_level <= ?")
		args = append(args, filter.MaxLevel)
	}

	// Required Level Range
	if filter.MinReqLevel > 0 {
		conditions = append(conditions, "required_level >= ?")
		args = append(args, filter.MinReqLevel)
	}
	if filter.MaxReqLevel > 0 {
		conditions = append(conditions, "required_level <= ?")
		args = append(args, filter.MaxReqLevel)
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) FROM items " + whereClause
	var totalCount int
	err := r.db.DB().QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("search count error: %w", err)
	}

	// Data query
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, icon_path
		FROM items
		%s
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add limit/offset args
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.DB().Query(dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("search data error: %w", err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		err := rows.Scan(
			&item.Entry, &item.Name, &item.Quality, &item.ItemLevel,
			&item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return &SearchResult{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}
