import json
import os
import re

DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')
BACKEND_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'database')

TYPE_MAP = {
    'tinyint': 'int',
    'smallint': 'int',
    'mediumint': 'int',
    'int': 'int',
    'bigint': 'int64',
    'float': 'float64',
    'double': 'float64',
    'decimal': 'float64',
    'varchar': 'string',
    'text': 'string',
    'longtext': 'string',
    'char': 'string',
    'date': 'string',
    'datetime': 'string',
    'timestamp': 'string'
}

def go_type(mysql_type):
    mysql_type_lower = mysql_type.lower()
    base_type = mysql_type_lower.split('(')[0].replace(' unsigned', '')
    
    if 'int' in base_type:  # tinyint, smallint, mediumint, int, bigint
        if 'bigint' in base_type:
            return 'int64'
        return 'int'
    elif base_type in ('float', 'double', 'decimal'):
        return 'float64'
    else:
        return 'string'

def to_pascal_case(snake_str):
    return ''.join(x.title() for x in snake_str.split('_'))

def generate_code():
    with open(os.path.join(DATA_DIR, 'schema_dump.json'), 'r') as f:
        schemas = json.load(f)

    # 1. Models
    models_content = "package models\n\n// Generated 1:1 MySQL models\n\n"
    for table, columns in schemas.items():
        struct_name = to_pascal_case(table) + "Full"  # Add suffix to avoid conflicts
        models_content += f"type {struct_name} struct {{\n"
        for col in columns:
            field_name = to_pascal_case(col['Field'])
            gotype = go_type(col['Type'])
            models_content += f"    {field_name} {gotype} `json:\"{col['Field']}\"`\n"
        models_content += "}\n\n"

    with open(os.path.join(BACKEND_DIR, 'models', 'generated_models.go'), 'w') as f:
        f.write(models_content)

    # 2. Schema
    schema_content = "package schema\n\nfunc GeneratedSchema() string {\n    return `\n"
    for table, columns in schemas.items():
        schema_content += f"    CREATE TABLE IF NOT EXISTS {table} (\n"
        pk = ""
        cols_def = []
        for col in columns:
            cname = col['Field']
            ctype = col['Type'].upper()
            # Simplify types for SQLite
            if 'INT' in ctype: sql_type = 'INTEGER'
            elif 'FLOAT' in ctype or 'DOUBLE' in ctype or 'DECIMAL' in ctype: sql_type = 'REAL'
            else: sql_type = 'TEXT'
            
            def_val = ""
            if col['Default'] is not None:
                d = col['Default']
                if sql_type == 'TEXT':
                    def_val = f" DEFAULT '{d}'"
                else:
                    def_val = f" DEFAULT {d}"
            
            cols_def.append(f"        {cname} {sql_type}{def_val}")
            
            if col['Key'] == 'PRI':
                pk = cname

        schema_content += ",\n".join(cols_def)
        
        # Add extra columns for item_template (icon_path from aowow)
        if table == 'item_template':
            schema_content += ",\n        icon_path TEXT DEFAULT ''"
        
        # Add extra columns for spell_template (iconName for spell icons)
        if table == 'spell_template':
            schema_content += ",\n        iconName TEXT DEFAULT ''"
        
        if pk:
            schema_content += f",\n        PRIMARY KEY ({pk})"
        schema_content += "\n    );\n\n"
    schema_content += "    `\n}\n"

    with open(os.path.join(BACKEND_DIR, 'schema', 'generated_schema.go'), 'w', encoding='utf-8') as f:
        f.write(schema_content)

    # 3. Importers
    importer_content = """package importers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "os"
    "shelllab/backend/database/models"
)

type GeneratedImporter struct {
    db *sql.DB
}

func NewGeneratedImporter(db *sql.DB) *GeneratedImporter {
    return &GeneratedImporter{db: db}
}

"""
    for table, columns in schemas.items():
        struct_name = to_pascal_case(table) + "Full"  # Match model name with suffix
        func_name = f"Import{to_pascal_case(table)}"
        col_names = [c['Field'] for c in columns]
        placeholders = ",".join(["?" for _ in col_names])
        cols_str = ",".join(col_names)
        
        # Build field accessors
        field_accessors = []
        for col in columns:
            fname = to_pascal_case(col['Field'])
            field_accessors.append(f"item.{fname}")
        
        importer_content += f"""
func (i *GeneratedImporter) {func_name}(jsonPath string) error {{
    file, err := os.Open(jsonPath)
    if err != nil {{
        return err
    }}
    defer file.Close()

    decoder := json.NewDecoder(file)
    // Expect array
    if _, err := decoder.Token(); err != nil {{
        return err
    }}

    tx, err := i.db.Begin()
    if err != nil {{
        return err
    }}
    defer tx.Rollback()

    query := `INSERT OR REPLACE INTO {table} ({cols_str}) VALUES ({placeholders})`
    stmt, err := tx.Prepare(query)
    if err != nil {{
        return err
    }}
    defer stmt.Close()

    count := 0
    for decoder.More() {{
        var item models.{struct_name}
        if err := decoder.Decode(&item); err != nil {{
            continue
        }}

        _, err = stmt.Exec({", ".join(field_accessors)})
        if err != nil {{
            // fmt.Printf("Error importing %s: %v\\n", "{table}", err)
            continue
        }}
        count++
    }}

    return tx.Commit()
}}
"""

    with open(os.path.join(BACKEND_DIR, 'importers', 'generated_importers.go'), 'w', encoding='utf-8') as f:
        f.write(importer_content)
        f.write("""
// ImportItemIcons loads icon paths from item_icons.json and updates item_template
func (i *GeneratedImporter) ImportItemIcons(jsonPath string) error {
	fmt.Printf("  -> Reading item icons from %s...\\n", jsonPath)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil // Icons are optional
	}

	var iconMap map[string]string
	if err := json.Unmarshal(data, &iconMap); err != nil {
		fmt.Printf("  ERROR parsing item_icons.json: %v\\n", err)
		return nil
	}

	fmt.Printf("  -> Updating database with %d icon mappings...\\n", len(iconMap))
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE item_template SET icon_path = ? WHERE display_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for displayIDStr, iconName := range iconMap {
		var displayID int
		fmt.Sscanf(displayIDStr, "%d", &displayID)
		if displayID > 0 {
			res, err := stmt.Exec(iconName, displayID)
			if err != nil {
				continue
			}
			if rows, _ := res.RowsAffected(); rows > 0 {
				count++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	fmt.Printf("  ✓ Successfully updated %d items with icons\\n", count)
	return nil
}

// SpellEnhanced represents a spell record from spells_enhanced.json
type SpellEnhanced struct {
	SpellIconId int    `json:"spellIconId"`
	IconName    string `json:"iconName"`
}

// ImportSpellIcons loads spell icons from spells_enhanced.json and updates spell_template
func (i *GeneratedImporter) ImportSpellIcons(jsonPath string) error {
	fmt.Printf("  -> Reading spell icons from %s...\\n", jsonPath)
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil // Optional
	}
	defer file.Close()

	var spells []SpellEnhanced
	if err := json.NewDecoder(file).Decode(&spells); err != nil {
		fmt.Printf("  ERROR parsing spells_enhanced.json: %v\\n", err)
		return nil
	}

	// Build unique icon map
	iconMap := make(map[int]string)
	for _, s := range spells {
		if s.SpellIconId > 0 && s.IconName != "" && s.IconName != "temp" {
			iconMap[s.SpellIconId] = s.IconName
		}
	}

	fmt.Printf("  -> Updating database with %d spell icon mappings...\\n", len(iconMap))
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE spell_template SET iconName = ? WHERE spellIconId = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for iconId, iconName := range iconMap {
		res, err := stmt.Exec(iconName, iconId)
		if err != nil {
			continue
		}
		if rows, _ := res.RowsAffected(); rows > 0 {
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	fmt.Printf("  ✓ Successfully updated %d spells with icons\\n", count)
	return nil
}
""")

    print("Generated models, schema, and importers.")

if __name__ == '__main__':
    generate_code()
