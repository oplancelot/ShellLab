package database

import "fmt"

// Creature represents a WoW NPC
type Creature struct {
	Entry     int    `json:"entry"`
	Name      string `json:"name"`
	Subname   string `json:"subname,omitempty"`
	LevelMin  int    `json:"levelMin"`
	LevelMax  int    `json:"levelMax"`
	HealthMin int    `json:"healthMin"`
	HealthMax int    `json:"healthMax"`
	ManaMin   int    `json:"manaMin"`
	ManaMax   int    `json:"manaMax"`
	Type      int    `json:"type"`
	TypeName  string `json:"typeName"`
	Rank      int    `json:"rank"`
	RankName  string `json:"rankName"`
	Faction   int    `json:"faction"`
	NPCFlags  int    `json:"npcFlags"`
}

// CreatureType represents a creature type category
type CreatureType struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetCreatureTypes returns all creature types with counts
func (r *ItemRepository) GetCreatureTypes() ([]*CreatureType, error) {
	rows, err := r.db.DB().Query(`
		SELECT creature_type, COUNT(*) as count
		FROM creatures
		GROUP BY creature_type
		ORDER BY creature_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*CreatureType
	for rows.Next() {
		t := &CreatureType{}
		if err := rows.Scan(&t.Type, &t.Count); err != nil {
			continue
		}
		t.Name = getCreatureTypeName(t.Type)
		types = append(types, t)
	}

	return types, nil
}

// GetCreaturesByType returns creatures filtered by type
func (r *ItemRepository) GetCreaturesByType(creatureType int, nameFilter string, limit, offset int) ([]*Creature, int, error) {
	// Build WHERE clause
	whereClause := "WHERE creature_type = ?"
	args := []interface{}{creatureType}

	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	// Count
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM creatures %s", whereClause)
	err := r.db.DB().QueryRow(countQuery, args...).Scan(&count)
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

	rows, err := r.db.DB().Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var creatures []*Creature
	for rows.Next() {
		c := &Creature{}
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
		c.TypeName = getCreatureTypeName(c.Type)
		c.RankName = getCreatureRankName(c.Rank)
		creatures = append(creatures, c)
	}

	return creatures, count, nil
}

// SearchCreatures searches for creatures by name
func (r *ItemRepository) SearchCreatures(query string, limit int) ([]*Creature, error) {
	rows, err := r.db.DB().Query(`
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

	var creatures []*Creature
	for rows.Next() {
		c := &Creature{}
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
		c.TypeName = getCreatureTypeName(c.Type)
		c.RankName = getCreatureRankName(c.Rank)
		creatures = append(creatures, c)
	}

	return creatures, nil
}

func getCreatureTypeName(t int) string {
	typeNames := map[int]string{
		0:  "None",
		1:  "Beast",
		2:  "Dragonkin",
		3:  "Demon",
		4:  "Elemental",
		5:  "Giant",
		6:  "Undead",
		7:  "Humanoid",
		8:  "Critter",
		9:  "Mechanical",
		10: "Not Specified",
		11: "Totem",
	}
	if name, ok := typeNames[t]; ok {
		return name
	}
	return "Unknown"
}

func getCreatureRankName(r int) string {
	rankNames := map[int]string{
		0: "Normal",
		1: "Elite",
		2: "Rare Elite",
		3: "Boss",
		4: "Rare",
	}
	if name, ok := rankNames[r]; ok {
		return name
	}
	return "Normal"
}
