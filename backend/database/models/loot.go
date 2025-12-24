package models

// LootItem represents an item drop with chance information
type LootItem struct {
	ItemID   int     `json:"itemId"`
	ItemName string  `json:"itemName"`
	Icon     string  `json:"icon"`
	Quality  int     `json:"quality"`
	Chance   float64 `json:"chance"`
	MinCount int     `json:"minCount"`
	MaxCount int     `json:"maxCount"`
}

// LootEntry represents a loot item with metadata (for AtlasLoot)
type LootEntry struct {
	ItemID     int    `json:"itemId"`
	ItemName   string `json:"itemName"`
	IconName   string `json:"iconName"`
	Quality    int    `json:"quality"`
	DropChance string `json:"dropChance,omitempty"`
}

// LootTemplateEntry represents a loot entry for JSON import
type LootTemplateEntry struct {
	Entry         int     `json:"entry"`
	Item          int     `json:"item"`
	Chance        float64 `json:"chance"`
	GroupID       int     `json:"groupId"`
	MinCountOrRef int     `json:"minCountOrRef"`
	MaxCount      int     `json:"maxCount"`
}
