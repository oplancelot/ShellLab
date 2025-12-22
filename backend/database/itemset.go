package database

// ItemSetBrowse represents an item set for browsing list
type ItemSetBrowse struct {
	ItemSetID  int    `json:"itemsetId"`
	Name       string `json:"name"`
	ItemIDs    []int  `json:"itemIds"`
	ItemCount  int    `json:"itemCount"`
	SkillID    int    `json:"skillId"`
	SkillLevel int    `json:"skillLevel"`
}

// ItemSetDetail includes items with their details
type ItemSetDetail struct {
	ItemSetID int        `json:"itemsetId"`
	Name      string     `json:"name"`
	Items     []*Item    `json:"items"`
	Bonuses   []SetBonus `json:"bonuses"`
}

// GetItemSets returns all item sets for browsing
func (r *ItemRepository) GetItemSets() ([]*ItemSetBrowse, error) {
	rows, err := r.db.DB().Query(`
		SELECT 
			itemset_id, name,
			item1, item2, item3, item4, item5, item6, item7, item8, item9, item10,
			skill_id, skill_level
		FROM itemsets
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*ItemSetBrowse
	for rows.Next() {
		set := &ItemSetBrowse{}
		var items [10]int
		err := rows.Scan(
			&set.ItemSetID, &set.Name,
			&items[0], &items[1], &items[2], &items[3], &items[4],
			&items[5], &items[6], &items[7], &items[8], &items[9],
			&set.SkillID, &set.SkillLevel,
		)
		if err != nil {
			continue
		}

		// Filter out zero item IDs
		for _, itemID := range items {
			if itemID > 0 {
				set.ItemIDs = append(set.ItemIDs, itemID)
			}
		}
		set.ItemCount = len(set.ItemIDs)

		sets = append(sets, set)
	}

	return sets, nil
}

// GetItemSetDetail returns detailed information about an item set
func (r *ItemRepository) GetItemSetDetail(itemSetID int) (*ItemSetDetail, error) {
	// First get the set info
	row := r.db.DB().QueryRow(`
		SELECT 
			itemset_id, name,
			item1, item2, item3, item4, item5, item6, item7, item8, item9, item10,
			spell1, spell2, spell3, spell4, spell5, spell6, spell7, spell8,
			bonus1, bonus2, bonus3, bonus4, bonus5, bonus6, bonus7, bonus8
		FROM itemsets
		WHERE itemset_id = ?
	`, itemSetID)

	var name string
	var items [10]int
	var spells [8]int
	var bonuses [8]int
	var setID int

	err := row.Scan(
		&setID, &name,
		&items[0], &items[1], &items[2], &items[3], &items[4],
		&items[5], &items[6], &items[7], &items[8], &items[9],
		&spells[0], &spells[1], &spells[2], &spells[3],
		&spells[4], &spells[5], &spells[6], &spells[7],
		&bonuses[0], &bonuses[1], &bonuses[2], &bonuses[3],
		&bonuses[4], &bonuses[5], &bonuses[6], &bonuses[7],
	)
	if err != nil {
		return nil, err
	}

	detail := &ItemSetDetail{
		ItemSetID: setID,
		Name:      name,
	}

	// Get item details for each item in the set
	for _, itemID := range items {
		if itemID > 0 {
			item, err := r.GetItemByID(itemID)
			if err == nil && item != nil {
				detail.Items = append(detail.Items, item)
			}
		}
	}

	// Build bonuses list
	for i := 0; i < 8; i++ {
		if spells[i] > 0 && bonuses[i] > 0 {
			detail.Bonuses = append(detail.Bonuses, SetBonus{
				Threshold: bonuses[i],
				SpellID:   spells[i],
			})
		}
	}

	return detail, nil
}
