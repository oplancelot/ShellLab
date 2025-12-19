# Item Data Import Pipeline

## 📋 Two-Stage Pipeline

### Stage 1: SQL → JSON

Extract item data from MySQL dump to JSON format.

```bash
go run scripts/db_import/main.go extract-sql
```

**Input:** `data/sql/tw_world.sql` (191MB MySQL dump)  
**Output:** `data/item_template.json` (60MB JSON)

### Stage 2: JSON → SQLite

Import JSON data into SQLite database with optional updates.

```bash
go run scripts/db_import/main.go import-items
```

**Input:**

- `data/item_template.json` (required - base data)
- `data/item_template_update.json` (optional - custom modifications)

**Output:** `data/shelllab.db` (items table populated)

## 🔄 Workflow

### Initial Setup (Full Pipeline)

```bash
# 1. Initialize database
go run scripts/db_import/main.go init

# 2. Extract from SQL (one-time, or when tw_world.sql updates)
go run scripts/db_import/main.go extract-sql

# 3. Import to database
go run scripts/db_import/main.go import-items

# 4. Import icon mappings
go run scripts/import_icons/main.go

# 5. Import AtlasLoot data
go run scripts/extract_atlasloot/main.go
```

### Quick Update (JSON → SQLite only)

If you only need to modify items without re-extracting from SQL:

```bash
# 1. Edit item_template.json or create item_template_update.json
nano data/item_template.json

# 2. Re-import
go run scripts/db_import/main.go import-items
```

## 📝 Custom Item Modifications

### Option 1: Direct JSON Edit

Edit `data/item_template.json` directly:

- Simple for small changes
- Changes will be overwritten if you re-run `extract-sql`

### Option 2: Update File (Recommended)

Create `data/item_template_update.json`:

- Preserves your modifications
- Automatically applied during import
- Won't be overwritten
- See `item_template_update.json.example` for format

```bash
# Copy example
cp data/item_template_update.json.example data/item_template_update.json

# Edit your modifications
nano data/item_template_update.json

# Import (automatically applies updates)
go run scripts/db_import/main.go import-items
```

### Update File Structure

```json
{
  "ENTRY_ID": {
    "entry": ENTRY_ID,
    "name": "Modified Item Name",
    "quality": 4,
    ...
  }
}
```

**Notes:**

- Use same entry ID to modify existing items
- Use new entry ID (e.g., 999999+) to add custom items
- Updates are applied with `INSERT OR REPLACE`, so existing items are overwritten

## 📊 Data Flow

```
tw_world.sql (191MB)
    ↓ [extract-sql]
item_template.json (60MB)
    ↓ [import-items]
    + item_template_update.json (optional)
    ↓
shelllab.db (items table)
```

## 🎯 Use Cases

### Scenario 1: Server Update

Official server updates tw_world.sql:

```bash
# Replace tw_world.sql with new version
# Re-extract and import
go run scripts/db_import/main.go extract-sql
go run scripts/db_import/main.go import-items
```

### Scenario 2: Custom Server

You maintain custom item modifications:

```bash
# One-time setup: extract base data
go run scripts/db_import/main.go extract-sql

# Ongoing: maintain updates in separate file
edit data/item_template_update.json
go run scripts/db_import/main.go import-items
```

### Scenario 3: Testing

Quickly test item changes:

```bash
# Edit JSON directly
edit data/item_template.json

# Re-import (fast)
go run scripts/db_import/main.go import-items
```

## 🔍 Verification

```bash
# Check database stats
go run scripts/db_import/main.go stats

# Output:
# === Database Statistics ===
# Items:      XXXXX
# Categories: XX
# DB Size:    X.XX MB
```

## ⚠️ Important Notes

1. **SQL Extraction is Slow**: The `extract-sql` command takes several minutes due to the 191MB file size. Only run when necessary.

2. **Updates are Additive**: `item_template_update.json` items will overwrite base items with the same entry ID.

3. **Preserve Custom Data**: Always use `item_template_update.json` for modifications you want to keep across SQL re-extractions.

4. **Backup Recommended**: Before major changes, backup your `shelllab.db` file.
