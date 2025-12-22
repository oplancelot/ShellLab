package database

type CreatureDetail struct {
	*Creature
	Loot         []*LootItem      `json:"loot"`
	StartsQuests []*QuestRelation `json:"startsQuests"`
	EndsQuests   []*QuestRelation `json:"endsQuests"`
}

type QuestRelationItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Level int    `json:"level"`
}

// GetCreatureDetail returns full details for a creature
func (r *ItemRepository) GetCreatureDetail(entry int) (*CreatureDetail, error) {
	// 1. Get Basic Info
	// Re-use functionality or duplicate query for simplicity?
	// Let's duplicate the single row fetch for efficiency.
	var c Creature
	var subname *string
	err := r.db.DB().QueryRow(`
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
	c.TypeName = getCreatureTypeName(c.Type)
	c.RankName = getCreatureRankName(c.Rank)

	detail := &CreatureDetail{
		Creature:     &c,
		Loot:         []*LootItem{},
		StartsQuests: []*QuestRelation{},
		EndsQuests:   []*QuestRelation{},
	}

	// 2. Get Loot (Reuse existing function logic is hard inside here without method call)
	// Calling the public method:
	loot, _ := r.GetCreatureLoot(entry)
	detail.Loot = loot

	// 3. Get Quests Started
	rowsS, err := r.db.DB().Query(`
        SELECT q.entry, IFNULL(q.title, '') FROM npc_quest_start n
        JOIN quests q ON n.quest = q.entry
        WHERE n.entry = ?
    `, entry)
	if err == nil {
		defer rowsS.Close()
		for rowsS.Next() {
			qr := &QuestRelation{Type: "quest"}
			rowsS.Scan(&qr.Entry, &qr.Name)
			detail.StartsQuests = append(detail.StartsQuests, qr)
		}
	}

	// 4. Get Quests Ended
	rowsE, err := r.db.DB().Query(`
        SELECT q.entry, IFNULL(q.title, '') FROM npc_quest_end n
        JOIN quests q ON n.quest = q.entry
        WHERE n.entry = ?
    `, entry)
	if err == nil {
		defer rowsE.Close()
		for rowsE.Next() {
			qr := &QuestRelation{Type: "quest"}
			rowsE.Scan(&qr.Entry, &qr.Name)
			detail.EndsQuests = append(detail.EndsQuests, qr)
		}
	}

	return detail, nil
}
