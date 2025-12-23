# Implementation Status

This document records the progress of ShellLab's local database (ETL workflow) and frontend features implementation.

## Database Import Status (ETL)

### Core Data Tables (MySQL -> JSON -> SQLite)

| Table         | Description    | Status  | Script                    |
| ------------- | -------------- | ------- | ------------------------- |
| **Items**     | Items          | ✅ Done | `db_import/main.go`       |
| **Icons**     | Item Icons     | ✅ Done | `import_icons/main.go`    |
| **Objects**   | Game Objects   | ✅ Done | `export_objects_mysql.py` |
| **Locks**     | Locks/Metadata | ✅ Done | `export_objects_mysql.py` |
| **Quests**    | Quests         | ✅ Done | `export_quests.py`        |
| **Creatures** | NPCs           | ✅ Done | `export_creatures.py`     |
| **Factions**  | Factions       | ✅ Done | `export_factions.py`      |
| **Spells**    | Spells         | ✅ Done | `export_spells.py`        |
| **Loot**      | Loot Tables    | ✅ Done | `export_loot.py`          |

### Pending Tables

| Table         | Description | Priority | Notes                       |
| ------------- | ----------- | -------- | --------------------------- |
| `item_set`    | Item Sets   | High     | Requires `aowow_itemset`    |
| `spell_icons` | Spell Icons | Medium   | Requires `aowow_spellicons` |
| `zones`       | Map Zones   | Medium   | Requires `aowow_zones`      |

---

## Feature Development Progress

### Backend API (Go)

- [x] **Items**: Get details, search, category suggestions
- [x] **Loot**: AtlasLoot hierarchy browsing, loot queries
- [x] **Objects**: Category browsing (lock-based), search
- [ ] **Quests**: Detail API, Search API
- [ ] **Creatures**: Detail API, Search API
- [ ] **Factions**: List API
- [ ] **Spells**: Detail API

### Frontend Pages (React)

- [x] **Loot Browser**: Complete AtlasLoot browsing interface
- [x] **Objects Browser**: Category browsing, details view (In Progress)
- [ ] **Quest Browser**: Quest list, detail page
- [ ] **Creature Browser**: NPC list, loot view
- [ ] **Faction Browser**: Faction list
- [ ] **Spell Browser**: Spell search

## Next Steps

1.  **Frontend**: Finalize Objects Browser UI details (icons, detailed info).
2.  **Frontend**: Develop Quests and Creatures browsers.
3.  **Backend**: Implement Item Sets import and API.
