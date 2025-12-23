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
	CREATE TABLE IF NOT EXISTS itemsets (
		itemset_id INTEGER PRIMARY KEY,
		name TEXT,
		item1 INTEGER DEFAULT 0, item2 INTEGER DEFAULT 0, item3 INTEGER DEFAULT 0, item4 INTEGER DEFAULT 0, item5 INTEGER DEFAULT 0,
		item6 INTEGER DEFAULT 0, item7 INTEGER DEFAULT 0, item8 INTEGER DEFAULT 0, item9 INTEGER DEFAULT 0, item10 INTEGER DEFAULT 0,
		spell1 INTEGER DEFAULT 0, spell2 INTEGER DEFAULT 0, spell3 INTEGER DEFAULT 0, spell4 INTEGER DEFAULT 0,
		spell5 INTEGER DEFAULT 0, spell6 INTEGER DEFAULT 0, spell7 INTEGER DEFAULT 0, spell8 INTEGER DEFAULT 0,
		bonus1 INTEGER DEFAULT 0, bonus2 INTEGER DEFAULT 0, bonus3 INTEGER DEFAULT 0, bonus4 INTEGER DEFAULT 0,
		bonus5 INTEGER DEFAULT 0, bonus6 INTEGER DEFAULT 0, bonus7 INTEGER DEFAULT 0, bonus8 INTEGER DEFAULT 0,
		skill_id INTEGER DEFAULT 0, skill_level INTEGER DEFAULT 0
	);

	-- Quests table
	CREATE TABLE IF NOT EXISTS quests (
		entry INTEGER PRIMARY KEY, 
		title TEXT, 
		min_level INTEGER, 
		quest_level INTEGER, 
		type INTEGER, 
		zone_or_sort INTEGER,
		details TEXT, 
		objectives TEXT, 
		offer_reward_text TEXT, 
		end_text TEXT,
		rew_xp INTEGER, 
		rew_money INTEGER, 
		rew_money_max_level INTEGER, 
		rew_spell INTEGER,
		rew_item1 INTEGER, rew_item2 INTEGER, rew_item3 INTEGER, rew_item4 INTEGER,
		rew_item_count1 INTEGER, rew_item_count2 INTEGER, rew_item_count3 INTEGER, rew_item_count4 INTEGER,
		rew_choice_item1 INTEGER, rew_choice_item2 INTEGER, rew_choice_item3 INTEGER, rew_choice_item4 INTEGER, rew_choice_item5 INTEGER, rew_choice_item6 INTEGER,
		rew_choice_item_count1 INTEGER, rew_choice_item_count2 INTEGER, rew_choice_item_count3 INTEGER, rew_choice_item_count4 INTEGER, rew_choice_item_count5 INTEGER, rew_choice_item_count6 INTEGER,
		rew_rep_faction1 INTEGER, rew_rep_faction2 INTEGER, rew_rep_faction3 INTEGER, rew_rep_faction4 INTEGER, rew_rep_faction5 INTEGER,
		rew_rep_value1 INTEGER, rew_rep_value2 INTEGER, rew_rep_value3 INTEGER, rew_rep_value4 INTEGER, rew_rep_value5 INTEGER,
		prev_quest_id INTEGER, next_quest_id INTEGER, exclusive_group INTEGER, next_quest_in_chain INTEGER,
		required_races INTEGER, required_classes INTEGER, src_item_id INTEGER
	);

	CREATE INDEX IF NOT EXISTS idx_quests_title ON quests(title);
	CREATE INDEX IF NOT EXISTS idx_quests_zone ON quests(zone_or_sort);

	-- Quest Relations
	CREATE TABLE IF NOT EXISTS npc_quest_start (entry INTEGER, quest INTEGER, PRIMARY KEY(entry, quest));
	CREATE TABLE IF NOT EXISTS npc_quest_end (entry INTEGER, quest INTEGER, PRIMARY KEY(entry, quest));
	CREATE TABLE IF NOT EXISTS go_quest_start (entry INTEGER, quest INTEGER, PRIMARY KEY(entry, quest));
	CREATE TABLE IF NOT EXISTS go_quest_end (entry INTEGER, quest INTEGER, PRIMARY KEY(entry, quest));

	-- Objects table (gameobject_template)
	CREATE TABLE IF NOT EXISTS objects (
		entry INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		type INTEGER NOT NULL,
		display_id INTEGER,
		size REAL,
		data0 INTEGER, data1 INTEGER, data2 INTEGER, data3 INTEGER,
		data4 INTEGER, data5 INTEGER, data6 INTEGER, data7 INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_objects_type ON objects(type);
	CREATE INDEX IF NOT EXISTS idx_objects_name ON objects(name);

	-- Locks table (aowow_lock)
	CREATE TABLE IF NOT EXISTS locks (
		id INTEGER PRIMARY KEY,
		type1 INTEGER, type2 INTEGER, type3 INTEGER, type4 INTEGER, type5 INTEGER,
		prop1 INTEGER, prop2 INTEGER, prop3 INTEGER, prop4 INTEGER, prop5 INTEGER,
		req1 INTEGER, req2 INTEGER, req3 INTEGER, req4 INTEGER, req5 INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_locks_prop1 ON locks(prop1);

	-- Loot tables
	CREATE TABLE IF NOT EXISTS creature_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL,
		groupid INTEGER,
		mincount_or_ref INTEGER,
		maxcount INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_creature_loot_entry ON creature_loot(entry);

	CREATE TABLE IF NOT EXISTS reference_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL,
		groupid INTEGER,
		mincount_or_ref INTEGER,
		maxcount INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_reference_loot_entry ON reference_loot(entry);
	
	CREATE TABLE IF NOT EXISTS gameobject_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL,
		groupid INTEGER,
		mincount_or_ref INTEGER,
		maxcount INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_gameobject_loot_entry ON gameobject_loot(entry);
	
	CREATE TABLE IF NOT EXISTS item_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL,
		groupid INTEGER,
		mincount_or_ref INTEGER,
		maxcount INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_item_loot_entry ON item_loot(entry);
	
	CREATE TABLE IF NOT EXISTS disenchant_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL,
		groupid INTEGER,
		mincount_or_ref INTEGER,
		maxcount INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_disenchant_loot_entry ON disenchant_loot(entry);

	-- Creatures table (creature_template)
	CREATE TABLE IF NOT EXISTS creatures (
		entry INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		subname TEXT,
		level_min INTEGER,
		level_max INTEGER,
		health_min INTEGER,
		health_max INTEGER,
		mana_min INTEGER,
		mana_max INTEGER,
		creature_type INTEGER,
		creature_rank INTEGER,
		faction INTEGER,
		npc_flags INTEGER,
		loot_id INTEGER,
		skin_loot_id INTEGER,
		pickpocket_loot_id INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_creatures_name ON creatures(name);
	CREATE INDEX IF NOT EXISTS idx_creatures_name ON creatures(name);
	CREATE INDEX IF NOT EXISTS idx_creatures_type ON creatures(creature_type);

	-- Factions table
	CREATE TABLE IF NOT EXISTS factions (
		id INTEGER PRIMARY KEY,
		name TEXT,
		description TEXT,
		side INTEGER,
		category_id INTEGER
	);

	-- Spells table
	CREATE TABLE IF NOT EXISTS spells (
		entry INTEGER PRIMARY KEY,
		name TEXT,
		description TEXT,
		effect_base_points1 INTEGER DEFAULT 0,
		effect_base_points2 INTEGER DEFAULT 0,
		effect_base_points3 INTEGER DEFAULT 0,
		effect_die_sides1 INTEGER DEFAULT 0,
		effect_die_sides2 INTEGER DEFAULT 0,
		effect_die_sides3 INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_spells_name ON spells(name);

	-- Quest Category Metadata
	CREATE TABLE IF NOT EXISTS quest_category_groups (
		id INTEGER PRIMARY KEY,
		name TEXT
	);

	CREATE TABLE IF NOT EXISTS quest_categories_enhanced (
		id INTEGER PRIMARY KEY, -- zone_or_sort
		group_id INTEGER,
		name TEXT,
		quest_count INTEGER DEFAULT 0
	);

	-- Spell Skill Metadata
	CREATE TABLE IF NOT EXISTS spell_skill_categories (
		id INTEGER PRIMARY KEY,
		name TEXT
	);

	CREATE TABLE IF NOT EXISTS spell_skills (
		id INTEGER PRIMARY KEY,
		category_id INTEGER,
		name TEXT
	);

	CREATE TABLE IF NOT EXISTS spell_skill_spells (
		skill_id INTEGER,
		spell_id INTEGER,
		PRIMARY KEY(skill_id, spell_id)
	);
	
	-- AtlasLoot Categories (for legacy API compatibility if needed)
	CREATE TABLE IF NOT EXISTS atlasloot_categories (
		id INTEGER PRIMARY KEY,
		name TEXT
	);
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
