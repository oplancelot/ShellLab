package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

// SQLiteDB wraps the SQLite database connection
type SQLiteDB struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// InitSchema creates the database schema if it doesn't exist
func (s *SQLiteDB) InitSchema() error {
	schema := `
	-- Items table (from item_template.json)
	CREATE TABLE IF NOT EXISTS items (
		entry INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		quality INTEGER DEFAULT 0,
		item_level INTEGER DEFAULT 0,
		required_level INTEGER DEFAULT 0,
		class INTEGER DEFAULT 0,
		subclass INTEGER DEFAULT 0,
		inventory_type INTEGER DEFAULT 0,
		display_id INTEGER DEFAULT 0,
		icon_path TEXT,
		buy_price INTEGER DEFAULT 0,
		sell_price INTEGER DEFAULT 0,
		allowable_class INTEGER DEFAULT -1,
		allowable_race INTEGER DEFAULT -1,
		max_stack INTEGER DEFAULT 1,
		bonding INTEGER DEFAULT 0,
		max_durability INTEGER DEFAULT 0,
		-- Stats
		stat_type1 INTEGER DEFAULT 0, stat_value1 INTEGER DEFAULT 0,
		stat_type2 INTEGER DEFAULT 0, stat_value2 INTEGER DEFAULT 0,
		stat_type3 INTEGER DEFAULT 0, stat_value3 INTEGER DEFAULT 0,
		stat_type4 INTEGER DEFAULT 0, stat_value4 INTEGER DEFAULT 0,
		stat_type5 INTEGER DEFAULT 0, stat_value5 INTEGER DEFAULT 0,
		stat_type6 INTEGER DEFAULT 0, stat_value6 INTEGER DEFAULT 0,
		stat_type7 INTEGER DEFAULT 0, stat_value7 INTEGER DEFAULT 0,
		stat_type8 INTEGER DEFAULT 0, stat_value8 INTEGER DEFAULT 0,
		stat_type9 INTEGER DEFAULT 0, stat_value9 INTEGER DEFAULT 0,
		stat_type10 INTEGER DEFAULT 0, stat_value10 INTEGER DEFAULT 0,
		-- Weapon
		delay INTEGER DEFAULT 0,
		dmg_min1 REAL DEFAULT 0,
		dmg_max1 REAL DEFAULT 0,
		dmg_type1 INTEGER DEFAULT 0,
		dmg_min2 REAL DEFAULT 0,
		dmg_max2 REAL DEFAULT 0,
		dmg_type2 INTEGER DEFAULT 0,
		-- Armor & Resistance
		armor INTEGER DEFAULT 0,
		holy_res INTEGER DEFAULT 0,
		fire_res INTEGER DEFAULT 0,
		nature_res INTEGER DEFAULT 0,
		frost_res INTEGER DEFAULT 0,
		shadow_res INTEGER DEFAULT 0,
		arcane_res INTEGER DEFAULT 0,
		-- Spells
		spell_id1 INTEGER DEFAULT 0, spell_trigger1 INTEGER DEFAULT 0,
		spell_id2 INTEGER DEFAULT 0, spell_trigger2 INTEGER DEFAULT 0,
		spell_id3 INTEGER DEFAULT 0, spell_trigger3 INTEGER DEFAULT 0,
		-- Set
		set_id INTEGER DEFAULT 0
	);

	-- Create indexes for common queries
	CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
	CREATE INDEX IF NOT EXISTS idx_items_quality ON items(quality);
	CREATE INDEX IF NOT EXISTS idx_items_class ON items(class);
	CREATE INDEX IF NOT EXISTS idx_items_subclass ON items(subclass);
	CREATE INDEX IF NOT EXISTS idx_items_inventory_type ON items(inventory_type);
	CREATE INDEX IF NOT EXISTS idx_items_item_level ON items(item_level);
	CREATE INDEX IF NOT EXISTS idx_items_required_level ON items(required_level);
	CREATE INDEX IF NOT EXISTS idx_items_set_id ON items(set_id);

	-- Categories table (AtlasLoot classification)
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		parent_id INTEGER,
		type TEXT NOT NULL, -- 'root', 'instance', 'boss', 'set', 'faction', 'pvp', 'crafting', 'world_event'
		sort_order INTEGER DEFAULT 0,
		FOREIGN KEY (parent_id) REFERENCES categories(id)
	);

	CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);
	CREATE INDEX IF NOT EXISTS idx_categories_type ON categories(type);
	CREATE INDEX IF NOT EXISTS idx_categories_key ON categories(key);

	-- Category-Item association (boss drops, set items, etc.)
	CREATE TABLE IF NOT EXISTS category_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		item_id INTEGER NOT NULL,
		drop_rate TEXT,
		sort_order INTEGER DEFAULT 0,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (item_id) REFERENCES items(entry)
	);

	CREATE INDEX IF NOT EXISTS idx_category_items_category ON category_items(category_id);
	CREATE INDEX IF NOT EXISTS idx_category_items_item ON category_items(item_id);

	-- Item Sets table
	CREATE TABLE IF NOT EXISTS item_sets (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		items_json TEXT -- JSON array of item IDs
	);

	-- Future: quests, spells, npcs, objects tables
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}

	// Initialize AtlasLoot schema
	return s.InitAtlasLootSchema()
}

// DB returns the underlying sql.DB for direct queries
func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}
