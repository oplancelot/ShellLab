package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"shelllab/backend/database/models"
)

// SpellImporter handles spell data imports
type SpellImporter struct {
	db *sql.DB
}

// NewSpellImporter creates a new spell importer
func NewSpellImporter(db *sql.DB) *SpellImporter {
	return &SpellImporter{db: db}
}

// ImportFromJSON imports spells from JSON into SQLite
func (s *SpellImporter) ImportFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var spells []models.SpellEntry
	if err := json.Unmarshal(data, &spells); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	tx, err := s.db.Begin()
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

	for _, sp := range spells {
		_, err := stmt.Exec(
			sp.Entry, sp.Name, sp.Description,
			sp.EffectBasePoints1, sp.EffectBasePoints2, sp.EffectBasePoints3,
			sp.EffectDieSides1, sp.EffectDieSides2, sp.EffectDieSides3,
		)
		if err != nil {
			continue
		}
	}
	return tx.Commit()
}

// CheckAndImport checks if spells table is empty and imports if JSON exists
func (s *SpellImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM spells").Scan(&count); err != nil {
		return nil
	}
	if count == 0 {
		path := fmt.Sprintf("%s/spells.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Importing Spells...")
			return s.ImportFromJSON(path)
		}
	}
	return nil
}
