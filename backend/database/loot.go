package database

import (
	"database/sql"
	"math"
)

type LootItem struct {
	ItemID   int     `json:"itemId"`
	ItemName string  `json:"itemName"`
	Icon     string  `json:"icon"`
	Quality  int     `json:"quality"`
	Chance   float64 `json:"chance"`
	MinCount int     `json:"minCount"`
	MaxCount int     `json:"maxCount"`
}

// GetCreatureLoot returns the flattened loot table for a creature
func (r *ItemRepository) GetCreatureLoot(creatureEntry int) ([]*LootItem, error) {
	// 1. Get loot_id from creatures table
	var lootID int
	err := r.db.DB().QueryRow("SELECT loot_id FROM creatures WHERE entry = ?", creatureEntry).Scan(&lootID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*LootItem{}, nil
		}
		return nil, err
	}

	if lootID == 0 {
		return []*LootItem{}, nil // No loot
	}

	// 2. Process loot recursively
	lootMap := make(map[int]*LootItem)
	if err := r.processLoot(lootID, 1.0, false, lootMap, 0); err != nil {
		return nil, err
	}

	// 3. Convert map to slice and enrich with item info
	var lootList []*LootItem
	for itemID, item := range lootMap {
		// Enrich with name, icon, quality
		var name, icon string
		var quality int
		err := r.db.DB().QueryRow("SELECT name, quality, icon FROM items WHERE entry = ?", itemID).Scan(&name, &quality, &icon)
		if err == nil {
			item.ItemName = name
			item.Quality = quality
			item.Icon = icon
			lootList = append(lootList, item)
		}
	}

	return lootList, nil
}

// processLoot recursively processes loot tables
// entry: loot table entry
// multiplier: chance multiplier (1.0 for root)
// isRef: whether querying reference_loot table
// results: map to store aggregated results
// depth: recursion depth limit
func (r *ItemRepository) processLoot(entry int, multiplier float64, isRef bool, results map[int]*LootItem, depth int) error {
	if depth > 10 {
		return nil // Prevent infinite recursion
	}

	tableName := "creature_loot"
	if isRef {
		tableName = "reference_loot"
	}

	rows, err := r.db.DB().Query(`
		SELECT item, chance, mincount_or_ref, maxcount, groupid
		FROM `+tableName+` 
		WHERE entry = ?
	`, entry)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var itemID, minOrRef, maxCount, groupID int
		var chance float64
		if err := rows.Scan(&itemID, &chance, &minOrRef, &maxCount, &groupID); err != nil {
			continue
		}

		// Calculate actual chance
		// Note: Negative chance means quest drop, we take absolute value for display.
		absChance := math.Abs(chance)
		currentChance := absChance * multiplier

		if minOrRef < 0 {
			// Reference
			refID := -minOrRef
			// Recursively process reference
			// For references, the chance defined here applies to the whole sub-table?
			// Yes, usually.
			if err := r.processLoot(refID, currentChance/100.0, true, results, depth+1); err != nil {
				return err
			}
		} else {
			// Item
			// Provide a floor for very small chances
			if currentChance < 0.0001 {
				currentChance = 0.0001
			}

			// If item already exists, maybe sum chances (e.g. from different groups)?
			// Or just take the max? Summing is safer for 'independent' drops, but complex.
			// Simplified: Sum chances.
			if existing, ok := results[itemID]; ok {
				existing.Chance += currentChance
				// Merge counts?
			} else {
				results[itemID] = &LootItem{
					ItemID:   itemID,
					Chance:   currentChance,
					MinCount: minOrRef, // Positive value is min count
					MaxCount: maxCount,
				}
			}
		}
	}
	return nil
}
