# ShellLab - World of Warcraft Database Browser

A comprehensive desktop application for browsing and exploring World of Warcraft (Turtle WoW) game data, built with Wails, Go, and React.

## Features

### Database Browser

- **Items**: Complete item database with detailed statistics
  - Search by name, class, subclass, and inventory slot
  - WoW-style tooltips with complete item information
  - Icon display with local cache and CDN fallback
- **AtlasLoot Integration**: Complete loot table browser

  - 7 categories: Instances, Sets, Factions, PvP, World Bosses, World Events, Crafting
  - Hierarchical navigation (Category → Module → Table → Items)
  - Drop chance information where available

- **Creatures**: Browse creature database

  - Search by name and type
  - Paginated results for performance
  - View creature loot tables

- **Quests**: Explore quest database

  - Browse by zone or quest category
  - View quest details and objectives

- **Spells**: Search spell database

  - Browse by class and skill category
  - View spell effects and icons

- **Game Objects**: Browse object database

  - Search by name and type
  - View object loot tables

- **Factions**: View faction database
  - Reputation and faction rewards

## Architecture

### Technology Stack

- **Backend**: Go 1.24 + Wails v2.11
- **Frontend**: React 18 + TypeScript + Vite
- **Database**: SQLite 3
- **Styling**: Tailwind CSS with custom WoW theme

### Data Pipeline

```
MySQL (Turtle WoW DB)
  ↓ Python export scripts
JSON Files (data/*.json)
  ↓ Go importers (auto-generated)
SQLite Database (shelllab.db)
  ↓ Go repositories
React Frontend
```

**Key Components**:

1. **Export Scripts** (`scripts/export_*.py`): Extract data from MySQL to JSON
2. **Code Generator** (`scripts/generate_go_code.py`): Auto-generate Go importers from MySQL schema
3. **Importers** (`backend/database/importers/`): Load JSON data into SQLite on first run
4. **Repositories** (`backend/database/repositories/`): Query layer for frontend
5. **Icon Service** (`backend/services/`): Downloads and caches item icons

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+
- Wails v2.11+

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/ShellLab.git
cd ShellLab

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run in development mode
wails dev
```

### First Run

On first startup, the application will:

1. Create SQLite database schema
2. Import all JSON data files from `data/` directory
3. Download missing item icons
4. This process takes ~1-2 minutes

## Development

### Database Schema

The application uses a SQLite database with 30+ tables:

**Core Tables**:

- `item_template`: Items (1:1 MySQL mapping)
- `creature_template`: Creatures
- `quest_template`: Quests
- `spell_template`: Spells
- `gameobject_template`: Objects

**AtlasLoot Tables**:

- `atlasloot_categories`: Categories
- `atlasloot_modules`: Modules
- `atlasloot_tables`: Loot tables
- `atlasloot_items`: Loot entries

**Loot Tables**:

- `creature_loot_template`
- `item_loot_template`
- `gameobject_loot_template`
- `reference_loot_template`
- `disenchant_loot_template`

### Data Update Workflow

1. **Export from MySQL** (when source data changes):

```bash
cd scripts
python export_all_data.py  # Exports all tables to JSON
```

2. **Regenerate Go Code** (when schema changes):

```bash
python generate_go_code.py  # Creates importers, models, schemas
```

3. **Rebuild Database**:

```bash
rm data/shelllab.db  # Remove old database
wails dev            # Reimport on startup
```

### Icon Management

Icons are automatically downloaded from:

1. **Wowhead CDN** (`wow.zamimg.com`) - Primary source
2. **Turtle WoW Database** (`database.turtlecraft.gg`) - Fallback
3. **Trinity AoWoW** (`aowow.trinitycore.info`) - Fallback

Icons are cached in `frontend/public/items/icons/` for offline use.

## Data Sources

- **Turtle-WoW Emulation Server Source Code**:https://github.com/brian8544/turtle-wow

## Key Technologies

- **Wails**: Go-powered desktop apps with web UI
- **SQLite**: Embedded database (no server needed)
- **Code Generation**: Python scripts auto-generate Go code
- **React Hooks**: Modern state management
- **Tailwind CSS**: Utility-first styling

## Future Enhancements

- Talent tree browser and calculator
- Equipment set manager
- Stat calculator and comparison
- DPS simulator
- Enchant and gem browser
- Character planner
- Export/import functionality

## Contributing

This project is for educational purposes and community use. Contributions welcome!

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

**Built with ❤️ for the Turtle WoW Community**
