package database

import "fmt"

// SpellSkillCategory represents a top-level category for spells
type SpellSkillCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SpellSkill represents a skill that contains spells
type SpellSkill struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"categoryId"`
	Name       string `json:"name"`
	SpellCount int    `json:"spellCount"`
}

// GetSpellSkillCategories returns all spell skill categories
func (r *ItemRepository) GetSpellSkillCategories() ([]*SpellSkillCategory, error) {
	rows, err := r.db.DB().Query(`
		SELECT id, name FROM spell_skill_categories ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*SpellSkillCategory
	for rows.Next() {
		c := &SpellSkillCategory{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetSpellSkillsByCategory returns all skills in a category with spell counts
func (r *ItemRepository) GetSpellSkillsByCategory(categoryID int) ([]*SpellSkill, error) {
	rows, err := r.db.DB().Query(`
		SELECT s.id, s.category_id, s.name, 
		       (SELECT COUNT(*) FROM spell_skill_spells ss WHERE ss.skill_id = s.id) as spell_count
		FROM spell_skills s 
		WHERE s.category_id = ?
		ORDER BY s.name
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*SpellSkill
	for rows.Next() {
		s := &SpellSkill{}
		if err := rows.Scan(&s.ID, &s.CategoryID, &s.Name, &s.SpellCount); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// GetSpellsBySkill returns all spells for a given skill
// GetSpellsBySkill returns all spells for a given skill
func (r *ItemRepository) GetSpellsBySkill(skillID int, nameFilter string) ([]*Spell, error) {
	whereClause := "WHERE ss.skill_id = ?"
	args := []interface{}{skillID}

	if nameFilter != "" {
		whereClause += " AND sp.name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	query := fmt.Sprintf(`
		SELECT sp.entry, sp.name, sp.description
		FROM spells sp
		INNER JOIN spell_skill_spells ss ON ss.spell_id = sp.entry
		%s
		ORDER BY sp.name
		LIMIT 10000
	`, whereClause)

	rows, err := r.db.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []*Spell
	for rows.Next() {
		s := &Spell{}
		var desc *string
		// Simplified scan: removed subname and icon_id which don't exist
		if err := rows.Scan(&s.Entry, &s.Name, &desc); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		if desc != nil {
			s.Description = *desc
		}
		spells = append(spells, s)
	}
	return spells, nil
}
