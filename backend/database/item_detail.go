package database

import "fmt"

type ItemDetail struct {
	*Item
	// Extra fields from items table not in Item struct
	DisplayID      int     `json:"displayId"`
	Flags          int     `json:"flags"`
	BuyCount       int     `json:"buyCount"`
	BuyPrice       int     `json:"buyPrice"`
	MaxCount       int     `json:"maxCount"`
	Stackable      int     `json:"stackable"`
	ContainerSlots int     `json:"containerSlots"`
	Material       int     `json:"material"`
	DmgMin2        float64 `json:"dmgMin2"`
	DmgMax2        float64 `json:"dmgMax2"`
	DmgType2       int     `json:"dmgType2"`

	DroppedBy  []*CreatureDrop `json:"droppedBy"`
	RewardFrom []*QuestReward  `json:"rewardFrom"`
}

type CreatureDrop struct {
	Entry    int     `json:"entry"`
	Name     string  `json:"name"`
	LevelMin int     `json:"levelMin"`
	LevelMax int     `json:"levelMax"`
	Chance   float64 `json:"chance"`
}

type QuestReward struct {
	Entry    int    `json:"entry"`
	Title    string `json:"title"`
	Level    int    `json:"level"`
	IsChoice bool   `json:"isChoice"`
}

// GetItemDetail returns full details for an item (sources, etc.)
func (r *ItemRepository) GetItemDetail(entry int) (*ItemDetail, error) {
	item := &Item{}
	detail := &ItemDetail{
		Item:       item,
		DroppedBy:  []*CreatureDrop{},
		RewardFrom: []*QuestReward{},
	}

	// Query 'items' table
	err := r.db.DB().QueryRow(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, icon_path, 
		display_id, description, flags, buy_count, buy_price, sell_price, 
		allowable_class, allowable_race, max_count, stackable, container_slots,
		stat_type1, stat_value1, stat_type2, stat_value2, stat_type3, stat_value3, stat_type4, stat_value4, stat_type5, stat_value5,
		stat_type6, stat_value6, stat_type7, stat_value7, stat_type8, stat_value8, stat_type9, stat_value9, stat_type10, stat_value10,
		delay, dmg_min1, dmg_max1, dmg_type1, dmg_min2, dmg_max2, dmg_type2,
		armor, holy_res, fire_res, nature_res, frost_res, shadow_res, arcane_res,
		spell_id1, spell_trigger1, spell_id2, spell_trigger2, spell_id3, spell_trigger3,
		bonding, max_durability, set_id, material
		FROM items WHERE entry = ?
	`, entry).Scan(
		&item.Entry, &item.Name, &item.Quality, &item.ItemLevel, &item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
		&detail.DisplayID, &item.Description, &detail.Flags, &detail.BuyCount, &detail.BuyPrice, &item.SellPrice,
		&item.AllowableClass, &item.AllowableRace, &detail.MaxCount, &detail.Stackable, &detail.ContainerSlots,
		&item.StatType1, &item.StatValue1, &item.StatType2, &item.StatValue2, &item.StatType3, &item.StatValue3, &item.StatType4, &item.StatValue4, &item.StatType5, &item.StatValue5,
		&item.StatType6, &item.StatValue6, &item.StatType7, &item.StatValue7, &item.StatType8, &item.StatValue8, &item.StatType9, &item.StatValue9, &item.StatType10, &item.StatValue10,
		&item.Delay, &item.DmgMin1, &item.DmgMax1, &item.DmgType1, &detail.DmgMin2, &detail.DmgMax2, &detail.DmgType2,
		&item.Armor, &item.HolyRes, &item.FireRes, &item.NatureRes, &item.FrostRes, &item.ShadowRes, &item.ArcaneRes,
		&item.SpellID1, &item.SpellTrigger1, &item.SpellID2, &item.SpellTrigger2, &item.SpellID3, &item.SpellTrigger3,
		&item.Bonding, &item.MaxDurability, &item.SetID, &detail.Material,
	)
	if err != nil {
		return nil, fmt.Errorf("item not found: %v", err)
	}

	// 2. Dropped By (Creature Loot)
	rows, err := r.db.DB().Query(`
		SELECT c.entry, c.name, c.level_min, c.level_max, cl.chance
		FROM creature_loot cl
		JOIN creatures c ON cl.entry = c.entry
		WHERE cl.item = ?
		ORDER BY cl.chance DESC
		LIMIT 50
	`, entry)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			drop := &CreatureDrop{}
			rows.Scan(&drop.Entry, &drop.Name, &drop.LevelMin, &drop.LevelMax, &drop.Chance)
			detail.DroppedBy = append(detail.DroppedBy, drop)
		}
	}

	// 3. Reward From (Quests)
	rowsQ, err := r.db.DB().Query(`
		SELECT entry, title, quest_level, 0 as is_choice
		FROM quests
		WHERE rew_item1=? OR rew_item2=? OR rew_item3=? OR rew_item4=?
		UNION
		SELECT entry, title, quest_level, 1 as is_choice
		FROM quests
		WHERE rew_choice_item1=? OR rew_choice_item2=? OR rew_choice_item3=? 
		OR rew_choice_item4=? OR rew_choice_item5=? OR rew_choice_item6=?
	`, entry, entry, entry, entry, entry, entry, entry, entry, entry, entry)

	if err == nil {
		defer rowsQ.Close()
		for rowsQ.Next() {
			r := &QuestReward{}
			rowsQ.Scan(&r.Entry, &r.Title, &r.Level, &r.IsChoice)
			detail.RewardFrom = append(detail.RewardFrom, r)
		}
	}

	return detail, nil
}
