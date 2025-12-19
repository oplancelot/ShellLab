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
	Class    int    `json:"class"`
	SubClass int    `json:"subClass"`
	Name     string `json:"name"`
}

// GetItemClasses returns all item classes with their subclasses
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
		defer subRows.Close()

		for subRows.Next() {
			var sc int
			if err := subRows.Scan(&sc); err != nil {
				continue
			}
			subName := getSubClassName(cls.Class, sc)
			cls.SubClasses = append(cls.SubClasses, &ItemSubClass{
				Class:    cls.Class,
				SubClass: sc,
				Name:     subName,
			})
		}
	}

	return classes, nil
}

// GetItemsByClass returns items filtered by class and subclass
func (r *ItemRepository) GetItemsByClass(class, subClass int, limit, offset int) ([]*Item, int, error) {
	// Count
	var count int
	err := r.db.DB().QueryRow(`
		SELECT COUNT(*) FROM items WHERE class = ? AND subclass = ?
	`, class, subClass).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Data
	rows, err := r.db.DB().Query(`
		SELECT entry, name, quality, item_level, required_level, class, subclass, inventory_type, icon_path
		FROM items
		WHERE class = ? AND subclass = ?
		ORDER BY quality DESC, item_level DESC
		LIMIT ? OFFSET ?
	`, class, subClass, limit, offset)
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
	// Simplified mapping
	if c == 2 { // Weapon
		switch sc {
		case 0:
			return "Axe 1H"
		case 1:
			return "Axe 2H"
		case 2:
			return "Bow"
		case 3:
			return "Gun"
		case 4:
			return "Mace 1H"
		case 5:
			return "Mace 2H"
		case 6:
			return "Polearm"
		case 7:
			return "Sword 1H"
		case 8:
			return "Sword 2H"
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
		case 18:
			return "Crossbow"
		case 19:
			return "Wand"
		case 20:
			return "Fishing Pole"
		}
	}
	if c == 4 { // Armor
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
		case 6:
			return "Shield"
		case 7:
			return "Libram"
		case 8:
			return "Idol"
		case 9:
			return "Totem"
		case 11:
			return "Relic"
		}
	}
	// Add other mappings if needed, or fallback
	return fmt.Sprintf("SubClass %d", sc)
}
