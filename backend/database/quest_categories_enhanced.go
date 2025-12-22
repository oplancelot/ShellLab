package database

import "fmt"

// QuestCategoryGroup represents a top-level quest category group
type QuestCategoryGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// QuestCategoryEnhanced represents an enhanced quest category with group info
type QuestCategoryEnhanced struct {
	ID         int    `json:"id"`
	GroupID    int    `json:"groupId"`
	Name       string `json:"name"`
	QuestCount int    `json:"questCount"`
}

// GetQuestCategoryGroups returns all quest category groups (Zones, Class Quests, etc.)
func (r *ItemRepository) GetQuestCategoryGroups() ([]*QuestCategoryGroup, error) {
	rows, err := r.db.DB().Query(`
		SELECT id, name FROM quest_category_groups ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*QuestCategoryGroup
	for rows.Next() {
		g := &QuestCategoryGroup{}
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			continue
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetQuestCategoriesByGroup returns all categories in a group with quest counts
func (r *ItemRepository) GetQuestCategoriesByGroup(groupID int) ([]*QuestCategoryEnhanced, error) {
	rows, err := r.db.DB().Query(`
		SELECT id, group_id, name, quest_count
		FROM quest_categories_enhanced 
		WHERE group_id = ?
		ORDER BY quest_count DESC, name
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*QuestCategoryEnhanced
	for rows.Next() {
		c := &QuestCategoryEnhanced{}
		if err := rows.Scan(&c.ID, &c.GroupID, &c.Name, &c.QuestCount); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetQuestsByEnhancedCategory returns quests for a given category (ZoneOrSort value)
func (r *ItemRepository) GetQuestsByEnhancedCategory(categoryID int) ([]*Quest, error) {
	rows, err := r.db.DB().Query(`
		SELECT entry, title, quest_level, min_level, type, zone_or_sort, reward_xp
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
		if err := rows.Scan(&q.Entry, &q.Title, &q.QuestLevel, &q.MinLevel, &q.Type, &q.ZoneOrSort, &q.RewardXP); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		quests = append(quests, q)
	}
	return quests, nil
}
