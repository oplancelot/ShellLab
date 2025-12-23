package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type CreatureTemplateEntry struct {
	Entry            int    `json:"entry"`
	Name             string `json:"name"`
	Subname          string `json:"subname"`
	LevelMin         int    `json:"level_min"`
	LevelMax         int    `json:"level_max"`
	HealthMin        int    `json:"health_min"`
	HealthMax        int    `json:"health_max"`
	ManaMin          int    `json:"mana_min"`
	ManaMax          int    `json:"mana_max"`
	CreatureType     int    `json:"creature_type"`
	CreatureRank     int    `json:"creature_rank"`
	Faction          int    `json:"faction"`
	NPCFlags         int    `json:"npc_flags"`
	LootID           int    `json:"loot_id"`
	SkinLootID       int    `json:"skinning_loot_id"`
	PickpocketLootID int    `json:"pickpocket_loot_id"`
}

// ImportCreaturesFromJSON imports creatures from JSON into SQLite
func (r *ItemRepository) ImportCreaturesFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", jsonPath, err)
	}

	var creatures []CreatureTemplateEntry
	if err := json.Unmarshal(data, &creatures); err != nil {
		return fmt.Errorf("failed to parse JSON %s: %w", jsonPath, err)
	}

	tx, err := r.db.DB().Begin()
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

	for _, c := range creatures {
		_, err := stmt.Exec(
			c.Entry, c.Name, c.Subname, c.LevelMin, c.LevelMax,
			c.HealthMin, c.HealthMax, c.ManaMin, c.ManaMax,
			c.CreatureType, c.CreatureRank, c.Faction, c.NPCFlags,
			c.LootID, c.SkinLootID, c.PickpocketLootID,
		)
		if err != nil {
			fmt.Printf("Warning: Failed to import creature %d: %v\n", c.Entry, err)
			continue
		}
	}

	return tx.Commit()
}

// ImportAllCreatures checks and imports creatures
func (r *ItemRepository) ImportAllCreatures(dataDir string) error {
	path := fmt.Sprintf("%s/creatures.json", dataDir)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Importing Creatures from %s...\n", path)
		if err := r.ImportCreaturesFromJSON(path); err != nil {
			fmt.Printf("Error importing creatures: %v\n", err)
			return err
		}
		fmt.Println("✓ Creatures imported successfully!")
	}
	return nil
}

// CheckAndImportCreatures checks if creatures table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportCreatures(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM creatures").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		return r.ImportAllCreatures(dataDir)
	}
	return nil
}
