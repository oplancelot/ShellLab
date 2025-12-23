package database

import "fmt"

// ItemClass represents a WoW item class (Weapon, Armor, etc.)
type ItemClass struct {
	Class      int             `json:"class"`
	Name       string          `json:"name"`
	SubClasses []*ItemSubClass `json:"subClasses,omitempty"`
}

// ItemSubClass represents a WoW item subclass (Axe, Bow, etc.)
type ItemSubClass struct {
	Class          int              `json:"class"`
	SubClass       int              `json:"subClass"`
	Name           string           `json:"name"`
	InventorySlots []*InventorySlot `json:"inventorySlots,omitempty"`
}

// InventorySlot represents a WoW inventory type (Head, Chest, etc.)
type InventorySlot struct {
	Class         int    `json:"class"`
	SubClass      int    `json:"subClass"`
	InventoryType int    `json:"inventoryType"`
	Name          string `json:"name"`
}

// GetItemClasses returns all item classes with their subclasses and inventory slots
func (r *ItemRepository) GetItemClasses() ([]*ItemClass, error) {
	// 1. Get distinct classes
	rows, err := r.db.DB().Query(`
		SELECT DISTINCT class
		FROM items
		ORDER BY class
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*ItemClass
	for rows.Next() {
		var c int
		if err := rows.Scan(&c); err != nil {
			continue
		}
		name := getClassName(c)
		classes = append(classes, &ItemClass{Class: c, Name: name})
	}

	// 2. Get subclasses for each class
	for _, cls := range classes {
		subRows, err := r.db.DB().Query(`
			SELECT DISTINCT subclass
			FROM items
			WHERE class = ?
			ORDER BY subclass
		`, cls.Class)
		if err != nil {
			continue
		}

		for subRows.Next() {
			var sc int
			if err := subRows.Scan(&sc); err != nil {
				continue
			}
			subName := getSubClassName(cls.Class, sc)
			subClass := &ItemSubClass{
				Class:    cls.Class,
				SubClass: sc,
				Name:     subName,
			}
			cls.SubClasses = append(cls.SubClasses, subClass)
		}
		subRows.Close()
	}

	// 3. Get inventory types for each subclass
	for _, cls := range classes {
		for _, sub := range cls.SubClasses {
			invRows, err := r.db.DB().Query(`
				SELECT DISTINCT inventory_type
				FROM items
				WHERE class = ? AND subclass = ?
				ORDER BY inventory_type
			`, cls.Class, sub.SubClass)
			if err != nil {
				continue
			}

			for invRows.Next() {
				var invType int
				if err := invRows.Scan(&invType); err != nil {
					continue
				}
				invName := getInventoryTypeName(invType)
				sub.InventorySlots = append(sub.InventorySlots, &InventorySlot{
					Class:         cls.Class,
					SubClass:      sub.SubClass,
					InventoryType: invType,
					Name:          invName,
				})
			}
			invRows.Close()
		}
	}

	return classes, nil
}

// GetItemsByClass returns items filtered by class and subclass
func (r *ItemRepository) GetItemsByClass(class, subClass int, nameFilter string, limit, offset int) ([]*Item, int, error) {
	// Build WHERE clause
	whereClause := "WHERE class = ? AND subclass = ?"
	args := []interface{}{class, subClass}

	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	// Count
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items %s", whereClause)
	err := r.db.DB().QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Data
	dataArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, icon_path
		FROM items
		%s
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	rows, err := r.db.DB().Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
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

	return items, count, nil
}

// GetItemsByClassAndSlot returns items filtered by class, subclass, and inventory type
func (r *ItemRepository) GetItemsByClassAndSlot(class, subClass, inventoryType int, nameFilter string, limit, offset int) ([]*Item, int, error) {
	// Build WHERE clause
	whereClause := "WHERE class = ? AND subclass = ? AND inventory_type = ?"
	args := []interface{}{class, subClass, inventoryType}

	// Add name filter if provided
	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	// Count
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items %s", whereClause)
	err := r.db.DB().QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Data - add limit and offset args
	dataArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, icon_path
		FROM items
		%s
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	rows, err := r.db.DB().Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
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

	return items, count, nil
}

// Helpers for names (simplified map, ideally DB or DBC would provide this)
func getClassName(c int) string {
	switch c {
	case 0:
		return "Consumable"
	case 1:
		return "Container"
	case 2:
		return "Weapon"
	case 3:
		return "Gem"
	case 4:
		return "Armor"
	case 5:
		return "Reagent"
	case 6:
		return "Projectile"
	case 7:
		return "Trade Goods"
	case 8:
		return "Generic"
	case 9:
		return "Recipe"
	case 10:
		return "Money"
	case 11:
		return "Quiver"
	case 12:
		return "Quest"
	case 13:
		return "Key"
	case 14:
		return "Permanent"
	case 15:
		return "Miscellaneous"
	default:
		return fmt.Sprintf("Class %d", c)
	}
}

func getSubClassName(c, sc int) string {
	switch c {
	case 0: // Consumable
		switch sc {
		case 0:
			return "Consumable"
		case 1:
			return "Potion"
		case 2:
			return "Elixir"
		case 3:
			return "Flask"
		case 4:
			return "Scroll"
		case 5:
			return "Food & Drink"
		case 6:
			return "Item Enhancement"
		case 7:
			return "Bandage"
		case 8:
			return "Other"
		}
	case 1: // Container
		switch sc {
		case 0:
			return "Bag"
		case 1:
			return "Soul Bag"
		case 2:
			return "Herb Bag"
		case 3:
			return "Enchanting Bag"
		case 4:
			return "Engineering Bag"
		case 5:
			return "Gem Bag"
		case 6:
			return "Mining Bag"
		case 7:
			return "Leatherworking Bag"
		}
	case 2: // Weapon
		switch sc {
		case 0:
			return "One-Handed Axe"
		case 1:
			return "Two-Handed Axe"
		case 2:
			return "Bow"
		case 3:
			return "Gun"
		case 4:
			return "One-Handed Mace"
		case 5:
			return "Two-Handed Mace"
		case 6:
			return "Polearm"
		case 7:
			return "One-Handed Sword"
		case 8:
			return "Two-Handed Sword"
		case 10:
			return "Staff"
		case 13:
			return "Fist Weapon"
		case 14:
			return "Miscellaneous"
		case 15:
			return "Dagger"
		case 16:
			return "Thrown"
		case 17:
			return "Spear"
		case 18:
			return "Crossbow"
		case 19:
			return "Wand"
		case 20:
			return "Fishing Pole"
		}
	case 3: // Gem
		switch sc {
		case 0:
			return "Red"
		case 1:
			return "Blue"
		case 2:
			return "Yellow"
		case 3:
			return "Purple"
		case 4:
			return "Green"
		case 5:
			return "Orange"
		case 6:
			return "Meta"
		case 7:
			return "Simple"
		case 8:
			return "Prismatic"
		}
	case 4: // Armor
		switch sc {
		case 0:
			return "Miscellaneous"
		case 1:
			return "Cloth"
		case 2:
			return "Leather"
		case 3:
			return "Mail"
		case 4:
			return "Plate"
		case 5:
			return "Buckler"
		case 6:
			return "Shield"
		case 7:
			return "Libram"
		case 8:
			return "Idol"
		case 9:
			return "Totem"
		case 10:
			return "Sigil"
		case 11:
			return "Relic"
		}
	case 5: // Reagent
		return "Reagent"
	case 6: // Projectile
		switch sc {
		case 2:
			return "Arrow"
		case 3:
			return "Bullet"
		}
	case 7: // Trade Goods
		switch sc {
		case 0:
			return "Trade Goods"
		case 1:
			return "Parts"
		case 2:
			return "Explosives"
		case 3:
			return "Devices"
		case 4:
			return "Jewelcrafting"
		case 5:
			return "Cloth"
		case 6:
			return "Leather"
		case 7:
			return "Metal & Stone"
		case 8:
			return "Meat"
		case 9:
			return "Herb"
		case 10:
			return "Elemental"
		case 11:
			return "Other"
		case 12:
			return "Enchanting"
		}
	case 9: // Recipe
		switch sc {
		case 0:
			return "Book"
		case 1:
			return "Leatherworking"
		case 2:
			return "Tailoring"
		case 3:
			return "Engineering"
		case 4:
			return "Blacksmithing"
		case 5:
			return "Cooking"
		case 6:
			return "Alchemy"
		case 7:
			return "First Aid"
		case 8:
			return "Enchanting"
		case 9:
			return "Fishing"
		case 10:
			return "Jewelcrafting"
		}
	case 11: // Quiver
		switch sc {
		case 2:
			return "Quiver"
		case 3:
			return "Ammo Pouch"
		}
	case 12: // Quest
		return "Quest Item"
	case 13: // Key
		switch sc {
		case 0:
			return "Key"
		case 1:
			return "Lockpick"
		}
	case 15: // Miscellaneous
		switch sc {
		case 0:
			return "Junk"
		case 1:
			return "Reagent"
		case 2:
			return "Pet"
		case 3:
			return "Holiday"
		case 4:
			return "Other"
		case 5:
			return "Mount"
		}
	}
	return fmt.Sprintf("SubClass %d", sc)
}

// getInventoryTypeName returns the inventory slot name
func getInventoryTypeName(invType int) string {
	switch invType {
	case 0:
		return "Non-Equipable"
	case 1:
		return "Head"
	case 2:
		return "Neck"
	case 3:
		return "Shoulder"
	case 4:
		return "Shirt"
	case 5:
		return "Chest"
	case 6:
		return "Waist"
	case 7:
		return "Legs"
	case 8:
		return "Feet"
	case 9:
		return "Wrists"
	case 10:
		return "Hands"
	case 11:
		return "Finger"
	case 12:
		return "Trinket"
	case 13:
		return "One-Hand"
	case 14:
		return "Shield"
	case 15:
		return "Ranged"
	case 16:
		return "Back"
	case 17:
		return "Two-Hand"
	case 18:
		return "Bag"
	case 19:
		return "Tabard"
	case 20:
		return "Robe"
	case 21:
		return "Main Hand"
	case 22:
		return "Off Hand"
	case 23:
		return "Holdable"
	case 24:
		return "Ammo"
	case 25:
		return "Thrown"
	case 26:
		return "Ranged Right"
	case 28:
		return "Relic"
	}
	return fmt.Sprintf("Slot %d", invType)
}
