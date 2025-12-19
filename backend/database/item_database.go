package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// ItemTemplate represents a complete item from item_template.json
type ItemTemplate struct {
	Entry          int    `json:"entry"`
	Class          int    `json:"class"`
	SubClass       int    `json:"subclass"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DisplayID      int    `json:"displayId"`
	Quality        int    `json:"quality"`
	Flags          int    `json:"flags"`
	BuyCount       int    `json:"buyCount"`
	BuyPrice       int    `json:"buyPrice"`
	SellPrice      int    `json:"sellPrice"`
	InventoryType  int    `json:"inventoryType"`
	AllowableClass int    `json:"allowableClass"`
	AllowableRace  int    `json:"allowableRace"`
	ItemLevel      int    `json:"itemLevel"`
	RequiredLevel  int    `json:"requiredLevel"`
	RequiredSkill  int    `json:"requiredSkill"`
	MaxCount       int    `json:"maxCount"`
	Stackable      int    `json:"stackable"`
	ContainerSlots int    `json:"containerSlots"`

	// Stats
	StatType1   int `json:"statType1"`
	StatValue1  int `json:"statValue1"`
	StatType2   int `json:"statType2"`
	StatValue2  int `json:"statValue2"`
	StatType3   int `json:"statType3"`
	StatValue3  int `json:"statValue3"`
	StatType4   int `json:"statType4"`
	StatValue4  int `json:"statValue4"`
	StatType5   int `json:"statType5"`
	StatValue5  int `json:"statValue5"`
	StatType6   int `json:"statType6"`
	StatValue6  int `json:"statValue6"`
	StatType7   int `json:"statType7"`
	StatValue7  int `json:"statValue7"`
	StatType8   int `json:"statType8"`
	StatValue8  int `json:"statValue8"`
	StatType9   int `json:"statType9"`
	StatValue9  int `json:"statValue9"`
	StatType10  int `json:"statType10"`
	StatValue10 int `json:"statValue10"`

	// Weapon
	Delay    int     `json:"delay"`
	DmgMin1  float64 `json:"dmgMin1"`
	DmgMax1  float64 `json:"dmgMax1"`
	DmgType1 int     `json:"dmgType1"`
	DmgMin2  float64 `json:"dmgMin2"`
	DmgMax2  float64 `json:"dmgMax2"`
	DmgType2 int     `json:"dmgType2"`

	// Armor & Resistance
	Armor     int `json:"armor"`
	HolyRes   int `json:"holyRes"`
	FireRes   int `json:"fireRes"`
	NatureRes int `json:"natureRes"`
	FrostRes  int `json:"frostRes"`
	ShadowRes int `json:"shadowRes"`
	ArcaneRes int `json:"arcaneRes"`

	// Spells
	SpellID1      int `json:"spellId1"`
	SpellTrigger1 int `json:"spellTrigger1"`
	SpellID2      int `json:"spellId2"`
	SpellTrigger2 int `json:"spellTrigger2"`
	SpellID3      int `json:"spellId3"`
	SpellTrigger3 int `json:"spellTrigger3"`

	// Other
	Bonding       int `json:"bonding"`
	MaxDurability int `json:"maxDurability"`
	SetID         int `json:"setId"`
	Material      int `json:"material"`
}

// SetBonus represents a set bonus
type SetBonus struct {
	Threshold int `json:"threshold"`
	SpellID   int `json:"spellId"`
}

// ItemSet represents an item set
type ItemSet struct {
	ID      int        `json:"id"`
	Name    string     `json:"name"`
	ItemIDs []int      `json:"itemIds"`
	Bonuses []SetBonus `json:"bonuses"`
}

// ItemSetInfo represents set info for tooltip
type ItemSetInfo struct {
	Name    string   `json:"name"`
	Items   []string `json:"items"`
	Bonuses []string `json:"bonuses"`
}

// TooltipData represents data for rendering a tooltip
type TooltipData struct {
	Entry         int          `json:"entry"`
	Name          string       `json:"name"`
	Quality       int          `json:"quality"`
	QualityName   string       `json:"qualityName"`
	ItemLevel     int          `json:"itemLevel"`
	RequiredLevel int          `json:"requiredLevel"`
	Binding       string       `json:"binding,omitempty"`
	SlotName      string       `json:"slotName,omitempty"`
	TypeName      string       `json:"typeName,omitempty"`
	Armor         int          `json:"armor,omitempty"`
	Stats         []string     `json:"stats,omitempty"`
	DamageText    string       `json:"damageText,omitempty"`
	SpeedText     string       `json:"speedText,omitempty"`
	DPS           string       `json:"dps,omitempty"`
	Resistances   []string     `json:"resistances,omitempty"`
	SpellEffects  []string     `json:"spellEffects,omitempty"`
	Description   string       `json:"description,omitempty"`
	SellPrice     string       `json:"sellPrice,omitempty"`
	SetName       string       `json:"setName,omitempty"`
	Durability    string       `json:"durability,omitempty"`
	Classes       string       `json:"classes,omitempty"`
	Races         string       `json:"races,omitempty"`
	SetInfo       *ItemSetInfo `json:"setInfo,omitempty"`
}

// ItemDef represents a single item's metadata (simplified for AtlasLoot)
type ItemDef struct {
	Entry          int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Quality        int    `json:"quality"`
	InventoryType  int    `json:"inventoryType"`
	AllowableClass int    `json:"-"`
	AllowableRace  int    `json:"-"`
	ItemLevel      int    `json:"itemLevel"`
	RequiredLevel  int    `json:"requiredLevel"`
	DisplayID      int    `json:"-"`
	IconPath       string `json:"icon"`
	SellPrice      int    `json:"sellPrice"`
	ItemLink       string `json:"itemLink"`
	Class          string `json:"class"`
	SubClass       string `json:"subClass"`
	EquipSlot      string `json:"equipSlot"`
	MaxStack       int    `json:"maxStack"`
}

// ItemDatabase is an in-memory item repository
type ItemDatabase struct {
	items         map[int]*ItemDef
	itemTemplates map[int]*ItemTemplate
	mu            sync.RWMutex
}

// NewItemDatabase creates a new item database
func NewItemDatabase() *ItemDatabase {
	db := &ItemDatabase{
		items:         make(map[int]*ItemDef),
		itemTemplates: make(map[int]*ItemTemplate),
	}
	db.loadMockData()
	return db
}

// GetItem retrieves an item by ID
func (db *ItemDatabase) GetItem(id int) (*ItemDef, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	item, found := db.items[id]
	return item, found
}

// GetItemTemplate retrieves full item template by ID
func (db *ItemDatabase) GetItemTemplate(id int) (*ItemTemplate, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	item, found := db.itemTemplates[id]
	return item, found
}

// GetTooltipData returns formatted tooltip data for an item
func (db *ItemDatabase) GetTooltipData(id int) *TooltipData {
	db.mu.RLock()
	defer db.mu.RUnlock()

	item, found := db.itemTemplates[id]
	if !found {
		return nil
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
	tooltip.Stats = parseStats(item)

	// Weapon damage
	if item.DmgMax1 > 0 && item.Delay > 0 {
		tooltip.DamageText = fmt.Sprintf("%.0f - %.0f Damage", item.DmgMin1, item.DmgMax1)
		speed := float64(item.Delay) / 1000.0
		tooltip.SpeedText = fmt.Sprintf("Speed %.2f", speed)
		dps := (item.DmgMin1 + item.DmgMax1) / 2.0 / speed
		tooltip.DPS = fmt.Sprintf("(%.1f damage per second)", dps)
	}

	// Resistances
	tooltip.Resistances = parseResistances(item)

	// Spell effects (simplified)
	tooltip.SpellEffects = parseSpellEffects(item)

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

	return tooltip
}

// AddItem adds or updates an item
func (db *ItemDatabase) AddItem(item *ItemDef) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[item.Entry] = item
}

// Count returns the number of items in the database
func (db *ItemDatabase) Count() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.items)
}

// TemplateCount returns the number of item templates
func (db *ItemDatabase) TemplateCount() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.itemTemplates)
}

// loadMockData loads some test items for POC phase
func (db *ItemDatabase) loadMockData() {
	mockItems := []*ItemDef{
		{
			Entry:         19019,
			Name:          "Thunderfury, Blessed Blade of the Windseeker",
			Quality:       5,
			ItemLevel:     80,
			RequiredLevel: 60,
			IconPath:      "inv_sword_39",
		},
		{
			Entry:         17076,
			Name:          "Bonereaver's Edge",
			Quality:       4,
			ItemLevel:     71,
			RequiredLevel: 60,
			IconPath:      "inv_axe_12",
		},
	}

	for _, item := range mockItems {
		db.items[item.Entry] = item
	}
}

// LoadFromJSON loads items from a simple JSON file (items.json)
func (db *ItemDatabase) LoadFromJSON(filepath string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var items map[string]*struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Quality int    `json:"quality"`
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, item := range items {
		db.items[item.ID] = &ItemDef{
			Entry:   item.ID,
			Name:    item.Name,
			Quality: item.Quality,
		}
	}

	return nil
}

// LoadItemTemplates loads full item data from item_template.json
func (db *ItemDatabase) LoadItemTemplates(filepath string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var items map[string]*ItemTemplate
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, item := range items {
		db.itemTemplates[item.Entry] = item

		// Also populate the simple items map
		db.items[item.Entry] = &ItemDef{
			Entry:         item.Entry,
			Name:          cleanItemName(item.Name),
			Description:   item.Description,
			Quality:       item.Quality,
			InventoryType: item.InventoryType,
			ItemLevel:     item.ItemLevel,
			RequiredLevel: item.RequiredLevel,
			SellPrice:     item.SellPrice,
			MaxStack:      item.Stackable,
		}
	}

	return nil
}

// === Helper Functions ===

func cleanItemName(name string) string {
	// Remove leading space and quotes from SQL extraction
	if len(name) > 2 && name[0] == ' ' && name[1] == '\'' {
		name = name[2:]
	}
	if len(name) > 0 && name[0] == '\'' {
		name = name[1:]
	}
	if len(name) > 0 && name[len(name)-1] == '\'' {
		name = name[:len(name)-1]
	}
	return name
}

func getQualityName(quality int) string {
	names := []string{"Poor", "Common", "Uncommon", "Rare", "Epic", "Legendary", "Artifact"}
	if quality >= 0 && quality < len(names) {
		return names[quality]
	}
	return "Unknown"
}

func getBindingText(bonding int) string {
	switch bonding {
	case 1:
		return "Binds when picked up"
	case 2:
		return "Binds when equipped"
	case 3:
		return "Binds when used"
	case 4:
		return "Quest Item"
	default:
		return ""
	}
}

func getSlotName(inventoryType int) string {
	slots := map[int]string{
		1:  "Head",
		2:  "Neck",
		3:  "Shoulder",
		4:  "Shirt",
		5:  "Chest",
		6:  "Waist",
		7:  "Legs",
		8:  "Feet",
		9:  "Wrist",
		10: "Hands",
		11: "Finger",
		12: "Trinket",
		13: "One-Hand",
		14: "Off Hand",
		15: "Ranged",
		16: "Back",
		17: "Two-Hand",
		18: "Bag",
		19: "Tabard",
		20: "Robe",
		21: "Main Hand",
		22: "Off Hand",
		23: "Held In Off-Hand",
		24: "Ammo",
		25: "Thrown",
		26: "Ranged",
		27: "Quiver",
		28: "Relic",
	}
	if name, ok := slots[inventoryType]; ok {
		return name
	}
	return ""
}

func getTypeName(class, subclass int) string {
	// Simplified type mapping
	switch class {
	case 2: // Weapon
		weaponTypes := map[int]string{
			0: "Axe", 1: "Axe", 2: "Bow", 3: "Gun", 4: "Mace", 5: "Mace",
			6: "Polearm", 7: "Sword", 8: "Sword", 10: "Staff", 13: "Fist Weapon",
			14: "Miscellaneous", 15: "Dagger", 16: "Thrown", 17: "Spear",
			18: "Crossbow", 19: "Wand", 20: "Fishing Pole",
		}
		if name, ok := weaponTypes[subclass]; ok {
			return name
		}
	case 4: // Armor
		armorTypes := map[int]string{
			0: "Misc", 1: "Cloth", 2: "Leather", 3: "Mail", 4: "Plate",
			5: "Buckler", 6: "Shield", 7: "Libram", 8: "Idol", 9: "Totem",
		}
		if name, ok := armorTypes[subclass]; ok {
			return name
		}
	}
	return ""
}

func parseStats(item *ItemTemplate) []string {
	var stats []string
	statPairs := [][2]int{
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

	for _, pair := range statPairs {
		if pair[1] != 0 {
			if name := getStatName(pair[0]); name != "" {
				prefix := "+"
				if pair[1] < 0 {
					prefix = ""
				}
				stats = append(stats, fmt.Sprintf("%s%d %s", prefix, pair[1], name))
			}
		}
	}

	return stats
}

func parseResistances(item *ItemTemplate) []string {
	var res []string
	if item.ArcaneRes != 0 {
		res = append(res, fmt.Sprintf("+%d Arcane Resistance", item.ArcaneRes))
	}
	if item.FireRes != 0 {
		res = append(res, fmt.Sprintf("+%d Fire Resistance", item.FireRes))
	}
	if item.FrostRes != 0 {
		res = append(res, fmt.Sprintf("+%d Frost Resistance", item.FrostRes))
	}
	if item.NatureRes != 0 {
		res = append(res, fmt.Sprintf("+%d Nature Resistance", item.NatureRes))
	}
	if item.ShadowRes != 0 {
		res = append(res, fmt.Sprintf("+%d Shadow Resistance", item.ShadowRes))
	}
	return res
}

func parseSpellEffects(item *ItemTemplate) []string {
	var effects []string
	// Simplified - would need spell database for full text
	spells := [][2]int{
		{item.SpellID1, item.SpellTrigger1},
		{item.SpellID2, item.SpellTrigger2},
		{item.SpellID3, item.SpellTrigger3},
	}

	for _, spell := range spells {
		if spell[0] > 0 {
			trigger := "Equip:"
			switch spell[1] {
			case 0:
				trigger = "Use:"
			case 1:
				trigger = "Equip:"
			case 2:
				trigger = "Chance on hit:"
			}
			effects = append(effects, fmt.Sprintf("%s Spell Effect #%d", trigger, spell[0]))
		}
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
		39: "Ranged Attack Power", 41: "Spell Healing", 42: "Spell Damage",
		43: "Mana Regeneration", 44: "Armor Penetration Rating", 45: "Spell Power",
		46: "Health Regeneration", 47: "Spell Penetration", 48: "Block Value",
	}
	return statNames[statType]
}

func formatMoney(copper int) string {
	gold := copper / 10000
	silver := (copper % 10000) / 100
	cop := copper % 100

	var parts []string
	if gold > 0 {
		parts = append(parts, fmt.Sprintf("%dg", gold))
	}
	if silver > 0 {
		parts = append(parts, fmt.Sprintf("%ds", silver))
	}
	if cop > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dc", cop))
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}
	return result
}
