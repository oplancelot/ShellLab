package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"shelllab/backend/database/models"
)

// LootImporter handles loot table data imports
type LootImporter struct {
	db *sql.DB
}

// NewLootImporter creates a new loot importer
func NewLootImporter(db *sql.DB) *LootImporter {
	return &LootImporter{db: db}
}

// ImportFromJSON imports a loot table from JSON into SQLite
func (l *LootImporter) ImportFromJSON(tableName, jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var entries []models.LootTemplateEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
		stmt.Exec(e.Entry, e.Item, e.Chance, e.GroupID, e.MinCountOrRef, e.MaxCount)
	}
	return tx.Commit()
}

// ImportAll imports all defined loot tables if they exist
func (l *LootImporter) ImportAll(dataDir string) error {
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
			fmt.Printf("Importing %s...\n", table)
			l.ImportFromJSON(table, path)
		}
	}
	return nil
}

// CheckAndImport checks if creature_loot table is empty and imports all loot if so
func (l *LootImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := l.db.QueryRow("SELECT COUNT(*) FROM creature_loot").Scan(&count); err != nil {
		return nil
	}
	if count == 0 {
		return l.ImportAll(dataDir)
	}
	return nil
}
