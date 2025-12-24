package repositories

import (
	"database/sql"
	"fmt"

	"shelllab/backend/database/models"
)

// QuestRepository handles quest-related database operations
type QuestRepository struct {
	db *sql.DB
}

// NewQuestRepository creates a new quest repository
func NewQuestRepository(db *sql.DB) *QuestRepository {
	return &QuestRepository{db: db}
}

// GetQuestCategories returns all quest categories (zones and sorts) with quest counts
func (r *QuestRepository) GetQuestCategories() ([]*models.QuestCategory, error) {
	rows, err := r.db.Query(`
		SELECT id, name FROM quest_categories ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make(map[int]*models.QuestCategory)
	var catList []*models.QuestCategory

	for rows.Next() {
		cat := &models.QuestCategory{}
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			continue
		}
		categories[cat.ID] = cat
		catList = append(catList, cat)
	}

	// Now count quests per category
	rows2, err := r.db.Query(`
		SELECT zone_or_sort, COUNT(*) 
		FROM quests 
		GROUP BY zone_or_sort
	`)
	if err != nil {
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
		}
	}

	// Filter out categories with 0 quests
	var activeCats []*models.QuestCategory
	for _, cat := range catList {
		if cat.Count > 0 {
			activeCats = append(activeCats, cat)
		}
	}

	return activeCats, nil
}

// GetQuestsByCategory returns quests filtered by category (zone or sort)
func (r *QuestRepository) GetQuestsByCategory(categoryID int) ([]*models.Quest, error) {
	rows, err := r.db.Query(`
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

	var quests []*models.Quest
	for rows.Next() {
		q := &models.Quest{}
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
func (r *QuestRepository) SearchQuests(query string) ([]*models.Quest, error) {
	rows, err := r.db.Query(`
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

	var quests []*models.Quest
	for rows.Next() {
		q := &models.Quest{}
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
			continue
		}
		if catName != nil {
			q.CategoryName = *catName
		}
		quests = append(quests, q)
	}
	return quests, nil
}

// GetQuestCategoryGroups returns all quest category groups
func (r *QuestRepository) GetQuestCategoryGroups() ([]*models.QuestCategoryGroup, error) {
	rows, err := r.db.Query(`
		SELECT id, name FROM quest_category_groups ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.QuestCategoryGroup
	for rows.Next() {
		g := &models.QuestCategoryGroup{}
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			continue
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetQuestCategoriesByGroup returns all categories in a group with quest counts
func (r *QuestRepository) GetQuestCategoriesByGroup(groupID int) ([]*models.QuestCategoryEnhanced, error) {
	rows, err := r.db.Query(`
		SELECT id, group_id, name, quest_count
		FROM quest_categories_enhanced 
		WHERE group_id = ?
		ORDER BY quest_count DESC, name
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.QuestCategoryEnhanced
	for rows.Next() {
		c := &models.QuestCategoryEnhanced{}
		if err := rows.Scan(&c.ID, &c.GroupID, &c.Name, &c.QuestCount); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetQuestsByEnhancedCategory returns quests for a given category (ZoneOrSort value)
func (r *QuestRepository) GetQuestsByEnhancedCategory(categoryID int, nameFilter string) ([]*models.Quest, error) {
	whereClause := "WHERE zone_or_sort = ?"
	args := []interface{}{categoryID}

	if nameFilter != "" {
		whereClause += " AND title LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	query := fmt.Sprintf(`
		SELECT entry, title, quest_level, min_level, type, zone_or_sort, rew_xp
		FROM quests 
		%s
		ORDER BY quest_level, title
		LIMIT 10000
	`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quests []*models.Quest
	for rows.Next() {
		q := &models.Quest{}
		if err := rows.Scan(&q.Entry, &q.Title, &q.QuestLevel, &q.MinLevel, &q.Type, &q.ZoneOrSort, &q.RewardXP); err != nil {
			continue
		}
		quests = append(quests, q)
	}
	return quests, nil
}

// GetQuestCount returns the total number of quests
func (r *QuestRepository) GetQuestCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM quests").Scan(&count)
	return count, err
}

// GetQuestDetail returns full quest information
func (r *QuestRepository) GetQuestDetail(entry int) (*models.QuestDetail, error) {
	row := r.db.QueryRow(`
		SELECT entry, title, details, objectives, offer_reward_text, end_text,
			quest_level, min_level, type, zone_or_sort,
			required_races, required_classes,
			rew_xp, rew_money, rew_spell,
			rew_item1, rew_item2, rew_item3, rew_item4,
			rew_item_count1, rew_item_count2, rew_item_count3, rew_item_count4,
			rew_choice_item1, rew_choice_item2, rew_choice_item3, rew_choice_item4, rew_choice_item5, rew_choice_item6,
			rew_choice_item_count1, rew_choice_item_count2, rew_choice_item_count3, rew_choice_item_count4, rew_choice_item_count5, rew_choice_item_count6,
			rew_rep_faction1, rew_rep_faction2, rew_rep_value1, rew_rep_value2,
			prev_quest_id, next_quest_id, exclusive_group, next_quest_in_chain
		FROM quests WHERE entry = ?
	`, entry)

	q := &models.QuestDetail{}
	var details, objectives, offerReward, endText *string
	var rewItems [4]int
	var rewItemCounts [4]int
	var rewChoiceItems [6]int
	var rewChoiceItemCounts [6]int
	var repFactions [2]int
	var repValues [2]int
	var prevQuestID, nextQuestID, exclusiveGroup, nextQuestInChain int

	err := row.Scan(
		&q.Entry, &q.Title, &details, &objectives, &offerReward, &endText,
		&q.QuestLevel, &q.MinLevel, &q.Type, &q.ZoneOrSort,
		&q.RequiredRaces, &q.RequiredClasses,
		&q.RewardXP, &q.RewardMoney, &q.RewardSpell,
		&rewItems[0], &rewItems[1], &rewItems[2], &rewItems[3],
		&rewItemCounts[0], &rewItemCounts[1], &rewItemCounts[2], &rewItemCounts[3],
		&rewChoiceItems[0], &rewChoiceItems[1], &rewChoiceItems[2], &rewChoiceItems[3], &rewChoiceItems[4], &rewChoiceItems[5],
		&rewChoiceItemCounts[0], &rewChoiceItemCounts[1], &rewChoiceItemCounts[2], &rewChoiceItemCounts[3], &rewChoiceItemCounts[4], &rewChoiceItemCounts[5],
		&repFactions[0], &repFactions[1], &repValues[0], &repValues[1],
		&prevQuestID, &nextQuestID, &exclusiveGroup, &nextQuestInChain,
	)
	if err != nil {
		return nil, err
	}

	if details != nil {
		q.Details = *details
	}
	if objectives != nil {
		q.Objectives = *objectives
	}
	if offerReward != nil {
		q.OfferRewardText = *offerReward
	}
	if endText != nil {
		q.EndText = *endText
	}

	// Process reward items
	for i := 0; i < 4; i++ {
		if rewItems[i] > 0 {
			item := &models.QuestItem{Entry: rewItems[i], Count: rewItemCounts[i]}
			var name, icon string
			var quality int
			r.db.QueryRow("SELECT name, icon_path, quality FROM items WHERE entry = ?", rewItems[i]).Scan(&name, &icon, &quality)
			item.Name = name
			item.Icon = icon
			item.Quality = quality
			q.RewardItems = append(q.RewardItems, item)
		}
	}

	// Process choice items
	for i := 0; i < 6; i++ {
		if rewChoiceItems[i] > 0 {
			item := &models.QuestItem{Entry: rewChoiceItems[i], Count: rewChoiceItemCounts[i]}
			var name, icon string
			var quality int
			r.db.QueryRow("SELECT name, icon_path, quality FROM items WHERE entry = ?", rewChoiceItems[i]).Scan(&name, &icon, &quality)
			item.Name = name
			item.Icon = icon
			item.Quality = quality
			q.ChoiceItems = append(q.ChoiceItems, item)
		}
	}

	// Process prev quests
	if prevQuestID != 0 {
		var title string
		r.db.QueryRow("SELECT title FROM quests WHERE entry = ?", prevQuestID).Scan(&title)
		q.PrevQuests = append(q.PrevQuests, &models.QuestSeriesItem{Entry: prevQuestID, Title: title})
	}

	// Process next/series quests
	if nextQuestInChain != 0 {
		var title string
		r.db.QueryRow("SELECT title FROM quests WHERE entry = ?", nextQuestInChain).Scan(&title)
		q.Series = append(q.Series, &models.QuestSeriesItem{Entry: nextQuestInChain, Title: title})
	}

	return q, nil
}
