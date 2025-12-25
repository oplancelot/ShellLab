// Package repositories contains database access layer implementations
package repositories

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"

	"shelllab/backend/database/helpers"
	"shelllab/backend/database/models"
)

// ItemRepository handles item-related database operations
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// SearchItems searches for items by name
func (r *ItemRepository) SearchItems(query string, limit int) ([]*models.Item, error) {
	rows, err := r.db.Query(`
		SELECT entry, name, quality, item_level, required_level, 
			class, subclass, inventory_type, COALESCE(icon_path, '')
		FROM item_template
		WHERE name LIKE ?
		ORDER BY length(name), name
		LIMIT ?
	`, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.Entry, &item.Name, &item.Quality, &item.ItemLevel,
			&item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

// GetItemByID retrieves a single item by ID
func (r *ItemRepository) GetItemByID(id int) (*models.Item, error) {
	item := &models.Item{}
	err := r.db.QueryRow(`
		SELECT entry, name, COALESCE(description, ''), quality, item_level, required_level,
			class, subclass, inventory_type, COALESCE(icon_path, ''), sell_price,
			allowable_class, allowable_race, bonding, max_durability, max_count, armor,
			stat_type1, stat_value1, stat_type2, stat_value2, stat_type3, stat_value3,
			stat_type4, stat_value4, stat_type5, stat_value5, stat_type6, stat_value6,
			stat_type7, stat_value7, stat_type8, stat_value8, stat_type9, stat_value9,
			stat_type10, stat_value10,
			delay, dmg_min1, dmg_max1, dmg_type1,
			dmg_min2, dmg_max2, dmg_type2,
			holy_res, fire_res, nature_res, frost_res, shadow_res, arcane_res,
			spellid_1, spelltrigger_1, spellid_2, spelltrigger_2, spellid_3, spelltrigger_3,
			set_id
		FROM item_template WHERE entry = ?
	`, id).Scan(
		&item.Entry, &item.Name, &item.Description, &item.Quality, &item.ItemLevel, &item.RequiredLevel,
		&item.Class, &item.SubClass, &item.InventoryType, &item.IconPath, &item.SellPrice,
		&item.AllowableClass, &item.AllowableRace, &item.Bonding, &item.MaxDurability, &item.MaxCount, &item.Armor,
		&item.StatType1, &item.StatValue1, &item.StatType2, &item.StatValue2, &item.StatType3, &item.StatValue3,
		&item.StatType4, &item.StatValue4, &item.StatType5, &item.StatValue5, &item.StatType6, &item.StatValue6,
		&item.StatType7, &item.StatValue7, &item.StatType8, &item.StatValue8, &item.StatType9, &item.StatValue9,
		&item.StatType10, &item.StatValue10,
		&item.Delay, &item.DmgMin1, &item.DmgMax1, &item.DmgType1,
		&item.DmgMin2, &item.DmgMax2, &item.DmgType2,
		&item.HolyRes, &item.FireRes, &item.NatureRes, &item.FrostRes, &item.ShadowRes, &item.ArcaneRes,
		&item.SpellID1, &item.SpellTrigger1, &item.SpellID2, &item.SpellTrigger2, &item.SpellID3, &item.SpellTrigger3,
		&item.SetID,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetItemCount returns the total number of items
func (r *ItemRepository) GetItemCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM item_template").Scan(&count)
	return count, err
}

// GetItemClasses returns all item classes with their subclasses and inventory slots
func (r *ItemRepository) GetItemClasses() ([]*models.ItemClass, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT class, subclass, inventory_type
		FROM item_template
		WHERE class IN (0,1,2,4,6,7,9,11,12,13,15)
		ORDER BY class, subclass, inventory_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classMap := make(map[int]*models.ItemClass)
	subclassMap := make(map[string]*models.ItemSubClass)

	for rows.Next() {
		var class, subclass, invType int
		if err := rows.Scan(&class, &subclass, &invType); err != nil {
			continue
		}

		// Ensure class exists
		if _, exists := classMap[class]; !exists {
			classMap[class] = &models.ItemClass{
				Class:      class,
				Name:       helpers.GetClassName(class),
				SubClasses: []*models.ItemSubClass{},
			}
		}

		// Ensure subclass exists
		subKey := fmt.Sprintf("%d-%d", class, subclass)
		if _, exists := subclassMap[subKey]; !exists {
			sc := &models.ItemSubClass{
				Class:          class,
				SubClass:       subclass,
				Name:           helpers.GetSubClassName(class, subclass),
				InventorySlots: []*models.InventorySlot{},
			}
			subclassMap[subKey] = sc
			classMap[class].SubClasses = append(classMap[class].SubClasses, sc)
		}

		// Add inventory slot if applicable (mainly for armor/weapons)
		if (class == 2 || class == 4) && invType > 0 {
			slot := &models.InventorySlot{
				Class:         class,
				SubClass:      subclass,
				InventoryType: invType,
				Name:          helpers.GetInventoryTypeName(invType),
			}
			subclassMap[subKey].InventorySlots = append(subclassMap[subKey].InventorySlots, slot)
		}
	}

	// Convert map to slice and sort
	var classes []*models.ItemClass
	for _, c := range classMap {
		classes = append(classes, c)
	}
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].Class < classes[j].Class
	})

	return classes, nil
}

// GetItemsByClass returns items filtered by class and subclass
func (r *ItemRepository) GetItemsByClass(class, subClass int, nameFilter string, limit, offset int) ([]*models.Item, int, error) {
	whereClause := "WHERE class = ? AND subclass = ?"
	args := []interface{}{class, subClass}

	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	// Count
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM item_template %s", whereClause)
	err := r.db.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Data
	dataArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, COALESCE(icon_path, '')
		FROM item_template %s
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.Entry, &item.Name, &item.Quality, &item.ItemLevel,
			&item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, count, nil
}

// GetItemsByClassAndSlot returns items filtered by class, subclass, and inventory type
func (r *ItemRepository) GetItemsByClassAndSlot(class, subClass, inventoryType int, nameFilter string, limit, offset int) ([]*models.Item, int, error) {
	whereClause := "WHERE class = ? AND subclass = ? AND inventory_type = ?"
	args := []interface{}{class, subClass, inventoryType}

	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	// Count
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM item_template %s", whereClause)
	err := r.db.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Data
	dataArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, COALESCE(icon_path, '')
		FROM item_template %s
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.Entry, &item.Name, &item.Quality, &item.ItemLevel,
			&item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, count, nil
}

// AdvancedSearch performs a multi-dimensional search on items
func (r *ItemRepository) AdvancedSearch(filter models.SearchFilter) (*models.SearchResult, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}

	var conditions []string
	var args []interface{}

	// Name filter
	if filter.Query != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+filter.Query+"%")
	}

	// Quality filter
	if len(filter.Quality) > 0 {
		placeholders := make([]string, len(filter.Quality))
		for i, q := range filter.Quality {
			placeholders[i] = "?"
			args = append(args, q)
		}
		conditions = append(conditions, fmt.Sprintf("quality IN (%s)", strings.Join(placeholders, ",")))
	}

	// Class filter
	if len(filter.Class) > 0 {
		placeholders := make([]string, len(filter.Class))
		for i, c := range filter.Class {
			placeholders[i] = "?"
			args = append(args, c)
		}
		conditions = append(conditions, fmt.Sprintf("class IN (%s)", strings.Join(placeholders, ",")))
	}

	// SubClass filter
	if len(filter.SubClass) > 0 {
		placeholders := make([]string, len(filter.SubClass))
		for i, sc := range filter.SubClass {
			placeholders[i] = "?"
			args = append(args, sc)
		}
		conditions = append(conditions, fmt.Sprintf("subclass IN (%s)", strings.Join(placeholders, ",")))
	}

	// InventoryType filter
	if len(filter.InventoryType) > 0 {
		placeholders := make([]string, len(filter.InventoryType))
		for i, it := range filter.InventoryType {
			placeholders[i] = "?"
			args = append(args, it)
		}
		conditions = append(conditions, fmt.Sprintf("inventory_type IN (%s)", strings.Join(placeholders, ",")))
	}

	// Level Range
	if filter.MinLevel > 0 {
		conditions = append(conditions, "item_level >= ?")
		args = append(args, filter.MinLevel)
	}
	if filter.MaxLevel > 0 {
		conditions = append(conditions, "item_level <= ?")
		args = append(args, filter.MaxLevel)
	}

	// Required Level Range
	if filter.MinReqLevel > 0 {
		conditions = append(conditions, "required_level >= ?")
		args = append(args, filter.MinReqLevel)
	}
	if filter.MaxReqLevel > 0 {
		conditions = append(conditions, "required_level <= ?")
		args = append(args, filter.MaxReqLevel)
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) FROM item_template " + whereClause
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("search count error: %w", err)
	}

	// Data query
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, COALESCE(icon_path, '')
		FROM item_template
		%s
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add limit/offset args
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("search data error: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.Entry, &item.Name, &item.Quality, &item.ItemLevel,
			&item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return &models.SearchResult{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

// GetItemSets returns all item sets for browsing
func (r *ItemRepository) GetItemSets() ([]*models.ItemSetBrowse, error) {
	rows, err := r.db.Query(`
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

	var sets []*models.ItemSetBrowse
	for rows.Next() {
		set := &models.ItemSetBrowse{}
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
func (r *ItemRepository) GetItemSetDetail(itemSetID int) (*models.ItemSetDetail, error) {
	row := r.db.QueryRow(`
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

	detail := &models.ItemSetDetail{
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
	var setBonuses []models.SetBonus
	for i := 0; i < 8; i++ {
		if spells[i] > 0 && bonuses[i] > 0 {
			setBonuses = append(setBonuses, models.SetBonus{
				Threshold: bonuses[i],
				SpellID:   spells[i],
			})
		}
	}

	// Sort bonuses by threshold (asc)
	sort.Slice(setBonuses, func(i, j int) bool {
		return setBonuses[i].Threshold < setBonuses[j].Threshold
	})

	// Resolve descriptions
	for i := range setBonuses {
		setBonuses[i].Description = r.resolveSpellText(setBonuses[i].SpellID)
	}

	detail.Bonuses = setBonuses

	return detail, nil
}

// GetTooltipData generates tooltip information for an item
func (r *ItemRepository) GetTooltipData(itemID int) (*models.TooltipData, error) {
	item, err := r.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}

	tooltip := &models.TooltipData{
		Entry:         item.Entry,
		Name:          item.Name,
		Quality:       item.Quality,
		ItemLevel:     item.ItemLevel,
		RequiredLevel: item.RequiredLevel,
		SellPrice:     item.SellPrice,
		Description:   item.Description,
	}

	// Unique
	if item.MaxCount == 1 {
		tooltip.Unique = true
	}

	// Binding
	tooltip.Binding = helpers.GetBondingName(item.Bonding)

	// Item Type and Slot
	tooltip.ItemType = strings.ReplaceAll(helpers.GetSubClassName(item.Class, item.SubClass), " (One-Handed)", "")
	tooltip.Slot = helpers.GetInventoryTypeName(item.InventoryType)

	// Armor
	if item.Armor > 0 {
		tooltip.Armor = item.Armor
	}

	// Weapon damage
	if item.DmgMin1 > 0 || item.DmgMax1 > 0 {
		tooltip.DamageRange = fmt.Sprintf("%.0f - %.0f Damage", item.DmgMin1, item.DmgMax1)
		if item.Delay > 0 {
			speed := float64(item.Delay) / 1000.0
			tooltip.AttackSpeed = fmt.Sprintf("Speed %.2f", speed)
			dps := (item.DmgMin1 + item.DmgMax1) / 2.0 / speed
			// Round half up for DPS
			dpsRounded := math.Round(dps*10) / 10
			tooltip.DPS = fmt.Sprintf("(%.1f damage per second)", dpsRounded)
		}
	}

	// Bonus Damage (e.g. Shadow Damage)
	if item.DmgMin2 > 0 || item.DmgMax2 > 0 {
		typeName := helpers.GetSchoolName(item.DmgType2)
		tooltip.Stats = append(tooltip.Stats, fmt.Sprintf("+%.0f - %.0f %s Damage", item.DmgMin2, item.DmgMax2, typeName))
	}

	// Stats
	statPairs := []struct{ t, v int }{
		{item.StatType1, item.StatValue1}, {item.StatType2, item.StatValue2},
		{item.StatType3, item.StatValue3}, {item.StatType4, item.StatValue4},
		{item.StatType5, item.StatValue5}, {item.StatType6, item.StatValue6},
		{item.StatType7, item.StatValue7}, {item.StatType8, item.StatValue8},
		{item.StatType9, item.StatValue9}, {item.StatType10, item.StatValue10},
	}
	for _, sp := range statPairs {
		if sp.t > 0 && sp.v != 0 {
			tooltip.Stats = append(tooltip.Stats, r.formatStat(sp.t, sp.v))
		}
	}

	// Resistances
	if item.HolyRes > 0 {
		tooltip.Resistances = append(tooltip.Resistances, fmt.Sprintf("+%d Holy Resistance", item.HolyRes))
	}
	if item.FireRes > 0 {
		tooltip.Resistances = append(tooltip.Resistances, fmt.Sprintf("+%d Fire Resistance", item.FireRes))
	}
	if item.NatureRes > 0 {
		tooltip.Resistances = append(tooltip.Resistances, fmt.Sprintf("+%d Nature Resistance", item.NatureRes))
	}
	if item.FrostRes > 0 {
		tooltip.Resistances = append(tooltip.Resistances, fmt.Sprintf("+%d Frost Resistance", item.FrostRes))
	}
	if item.ShadowRes > 0 {
		tooltip.Resistances = append(tooltip.Resistances, fmt.Sprintf("+%d Shadow Resistance", item.ShadowRes))
	}
	if item.ArcaneRes > 0 {
		tooltip.Resistances = append(tooltip.Resistances, fmt.Sprintf("+%d Arcane Resistance", item.ArcaneRes))
	}

	// Durability
	if item.MaxDurability > 0 {
		tooltip.Durability = fmt.Sprintf("Durability %d / %d", item.MaxDurability, item.MaxDurability)
	}

	// Spell Effects
	spellPairs := []struct{ id, trigger int }{
		{item.SpellID1, item.SpellTrigger1},
		{item.SpellID2, item.SpellTrigger2},
		{item.SpellID3, item.SpellTrigger3},
	}
	for _, sp := range spellPairs {
		if sp.id > 0 {
			effect := r.formatSpellEffect(sp.id, sp.trigger)
			if effect != "" {
				tooltip.Effects = append(tooltip.Effects, effect)
			}
		}
	}

	// Set Info
	if item.SetID > 0 {
		var setInfo models.ItemSetInfo

		var setID, skillID, skillLevel int
		var item1, item2, item3, item4, item5, item6, item7, item8, item9, item10 int
		var spell1, spell2, spell3, spell4, spell5, spell6, spell7, spell8 int
		var bonus1, bonus2, bonus3, bonus4, bonus5, bonus6, bonus7, bonus8 int

		err := r.db.QueryRow(`
			SELECT itemset_id, COALESCE(name, ''),
				item1, item2, item3, item4, item5, item6, item7, item8, item9, item10,
				spell1, spell2, spell3, spell4, spell5, spell6, spell7, spell8,
				bonus1, bonus2, bonus3, bonus4, bonus5, bonus6, bonus7, bonus8,
				skill_id, skill_level
			FROM itemsets WHERE itemset_id = ?
		`, item.SetID).Scan(
			&setID, &setInfo.Name,
			&item1, &item2, &item3, &item4, &item5, &item6, &item7, &item8, &item9, &item10,
			&spell1, &spell2, &spell3, &spell4, &spell5, &spell6, &spell7, &spell8,
			&bonus1, &bonus2, &bonus3, &bonus4, &bonus5, &bonus6, &bonus7, &bonus8,
			&skillID, &skillLevel,
		)

		if err == nil {
			// Process items
			itemIDs := []int{item1, item2, item3, item4, item5, item6, item7, item8, item9, item10}
			for _, id := range itemIDs {
				if id > 0 {
					var itemName string
					r.db.QueryRow("SELECT name FROM item_template WHERE entry = ?", id).Scan(&itemName)
					setInfo.Items = append(setInfo.Items, itemName)
				}
			}

			// Process bonuses
			bonuses := []struct{ spell, threshold int }{
				{spell1, bonus1}, {spell2, bonus2}, {spell3, bonus3}, {spell4, bonus4},
				{spell5, bonus5}, {spell6, bonus6}, {spell7, bonus7}, {spell8, bonus8},
			}
			// Sort bonuses by threshold (asc)
			sort.Slice(bonuses, func(i, j int) bool {
				return bonuses[i].threshold < bonuses[j].threshold
			})

			for _, b := range bonuses {
				if b.spell > 0 && b.threshold > 0 {
					description := r.resolveSpellText(b.spell)
					if description != "" {
						setInfo.Bonuses = append(setInfo.Bonuses, fmt.Sprintf("(%d) Set: %s", b.threshold, description))
					}
				}
			}

			tooltip.SetInfo = &setInfo
		}
	}

	return tooltip, nil
}

// resolveSpellText fetches and formats spell description with parameters
func (r *ItemRepository) resolveSpellText(spellID int) string {
	// Get spell data including effect values and duration index
	var name, description string
	var bp1, bp2, bp3, ds1, ds2, ds3 int
	var durationIndex int

	// Try to query with durationIndex, fallback if column doesn't exist yet
	err := r.db.QueryRow(`
		SELECT COALESCE(name, ''), COALESCE(description, ''),
			effectBasePoints1, effectBasePoints2, effectBasePoints3,
			effectDieSides1, effectDieSides2, effectDieSides3,
			durationIndex
		FROM spell_template WHERE entry = ?
	`, spellID).Scan(&name, &description, &bp1, &bp2, &bp3, &ds1, &ds2, &ds3, &durationIndex)

	if err != nil {
		// Fallback for old schema (without durationIndex)
		err = r.db.QueryRow(`
			SELECT COALESCE(name, ''), COALESCE(description, ''),
				effectBasePoints1, effectBasePoints2, effectBasePoints3,
				effectDieSides1, effectDieSides2, effectDieSides3
			FROM spell_template WHERE entry = ?
		`, spellID).Scan(&name, &description, &bp1, &bp2, &bp3, &ds1, &ds2, &ds3)

		if err != nil {
			return ""
		}
		durationIndex = 0
	}

	// Use description if available, otherwise use name
	text := description
	if text == "" {
		text = name
	}
	if text == "" {
		return ""
	}

	// Calculate actual values (base_points + 1 for typical spells, or base_points + die_sides for ranges)
	v1 := bp1 + 1
	v2 := bp2 + 1
	v3 := bp3 + 1
	if ds1 > 1 {
		v1 = bp1 + ds1
	}
	if ds2 > 1 {
		v2 = bp2 + ds2
	}
	if ds3 > 1 {
		v3 = bp3 + ds3
	}

	// Get Duration
	var durationText string
	if durationIndex > 0 {
		var durationBase int
		r.db.QueryRow("SELECT duration_base FROM spell_durations WHERE id = ?", durationIndex).Scan(&durationBase)
		if durationBase > 0 {
			if durationBase < 0 {
				durationBase = -durationBase
			}
			// Duration is usually in ms
			seconds := durationBase / 1000
			if seconds < 60 {
				durationText = fmt.Sprintf("%d sec", seconds)
			} else if seconds < 3600 {
				durationText = fmt.Sprintf("%d min", seconds/60)
			} else {
				durationText = fmt.Sprintf("%d hr", seconds/3600)
			}
		}
	}
	if durationText == "" {
		durationText = "duration" // fallback
	}

	// Replace placeholders
	text = strings.ReplaceAll(text, "$d", durationText)
	text = strings.ReplaceAll(text, "$s1", fmt.Sprintf("%d", v1))
	text = strings.ReplaceAll(text, "$s2", fmt.Sprintf("%d", v2))
	text = strings.ReplaceAll(text, "$s3", fmt.Sprintf("%d", v3))
	// Over-time effects (damage/healing over duration)
	text = strings.ReplaceAll(text, "$o1", fmt.Sprintf("%d", v1))
	text = strings.ReplaceAll(text, "$o2", fmt.Sprintf("%d", v2))
	text = strings.ReplaceAll(text, "$o3", fmt.Sprintf("%d", v3))
	// Also handle ${} format
	text = strings.ReplaceAll(text, "${s1}", fmt.Sprintf("%d", v1))
	text = strings.ReplaceAll(text, "${s2}", fmt.Sprintf("%d", v2))
	text = strings.ReplaceAll(text, "${s3}", fmt.Sprintf("%d", v3))
	text = strings.ReplaceAll(text, "${o1}", fmt.Sprintf("%d", v1))
	text = strings.ReplaceAll(text, "${o2}", fmt.Sprintf("%d", v2))
	text = strings.ReplaceAll(text, "${o3}", fmt.Sprintf("%d", v3))

	return text
}

// formatSpellEffect returns a formatted spell effect string with trigger prefix
func (r *ItemRepository) formatSpellEffect(spellID, trigger int) string {
	text := r.resolveSpellText(spellID)
	if text == "" {
		return ""
	}

	// Format based on trigger type
	var prefix string
	switch trigger {
	case 0: // Use
		prefix = "Use:"
	case 1: // On Equip
		prefix = "Equip:"
	case 2: // Chance on Hit
		prefix = "Chance on hit:"
	case 4: // Soulstone
		prefix = "Use:"
	case 5: // Use with no delay
		prefix = "Use:"
	case 6: // Learn spell
		prefix = "Use:"
	default:
		prefix = "Equip:"
	}

	return fmt.Sprintf("%s %s", prefix, text)
}

// GetItemDetail returns full item information with drop sources
func (r *ItemRepository) GetItemDetail(entry int) (*models.ItemDetail, error) {
	item, err := r.GetItemByID(entry)
	if err != nil {
		return nil, err
	}

	detail := &models.ItemDetail{Item: item}

	// Get dropped by creatures
	rows, err := r.db.Query(`
		SELECT c.entry, c.name, c.level_min, c.level_max, cl.chance
		FROM creature_loot_template cl
		JOIN creature_template c ON cl.entry = c.loot_id
		WHERE cl.item = ?
		ORDER BY cl.chance DESC
		LIMIT 20
	`, entry)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			drop := &models.CreatureDrop{}
			rows.Scan(&drop.Entry, &drop.Name, &drop.LevelMin, &drop.LevelMax, &drop.Chance)
			detail.DroppedBy = append(detail.DroppedBy, drop)
		}
	}

	// Get quest rewards
	rows2, err := r.db.Query(`
		SELECT entry, Title, QuestLevel, 0 as is_choice
		FROM quest_template
		WHERE RewItemId1 = ? OR RewItemId2 = ? OR RewItemId3 = ? OR RewItemId4 = ?
		UNION
		SELECT entry, Title, QuestLevel, 1 as is_choice
		FROM quest_template
		WHERE RewChoiceItemId1 = ? OR RewChoiceItemId2 = ? OR RewChoiceItemId3 = ? 
		   OR RewChoiceItemId4 = ? OR RewChoiceItemId5 = ? OR RewChoiceItemId6 = ?
		LIMIT 20
	`, entry, entry, entry, entry, entry, entry, entry, entry, entry, entry)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			reward := &models.QuestReward{}
			var isChoice int
			rows2.Scan(&reward.Entry, &reward.Title, &reward.Level, &isChoice)
			reward.IsChoice = isChoice == 1
			detail.RewardFrom = append(detail.RewardFrom, reward)
		}
	}

	// Get contains (if item is a container)
	rows3, err := r.db.Query(`
		SELECT i.entry, i.name, i.quality, COALESCE(i.icon_path, ''), il.chance, il.mincount_or_ref, il.maxcount
		FROM item_loot_template il
		JOIN item_template i ON il.item = i.entry
		WHERE il.entry = ?
		ORDER BY il.chance DESC
	`, entry)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			drop := &models.ItemDrop{}
			rows3.Scan(&drop.Entry, &drop.Name, &drop.Quality, &drop.IconPath, &drop.Chance, &drop.MinCount, &drop.MaxCount)
			detail.Contains = append(detail.Contains, drop)
		}
	}

	return detail, nil
}

// formatStat returns a formatted stat string
func (r *ItemRepository) formatStat(statType, value int) string {
	statNames := map[int]string{
		0: "Mana", 1: "Health", 3: "Agility", 4: "Strength",
		5: "Intellect", 6: "Spirit", 7: "Stamina",
		12: "Defense Rating", 13: "Dodge Rating", 14: "Parry Rating",
		15: "Shield Block Rating", 16: "Melee Hit Rating", 17: "Ranged Hit Rating",
		18: "Spell Hit Rating", 19: "Melee Critical Rating", 20: "Ranged Critical Rating",
		21: "Spell Critical Rating", 35: "Resilience Rating", 36: "Haste Rating",
		37: "Expertise Rating", 38: "Attack Power", 39: "Ranged Attack Power",
		41: "Spell Healing", 42: "Spell Damage", 43: "Mana Regeneration",
		44: "Armor Penetration Rating", 45: "Spell Power",
	}
	name := statNames[statType]
	if name == "" {
		name = fmt.Sprintf("Unknown Stat %d", statType)
	}
	if value > 0 {
		return fmt.Sprintf("+%d %s", value, name)
	}
	return fmt.Sprintf("%d %s", value, name)
}
