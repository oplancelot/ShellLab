// Package schema contains database schema definitions
package schema

// CoreSchema returns the SQL statements for core tables
func CoreSchema() string {
	return `
	-- Items table
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
		buy_price INTEGER DEFAULT 0,
		sell_price INTEGER DEFAULT 0,
		allowable_class INTEGER DEFAULT -1,
		allowable_race INTEGER DEFAULT -1,
		max_stack INTEGER DEFAULT 1,
		bonding INTEGER DEFAULT 0,
		max_durability INTEGER DEFAULT 0,
		flags INTEGER DEFAULT 0,
		buy_count INTEGER DEFAULT 1,
		max_count INTEGER DEFAULT 0,
		stackable INTEGER DEFAULT 1,
		container_slots INTEGER DEFAULT 0,
		material INTEGER DEFAULT 0,
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
		-- Damage
		delay INTEGER DEFAULT 0,
		dmg_min1 REAL DEFAULT 0, dmg_max1 REAL DEFAULT 0, dmg_type1 INTEGER DEFAULT 0,
		dmg_min2 REAL DEFAULT 0, dmg_max2 REAL DEFAULT 0, dmg_type2 INTEGER DEFAULT 0,
		-- Armor & Resistances
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
		set_id INTEGER DEFAULT 0,
		-- Icon (from aowow data)
		icon TEXT DEFAULT '',
		icon_path TEXT DEFAULT ''
	);

	CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
	CREATE INDEX IF NOT EXISTS idx_items_quality ON items(quality);
	CREATE INDEX IF NOT EXISTS idx_items_class ON items(class);
	CREATE INDEX IF NOT EXISTS idx_items_subclass ON items(subclass);
	CREATE INDEX IF NOT EXISTS idx_items_set_id ON items(set_id);

	-- Item Sets
	CREATE TABLE IF NOT EXISTS itemsets (
		itemset_id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		item1 INTEGER DEFAULT 0, item2 INTEGER DEFAULT 0, item3 INTEGER DEFAULT 0,
		item4 INTEGER DEFAULT 0, item5 INTEGER DEFAULT 0, item6 INTEGER DEFAULT 0,
		item7 INTEGER DEFAULT 0, item8 INTEGER DEFAULT 0, item9 INTEGER DEFAULT 0,
		item10 INTEGER DEFAULT 0,
		skill_id INTEGER DEFAULT 0,
		skill_level INTEGER DEFAULT 0,
		bonus1 INTEGER DEFAULT 0, bonus2 INTEGER DEFAULT 0, bonus3 INTEGER DEFAULT 0,
		bonus4 INTEGER DEFAULT 0, bonus5 INTEGER DEFAULT 0, bonus6 INTEGER DEFAULT 0,
		bonus7 INTEGER DEFAULT 0, bonus8 INTEGER DEFAULT 0,
		spell1 INTEGER DEFAULT 0, spell2 INTEGER DEFAULT 0, spell3 INTEGER DEFAULT 0,
		spell4 INTEGER DEFAULT 0, spell5 INTEGER DEFAULT 0, spell6 INTEGER DEFAULT 0,
		spell7 INTEGER DEFAULT 0, spell8 INTEGER DEFAULT 0
	);

	-- Creatures
	CREATE TABLE IF NOT EXISTS creatures (
		entry INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		subname TEXT,
		level_min INTEGER DEFAULT 1,
		level_max INTEGER DEFAULT 1,
		health_min INTEGER DEFAULT 1,
		health_max INTEGER DEFAULT 1,
		mana_min INTEGER DEFAULT 0,
		mana_max INTEGER DEFAULT 0,
		creature_type INTEGER DEFAULT 0,
		creature_rank INTEGER DEFAULT 0,
		faction INTEGER DEFAULT 0,
		npc_flags INTEGER DEFAULT 0,
		loot_id INTEGER DEFAULT 0,
		skin_loot_id INTEGER DEFAULT 0,
		pickpocket_loot_id INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_creatures_name ON creatures(name);
	CREATE INDEX IF NOT EXISTS idx_creatures_type ON creatures(creature_type);
	CREATE INDEX IF NOT EXISTS idx_creatures_loot ON creatures(loot_id);

	-- Quests
	CREATE TABLE IF NOT EXISTS quests (
		entry INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		min_level INTEGER DEFAULT 0,
		quest_level INTEGER DEFAULT 0,
		type INTEGER DEFAULT 0,
		zone_or_sort INTEGER DEFAULT 0,
		details TEXT,
		objectives TEXT,
		offer_reward_text TEXT,
		end_text TEXT,
		rew_xp INTEGER DEFAULT 0,
		rew_money INTEGER DEFAULT 0,
		rew_money_max_level INTEGER DEFAULT 0,
		rew_spell INTEGER DEFAULT 0,
		rew_item1 INTEGER DEFAULT 0, rew_item2 INTEGER DEFAULT 0,
		rew_item3 INTEGER DEFAULT 0, rew_item4 INTEGER DEFAULT 0,
		rew_item_count1 INTEGER DEFAULT 0, rew_item_count2 INTEGER DEFAULT 0,
		rew_item_count3 INTEGER DEFAULT 0, rew_item_count4 INTEGER DEFAULT 0,
		rew_choice_item1 INTEGER DEFAULT 0, rew_choice_item2 INTEGER DEFAULT 0,
		rew_choice_item3 INTEGER DEFAULT 0, rew_choice_item4 INTEGER DEFAULT 0,
		rew_choice_item5 INTEGER DEFAULT 0, rew_choice_item6 INTEGER DEFAULT 0,
		rew_choice_item_count1 INTEGER DEFAULT 0, rew_choice_item_count2 INTEGER DEFAULT 0,
		rew_choice_item_count3 INTEGER DEFAULT 0, rew_choice_item_count4 INTEGER DEFAULT 0,
		rew_choice_item_count5 INTEGER DEFAULT 0, rew_choice_item_count6 INTEGER DEFAULT 0,
		rew_rep_faction1 INTEGER DEFAULT 0, rew_rep_faction2 INTEGER DEFAULT 0,
		rew_rep_faction3 INTEGER DEFAULT 0, rew_rep_faction4 INTEGER DEFAULT 0,
		rew_rep_faction5 INTEGER DEFAULT 0,
		rew_rep_value1 INTEGER DEFAULT 0, rew_rep_value2 INTEGER DEFAULT 0,
		rew_rep_value3 INTEGER DEFAULT 0, rew_rep_value4 INTEGER DEFAULT 0,
		rew_rep_value5 INTEGER DEFAULT 0,
		prev_quest_id INTEGER DEFAULT 0,
		next_quest_id INTEGER DEFAULT 0,
		exclusive_group INTEGER DEFAULT 0,
		next_quest_in_chain INTEGER DEFAULT 0,
		required_races INTEGER DEFAULT 0,
		required_classes INTEGER DEFAULT 0,
		src_item_id INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_quests_title ON quests(title);
	CREATE INDEX IF NOT EXISTS idx_quests_zone ON quests(zone_or_sort);

	-- Quest Categories
	CREATE TABLE IF NOT EXISTS quest_categories (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	-- Quest Category Groups
	CREATE TABLE IF NOT EXISTS quest_category_groups (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	-- Quest Categories Enhanced
	CREATE TABLE IF NOT EXISTS quest_categories_enhanced (
		id INTEGER PRIMARY KEY,
		group_id INTEGER DEFAULT 0,
		name TEXT NOT NULL,
		quest_count INTEGER DEFAULT 0
	);

	-- NPC Quest Start/End
	CREATE TABLE IF NOT EXISTS npc_quest_start (
		entry INTEGER,
		quest INTEGER,
		PRIMARY KEY (entry, quest)
	);

	CREATE TABLE IF NOT EXISTS npc_quest_end (
		entry INTEGER,
		quest INTEGER,
		PRIMARY KEY (entry, quest)
	);

	-- Spells
	CREATE TABLE IF NOT EXISTS spells (
		entry INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		effect_base_points1 INTEGER DEFAULT 0,
		effect_base_points2 INTEGER DEFAULT 0,
		effect_base_points3 INTEGER DEFAULT 0,
		effect_die_sides1 INTEGER DEFAULT 0,
		effect_die_sides2 INTEGER DEFAULT 0,
		effect_die_sides3 INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_spells_name ON spells(name);

	-- Spell Skill Categories
	CREATE TABLE IF NOT EXISTS spell_skill_categories (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	-- Spell Skills
	CREATE TABLE IF NOT EXISTS spell_skills (
		id INTEGER PRIMARY KEY,
		category_id INTEGER DEFAULT 0,
		name TEXT NOT NULL
	);

	-- Spell Skill Spells
	CREATE TABLE IF NOT EXISTS spell_skill_spells (
		skill_id INTEGER,
		spell_id INTEGER,
		PRIMARY KEY (skill_id, spell_id)
	);

	-- Loot Tables
	CREATE TABLE IF NOT EXISTS creature_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL DEFAULT 0,
		groupid INTEGER DEFAULT 0,
		mincount_or_ref INTEGER DEFAULT 1,
		maxcount INTEGER DEFAULT 1,
		PRIMARY KEY (entry, item)
	);

	CREATE TABLE IF NOT EXISTS reference_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL DEFAULT 0,
		groupid INTEGER DEFAULT 0,
		mincount_or_ref INTEGER DEFAULT 1,
		maxcount INTEGER DEFAULT 1,
		PRIMARY KEY (entry, item)
	);

	CREATE TABLE IF NOT EXISTS gameobject_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL DEFAULT 0,
		groupid INTEGER DEFAULT 0,
		mincount_or_ref INTEGER DEFAULT 1,
		maxcount INTEGER DEFAULT 1,
		PRIMARY KEY (entry, item)
	);

	CREATE TABLE IF NOT EXISTS item_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL DEFAULT 0,
		groupid INTEGER DEFAULT 0,
		mincount_or_ref INTEGER DEFAULT 1,
		maxcount INTEGER DEFAULT 1,
		PRIMARY KEY (entry, item)
	);

	CREATE TABLE IF NOT EXISTS disenchant_loot (
		entry INTEGER,
		item INTEGER,
		chance REAL DEFAULT 0,
		groupid INTEGER DEFAULT 0,
		mincount_or_ref INTEGER DEFAULT 1,
		maxcount INTEGER DEFAULT 1,
		PRIMARY KEY (entry, item)
	);

	CREATE INDEX IF NOT EXISTS idx_creature_loot_item ON creature_loot(item);
	CREATE INDEX IF NOT EXISTS idx_reference_loot_entry ON reference_loot(entry);

	-- Game Objects
	CREATE TABLE IF NOT EXISTS objects (
		entry INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		type INTEGER DEFAULT 0,
		display_id INTEGER DEFAULT 0,
		size REAL DEFAULT 1.0,
		data0 INTEGER DEFAULT 0, data1 INTEGER DEFAULT 0, data2 INTEGER DEFAULT 0,
		data3 INTEGER DEFAULT 0, data4 INTEGER DEFAULT 0, data5 INTEGER DEFAULT 0,
		data6 INTEGER DEFAULT 0, data7 INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_objects_name ON objects(name);
	CREATE INDEX IF NOT EXISTS idx_objects_type ON objects(type);

	-- Locks (for object requirements)
	CREATE TABLE IF NOT EXISTS locks (
		id INTEGER PRIMARY KEY,
		type1 INTEGER DEFAULT 0, type2 INTEGER DEFAULT 0, type3 INTEGER DEFAULT 0,
		type4 INTEGER DEFAULT 0, type5 INTEGER DEFAULT 0,
		prop1 INTEGER DEFAULT 0, prop2 INTEGER DEFAULT 0, prop3 INTEGER DEFAULT 0,
		prop4 INTEGER DEFAULT 0, prop5 INTEGER DEFAULT 0,
		req1 INTEGER DEFAULT 0, req2 INTEGER DEFAULT 0, req3 INTEGER DEFAULT 0,
		req4 INTEGER DEFAULT 0, req5 INTEGER DEFAULT 0
	);

	-- Factions
	CREATE TABLE IF NOT EXISTS factions (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		side INTEGER DEFAULT 0,
		category_id INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_factions_name ON factions(name);

	-- Categories (for general categorization)
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		parent_id INTEGER,
		type TEXT,
		sort_order INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS category_items (
		category_id INTEGER,
		item_id INTEGER,
		drop_rate TEXT,
		sort_order INTEGER DEFAULT 0,
		PRIMARY KEY (category_id, item_id)
	);
	`
}
