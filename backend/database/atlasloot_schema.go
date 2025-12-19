package database

// InitAtlasLootSchema creates AtlasLoot-specific tables
func (s *SQLiteDB) InitAtlasLootSchema() error {
	schema := `
	-- AtlasLoot Categories (top level: Instances, Sets, Factions, etc.)
	CREATE TABLE IF NOT EXISTS atlasloot_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		display_name TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0
	);

	-- AtlasLoot Modules (e.g., Molten Core for Instances category)
	CREATE TABLE IF NOT EXISTS atlasloot_modules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		display_name TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		FOREIGN KEY (category_id) REFERENCES atlasloot_categories(id)
	);

	-- AtlasLoot Tables (e.g., Ragnaros boss table)
	CREATE TABLE IF NOT EXISTS atlasloot_tables (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		module_id INTEGER NOT NULL,
		table_key TEXT NOT NULL,
		display_name TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		FOREIGN KEY (module_id) REFERENCES atlasloot_modules(id)
	);

	-- AtlasLoot Items (actual loot entries)
	CREATE TABLE IF NOT EXISTS atlasloot_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		table_id INTEGER NOT NULL,
		item_id INTEGER NOT NULL,
		drop_chance TEXT,
		sort_order INTEGER DEFAULT 0,
		FOREIGN KEY (table_id) REFERENCES atlasloot_tables(id),
		FOREIGN KEY (item_id) REFERENCES items(entry)
	);

	CREATE INDEX IF NOT EXISTS idx_atlasloot_modules_category ON atlasloot_modules(category_id);
	CREATE INDEX IF NOT EXISTS idx_atlasloot_tables_module ON atlasloot_tables(module_id);
	CREATE INDEX IF NOT EXISTS idx_atlasloot_items_table ON atlasloot_items(table_id);
	CREATE INDEX IF NOT EXISTS idx_atlasloot_items_item ON atlasloot_items(item_id);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}

	// Initialize locale schema
	return s.InitAtlasLootLocaleSchema()
}

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
