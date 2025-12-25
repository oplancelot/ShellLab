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
		SELECT ZoneOrSort, COUNT(*) 
		FROM quest_template 
		GROUP BY ZoneOrSort
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
		SELECT entry, IFNULL(Title,''), IFNULL(QuestLevel,0), IFNULL(MinLevel,0), 
			IFNULL(Type,0), IFNULL(ZoneOrSort,0),
			IFNULL(RewXP,0), IFNULL(RewOrReqMoney,0),
			IFNULL(RequiredRaces,0), IFNULL(RequiredClasses,0), IFNULL(SrcItemId,0),
			IFNULL(PrevQuestId,0), IFNULL(NextQuestId,0), IFNULL(ExclusiveGroup,0), IFNULL(NextQuestInChain,0)
		FROM quest_template
		WHERE ZoneOrSort = ?
		ORDER BY QuestLevel, Title
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
		SELECT q.entry, IFNULL(q.Title,''), IFNULL(q.QuestLevel,0), IFNULL(q.MinLevel,0), 
			IFNULL(q.Type,0), IFNULL(q.ZoneOrSort,0),
			IFNULL(q.RewXP,0), IFNULL(q.RewOrReqMoney,0),
			IFNULL(q.RequiredRaces,0), IFNULL(q.RequiredClasses,0), IFNULL(q.SrcItemId,0),
			IFNULL(q.PrevQuestId,0), IFNULL(q.NextQuestId,0), IFNULL(q.ExclusiveGroup,0), IFNULL(q.NextQuestInChain,0),
			c.name
		FROM quest_template q
		LEFT JOIN quest_categories c ON q.ZoneOrSort = c.id
		WHERE q.Title LIKE ?
		ORDER BY length(q.Title), q.Title
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
		SELECT qce.id, qce.group_id, qce.name, 
			COALESCE((SELECT COUNT(*) FROM quest_template WHERE ZoneOrSort = qce.id), 0) as quest_count
		FROM quest_categories_enhanced qce
		WHERE qce.group_id = ?
		ORDER BY quest_count DESC, qce.name
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
	whereClause := "WHERE ZoneOrSort = ?"
	args := []interface{}{categoryID}

	if nameFilter != "" {
		whereClause += " AND title LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	query := fmt.Sprintf(`
		SELECT entry, Title, QuestLevel, MinLevel, Type, ZoneOrSort, RewXP
		FROM quest_template 
		%s
		ORDER BY QuestLevel, Title
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
	err := r.db.QueryRow("SELECT COUNT(*) FROM quest_template").Scan(&count)
	return count, err
}

// GetQuestDetail returns full quest information
func (r *QuestRepository) GetQuestDetail(entry int) (*models.QuestDetail, error) {
	row := r.db.QueryRow(`
		SELECT entry, Title, Details, Objectives, OfferRewardText, EndText,
			QuestLevel, MinLevel, Type, ZoneOrSort,
			RequiredRaces, RequiredClasses,
			RewXP, RewOrReqMoney, RewSpell,
			RewItemId1, RewItemId2, RewItemId3, RewItemId4,
			RewItemCount1, RewItemCount2, RewItemCount3, RewItemCount4,
			RewChoiceItemId1, RewChoiceItemId2, RewChoiceItemId3, RewChoiceItemId4, RewChoiceItemId5, RewChoiceItemId6,
			RewChoiceItemCount1, RewChoiceItemCount2, RewChoiceItemCount3, RewChoiceItemCount4, RewChoiceItemCount5, RewChoiceItemCount6,
			RewRepFaction1, RewRepFaction2, RewRepValue1, RewRepValue2,
			PrevQuestId, NextQuestId, ExclusiveGroup, NextQuestInChain
		FROM quest_template WHERE entry = ?
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
			r.db.QueryRow("SELECT name, COALESCE(icon_path, ''), quality FROM item_template WHERE entry = ?", rewItems[i]).Scan(&name, &icon, &quality)
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
			r.db.QueryRow("SELECT name, COALESCE(icon_path, ''), quality FROM item_template WHERE entry = ?", rewChoiceItems[i]).Scan(&name, &icon, &quality)
			item.Name = name
			item.Icon = icon
			item.Quality = quality
			q.ChoiceItems = append(q.ChoiceItems, item)
		}
	}

	// Process prev quests
	if prevQuestID != 0 {
		var title string
		r.db.QueryRow("SELECT Title FROM quest_template WHERE entry = ?", prevQuestID).Scan(&title)
		q.PrevQuests = append(q.PrevQuests, &models.QuestSeriesItem{Entry: prevQuestID, Title: title})
	}

	// Process next/series quests
	if nextQuestInChain != 0 {
		var title string
		r.db.QueryRow("SELECT Title FROM quest_template WHERE entry = ?", nextQuestInChain).Scan(&title)
		q.Series = append(q.Series, &models.QuestSeriesItem{Entry: nextQuestInChain, Title: title})
	}

	return q, nil
}
