package database

// Faction represents a WoW reputation faction
type Faction struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Side        int    `json:"side"` // 1=Alliance, 2=Horde, 3=Both? (Need to check DB values)
	CategoryId  int    `json:"categoryId"`
}

// GetFactions returns all factions ordered by side and name
func (r *ItemRepository) GetFactions() ([]*Faction, error) {
	rows, err := r.db.DB().Query(`
		SELECT id, name, description, side, category_id
		FROM factions
		ORDER BY side, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var factions []*Faction
	for rows.Next() {
		f := &Faction{}
		var desc *string
		if err := rows.Scan(&f.ID, &f.Name, &desc, &f.Side, &f.CategoryId); err != nil {
			continue
		}
		if desc != nil {
			f.Description = *desc
		}
		factions = append(factions, f)
	}
	return factions, nil
}
