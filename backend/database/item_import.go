package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type ItemTemplateEntry struct {
	Entry          int     `json:"entry"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Quality        int     `json:"quality"`
	ItemLevel      int     `json:"item_level"`
	RequiredLevel  int     `json:"required_level"`
	Class          int     `json:"class"`
	Subclass       int     `json:"subclass"`
	InventoryType  int     `json:"inventory_type"`
	DisplayID      int     `json:"display_id"`
	BuyPrice       int     `json:"buy_price"`
	SellPrice      int     `json:"sell_price"`
	AllowableClass int     `json:"allowable_class"`
	AllowableRace  int     `json:"allowable_race"`
	Stackable      int     `json:"stackable"`
	Bonding        int     `json:"bonding"`
	MaxDurability  int     `json:"max_durability"`
	ContainerSlots int     `json:"container_slots"`
	StatType1      int     `json:"stat_type1"`
	StatValue1     int     `json:"stat_value1"`
	StatType2      int     `json:"stat_type2"`
	StatValue2     int     `json:"stat_value2"`
	StatType3      int     `json:"stat_type3"`
	StatValue3     int     `json:"stat_value3"`
	StatType4      int     `json:"stat_type4"`
	StatValue4     int     `json:"stat_value4"`
	StatType5      int     `json:"stat_type5"`
	StatValue5     int     `json:"stat_value5"`
	StatType6      int     `json:"stat_type6"`
	StatValue6     int     `json:"stat_value6"`
	StatType7      int     `json:"stat_type7"`
	StatValue7     int     `json:"stat_value7"`
	StatType8      int     `json:"stat_type8"`
	StatValue8     int     `json:"stat_value8"`
	StatType9      int     `json:"stat_type9"`
	StatValue9     int     `json:"stat_value9"`
	StatType10     int     `json:"stat_type10"`
	StatValue10    int     `json:"stat_value10"`
	Delay          int     `json:"delay"`
	DmgMin1        float64 `json:"dmg_min1"`
	DmgMax1        float64 `json:"dmg_max1"`
	DmgType1       int     `json:"dmg_type1"`
	DmgMin2        float64 `json:"dmg_min2"`
	DmgMax2        float64 `json:"dmg_max2"`
	DmgType2       int     `json:"dmg_type2"`
	Armor          int     `json:"armor"`
	HolyRes        int     `json:"holy_res"`
	FireRes        int     `json:"fire_res"`
	NatureRes      int     `json:"nature_res"`
	FrostRes       int     `json:"frost_res"`
	ShadowRes      int     `json:"shadow_res"`
	ArcaneRes      int     `json:"arcane_res"`
	SpellID1       int     `json:"spellid_1"`
	SpellTrigger1  int     `json:"spelltrigger_1"`
	SpellID2       int     `json:"spellid_2"`
	SpellTrigger2  int     `json:"spelltrigger_2"`
	SpellID3       int     `json:"spellid_3"`
	SpellTrigger3  int     `json:"spelltrigger_3"`
	ItemSet        int     `json:"set_id"`
}

// ImportItemsFromJSON imports items from JSON into SQLite
func (r *ItemRepository) ImportItemsFromJSON(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read opening bracket
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use REPLACE INTO to update existing or insert new
	stmt, err := tx.Prepare(`
		REPLACE INTO items (
			entry, name, description, quality, item_level, required_level,
			class, subclass, inventory_type, display_id,
			buy_price, sell_price, allowable_class, allowable_race,
			max_stack, bonding, max_durability,
			stat_type1, stat_value1, stat_type2, stat_value2,
			stat_type3, stat_value3, stat_type4, stat_value4,
			stat_type5, stat_value5, stat_type6, stat_value6,
			stat_type7, stat_value7, stat_type8, stat_value8,
			stat_type9, stat_value9, stat_type10, stat_value10,
			delay, dmg_min1, dmg_max1, dmg_type1,
			dmg_min2, dmg_max2, dmg_type2,
			armor, holy_res, fire_res, nature_res, frost_res, shadow_res, arcane_res,
			spell_id1, spell_trigger1, spell_id2, spell_trigger2,
			spell_id3, spell_trigger3,
			set_id
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?,
			?
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for decoder.More() {
		var item ItemTemplateEntry
		if err := decoder.Decode(&item); err != nil {
			fmt.Printf("Error decoding item: %v\n", err)
			continue
		}

		_, err := stmt.Exec(
			item.Entry, item.Name, item.Description, item.Quality, item.ItemLevel, item.RequiredLevel,
			item.Class, item.Subclass, item.InventoryType, item.DisplayID,
			item.BuyPrice, item.SellPrice, item.AllowableClass, item.AllowableRace,
			item.Stackable, item.Bonding, item.MaxDurability,
			item.StatType1, item.StatValue1, item.StatType2, item.StatValue2,
			item.StatType3, item.StatValue3, item.StatType4, item.StatValue4,
			item.StatType5, item.StatValue5, item.StatType6, item.StatValue6,
			item.StatType7, item.StatValue7, item.StatType8, item.StatValue8,
			item.StatType9, item.StatValue9, item.StatType10, item.StatValue10,
			item.Delay, item.DmgMin1, item.DmgMax1, item.DmgType1,
			item.DmgMin2, item.DmgMax2, item.DmgType2,
			item.Armor, item.HolyRes, item.FireRes, item.NatureRes, item.FrostRes, item.ShadowRes, item.ArcaneRes,
			item.SpellID1, item.SpellTrigger1, item.SpellID2, item.SpellTrigger2,
			item.SpellID3, item.SpellTrigger3,
			item.ItemSet,
		)
		if err != nil {
			fmt.Printf("Error inserting item %d: %v\n", item.Entry, err)
			continue
		}
		count++
		if count%1000 == 0 {
			fmt.Printf("Imported %d items...\n", count)
		}
	}

	return tx.Commit()
}

// CheckAndImportItems checks if items table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportItems(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if err != nil {
		// Table might not exist yet, schema init should handle it, but just in case
		return err
	}

	if count == 0 {
		path := fmt.Sprintf("%s/item_template.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Items table is empty. Importing from item_template.json...")
			return r.ImportItemsFromJSON(path)
		}
	}
	return nil
}
