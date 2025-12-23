package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type ItemSetEntry struct {
	ID         int    `json:"itemsetID"`
	Name       string `json:"name_loc0"`
	Item1      int    `json:"item1"`
	Item2      int    `json:"item2"`
	Item3      int    `json:"item3"`
	Item4      int    `json:"item4"`
	Item5      int    `json:"item5"`
	Item6      int    `json:"item6"`
	Item7      int    `json:"item7"`
	Item8      int    `json:"item8"`
	Item9      int    `json:"item9"`
	Item10     int    `json:"item10"`
	SkillID    int    `json:"skillID"`
	SkillLevel int    `json:"skilllevel"`
	Bonus1     int    `json:"bonus1"`
	Bonus2     int    `json:"bonus2"`
	Bonus3     int    `json:"bonus3"`
	Bonus4     int    `json:"bonus4"`
	Bonus5     int    `json:"bonus5"`
	Bonus6     int    `json:"bonus6"`
	Bonus7     int    `json:"bonus7"`
	Bonus8     int    `json:"bonus8"`
	Spell1     int    `json:"spell1"`
	Spell2     int    `json:"spell2"`
	Spell3     int    `json:"spell3"`
	Spell4     int    `json:"spell4"`
	Spell5     int    `json:"spell5"`
	Spell6     int    `json:"spell6"`
	Spell7     int    `json:"spell7"`
	Spell8     int    `json:"spell8"`
}

func (r *ItemRepository) ImportItemSets(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open item sets JSON: %w", err)
	}
	defer file.Close()

	var sets []ItemSetEntry
	if err := json.NewDecoder(file).Decode(&sets); err != nil {
		return fmt.Errorf("failed to decode item sets JSON: %w", err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		REPLACE INTO itemsets (
			itemset_id, name,
			item1, item2, item3, item4, item5, item6, item7, item8, item9, item10,
			skill_id, skill_level,
			bonus1, bonus2, bonus3, bonus4, bonus5, bonus6, bonus7, bonus8,
			spell1, spell2, spell3, spell4, spell5, spell6, spell7, spell8
		) VALUES (
			?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range sets {
		_, err := stmt.Exec(
			s.ID, s.Name,
			s.Item1, s.Item2, s.Item3, s.Item4, s.Item5, s.Item6, s.Item7, s.Item8, s.Item9, s.Item10,
			s.SkillID, s.SkillLevel,
			s.Bonus1, s.Bonus2, s.Bonus3, s.Bonus4, s.Bonus5, s.Bonus6, s.Bonus7, s.Bonus8,
			s.Spell1, s.Spell2, s.Spell3, s.Spell4, s.Spell5, s.Spell6, s.Spell7, s.Spell8,
		)
		if err != nil {
			fmt.Printf("Error importing item set %d: %v\n", s.ID, err)
			continue
		}
	}

	return tx.Commit()
}

// CheckAndImportItemSets checks if itemsets table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportItemSets(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM itemsets").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		path := fmt.Sprintf("%s/item_sets.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Importing Item Sets from JSON...")
			return r.ImportItemSets(path)
		}
	}
	return nil
}
