package database

import "fmt"

// Quest represents a WoW quest
type Quest struct {
	Entry            int    `json:"entry"`
	Title            string `json:"title"`
	QuestLevel       int    `json:"questLevel"`
	MinLevel         int    `json:"minLevel"`
	Type             int    `json:"type"`
	ZoneOrSort       int    `json:"zoneOrSort"`
	CategoryName     string `json:"categoryName"`
	RequiredRaces    int    `json:"requiredRaces"`
	RequiredClasses  int    `json:"requiredClasses"`
	SrcItem          int    `json:"srcItemId"`
	RewardXP         int    `json:"rewardXp"`
	RewardMoney      int    `json:"rewardMoney"`
	PrevQuestID      int    `json:"prevQuestId"`
	NextQuestID      int    `json:"nextQuestId"`
	ExclusiveGroup   int    `json:"exclusiveGroup"`
	NextQuestInChain int    `json:"nextQuestInChain"`
}

// QuestCategory represents a zone or category for quests
type QuestCategory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetQuestCategories returns all quest categories (zones and sorts) with quest counts
func (r *ItemRepository) GetQuestCategories() ([]*QuestCategory, error) {
	// First get categories from our table
	rows, err := r.db.DB().Query(`
		SELECT id, name FROM quest_categories ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make(map[int]*QuestCategory)
	var catList []*QuestCategory

	for rows.Next() {
		cat := &QuestCategory{}
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			continue
		}
		categories[cat.ID] = cat
		catList = append(catList, cat)
	}

	// Now count quests per category
	rows2, err := r.db.DB().Query(`
		SELECT zone_or_sort, COUNT(*) 
		FROM quests 
		GROUP BY zone_or_sort
	`)
	if err != nil {
		// Just return categories without counts if this fails? better to log error
		return catList, nil
	}
	defer rows2.Close()

	for rows2.Next() {
		var zoneID, count int
		if err := rows2.Scan(&zoneID, &count); err != nil {
			continue
		}
		if cat, exists := categories[zoneID]; exists {
			cat.Count = count
		} else {
			// Category exists in quests but not in category table, create generic one
			// Prefer not to show unknown categories in production, but for debug it's useful.
			// Currently ignoring to keep list clean, or we can add "Unknown (ID)" categories.
		}
	}

	// Filter out categories with 0 quests
	var activeCats []*QuestCategory
	for _, cat := range catList {
		if cat.Count > 0 {
			activeCats = append(activeCats, cat)
		}
	}

	return activeCats, nil
}

// GetQuestsByCategory returns quests filtered by category (zone or sort)
func (r *ItemRepository) GetQuestsByCategory(categoryID int) ([]*Quest, error) {
	rows, err := r.db.DB().Query(`
		SELECT entry, IFNULL(title,''), IFNULL(quest_level,0), IFNULL(min_level,0), 
			IFNULL(type,0), IFNULL(zone_or_sort,0),
			IFNULL(rew_xp,0), IFNULL(rew_money,0),
			IFNULL(required_races,0), IFNULL(required_classes,0), IFNULL(src_item_id,0),
			IFNULL(prev_quest_id,0), IFNULL(next_quest_id,0), IFNULL(exclusive_group,0), IFNULL(next_quest_in_chain,0)
		FROM quests
		WHERE zone_or_sort = ?
		ORDER BY quest_level, title
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quests []*Quest
	for rows.Next() {
		q := &Quest{}
		err := rows.Scan(
			&q.Entry, &q.Title, &q.QuestLevel, &q.MinLevel,
			&q.Type, &q.ZoneOrSort,
			&q.RewardXP, &q.RewardMoney,
			&q.RequiredRaces, &q.RequiredClasses, &q.SrcItem,
			&q.PrevQuestID, &q.NextQuestID, &q.ExclusiveGroup, &q.NextQuestInChain,
		)
		if err != nil {
			fmt.Printf("Error scanning quest list: %v\n", err)
			continue
		}
		quests = append(quests, q)
	}
	return quests, nil
}

// SearchQuests searches for quests by title
func (r *ItemRepository) SearchQuests(query string) ([]*Quest, error) {
	rows, err := r.db.DB().Query(`
		SELECT q.entry, IFNULL(q.title,''), IFNULL(q.quest_level,0), IFNULL(q.min_level,0), 
			IFNULL(q.type,0), IFNULL(q.zone_or_sort,0),
			IFNULL(q.rew_xp,0), IFNULL(q.rew_money,0),
			IFNULL(q.required_races,0), IFNULL(q.required_classes,0), IFNULL(q.src_item_id,0),
			IFNULL(q.prev_quest_id,0), IFNULL(q.next_quest_id,0), IFNULL(q.exclusive_group,0), IFNULL(q.next_quest_in_chain,0),
			c.name
		FROM quests q
		LEFT JOIN quest_categories c ON q.zone_or_sort = c.id
		WHERE q.title LIKE ?
		ORDER BY length(q.title), q.title
		LIMIT 50
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quests []*Quest
	for rows.Next() {
		q := &Quest{}
		var catName *string
		err := rows.Scan(
			&q.Entry, &q.Title, &q.QuestLevel, &q.MinLevel,
			&q.Type, &q.ZoneOrSort,
			&q.RewardXP, &q.RewardMoney,
			&q.RequiredRaces, &q.RequiredClasses, &q.SrcItem,
			&q.PrevQuestID, &q.NextQuestID, &q.ExclusiveGroup, &q.NextQuestInChain,
			&catName,
		)
		if err != nil {
			fmt.Printf("Error scanning quest search: %v\n", err)
			continue
		}
		if catName != nil {
			q.CategoryName = *catName
		}
		quests = append(quests, q)
	}
	return quests, nil
}
