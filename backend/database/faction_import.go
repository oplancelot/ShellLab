package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type FactionEntry struct {
	FactionID   int    `json:"factionID"`
	Name        string `json:"name_loc0"`
	Description string `json:"description1_loc0"`
	Side        int    `json:"side"`
	Team        int    `json:"team"`
}

// ImportFactionsFromJSON imports factions from JSON into SQLite
func (r *ItemRepository) ImportFactionsFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", jsonPath, err)
	}

	var factions []FactionEntry
	if err := json.Unmarshal(data, &factions); err != nil {
		return fmt.Errorf("failed to parse JSON %s: %w", jsonPath, err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM factions")

	stmt, err := tx.Prepare("INSERT INTO factions (id, name, description, side, category_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, f := range factions {
		_, err := stmt.Exec(f.FactionID, f.Name, f.Description, f.Side, f.Team)
		if err != nil {
			continue
		}
	}

	return tx.Commit()
}

// ImportAllFactions checks and imports
func (r *ItemRepository) ImportAllFactions(dataDir string) error {
	path := fmt.Sprintf("%s/factions.json", dataDir)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Importing Factions from %s...\n", path)
		if err := r.ImportFactionsFromJSON(path); err != nil {
			fmt.Printf("Error importing factions: %v\n", err)
			return err
		}
		fmt.Println("✓ Factions imported successfully!")
	}
	return nil
}

// CheckAndImportFactions checks if factions table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportFactions(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM factions").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		return r.ImportAllFactions(dataDir)
	}
	return nil
}
