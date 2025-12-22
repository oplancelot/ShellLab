package database

import "fmt"

// Spell represents a WoW spell
type Spell struct {
	Entry       int    `json:"entry"`
	Name        string `json:"name"`
	SubName     string `json:"subname"` // Rank or subtext
	Description string `json:"description"`
	IconID      int    `json:"iconId"`
}

// SearchSpells searches for spells by name
func (r *ItemRepository) SearchSpells(query string) ([]*Spell, error) {
	rows, err := r.db.DB().Query(`
		SELECT entry, name, subname, description, icon_id
		FROM spells
		WHERE name LIKE ?
		ORDER BY length(name), name
		LIMIT 100
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []*Spell
	for rows.Next() {
		s := &Spell{}
		var subname, desc *string
		if err := rows.Scan(&s.Entry, &s.Name, &subname, &desc, &s.IconID); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		if subname != nil {
			s.SubName = *subname
		}
		if desc != nil {
			s.Description = *desc
		}
		spells = append(spells, s)
	}
	return spells, nil
}
