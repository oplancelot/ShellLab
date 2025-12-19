package database

import (
	"fmt"
)

// AtlasLootRepository handles AtlasLoot data queries
type AtlasLootRepository struct {
	db *SQLiteDB
}

// NewAtlasLootRepository creates a new repository
func NewAtlasLootRepository(db *SQLiteDB) *AtlasLootRepository {
	return &AtlasLootRepository{db: db}
}

// GetCategories returns all category names
func (r *AtlasLootRepository) GetCategories() ([]string, error) {
	rows, err := r.db.DB().Query(`
		SELECT display_name 
		FROM atlasloot_categories 
		ORDER BY sort_order, display_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		categories = append(categories, name)
	}
	return categories, nil
}

// GetModules returns module names for a category
func (r *AtlasLootRepository) GetModules(categoryName string) ([]string, error) {
	rows, err := r.db.DB().Query(`
		SELECT m.display_name
		FROM atlasloot_modules m
		JOIN atlasloot_categories c ON m.category_id = c.id
		WHERE c.display_name = ?
		ORDER BY m.sort_order, m.display_name
	`, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		modules = append(modules, name)
	}
	return modules, nil
}

// AtlasTable represents a loot table reference
type AtlasTable struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
}

// GetTables returns table references for a module
func (r *AtlasLootRepository) GetTables(categoryName, moduleName string) ([]AtlasTable, error) {
	rows, err := r.db.DB().Query(`
		SELECT t.table_key, t.display_name
		FROM atlasloot_tables t
		JOIN atlasloot_modules m ON t.module_id = m.id
		JOIN atlasloot_categories c ON m.category_id = c.id
		WHERE c.display_name = ? AND m.display_name = ?
		ORDER BY t.sort_order, t.display_name
	`, categoryName, moduleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []AtlasTable
	for rows.Next() {
		var t AtlasTable
		if err := rows.Scan(&t.Key, &t.DisplayName); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// LootEntry represents a loot item with metadata
type LootEntry struct {
	ItemID     int    `json:"itemId"`
	ItemName   string `json:"itemName"`
	IconName   string `json:"iconName"`
	Quality    int    `json:"quality"`
	DropChance string `json:"dropChance,omitempty"`
}

// GetLootItems returns items for a specific table
func (r *AtlasLootRepository) GetLootItems(categoryName, moduleName, tableKey string) ([]*LootEntry, error) {
	rows, err := r.db.DB().Query(`
		SELECT 
			al.item_id,
			al.drop_chance,
			al.sort_order,
			i.name,
			i.icon_path,
			i.quality
		FROM atlasloot_items al
		JOIN atlasloot_tables t ON al.table_id = t.id
		JOIN atlasloot_modules m ON t.module_id = m.id
		JOIN atlasloot_categories c ON m.category_id = c.id
		JOIN items i ON al.item_id = i.entry
		WHERE c.display_name = ? AND m.display_name = ? AND t.table_key = ?
		ORDER BY al.sort_order, i.name
	`, categoryName, moduleName, tableKey)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*LootEntry
	for rows.Next() {
		var entry LootEntry
		var sortOrder int
		if err := rows.Scan(
			&entry.ItemID,
			&entry.DropChance,
			&sortOrder,
			&entry.ItemName,
			&entry.IconName,
			&entry.Quality,
		); err != nil {
			return nil, err
		}
		items = append(items, &entry)
	}
	return items, nil
}

// InsertCategory inserts a new category
func (r *AtlasLootRepository) InsertCategory(name, displayName string, sortOrder int) (int64, error) {
	_, err := r.db.DB().Exec(`
		INSERT OR IGNORE INTO atlasloot_categories (name, display_name, sort_order)
		VALUES (?, ?, ?)
	`, name, displayName, sortOrder)
	if err != nil {
		return 0, err
	}

	// Get the ID (either newly inserted or existing)
	var id int64
	err = r.db.DB().QueryRow(`
		SELECT id FROM atlasloot_categories WHERE name = ?
	`, name).Scan(&id)
	return id, err
}

// InsertModule inserts a new module
func (r *AtlasLootRepository) InsertModule(categoryID int, name, displayName string, sortOrder int) (int64, error) {
	result, err := r.db.DB().Exec(`
		INSERT INTO atlasloot_modules (category_id, name, display_name, sort_order)
		VALUES (?, ?, ?, ?)
	`, categoryID, name, displayName, sortOrder)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertTable inserts a new table
func (r *AtlasLootRepository) InsertTable(moduleID int, tableKey, displayName string, sortOrder int) (int64, error) {
	result, err := r.db.DB().Exec(`
		INSERT INTO atlasloot_tables (module_id, table_key, display_name, sort_order)
		VALUES (?, ?, ?, ?)
	`, moduleID, tableKey, displayName, sortOrder)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertItem inserts a loot item
func (r *AtlasLootRepository) InsertItem(tableID int, itemID int, dropChance string, sortOrder int) error {
	_, err := r.db.DB().Exec(`
		INSERT INTO atlasloot_items (table_id, item_id, drop_chance, sort_order)
		VALUES (?, ?, ?, ?)
	`, tableID, itemID, dropChance, sortOrder)
	return err
}

// ClearAllData removes all AtlasLoot data (for reimport)
func (r *AtlasLootRepository) ClearAllData() error {
	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tables := []string{
		"atlasloot_items",
		"atlasloot_tables",
		"atlasloot_modules",
		"atlasloot_categories",
	}

	for _, table := range tables {
		res, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return err
		}
		count, _ := res.RowsAffected()
		fmt.Printf("  Deleted %d rows from %s\n", count, table)
	}

	return tx.Commit()
}

// GetStats returns statistics about AtlasLoot data
func (r *AtlasLootRepository) GetStats() (map[string]int, error) {
	stats := make(map[string]int)

	tables := map[string]string{
		"categories": "atlasloot_categories",
		"modules":    "atlasloot_modules",
		"tables":     "atlasloot_tables",
		"items":      "atlasloot_items",
	}

	for key, table := range tables {
		var count int
		err := r.db.DB().QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			return nil, err
		}
		stats[key] = count
	}

	return stats, nil
}
