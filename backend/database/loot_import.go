package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type LootTemplateEntry struct {
	Entry         int     `json:"entry"`
	Item          int     `json:"item"`
	Chance        float64 `json:"chance"`
	GroupID       int     `json:"groupId"`
	MinCountOrRef int     `json:"minCountOrRef"`
	MaxCount      int     `json:"maxCount"`
}

// ImportLootFromJSON imports a loot table from JSON into SQLite
// tableName: creature_loot, reference_loot, etc.
// jsonPath: path to JSON file
func (r *ItemRepository) ImportLootFromJSON(tableName, jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", jsonPath, err)
	}

	var entries []LootTemplateEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse JSON %s: %w", jsonPath, err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear table
	tx.Exec(fmt.Sprintf("DELETE FROM %s", tableName))

	stmt, err := tx.Prepare(fmt.Sprintf(
		"INSERT INTO %s (entry, item, chance, groupid, mincount_or_ref, maxcount) VALUES (?, ?, ?, ?, ?, ?)",
		tableName,
	))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range entries {
		_, err := stmt.Exec(e.Entry, e.Item, e.Chance, e.GroupID, e.MinCountOrRef, e.MaxCount)
		if err != nil {
			continue // Skip errors
		}
	}

	return tx.Commit()
}

// ImportAllLoot imports all defined loot tables if they exist
func (r *ItemRepository) ImportAllLoot(dataDir string) error {
	lootFiles := map[string]string{
		"creature_loot":   "creature_loot.json",
		"reference_loot":  "reference_loot.json",
		"gameobject_loot": "gameobject_loot.json",
		"item_loot":       "item_loot.json",
		"disenchant_loot": "disenchant_loot.json",
	}

	for table, file := range lootFiles {
		path := fmt.Sprintf("%s/%s", dataDir, file)
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Importing %s from %s...\n", table, path)
			if err := r.ImportLootFromJSON(table, path); err != nil {
				fmt.Printf("Error importing %s: %v\n", table, err)
			}
		}
	}
	return nil
}

// CheckAndImportLoot checks if creature_loot table is empty and imports all loot if so
func (r *ItemRepository) CheckAndImportLoot(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM creature_loot").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		return r.ImportAllLoot(dataDir)
	}
	return nil
}
