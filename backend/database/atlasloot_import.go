package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type AtlasLootImportItem struct {
	ID       int    `json:"id"`
	DropRate string `json:"drop_rate"`
}

type AtlasLootImportTable struct {
	Key   string                `json:"key"`
	Name  string                `json:"name"`
	Items []AtlasLootImportItem `json:"items"`
}

type AtlasLootImportModule struct {
	Key    string                 `json:"key"`
	Name   string                 `json:"name"`
	Tables []AtlasLootImportTable `json:"tables"`
}

type AtlasLootImportCategory struct {
	Key     string                  `json:"key"`
	Name    string                  `json:"name"`
	Sort    int                     `json:"sort"`
	Modules []AtlasLootImportModule `json:"modules"`
}

// CheckAndImportAtlasLoot checks if AtlasLoot data exists and imports it
func (r *ItemRepository) CheckAndImportAtlasLoot(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM atlasloot_categories").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		path := fmt.Sprintf("%s/atlasloot.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("AtlasLoot table is empty. Importing from atlasloot.json...")
			return r.ImportAtlasLootFromJSON(path)
		}
	}
	return nil
}

// ImportAtlasLootFromJSON imports AtlasLoot structure from JSON
func (r *ItemRepository) ImportAtlasLootFromJSON(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var categories []AtlasLootImportCategory
	if err := json.NewDecoder(file).Decode(&categories); err != nil {
		return fmt.Errorf("failed to decode AtlasLoot JSON: %w", err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Clear existing (just in case)
	// Actually no need if we check count == 0, but safe to clear related tables if re-importing
	// Since we replace, let's keep it simple.

	// Prepare statements
	stmtCat, err := tx.Prepare("INSERT INTO atlasloot_categories (key, name, type, sort_order) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtCat.Close()

	stmtMod, err := tx.Prepare("INSERT INTO atlasloot_modules (category_id, name, display_name, sort_order) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtMod.Close()

	stmtTbl, err := tx.Prepare("INSERT INTO atlasloot_tables (module_id, table_key, display_name, sort_order) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtTbl.Close()

	// Items insert. Note: atlasloot_items table? or category_items?
	// Existing schema uses category_items connected to categories/tables?
	// Let's check `atlasloot_schema.go` or `sqlite_db.go`.
	// New schema in `InitAtlasLootSchema` uses `atlasloot_items`.

	stmtItem, err := tx.Prepare("INSERT INTO atlasloot_items (table_id, item_id, drop_rate, sort_order) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtItem.Close()

	for _, cat := range categories {
		// Category
		res, err := stmtCat.Exec(cat.Key, cat.Name, "root", cat.Sort)
		if err != nil {
			return fmt.Errorf("error inserting category %s: %w", cat.Name, err)
		}
		catID, _ := res.LastInsertId()

		for j, mod := range cat.Modules {
			// Module
			res, err := stmtMod.Exec(catID, mod.Key, mod.Name, j)
			if err != nil {
				return fmt.Errorf("error inserting module %s: %w", mod.Name, err)
			}
			modID, _ := res.LastInsertId()

			for k, tbl := range mod.Tables {
				// Table
				res, err := stmtTbl.Exec(modID, tbl.Key, tbl.Name, k)
				if err != nil {
					return fmt.Errorf("error inserting table %s: %w", tbl.Name, err)
				}
				tblID, _ := res.LastInsertId()

				for l, item := range tbl.Items {
					if item.ID <= 0 {
						continue
					}
					_, err := stmtItem.Exec(tblID, item.ID, item.DropRate, l)
					if err != nil {
						continue
					}
				}
			}
		}
	}

	return tx.Commit()
}
