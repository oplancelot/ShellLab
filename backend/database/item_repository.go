package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ItemRepository handles item-related database operations
type ItemRepository struct {
	db *SQLiteDB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *SQLiteDB) *ItemRepository {
	return &ItemRepository{db: db}
}

// ImportFromJSON imports items from item_template.json into SQLite
func (r *ItemRepository) ImportFromJSON(jsonPath string) (int, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var items map[string]map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Begin transaction for bulk insert
	tx, err := r.db.DB().Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO items (
			entry, name, description, quality, item_level, required_level,
			class, subclass, inventory_type, display_id, icon_path,
			buy_price, sell_price, allowable_class, allowable_race, max_stack,
			bonding, max_durability,
			stat_type1, stat_value1, stat_type2, stat_value2,
			stat_type3, stat_value3, stat_type4, stat_value4,
			stat_type5, stat_value5, stat_type6, stat_value6,
			stat_type7, stat_value7, stat_type8, stat_value8,
			stat_type9, stat_value9, stat_type10, stat_value10,
			delay, dmg_min1, dmg_max1, dmg_type1, dmg_min2, dmg_max2, dmg_type2,
			armor, holy_res, fire_res, nature_res, frost_res, shadow_res, arcane_res,
			spell_id1, spell_trigger1, spell_id2, spell_trigger2, spell_id3, spell_trigger3,
			set_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	count := 0
	for _, item := range items {
		_, err := stmt.Exec(
			getInt(item, "entry"),
			cleanName(getString(item, "name")),
			cleanName(getString(item, "description")),
			getInt(item, "quality"),
			getInt(item, "itemLevel"),
			getInt(item, "requiredLevel"),
			getInt(item, "class"),
			getInt(item, "subclass"),
			getInt(item, "inventoryType"),
			getInt(item, "displayId"),
			"", // icon_path - will be populated later
			getInt(item, "buyPrice"),
			getInt(item, "sellPrice"),
			getInt(item, "allowableClass"),
			getInt(item, "allowableRace"),
			getInt(item, "stackable"),
			getInt(item, "bonding"),
			getInt(item, "maxDurability"),
			getInt(item, "statType1"), getInt(item, "statValue1"),
			getInt(item, "statType2"), getInt(item, "statValue2"),
			getInt(item, "statType3"), getInt(item, "statValue3"),
			getInt(item, "statType4"), getInt(item, "statValue4"),
			getInt(item, "statType5"), getInt(item, "statValue5"),
			getInt(item, "statType6"), getInt(item, "statValue6"),
			getInt(item, "statType7"), getInt(item, "statValue7"),
			getInt(item, "statType8"), getInt(item, "statValue8"),
			getInt(item, "statType9"), getInt(item, "statValue9"),
			getInt(item, "statType10"), getInt(item, "statValue10"),
			getInt(item, "delay"),
			getFloat(item, "dmgMin1"), getFloat(item, "dmgMax1"), getInt(item, "dmgType1"),
			getFloat(item, "dmgMin2"), getFloat(item, "dmgMax2"), getInt(item, "dmgType2"),
			getInt(item, "armor"),
			getInt(item, "holyRes"), getInt(item, "fireRes"), getInt(item, "natureRes"),
			getInt(item, "frostRes"), getInt(item, "shadowRes"), getInt(item, "arcaneRes"),
			getInt(item, "spellId1"), getInt(item, "spellTrigger1"),
			getInt(item, "spellId2"), getInt(item, "spellTrigger2"),
			getInt(item, "spellId3"), getInt(item, "spellTrigger3"),
			getInt(item, "setId"),
		)
		if err != nil {
			fmt.Printf("Warning: failed to insert item %v: %v\n", getInt(item, "entry"), err)
			continue
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return count, nil
}

// SearchItems searches for items by name
func (r *ItemRepository) SearchItems(query string, limit int) ([]*Item, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.db.DB().Query(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, icon_path
		FROM items
		WHERE name LIKE ?
		ORDER BY quality DESC, item_level DESC
		LIMIT ?
	`, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
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
func (r *ItemRepository) GetItemByID(id int) (*Item, error) {
	row := r.db.DB().QueryRow(`
		SELECT entry, name, description, quality, item_level, required_level,
			class, subclass, inventory_type, icon_path, sell_price, bonding, max_durability,
			stat_type1, stat_value1, stat_type2, stat_value2,
			stat_type3, stat_value3, stat_type4, stat_value4,
			stat_type5, stat_value5, stat_type6, stat_value6,
			stat_type7, stat_value7, stat_type8, stat_value8,
			stat_type9, stat_value9, stat_type10, stat_value10,
			delay, dmg_min1, dmg_max1, dmg_type1,
			armor, holy_res, fire_res, nature_res, frost_res, shadow_res, arcane_res,
			spell_id1, spell_trigger1, spell_id2, spell_trigger2, spell_id3, spell_trigger3,
			set_id
		FROM items WHERE entry = ?
	`, id)

	item := &Item{}
	err := row.Scan(
		&item.Entry, &item.Name, &item.Description, &item.Quality, &item.ItemLevel, &item.RequiredLevel,
		&item.Class, &item.SubClass, &item.InventoryType, &item.IconPath, &item.SellPrice, &item.Bonding, &item.MaxDurability,
		&item.StatType1, &item.StatValue1, &item.StatType2, &item.StatValue2,
		&item.StatType3, &item.StatValue3, &item.StatType4, &item.StatValue4,
		&item.StatType5, &item.StatValue5, &item.StatType6, &item.StatValue6,
		&item.StatType7, &item.StatValue7, &item.StatType8, &item.StatValue8,
		&item.StatType9, &item.StatValue9, &item.StatType10, &item.StatValue10,
		&item.Delay, &item.DmgMin1, &item.DmgMax1, &item.DmgType1,
		&item.Armor, &item.HolyRes, &item.FireRes, &item.NatureRes, &item.FrostRes, &item.ShadowRes, &item.ArcaneRes,
		&item.SpellID1, &item.SpellTrigger1, &item.SpellID2, &item.SpellTrigger2, &item.SpellID3, &item.SpellTrigger3,
		&item.SetID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetItemCount returns the total number of items
func (r *ItemRepository) GetItemCount() (int, error) {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	return count, err
}

// Item represents a database item record
type Item struct {
	Entry         int    `json:"entry"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Quality       int    `json:"quality"`
	ItemLevel     int    `json:"itemLevel"`
	RequiredLevel int    `json:"requiredLevel"`
	Class         int    `json:"class"`
	SubClass      int    `json:"subClass"`
	InventoryType int    `json:"inventoryType"`
	IconPath      string `json:"iconPath,omitempty"`
	SellPrice     int    `json:"sellPrice,omitempty"`
	Bonding       int    `json:"bonding,omitempty"`
	MaxDurability int    `json:"maxDurability,omitempty"`
	// Stats
	StatType1   int `json:"statType1,omitempty"`
	StatValue1  int `json:"statValue1,omitempty"`
	StatType2   int `json:"statType2,omitempty"`
	StatValue2  int `json:"statValue2,omitempty"`
	StatType3   int `json:"statType3,omitempty"`
	StatValue3  int `json:"statValue3,omitempty"`
	StatType4   int `json:"statType4,omitempty"`
	StatValue4  int `json:"statValue4,omitempty"`
	StatType5   int `json:"statType5,omitempty"`
	StatValue5  int `json:"statValue5,omitempty"`
	StatType6   int `json:"statType6,omitempty"`
	StatValue6  int `json:"statValue6,omitempty"`
	StatType7   int `json:"statType7,omitempty"`
	StatValue7  int `json:"statValue7,omitempty"`
	StatType8   int `json:"statType8,omitempty"`
	StatValue8  int `json:"statValue8,omitempty"`
	StatType9   int `json:"statType9,omitempty"`
	StatValue9  int `json:"statValue9,omitempty"`
	StatType10  int `json:"statType10,omitempty"`
	StatValue10 int `json:"statValue10,omitempty"`
	// Weapon
	Delay    int     `json:"delay,omitempty"`
	DmgMin1  float64 `json:"dmgMin1,omitempty"`
	DmgMax1  float64 `json:"dmgMax1,omitempty"`
	DmgType1 int     `json:"dmgType1,omitempty"`
	// Armor & Resistance
	Armor     int `json:"armor,omitempty"`
	HolyRes   int `json:"holyRes,omitempty"`
	FireRes   int `json:"fireRes,omitempty"`
	NatureRes int `json:"natureRes,omitempty"`
	FrostRes  int `json:"frostRes,omitempty"`
	ShadowRes int `json:"shadowRes,omitempty"`
	ArcaneRes int `json:"arcaneRes,omitempty"`
	// Spells
	SpellID1      int `json:"spellId1,omitempty"`
	SpellTrigger1 int `json:"spellTrigger1,omitempty"`
	SpellID2      int `json:"spellId2,omitempty"`
	SpellTrigger2 int `json:"spellTrigger2,omitempty"`
	SpellID3      int `json:"spellId3,omitempty"`
	SpellTrigger3 int `json:"spellTrigger3,omitempty"`
	// Set
	SetID int `json:"setId,omitempty"`

	// Contextual
	DropRate string `json:"dropRate,omitempty"`
}

// Helper functions for JSON parsing
func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return 0
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		}
	}
	return 0
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func cleanName(name string) string {
	// Remove leading space and quotes from SQL extraction
	name = strings.TrimSpace(name)
	if len(name) > 0 && name[0] == '\'' {
		name = name[1:]
	}
	if len(name) > 0 && name[len(name)-1] == '\'' {
		name = name[:len(name)-1]
	}
	return name
}

// GetTooltipData returns formatted tooltip data for an item
func (r *ItemRepository) GetTooltipData(id int) (*TooltipData, error) {
	item, err := r.GetItemByID(id)
	if err != nil {
		return nil, err
	}

	tooltip := &TooltipData{
		Entry:         item.Entry,
		Name:          cleanItemName(item.Name),
		Quality:       item.Quality,
		QualityName:   getQualityName(item.Quality),
		ItemLevel:     item.ItemLevel,
		RequiredLevel: item.RequiredLevel,
	}

	// Binding
	tooltip.Binding = getBindingText(item.Bonding)

	// Slot & Type
	tooltip.SlotName = getSlotName(item.InventoryType)
	tooltip.TypeName = getTypeName(item.Class, item.SubClass)

	// Armor
	if item.Armor > 0 {
		tooltip.Armor = item.Armor
	}

	// Stats
	tooltip.Stats = parseItemStats(item)

	// Weapon damage
	if item.DmgMax1 > 0 && item.Delay > 0 {
		tooltip.DamageText = fmt.Sprintf("%.0f - %.0f Damage", item.DmgMin1, item.DmgMax1)
		speed := float64(item.Delay) / 1000.0
		tooltip.SpeedText = fmt.Sprintf("Speed %.2f", speed)
		dps := (item.DmgMin1 + item.DmgMax1) / 2.0 / speed
		tooltip.DPS = fmt.Sprintf("(%.1f damage per second)", dps)
	}

	// Resistances
	tooltip.Resistances = parseItemResistances(item)

	// Spell effects (simplified)
	tooltip.SpellEffects = parseItemSpellEffects(item)

	// Description
	if item.Description != "" && item.Description != " ''" {
		tooltip.Description = cleanItemName(item.Description)
	}

	// Sell price
	if item.SellPrice > 0 {
		tooltip.SellPrice = formatMoney(item.SellPrice)
	}

	// Durability
	if item.MaxDurability > 0 {
		tooltip.Durability = fmt.Sprintf("Durability %d / %d", item.MaxDurability, item.MaxDurability)
	}

	return tooltip, nil
}

func parseItemStats(item *Item) []string {
	var stats []string
	statTypes := []struct {
		Type  int
		Value int
	}{
		{item.StatType1, item.StatValue1},
		{item.StatType2, item.StatValue2},
		{item.StatType3, item.StatValue3},
		{item.StatType4, item.StatValue4},
		{item.StatType5, item.StatValue5},
		{item.StatType6, item.StatValue6},
		{item.StatType7, item.StatValue7},
		{item.StatType8, item.StatValue8},
		{item.StatType9, item.StatValue9},
		{item.StatType10, item.StatValue10},
	}

	for _, stat := range statTypes {
		if stat.Type > 0 && stat.Value != 0 {
			statName := getStatName(stat.Type)
			if statName != "" {
				if stat.Value > 0 {
					stats = append(stats, fmt.Sprintf("+%d %s", stat.Value, statName))
				} else {
					stats = append(stats, fmt.Sprintf("%d %s", stat.Value, statName))
				}
			}
		}
	}
	return stats
}

func parseItemResistances(item *Item) []string {
	var res []string
	if item.HolyRes > 0 {
		res = append(res, fmt.Sprintf("+%d Holy Resistance", item.HolyRes))
	}
	if item.FireRes > 0 {
		res = append(res, fmt.Sprintf("+%d Fire Resistance", item.FireRes))
	}
	if item.NatureRes > 0 {
		res = append(res, fmt.Sprintf("+%d Nature Resistance", item.NatureRes))
	}
	if item.FrostRes > 0 {
		res = append(res, fmt.Sprintf("+%d Frost Resistance", item.FrostRes))
	}
	if item.ShadowRes > 0 {
		res = append(res, fmt.Sprintf("+%d Shadow Resistance", item.ShadowRes))
	}
	if item.ArcaneRes > 0 {
		res = append(res, fmt.Sprintf("+%d Arcane Resistance", item.ArcaneRes))
	}
	return res
}

func parseItemSpellEffects(item *Item) []string {
	var effects []string
	// Simplified - would need spell data to format properly
	if item.SpellID1 > 0 {
		effects = append(effects, fmt.Sprintf("Effect from spell %d", item.SpellID1))
	}
	if item.SpellID2 > 0 {
		effects = append(effects, fmt.Sprintf("Effect from spell %d", item.SpellID2))
	}
	if item.SpellID3 > 0 {
		effects = append(effects, fmt.Sprintf("Effect from spell %d", item.SpellID3))
	}
	return effects
}

func getStatName(statType int) string {
	statNames := map[int]string{
		0: "Mana", 1: "Health", 3: "Agility", 4: "Strength",
		5: "Intellect", 6: "Spirit", 7: "Stamina",
		12: "Defense Rating", 13: "Dodge Rating", 14: "Parry Rating",
		15: "Block Rating", 16: "Hit Melee Rating", 17: "Hit Ranged Rating",
		18: "Hit Spell Rating", 19: "Crit Melee Rating", 20: "Crit Ranged Rating",
		21: "Crit Spell Rating", 31: "Hit Rating", 32: "Crit Rating",
		35: "Resilience Rating", 36: "Haste Rating", 38: "Attack Power",
		44: "Armor Penetration Rating", 45: "Spell Power",
	}
	if name, ok := statNames[statType]; ok {
		return name
	}
	return ""
}
