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
			effect_die_sides1, effect_die_sides2, effect_die_sides3,
			duration_index,
			icon_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
			sp.DurationIndex,
			sp.IconName,
		)
		if err != nil {
			continue
		}
	}
	return tx.Commit()
}

// ImportDurationsFromJSON imports spell durations from JSON into SQLite
func (s *SpellImporter) ImportDurationsFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var durations []models.SpellDurationEntry
	if err := json.Unmarshal(data, &durations); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM spell_durations")

	stmt, err := tx.Prepare(`
		INSERT INTO spell_durations (
			id, duration_base, duration_per_level, max_duration
		) VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range durations {
		_, err := stmt.Exec(d.ID, d.DurationBase, d.DurationPerLevel, d.MaxDuration)
		if err != nil {
			continue
		}
	}
	return tx.Commit()
}

// CheckAndImport checks if spells table is empty and imports if JSON exists
func (s *SpellImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM spells").Scan(&count); err == nil && count == 0 {
		// Try enhanced spells first
		path := fmt.Sprintf("%s/spells_enhanced.json", dataDir)
		if _, err := os.Stat(path); err != nil {
			path = fmt.Sprintf("%s/spells.json", dataDir)
		}

		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Importing Spells from %s...\n", path)
			if err := s.ImportFromJSON(path); err != nil {
				return err
			}
		}
	}

	// Check durations
	if err := s.db.QueryRow("SELECT COUNT(*) FROM spell_durations").Scan(&count); err == nil && count == 0 {
		path := fmt.Sprintf("%s/spell_durations.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Importing Spell Durations...")
			if err := s.ImportDurationsFromJSON(path); err != nil {
				fmt.Printf("Error importing durations: %v\n", err) // Log but continue
			}
		}
	}

	return nil
}
