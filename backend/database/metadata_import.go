package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type ZoneEntry struct {
	AreaID int    `json:"areatableID"`
	MapID  int    `json:"mapID"`
	Name   string `json:"name_loc0"`
}

type SkillEntry struct {
	ID         int    `json:"skillID"`
	CategoryID int    `json:"categoryID"`
	Name       string `json:"name_loc0"`
}

type SkillLineAbilityEntry struct {
	SkillID int `json:"skillID"`
	SpellID int `json:"spellID"`
}

// CheckAndImportMetadata handles all metadata imports
func (r *ItemRepository) CheckAndImportMetadata(dataDir string) error {
	// Init static groups first
	r.initStaticMetadata()

	// 1. Import Spell Skills
	var skillCount int
	r.db.DB().QueryRow("SELECT COUNT(*) FROM spell_skills").Scan(&skillCount)
	if skillCount == 0 {
		if err := r.importSkills(dataDir); err != nil {
			fmt.Printf("Warning: Failed to import skills: %v\n", err)
		} else {
			fmt.Println("✓ Spell Skills imported")
		}
	}

	// 2. Import Quest Zones
	var questCatCount int
	r.db.DB().QueryRow("SELECT COUNT(*) FROM quest_categories_enhanced").Scan(&questCatCount)
	if questCatCount == 0 {
		if err := r.importQuestZones(dataDir); err != nil {
			fmt.Printf("Warning: Failed to import quest zones: %v\n", err)
		} else {
			fmt.Println("✓ Quest Zones imported")
		}
	}

	return nil
}

func (r *ItemRepository) initStaticMetadata() {
	// Quest Groups
	groups := []struct {
		ID   int
		Name string
	}{
		{0, "Eastern Kingdoms"},
		{1, "Kalimdor"},
		{2, "Dungeons"},
		{3, "Raids"},
		{4, "Classes"},
		{5, "Professions"},
		{6, "Battlegrounds"},
		{7, "Misc"},
	}
	r.db.DB().Exec("DELETE FROM quest_category_groups") // Ensure clean state or use INSERT OR IGNORE
	for _, g := range groups {
		r.db.DB().Exec("INSERT OR IGNORE INTO quest_category_groups (id, name) VALUES (?, ?)", g.ID, g.Name)
	}

	// Spell Categories
	spellCats := []struct {
		ID   int
		Name string
	}{
		{6, "Weapon Skills"},
		{8, "Armor Proficiencies"},
		{10, "Languages"},
		{7, "Class Skills"},
		{9, "Professions"},
		{11, "Racial Traits"},
	}
	r.db.DB().Exec("DELETE FROM spell_skill_categories")
	for _, c := range spellCats {
		r.db.DB().Exec("INSERT OR IGNORE INTO spell_skill_categories (id, name) VALUES (?, ?)", c.ID, c.Name)
	}
}

func (r *ItemRepository) importSkills(dataDir string) error {
	// 1. Skills
	file, err := os.Open(fmt.Sprintf("%s/skills.json", dataDir))
	if err != nil {
		return err
	}
	defer file.Close()

	var skills []SkillEntry
	if err := json.NewDecoder(file).Decode(&skills); err != nil {
		return err
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	skillStmt, _ := tx.Prepare("REPLACE INTO spell_skills (id, category_id, name) VALUES (?, ?, ?)")
	defer skillStmt.Close()

	for _, s := range skills {
		skillStmt.Exec(s.ID, s.CategoryID, s.Name)
	}

	// 2. Skill Line Abilities
	file2, err := os.Open(fmt.Sprintf("%s/skill_line_abilities.json", dataDir))
	if err != nil {
		return err
	}
	defer file2.Close()

	var abilities []SkillLineAbilityEntry
	if err := json.NewDecoder(file2).Decode(&abilities); err != nil {
		return err
	}

	abilityStmt, _ := tx.Prepare("REPLACE INTO spell_skill_spells (skill_id, spell_id) VALUES (?, ?)")
	defer abilityStmt.Close()

	for _, a := range abilities {
		abilityStmt.Exec(a.SkillID, a.SpellID)
	}

	return tx.Commit()
}

func (r *ItemRepository) importQuestZones(dataDir string) error {
	file, err := os.Open(fmt.Sprintf("%s/zones.json", dataDir))
	if err != nil {
		return err
	}
	defer file.Close()

	var zones []ZoneEntry
	if err := json.NewDecoder(file).Decode(&zones); err != nil {
		return err
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, _ := tx.Prepare("REPLACE INTO quest_categories_enhanced (id, group_id, name) VALUES (?, ?, ?)")
	defer stmt.Close()

	for _, z := range zones {
		groupID := 7 // Misc default
		if z.MapID == 0 {
			groupID = 0 // WK
		} else if z.MapID == 1 {
			groupID = 1 // Kalimdor
		} else {
			// Simple heuristic for instances (mapID usually > 1)
			groupID = 2 // Dungeons
		}
		stmt.Exec(z.AreaID, groupID, z.Name)
	}

	return tx.Commit()
}
