package database

// Quest represents a WoW quest
type Quest struct {
	Entry           int    `json:"entry"`
	Title           string `json:"title"`
	QuestLevel      int    `json:"questLevel"`
	MinLevel        int    `json:"minLevel"`
	MaxLevel        int    `json:"maxLevel"`
	Type            int    `json:"type"`
	ZoneOrSort      int    `json:"zoneOrSort"`
	CategoryName    string `json:"categoryName"`
	RequiredRaces   int    `json:"requiredRaces"`
	RequiredClasses int    `json:"requiredClasses"`
	SrcItem         int    `json:"srcItemId"`
	RewardXP        int    `json:"rewardXp"`
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
		SELECT entry, title, quest_level, min_level, max_level, 
			quest_type, zone_or_sort, required_races, required_classes,
			src_item_id, reward_xp
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
			&q.Entry, &q.Title, &q.QuestLevel, &q.MinLevel, &q.MaxLevel,
			&q.Type, &q.ZoneOrSort, &q.RequiredRaces, &q.RequiredClasses,
			&q.SrcItem, &q.RewardXP,
		)
		if err != nil {
			continue
		}
		quests = append(quests, q)
	}
	return quests, nil
}

// SearchQuests searches for quests by title
func (r *ItemRepository) SearchQuests(query string) ([]*Quest, error) {
	rows, err := r.db.DB().Query(`
		SELECT q.entry, q.title, q.quest_level, q.min_level, q.max_level, 
			q.quest_type, q.zone_or_sort, q.required_races, q.required_classes,
			q.src_item_id, q.reward_xp, c.name
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
			&q.Entry, &q.Title, &q.QuestLevel, &q.MinLevel, &q.MaxLevel,
			&q.Type, &q.ZoneOrSort, &q.RequiredRaces, &q.RequiredClasses,
			&q.SrcItem, &q.RewardXP, &catName,
		)
		if err != nil {
			continue
		}
		if catName != nil {
			q.CategoryName = *catName
		}
		quests = append(quests, q)
	}
	return quests, nil
}
