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

	stmtCat, _ := tx.Prepare("INSERT INTO atlasloot_categories (name, display_name, sort_order) VALUES (?, ?, ?)")
	defer stmtCat.Close()
	stmtMod, _ := tx.Prepare("INSERT INTO atlasloot_modules (category_id, name, display_name, sort_order) VALUES (?, ?, ?, ?)")
	defer stmtMod.Close()
	stmtTbl, _ := tx.Prepare("INSERT INTO atlasloot_tables (module_id, table_key, display_name, sort_order) VALUES (?, ?, ?, ?)")
	defer stmtTbl.Close()
	stmtItem, _ := tx.Prepare("INSERT INTO atlasloot_items (table_id, item_id, drop_chance, sort_order) VALUES (?, ?, ?, ?)")
	defer stmtItem.Close()

	for _, cat := range categories {
		res, err := stmtCat.Exec(cat.Key, cat.Name, cat.Sort)
		if err != nil {
			continue
		}
		catID, _ := res.LastInsertId()

		for j, mod := range cat.Modules {
			res, err := stmtMod.Exec(catID, mod.Key, mod.Name, j)
			if err != nil {
				continue
			}
			modID, _ := res.LastInsertId()

			for k, tbl := range mod.Tables {
				res, err := stmtTbl.Exec(modID, tbl.Key, tbl.Name, k)
				if err != nil {
					continue
				}
				tblID, _ := res.LastInsertId()

				for l, item := range tbl.Items {
					if item.ID <= 0 {
						continue
					}
					stmtItem.Exec(tblID, item.ID, item.DropRate, l)
				}
			}
		}
	}
	return tx.Commit()
}

// CheckAndImport checks if AtlasLoot data exists and imports it
func (a *AtlasLootImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := a.db.QueryRow("SELECT COUNT(*) FROM atlasloot_categories").Scan(&count); err != nil {
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
