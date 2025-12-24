// Package importers contains data import logic
package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"shelllab/backend/database/models"
)

// ItemImporter handles item data imports
type ItemImporter struct {
	db *sql.DB
}

// NewItemImporter creates a new item importer
func NewItemImporter(db *sql.DB) *ItemImporter {
	return &ItemImporter{db: db}
}

// ImportFromJSON imports items from JSON into SQLite
func (i *ItemImporter) ImportFromJSON(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
			spell_id3, spell_trigger3, set_id
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
			?, ?, ?
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for decoder.More() {
		var item models.ItemTemplateEntry
		if err := decoder.Decode(&item); err != nil {
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
			item.SpellID3, item.SpellTrigger3, item.ItemSet,
		)
		if err != nil {
			continue
		}
		count++
		if count%1000 == 0 {
			fmt.Printf("Imported %d items...\n", count)
		}
	}
	return tx.Commit()
}

// CheckAndImport checks if items table is empty and imports if JSON exists
func (i *ItemImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := i.db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		path := fmt.Sprintf("%s/item_template.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Items table is empty. Importing from item_template.json...")
			return i.ImportFromJSON(path)
		}
	}
	return nil
}
