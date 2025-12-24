package repositories

import (
	"database/sql"
	"fmt"

	"shelllab/backend/database/helpers"
	"shelllab/backend/database/models"
)

// CreatureRepository handles creature-related database operations
type CreatureRepository struct {
	db *sql.DB
}

// NewCreatureRepository creates a new creature repository
func NewCreatureRepository(db *sql.DB) *CreatureRepository {
	return &CreatureRepository{db: db}
}

// GetCreatureTypes returns all creature types with counts
func (r *CreatureRepository) GetCreatureTypes() ([]*models.CreatureType, error) {
	rows, err := r.db.Query(`
		SELECT creature_type, COUNT(*) as count
		FROM creatures
		GROUP BY creature_type
		ORDER BY creature_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*models.CreatureType
	for rows.Next() {
		t := &models.CreatureType{}
		if err := rows.Scan(&t.Type, &t.Count); err != nil {
			continue
		}
		t.Name = helpers.GetCreatureTypeName(t.Type)
		types = append(types, t)
	}

	return types, nil
}

// GetCreaturesByType returns creatures filtered by type
func (r *CreatureRepository) GetCreaturesByType(creatureType int, nameFilter string, limit, offset int) ([]*models.Creature, int, error) {
	whereClause := "WHERE creature_type = ?"
	args := []interface{}{creatureType}

	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	// Count
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM creatures %s", whereClause)
	err := r.db.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Data
	dataArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, subname, level_min, level_max, 
			health_min, health_max, mana_min, mana_max,
			creature_type, creature_rank, faction, npc_flags
		FROM creatures
		%s
		ORDER BY level_max DESC, name
		LIMIT ? OFFSET ?
	`, whereClause)

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var creatures []*models.Creature
	for rows.Next() {
		c := &models.Creature{}
		var subname *string
		err := rows.Scan(
			&c.Entry, &c.Name, &subname, &c.LevelMin, &c.LevelMax,
			&c.HealthMin, &c.HealthMax, &c.ManaMin, &c.ManaMax,
			&c.Type, &c.Rank, &c.Faction, &c.NPCFlags,
		)
		if err != nil {
			continue
		}
		if subname != nil {
			c.Subname = *subname
		}
		c.TypeName = helpers.GetCreatureTypeName(c.Type)
		c.RankName = helpers.GetCreatureRankName(c.Rank)
		creatures = append(creatures, c)
	}

	return creatures, count, nil
}

// SearchCreatures searches for creatures by name
func (r *CreatureRepository) SearchCreatures(query string, limit int) ([]*models.Creature, error) {
	rows, err := r.db.Query(`
		SELECT entry, name, subname, level_min, level_max, 
			health_min, health_max, mana_min, mana_max,
			creature_type, creature_rank, faction, npc_flags
		FROM creatures
		WHERE name LIKE ?
		ORDER BY length(name), name
		LIMIT ?
	`, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creatures []*models.Creature
	for rows.Next() {
		c := &models.Creature{}
		var subname *string
		err := rows.Scan(
			&c.Entry, &c.Name, &subname, &c.LevelMin, &c.LevelMax,
			&c.HealthMin, &c.HealthMax, &c.ManaMin, &c.ManaMax,
			&c.Type, &c.Rank, &c.Faction, &c.NPCFlags,
		)
		if err != nil {
			continue
		}
		if subname != nil {
			c.Subname = *subname
		}
		c.TypeName = helpers.GetCreatureTypeName(c.Type)
		c.RankName = helpers.GetCreatureRankName(c.Rank)
		creatures = append(creatures, c)
	}

	return creatures, nil
}

// GetCreatureByID retrieves a single creature by ID
func (r *CreatureRepository) GetCreatureByID(entry int) (*models.Creature, error) {
	c := &models.Creature{}
	var subname *string
	err := r.db.QueryRow(`
		SELECT entry, name, subname, level_min, level_max, 
			health_min, health_max, mana_min, mana_max,
			creature_type, creature_rank, faction, npc_flags
		FROM creatures WHERE entry = ?
	`, entry).Scan(
		&c.Entry, &c.Name, &subname, &c.LevelMin, &c.LevelMax,
		&c.HealthMin, &c.HealthMax, &c.ManaMin, &c.ManaMax,
		&c.Type, &c.Rank, &c.Faction, &c.NPCFlags,
	)
	if err != nil {
		return nil, err
	}
	if subname != nil {
		c.Subname = *subname
	}
	c.TypeName = helpers.GetCreatureTypeName(c.Type)
	c.RankName = helpers.GetCreatureRankName(c.Rank)
	return c, nil
}

// GetCreatureCount returns the total number of creatures
func (r *CreatureRepository) GetCreatureCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM creatures").Scan(&count)
	return count, err
}

// GetCreatureDetail returns full creature information with loot and quests
func (r *CreatureRepository) GetCreatureDetail(entry int) (*models.CreatureDetail, error) {
	creature, err := r.GetCreatureByID(entry)
	if err != nil {
		return nil, err
	}

	detail := &models.CreatureDetail{Creature: creature}

	// Get loot
	lootRepo := NewLootRepository(r.db)
	loot, err := lootRepo.GetCreatureLoot(entry)
	if err == nil {
		detail.Loot = loot
	}

	// Get quests this creature starts
	rows, err := r.db.Query(`
		SELECT q.entry, q.title
		FROM npc_quest_start nqs
		JOIN quests q ON nqs.quest = q.entry
		WHERE nqs.entry = ?
	`, entry)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			qr := &models.QuestRelation{Type: "quest"}
			rows.Scan(&qr.Entry, &qr.Name)
			detail.StartsQuests = append(detail.StartsQuests, qr)
		}
	}

	// Get quests this creature ends
	rows2, err := r.db.Query(`
		SELECT q.entry, q.title
		FROM npc_quest_end nqe
		JOIN quests q ON nqe.quest = q.entry
		WHERE nqe.entry = ?
	`, entry)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			qr := &models.QuestRelation{Type: "quest"}
			rows2.Scan(&qr.Entry, &qr.Name)
			detail.EndsQuests = append(detail.EndsQuests, qr)
		}
	}

	return detail, nil
}
