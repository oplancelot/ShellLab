package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type SpellEntry struct {
	Entry             int    `json:"entry"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	EffectBasePoints1 int    `json:"effectBasePoints1"`
	EffectBasePoints2 int    `json:"effectBasePoints2"`
	EffectBasePoints3 int    `json:"effectBasePoints3"`
	EffectDieSides1   int    `json:"effectDieSides1"`
	EffectDieSides2   int    `json:"effectDieSides2"`
	EffectDieSides3   int    `json:"effectDieSides3"`
}

// ImportSpellsFromJSON imports spells from JSON into SQLite
func (r *ItemRepository) ImportSpellsFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", jsonPath, err)
	}

	var spells []SpellEntry
	if err := json.Unmarshal(data, &spells); err != nil {
		return fmt.Errorf("failed to parse JSON %s: %w", jsonPath, err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM spells")

	stmt, err := tx.Prepare(`
		INSERT INTO spells (
			entry, name, description,
			effect_base_points1, effect_base_points2, effect_base_points3,
			effect_die_sides1, effect_die_sides2, effect_die_sides3
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range spells {
		_, err := stmt.Exec(
			s.Entry, s.Name, s.Description,
			s.EffectBasePoints1, s.EffectBasePoints2, s.EffectBasePoints3,
			s.EffectDieSides1, s.EffectDieSides2, s.EffectDieSides3,
		)
		if err != nil {
			continue
		}
	}

	return tx.Commit()
}

// ImportAllSpells checks and imports spells
func (r *ItemRepository) ImportAllSpells(dataDir string) error {
	path := fmt.Sprintf("%s/spells.json", dataDir)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Importing Spells from %s...\n", path)
		if err := r.ImportSpellsFromJSON(path); err != nil {
			fmt.Printf("Error importing spells: %v\n", err)
			return err
		}
		fmt.Println("✓ Spells imported successfully!")
	}
	return nil
}

// CheckAndImportSpells checks if spells table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportSpells(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM spells").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		return r.ImportAllSpells(dataDir)
	}
	return nil
}
