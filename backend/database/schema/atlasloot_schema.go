package schema

// AtlasLootSchema returns the SQL statements for AtlasLoot tables
func AtlasLootSchema() string {
	return `
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
}

// LocaleSchema returns the SQL statements for locale tables
func LocaleSchema() string {
	return `
	-- AtlasLoot Locale table for multi-language support
	CREATE TABLE IF NOT EXISTS atlasloot_locale (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		locale_key TEXT NOT NULL,
		language TEXT NOT NULL,  -- 'en', 'cn', 'de', 'fr', 'es', etc.
		text TEXT NOT NULL,
		UNIQUE(locale_key, language)
	);

	CREATE INDEX IF NOT EXISTS idx_atlasloot_locale_key ON atlasloot_locale(locale_key);
	CREATE INDEX IF NOT EXISTS idx_atlasloot_locale_lang ON atlasloot_locale(language);
	`
}
