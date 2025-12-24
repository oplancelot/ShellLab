package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"shelllab/backend/database/models"
)

// CreatureImporter handles creature data imports
type CreatureImporter struct {
	db *sql.DB
}

// NewCreatureImporter creates a new creature importer
func NewCreatureImporter(db *sql.DB) *CreatureImporter {
	return &CreatureImporter{db: db}
}

// ImportFromJSON imports creatures from JSON into SQLite
func (c *CreatureImporter) ImportFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var creatures []models.CreatureTemplateEntry
	if err := json.Unmarshal(data, &creatures); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM creatures")

	stmt, err := tx.Prepare(`
		INSERT INTO creatures (
			entry, name, subname, level_min, level_max,
			health_min, health_max, mana_min, mana_max,
			creature_type, creature_rank, faction, npc_flags,
			loot_id, skin_loot_id, pickpocket_loot_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, cr := range creatures {
		_, err := stmt.Exec(
			cr.Entry, cr.Name, cr.Subname, cr.LevelMin, cr.LevelMax,
			cr.HealthMin, cr.HealthMax, cr.ManaMin, cr.ManaMax,
			cr.CreatureType, cr.CreatureRank, cr.Faction, cr.NPCFlags,
			cr.LootID, cr.SkinLootID, cr.PickpocketLootID,
		)
		if err != nil {
			continue
		}
	}
	return tx.Commit()
}

// CheckAndImport checks if creatures table is empty and imports if JSON exists
func (c *CreatureImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := c.db.QueryRow("SELECT COUNT(*) FROM creatures").Scan(&count); err != nil {
		return nil
	}
	if count == 0 {
		path := fmt.Sprintf("%s/creatures.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Importing Creatures...")
			return c.ImportFromJSON(path)
		}
	}
	return nil
}
