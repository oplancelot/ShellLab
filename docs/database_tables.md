# MySQL Database Schema Documentation

## Database Overview

| Database   | Tables | Purpose                                               |
| ---------- | ------ | ----------------------------------------------------- |
| `aowow`    | 22     | AOWOW Website Data (Icons, Spells, Factions metadata) |
| `tw_world` | 180+   | Turtle WoW World Data (Items, NPCs, Quests, etc.)     |

---

## aowow Database (22 Tables)

### Character Related

| Table Name          | Purpose          |
| ------------------- | ---------------- |
| `aowow_char_titles` | Character Titles |

### Faction Related

| Table Name              | Purpose                       |
| ----------------------- | ----------------------------- |
| `aowow_factions`        | Faction Data                  |
| `aowow_factiontemplate` | Faction Templates (Relations) |

### Item Related

| Table Name              | Purpose                              |
| ----------------------- | ------------------------------------ |
| `aowow_icons`           | Item/Spell Icon Name Mappings        |
| `aowow_itemenchantment` | Item Enchantment Effects             |
| `aowow_itemset`         | Item Set Data (components & bonuses) |

### Skill Related

| Table Name                 | Purpose                      |
| -------------------------- | ---------------------------- |
| `aowow_skill`              | Skill Data                   |
| `aowow_skill_line_ability` | Skill Line Ability Relations |

### Spell Related

| Table Name              | Purpose         |
| ----------------------- | --------------- |
| `aowow_spell`           | Spell Base Data |
| `aowow_spellcasttimes`  | Cast Times      |
| `aowow_spelldispeltype` | Dispel Types    |
| `aowow_spellduration`   | Spell Durations |
| `aowow_spellicons`      | Spell Icons     |
| `aowow_spellmechanic`   | Spell Mechanics |
| `aowow_spellradius`     | Spell Radius    |
| `aowow_spellrange`      | Spell Range     |

### Others

| Table Name             | Purpose         |
| ---------------------- | --------------- |
| `aowow_comments`       | User Comments   |
| `aowow_comments_rates` | Comment Ratings |
| `aowow_lock`           | Lock Data       |
| `aowow_news`           | News            |
| `aowow_resistances`    | Resistance Data |
| `aowow_zones`          | Zone/Map Data   |

---

## tw_world Database (Main Tables)

### Item Related

| Table Name                   | Purpose                   |
| ---------------------------- | ------------------------- |
| `item_template`              | Item Template (Base data) |
| `item_display_info`          | Item Display Info         |
| `item_enchantment_template`  | Item Enchantment Template |
| `item_loot_template`         | Item Loot Template        |
| `item_required_target`       | Item Required Target      |
| `item_transmogrify_template` | Transmogrify Template     |
| `locales_item`               | Item Localization         |

### NPC/Creature Related

| Table Name                    | Purpose                     |
| ----------------------------- | --------------------------- |
| `creature_template`           | Creature Template           |
| `creature`                    | Creature Instances (Spawns) |
| `creature_addon`              | Creature Addon Data (Auras) |
| `creature_ai_events`          | Creature AI Events          |
| `creature_ai_scripts`         | Creature AI Scripts         |
| `creature_equip_template`     | Creature Equipment Template |
| `creature_groups`             | Creature Groups             |
| `creature_involvedrelation`   | Quest Relations (Finisher)  |
| `creature_loot_template`      | Creature Loot Template      |
| `creature_movement`           | Creature Movement Paths     |
| `creature_movement_template`  | Creature Movement Template  |
| `creature_onkill_reputation`  | On-Kill Reputation          |
| `creature_questrelation`      | Quest Relations (Starter)   |
| `creature_spells`             | Creature Spells             |
| `creature_display_info_addon` | Creature Display Info Addon |
| `locales_creature`            | Creature Localization       |

### Quest Related

| Table Name             | Purpose              |
| ---------------------- | -------------------- |
| `quest_template`       | Quest Template       |
| `quest_cast_objective` | Quest Cast Objective |
| `quest_end_scripts`    | Quest End Scripts    |
| `quest_greeting`       | Quest Greeting       |
| `quest_start_scripts`  | Quest Start Scripts  |
| `locales_quest`        | Quest Localization   |

### GameObject Related

| Table Name                      | Purpose                     |
| ------------------------------- | --------------------------- |
| `gameobject_template`           | GameObject Template         |
| `gameobject`                    | GameObject Instances        |
| `gameobject_involvedrelation`   | GameObject Quest (Finisher) |
| `gameobject_loot_template`      | GameObject Loot Template    |
| `gameobject_questrelation`      | GameObject Quest (Starter)  |
| `gameobject_scripts`            | GameObject Scripts          |
| `gameobject_display_info_addon` | GameObject Display Info     |
| `locales_gameobject`            | GameObject Localization     |

### Spell Related

| Table Name              | Purpose               |
| ----------------------- | --------------------- |
| `spell_template`        | Spell Template        |
| `spell_affect`          | Spell Affects         |
| `spell_area`            | Area Spells           |
| `spell_chain`           | Spell Chain (Ranks)   |
| `spell_disabled`        | Disabled Spells       |
| `spell_effect_mod`      | Spell Effect Mods     |
| `spell_elixir`          | Elixir Types          |
| `spell_group`           | Spell Groups          |
| `spell_learn_spell`     | Learn Spell Relations |
| `spell_mod`             | Spell Mods            |
| `spell_proc_event`      | Spell Proc Events     |
| `spell_scripts`         | Spell Scripts         |
| `spell_target_position` | Spell Target Position |
| `locales_spell`         | Spell Localization    |

### Loot Templates

| Table Name                    | Purpose               |
| ----------------------------- | --------------------- |
| `creature_loot_template`      | Creature Loot         |
| `gameobject_loot_template`    | GameObject Loot       |
| `item_loot_template`          | Item Loot (Container) |
| `disenchant_loot_template`    | Disenchant Loot       |
| `fishing_loot_template`       | Fishing Loot          |
| `mail_loot_template`          | Mail Loot             |
| `pickpocketing_loot_template` | Pickpocket Loot       |
| `reference_loot_template`     | Reference Loot        |
| `skinning_loot_template`      | Skinning Loot         |

### NPC Interaction

| Table Name             | Purpose            |
| ---------------------- | ------------------ |
| `npc_gossip`           | NPC Gossip         |
| `npc_text`             | NPC Text           |
| `npc_trainer`          | NPC Trainer        |
| `npc_trainer_template` | Trainer Template   |
| `npc_vendor`           | NPC Vendor         |
| `npc_vendor_template`  | Vendor Template    |
| `gossip_menu`          | Gossip Menu        |
| `gossip_menu_option`   | Gossip Menu Option |
| `gossip_scripts`       | Gossip Scripts     |

### Map/Area Related

| Table Name               | Purpose           |
| ------------------------ | ----------------- |
| `area_template`          | Area Template     |
| `areatrigger_teleport`   | Teleport Triggers |
| `areatrigger_template`   | Trigger Template  |
| `areatrigger_tavern`     | Tavern Triggers   |
| `map_template`           | Map Template      |
| `game_graveyard_zone`    | Graveyards        |
| `world_safe_locs_facing` | Safe Locations    |
| `locales_area`           | Area Localization |

### Battlegrounds

| Table Name                | Purpose            |
| ------------------------- | ------------------ |
| `battleground_events`     | BG Events          |
| `battleground_template`   | BG Template        |
| `battlemaster_entry`      | Battlemaster Entry |
| `creature_battleground`   | BG Creatures       |
| `gameobject_battleground` | BG GameObjects     |

### Faction/Reputation

| Table Name                      | Purpose                |
| ------------------------------- | ---------------------- |
| `faction`                       | Faction                |
| `faction_template`              | Faction Template       |
| `reputation_reward_rate`        | Rep Reward Rate        |
| `reputation_spillover_template` | Rep Spillover Template |
| `locales_faction`               | Faction Localization   |

### Player Related

| Table Name                | Purpose             |
| ------------------------- | ------------------- |
| `player_classlevelstats`  | Class Level Stats   |
| `player_levelstats`       | Level Stats         |
| `player_xp_for_level`     | XP for Level        |
| `playercreateinfo`        | Create Info         |
| `playercreateinfo_action` | Create Actions      |
| `playercreateinfo_item`   | Create Items        |
| `playercreateinfo_spell`  | Create Spells       |
| `player_factionchange_*`  | Faction Change Data |

### Skill Related

| Table Name                 | Purpose            |
| -------------------------- | ------------------ |
| `skill_fishing_base_level` | Fishing Base Level |
| `skill_line_ability`       | Skill Line Ability |

### Pet Related

| Table Name            | Purpose           |
| --------------------- | ----------------- |
| `pet_levelstats`      | Pet Level Stats   |
| `pet_name_generation` | Pet Name Gen      |
| `pet_spell_data`      | Pet Spell Data    |
| `petcreateinfo_spell` | Pet Create Spells |
| `collection_pet`      | Pet Collection    |
| `collection_mount`    | Mount Collection  |

### Transport

| Table Name              | Purpose            |
| ----------------------- | ------------------ |
| `taxi_nodes`            | Flight Paths       |
| `taxi_path_transitions` | Flight Transitions |
| `transports`            | Transports (Ships) |
| `game_tele`             | Teleport Locations |
| `locales_taxi_node`     | Flight Path Loc    |

### Scripts

| Table Name             | Purpose            |
| ---------------------- | ------------------ |
| `event_scripts`        | Event Scripts      |
| `generic_scripts`      | Generic Scripts    |
| `script_texts`         | Script Texts       |
| `script_waypoint`      | Script Waypoints   |
| `scripted_areatrigger` | Scripted Triggers  |
| `scripted_event_id`    | Scripted Event IDs |

### Game Events

| Table Name                 | Purpose             |
| -------------------------- | ------------------- |
| `game_event_creature`      | Event Creatures     |
| `game_event_creature_data` | Event Creature Data |
| `game_event_gameobject`    | Event Objects       |
| `game_event_mail`          | Event Mail          |
| `game_event_quest`         | Event Quests        |

### Pools

| Table Name                 | Purpose            |
| -------------------------- | ------------------ |
| `pool_creature`            | Creature Pool      |
| `pool_creature_template`   | Creature Pool Tmpl |
| `pool_gameobject`          | GameObject Pool    |
| `pool_gameobject_template` | Object Pool Tmpl   |
| `pool_pool`                | Pool of Pools      |
| `pool_template`            | Pool Template      |

### Shop

| Table Name        | Purpose         |
| ----------------- | --------------- |
| `shop_categories` | Shop Categories |
| `shop_items`      | Shop Items      |

### Others

| Table Name           | Purpose           |
| -------------------- | ----------------- |
| `autobroadcast`      | Autobroadcast     |
| `broadcast_text`     | Broadcast Text    |
| `conditions`         | Conditions        |
| `exploration_basexp` | Exploration XP    |
| `game_weather`       | Weather           |
| `mangos_string`      | MaNGOS Strings    |
| `page_text`          | Page Text (Books) |
| `points_of_interest` | POI               |
| `reserved_name`      | Reserved Names    |
| `sound_entries`      | Sound Entries     |
| `variables`          | Server Variables  |
| `warden_checks`      | Warden Checks     |
| `warden_scans`       | Warden Scans      |

---

## ShellLab Database Architecture

This project uses an ETL (Extract-Transform-Load) pipeline to export data from the source (MySQL) to JSON, and then import it into the local SQLite database.

### ETL Process Status

| Module        | Source Table (MySQL)       | Intermediate File (JSON) | Target Table (SQLite) | Status  | Import Script             |
| ------------- | -------------------------- | ------------------------ | --------------------- | ------- | ------------------------- |
| **Items**     | `item_template`            | `item_template.json`     | `items`               | ✅ Done | `db_import/main.go`       |
| **Objects**   | `gameobject_template`      | `objects.json`           | `objects`             | ✅ Done | `export_objects_mysql.py` |
| **Locks**     | `aowow_lock`               | `locks.json`             | `locks`               | ✅ Done | `export_objects_mysql.py` |
| **Quests**    | `quest_template`           | `quests.json`            | `quests`              | ✅ Done | `export_quests.py`        |
| **Creatures** | `creature_template`        | `creatures.json`         | `creatures`           | ✅ Done | `export_creatures.py`     |
| **Factions**  | `aowow_factions`           | `factions.json`          | `factions`            | ✅ Done | `export_factions.py`      |
| **Spells**    | `spell_template`           | `spells.json`            | `spells`              | ✅ Done | `export_spells.py`        |
| **Loot**      | `creature_loot_template`   | `creature_loot.json`     | `creature_loot`       | ✅ Done | `export_loot.py`          |
|               | `reference_loot_template`  | `reference_loot.json`    | `reference_loot`      | ✅ Done | `export_loot.py`          |
|               | `gameobject_loot_template` | `gameobject_loot.json`   | `gameobject_loot`     | ✅ Done | `export_loot.py`          |
|               | `item_loot_template`       | `item_loot.json`         | `item_loot`           | ✅ Done | `export_loot.py`          |
|               | `disenchant_loot_template` | `disenchant_loot.json`   | `disenchant_loot`     | ✅ Done | `export_loot.py`          |

### SQLite Table Schema Definitions

Below are the main table structures for the local SQLite database (`data/shelllab.db`).

#### 1. Objects

```sql
CREATE TABLE objects (
    entry INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type INTEGER NOT NULL,
    display_id INTEGER,
    size REAL,
    data0-7 INTEGER  -- Contains critical data like LockID
);
```

#### 2. Locks

Used to derive special Object categories (e.g. Herb, Mine).

```sql
CREATE TABLE locks (
    id INTEGER PRIMARY KEY,
    type1-5 INTEGER,
    prop1-5 INTEGER, -- Key properties (1=Lockpicking, 2=Herbalism, 3=Mining)
    req1-5 INTEGER
);
```

#### 3. Quests

```sql
CREATE TABLE quests (
    entry INTEGER PRIMARY KEY,
    title TEXT,
    min_level INTEGER,
    quest_level INTEGER,
    ... -- Contains detailed text, rewards, objectives, etc.
);
```

#### 4. Creatures

```sql
CREATE TABLE creatures (
    entry INTEGER PRIMARY KEY,
    name TEXT,
    subname TEXT,
    level_min/max INTEGER,
    health_min/max INTEGER,
    creature_type INTEGER,
    creature_rank INTEGER,
    loot_id INTEGER,
    skin_loot_id INTEGER,
    pickpocket_loot_id INTEGER
    ...
);
```

#### 5. Loot Tables

All loot tables share the same structure:

```sql
CREATE TABLE *_loot (
    entry INTEGER,          -- Related ID (e.g. creature.loot_id)
    item INTEGER,           -- Item ID (FK: items.entry)
    chance REAL,            -- Drop rate
    groupid INTEGER,
    mincount_or_ref INTEGER, -- Positive=Count, Negative=Reference to other loot table
    maxcount INTEGER
);
```

#### 6. Spells

Used for displaying spell details and calculating set/item effects.

```sql
CREATE TABLE spells (
    entry INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    effect_base_points1-3 INTEGER,
    effect_die_sides1-3 INTEGER
);
```

#### 7. Factions

```sql
CREATE TABLE factions (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    side INTEGER,      -- Alliance/Horde/Neutral
    category_id INTEGER
);
```
