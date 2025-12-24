package models

// AtlasLootCategory represents a top-level category
type AtlasLootCategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	SortOrder   int    `json:"sortOrder"`
}

// AtlasLootModule represents a module within a category
type AtlasLootModule struct {
	ID          int    `json:"id"`
	CategoryID  int    `json:"categoryId"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	SortOrder   int    `json:"sortOrder"`
}

// AtlasLootTable represents a boss/loot table
type AtlasLootTable struct {
	ID          int    `json:"id"`
	ModuleID    int    `json:"moduleId"`
	TableKey    string `json:"tableKey"`
	DisplayName string `json:"displayName"`
	SortOrder   int    `json:"sortOrder"`
}

// AtlasLootItem represents a loot entry
type AtlasLootItem struct {
	ID         int    `json:"id"`
	TableID    int    `json:"tableId"`
	ItemID     int    `json:"itemId"`
	DropChance string `json:"dropChance,omitempty"`
	SortOrder  int    `json:"sortOrder"`
}

// AtlasTable represents a loot table reference
type AtlasTable struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
}

// AtlasLoot Import Types

// AtlasLootImportItem represents an item for import
type AtlasLootImportItem struct {
	ID       int    `json:"id"`
	DropRate string `json:"drop_rate"`
}

// AtlasLootImportTable represents a table for import
type AtlasLootImportTable struct {
	Key   string                `json:"key"`
	Name  string                `json:"name"`
	Items []AtlasLootImportItem `json:"items"`
}

// AtlasLootImportModule represents a module for import
type AtlasLootImportModule struct {
	Key    string                 `json:"key"`
	Name   string                 `json:"name"`
	Tables []AtlasLootImportTable `json:"tables"`
}

// AtlasLootImportCategory represents a category for import
type AtlasLootImportCategory struct {
	Key     string                  `json:"key"`
	Name    string                  `json:"name"`
	Sort    int                     `json:"sort"`
	Modules []AtlasLootImportModule `json:"modules"`
}
