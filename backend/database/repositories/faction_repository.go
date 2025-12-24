package repositories

import (
	"database/sql"

	"shelllab/backend/database/models"
)

// FactionRepository handles faction-related database operations
type FactionRepository struct {
	db *sql.DB
}

// NewFactionRepository creates a new faction repository
func NewFactionRepository(db *sql.DB) *FactionRepository {
	return &FactionRepository{db: db}
}

// GetFactions returns all factions ordered by side and name
func (r *FactionRepository) GetFactions() ([]*models.Faction, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, side, category_id
		FROM factions
		ORDER BY side, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var factions []*models.Faction
	for rows.Next() {
		f := &models.Faction{}
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
