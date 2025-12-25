package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"shelllab/backend/database/models"
)

// AtlasLootImporter handles AtlasLoot data imports
type AtlasLootImporter struct {
	db *sql.DB
}

// NewAtlasLootImporter creates a new AtlasLoot importer
func NewAtlasLootImporter(db *sql.DB) *AtlasLootImporter {
	return &AtlasLootImporter{db: db}
}

// ImportFromJSON imports AtlasLoot structure from JSON
func (a *AtlasLootImporter) ImportFromJSON(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var categories []models.AtlasLootImportCategory
	if err := json.NewDecoder(file).Decode(&categories); err != nil {
		return fmt.Errorf("failed to decode AtlasLoot JSON: %w", err)
	}

	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear all existing data first to avoid constraint violations
	fmt.Println("  Clearing existing AtlasLoot data...")
	tx.Exec("DELETE FROM atlasloot_items")
	tx.Exec("DELETE FROM atlasloot_tables")
	tx.Exec("DELETE FROM atlasloot_modules")
	tx.Exec("DELETE FROM atlasloot_categories")

	stmtCat, _ := tx.Prepare("INSERT INTO atlasloot_categories (name, display_name, sort_order) VALUES (?, ?, ?)")
	defer stmtCat.Close()
	stmtMod, _ := tx.Prepare("INSERT INTO atlasloot_modules (category_id, name, display_name, sort_order) VALUES (?, ?, ?, ?)")
	defer stmtMod.Close()
	stmtTbl, _ := tx.Prepare("INSERT INTO atlasloot_tables (module_id, table_key, display_name, sort_order) VALUES (?, ?, ?, ?)")
	defer stmtTbl.Close()
	stmtItem, _ := tx.Prepare("INSERT INTO atlasloot_items (table_id, item_id, drop_chance, sort_order) VALUES (?, ?, ?, ?)")
	defer stmtItem.Close()

	for i, cat := range categories {
		fmt.Printf("  [%d/%d] Category: %s (modules: %d)\n", i+1, len(categories), cat.Name, len(cat.Modules))
		res, err := stmtCat.Exec(cat.Key, cat.Name, cat.Sort)
		if err != nil {
			fmt.Printf("    ERROR inserting category: %v\n", err)
			continue
		}
		catID, _ := res.LastInsertId()

		for j, mod := range cat.Modules {
			fmt.Printf("    Module: %s (tables: %d)\n", mod.Name, len(mod.Tables))
			res, err := stmtMod.Exec(catID, mod.Key, mod.Name, j)
			if err != nil {
				fmt.Printf("      ERROR inserting module: %v\n", err)
				continue
			}
			modID, _ := res.LastInsertId()

			for k, tbl := range mod.Tables {
				itemCount := len(tbl.Items)
				if itemCount > 0 {
					fmt.Printf("      Table: %s (items: %d)\n", tbl.Name, itemCount)
				}
				res, err := stmtTbl.Exec(modID, tbl.Key, tbl.Name, k)
				if err != nil {
					fmt.Printf("        ERROR inserting table: %v\n", err)
					continue
				}
				tblID, _ := res.LastInsertId()

				insertedItems := 0
				for l, item := range tbl.Items {
					if item.ID <= 0 {
						continue
					}
					if _, err := stmtItem.Exec(tblID, item.ID, item.DropRate, l); err != nil {
						fmt.Printf("        ERROR inserting item %d: %v\n", item.ID, err)
					} else {
						insertedItems++
					}
				}
				if insertedItems > 0 && insertedItems != itemCount {
					fmt.Printf("        Inserted %d/%d items (skipped %d with ID<=0)\n", insertedItems, itemCount, itemCount-insertedItems)
				}
			}
		}
	}
	fmt.Println("Committing transaction...")
	return tx.Commit()
}

// CheckAndImport checks if AtlasLoot data exists and imports it
func (a *AtlasLootImporter) CheckAndImport(dataDir string) error {
	var count int
	// Check items table - the most important one
	if err := a.db.QueryRow("SELECT COUNT(*) FROM atlasloot_items").Scan(&count); err != nil {
		return nil
	}
	if count == 0 {
		path := fmt.Sprintf("%s/atlasloot.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Importing AtlasLoot...")
			return a.ImportFromJSON(path)
		}
	}
	return nil
}
