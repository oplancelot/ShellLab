package repositories

import (
	"database/sql"
	"fmt"

	"shelllab/backend/database/models"
)

// SpellRepository handles spell-related database operations
type SpellRepository struct {
	db *sql.DB
}

// NewSpellRepository creates a new spell repository
func NewSpellRepository(db *sql.DB) *SpellRepository {
	return &SpellRepository{db: db}
}

// SearchSpells searches for spells by name
func (r *SpellRepository) SearchSpells(query string) ([]*models.Spell, error) {
	rows, err := r.db.Query(`
		SELECT entry, name, description, iconName
		FROM spell_template
		WHERE name LIKE ?
		ORDER BY length(name), name
		LIMIT 100
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []*models.Spell
	for rows.Next() {
		s := &models.Spell{}
		var desc *string
		if err := rows.Scan(&s.Entry, &s.Name, &desc, &s.Icon); err != nil {
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

// GetSpellSkillCategories returns all spell skill categories
func (r *SpellRepository) GetSpellSkillCategories() ([]*models.SpellSkillCategory, error) {
	rows, err := r.db.Query(`
		SELECT id, name FROM spell_skill_categories ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.SpellSkillCategory
	for rows.Next() {
		c := &models.SpellSkillCategory{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetSpellSkillsByCategory returns all skills in a category with spell counts
func (r *SpellRepository) GetSpellSkillsByCategory(categoryID int) ([]*models.SpellSkill, error) {
	rows, err := r.db.Query(`
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

	var skills []*models.SpellSkill
	for rows.Next() {
		s := &models.SpellSkill{}
		if err := rows.Scan(&s.ID, &s.CategoryID, &s.Name, &s.SpellCount); err != nil {
			continue
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// GetSpellsBySkill returns all spells for a given skill
func (r *SpellRepository) GetSpellsBySkill(skillID int, nameFilter string) ([]*models.Spell, error) {
	whereClause := "WHERE ss.skill_id = ?"
	args := []interface{}{skillID}

	if nameFilter != "" {
		whereClause += " AND sp.name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	query := fmt.Sprintf(`
		SELECT sp.entry, sp.name, sp.description, sp.iconName
		FROM spell_template sp
		INNER JOIN spell_skill_spells ss ON ss.spell_id = sp.entry
		%s
		ORDER BY sp.name
		LIMIT 10000
	`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []*models.Spell
	for rows.Next() {
		s := &models.Spell{}
		var desc *string
		if err := rows.Scan(&s.Entry, &s.Name, &desc, &s.Icon); err != nil {
			continue
		}
		if desc != nil {
			s.Description = *desc
		}
		spells = append(spells, s)
	}
	return spells, nil
}

// GetSpellByID retrieves a single spell by ID
func (r *SpellRepository) GetSpellByID(entry int) (*models.Spell, error) {
	s := &models.Spell{}
	var desc *string
	err := r.db.QueryRow(`
		SELECT entry, name, description, iconName
		FROM spell_template WHERE entry = ?
	`, entry).Scan(&s.Entry, &s.Name, &desc, &s.Icon)
	if err != nil {
		return nil, err
	}
	if desc != nil {
		s.Description = *desc
	}
	return s, nil
}

// GetSpellDescription retrieves spell description and base points
func (r *SpellRepository) GetSpellDescription(spellID int) (string, []int) {
	var desc string
	var bp1, bp2, bp3 int
	err := r.db.QueryRow(`
		SELECT description, effectBasePoints1, effectBasePoints2, effectBasePoints3
		FROM spell_template WHERE entry = ?
	`, spellID).Scan(&desc, &bp1, &bp2, &bp3)
	if err != nil {
		return "", nil
	}
	return desc, []int{bp1, bp2, bp3}
}
