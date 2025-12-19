package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"shelllab/backend/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: db_import <command>")
		fmt.Println("\nCommands:")
		fmt.Println("  init         - Initialize the database schema")
		fmt.Println("  extract-sql  - Extract item_template.json from tw_world.sql")
		fmt.Println("  import-items - Import items from item_template.json (+ updates)")
		fmt.Println("  stats        - Show database statistics")
		fmt.Println("\nTwo-Stage Pipeline:")
		fmt.Println("  Stage 1: SQL → JSON")
		fmt.Println("    extract-sql converts tw_world.sql to item_template.json")
		fmt.Println("  Stage 2: JSON → SQLite")
		fmt.Println("    import-items loads item_template.json into database")
		fmt.Println("    Automatically applies item_template_update.json if exists")
		os.Exit(1)
	}

	// Get paths
	dataDir := "data"
	dbPath := filepath.Join(dataDir, "shelllab.db")
	jsonPath := filepath.Join(dataDir, "item_template.json")
	updateJsonPath := filepath.Join(dataDir, "item_template_update.json")
	sqlPath := filepath.Join(dataDir, "sql", "tw_world.sql")

	command := os.Args[1]

	switch command {
	case "init":
		initDatabase(dbPath)
	case "extract-sql":
		extractFromSQL(sqlPath, jsonPath)
	case "import-items":
		importItems(dbPath, jsonPath, updateJsonPath)
	case "import-spells":
		importSpells(dbPath, "data/sql/aowow.sql")
	case "import-itemsets":
		importItemSets(dbPath, "data/sql/aowow.sql")
	case "stats":
		showStats(dbPath)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func initDatabase(dbPath string) {
	fmt.Printf("Initializing database at %s...\n", dbPath)

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		fmt.Printf("ERROR: Failed to initialize schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Database schema initialized successfully!")
}

func extractFromSQL(sqlPath, jsonPath string) {
	fmt.Println("===== Extracting item_template from tw_world.sql =====")
	fmt.Printf("Reading: %s\n", sqlPath)
	fmt.Printf("Output:  %s\n\n", jsonPath)

	// Check if SQL file exists
	if _, err := os.Stat(sqlPath); os.IsNotExist(err) {
		fmt.Printf("ERROR: SQL file not found: %s\n", sqlPath)
		os.Exit(1)
	}

	file, err := os.Open(sqlPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open SQL file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 100*1024*1024)

	items := make(map[string]map[string]interface{})

	insertRegex := regexp.MustCompile(`^INSERT INTO \x60?item_template\x60?`)
	lineCount := 0
	inInsert := false
	var valueBuffer strings.Builder

	fmt.Println("Parsing SQL file (this may take a while)...")

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if insertRegex.MatchString(line) {
			inInsert = true
			// Find VALUES section
			if idx := strings.Index(line, "VALUES"); idx != -1 {
				valueBuffer.WriteString(line[idx+6:])
			}
			continue
		}

		if inInsert {
			valueBuffer.WriteString(" ")
			valueBuffer.WriteString(line)

			// Check if finished
			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				// Parse all values
				parseInsertValues(valueBuffer.String(), items)
				valueBuffer.Reset()
				inInsert = false
			}
		}

		if lineCount%500000 == 0 {
			fmt.Printf("Processed %d lines, found %d items...\n", lineCount, len(items))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR: Failed to read SQL file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Extracted %d items\n", len(items))

	// 导出为JSON
	fmt.Println("Writing JSON file...")
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Printf("ERROR: Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		fmt.Printf("ERROR: Failed to write JSON file: %v\n", err)
		os.Exit(1)
	}

	fileInfo, _ := os.Stat(jsonPath)
	sizeMB := float64(fileInfo.Size()) / 1024 / 1024

	fmt.Printf("✓ Exported to %s (%.2f MB)\n", jsonPath, sizeMB)
	fmt.Println("\n===== Complete =====")
	fmt.Println("Next step: go run scripts/db_import/main.go import-items")
}

func parseInsertValues(valueStr string, items map[string]map[string]interface{}) {
	// 正则匹配每一行的值: (value1, value2, ...)
	rowRegex := regexp.MustCompile(`\(([^)]*(?:\([^)]*\)[^)]*)*)\)`)
	rows := rowRegex.FindAllStringSubmatch(valueStr, -1)

	for _, row := range rows {
		if len(row) < 2 {
			continue
		}

		values := parseSQLRow(row[1])
		if len(values) < 80 { // item_template should have many fields
			continue
		}

		entry := getStr(values, 0)
		if entry == "" || entry == "0" {
			continue
		}

		item := map[string]interface{}{
			"entry":             parseInt(values[0]),
			"class":             parseInt(values[1]),
			"subclass":          parseInt(values[2]),
			"name":              cleanStr(values[3]),
			"description":       cleanStr(values[4]),
			"displayId":         parseInt(values[5]),
			"quality":           parseInt(values[6]),
			"flags":             parseInt(values[7]),
			"buyCount":          parseInt(values[8]),
			"buyPrice":          parseInt(values[9]),
			"sellPrice":         parseInt(values[10]),
			"inventoryType":     parseInt(values[11]),
			"allowableClass":    parseInt(values[12]),
			"allowableRace":     parseInt(values[13]),
			"itemLevel":         parseInt(values[14]),
			"requiredLevel":     parseInt(values[15]),
			"requiredSkill":     parseInt(values[16]),
			"requiredSkillRank": parseInt(values[17]),
			"requiredSpell":     parseInt(values[18]),
			"stackable":         parseInt(values[24]),

			// Stats (26-45)
			"statType1":   parseInt(values[26]),
			"statValue1":  parseInt(values[27]),
			"statType2":   parseInt(values[28]),
			"statValue2":  parseInt(values[29]),
			"statType3":   parseInt(values[30]),
			"statValue3":  parseInt(values[31]),
			"statType4":   parseInt(values[32]),
			"statValue4":  parseInt(values[33]),
			"statType5":   parseInt(values[34]),
			"statValue5":  parseInt(values[35]),
			"statType6":   parseInt(values[36]),
			"statValue6":  parseInt(values[37]),
			"statType7":   parseInt(values[38]),
			"statValue7":  parseInt(values[39]),
			"statType8":   parseInt(values[40]),
			"statValue8":  parseInt(values[41]),
			"statType9":   parseInt(values[42]),
			"statValue9":  parseInt(values[43]),
			"statType10":  parseInt(values[44]),
			"statValue10": parseInt(values[45]),

			// Damage & Speed
			"delay":    parseInt(values[46]),
			"dmgMin1":  parseFloat(values[49]),
			"dmgMax1":  parseFloat(values[50]),
			"dmgType1": parseInt(values[51]),
			"dmgMin2":  parseFloat(values[52]),
			"dmgMax2":  parseFloat(values[53]),
			"dmgType2": parseInt(values[54]),

			// Defenses
			"armor":     parseInt(values[65]),
			"holyRes":   parseInt(values[66]),
			"fireRes":   parseInt(values[67]),
			"natureRes": parseInt(values[68]),
			"frostRes":  parseInt(values[69]),
			"shadowRes": parseInt(values[70]),
			"arcaneRes": parseInt(values[71]),

			// Spells
			"spellId1":      parseInt(values[72]),
			"spellTrigger1": parseInt(values[73]),
			"spellId2":      parseInt(values[79]),
			"spellTrigger2": parseInt(values[80]),
			"spellId3":      parseInt(values[86]),
			"spellTrigger3": parseInt(values[87]),

			// Other
			"bonding":       parseInt(values[107]),
			"setId":         parseInt(values[116]),
			"maxDurability": parseInt(values[117]),
		}

		items[entry] = item
	}
}

func parseSQLRow(row string) []string {
	var result []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)
	escaped := false

	for _, ch := range row {
		if escaped {
			current.WriteRune(ch)
			escaped = false
			continue
		}

		if ch == '\\' && inQuote {
			escaped = true
			current.WriteRune(ch)
			continue
		}

		if (ch == '\'' || ch == '"') && !inQuote {
			inQuote = true
			quoteChar = ch
			current.WriteRune(ch)
			continue
		}

		if ch == quoteChar && inQuote {
			inQuote = false
			quoteChar = 0
			current.WriteRune(ch)
			continue
		}

		if ch == ',' && !inQuote {
			result = append(result, strings.TrimSpace(current.String()))
			current.Reset()
			continue
		}

		current.WriteRune(ch)
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

func getStr(values []string, index int) string {
	if index < len(values) {
		return values[index]
	}
	return ""
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" || s == "NULL" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "NULL" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func cleanStr(s string) string {
	s = strings.TrimSpace(s)
	// Remove quotes
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			s = s[1 : len(s)-1]
		}
	}
	// Handle escapes
	s = strings.ReplaceAll(s, "\\'", "'")
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}

func importItems(dbPath, jsonPath, updateJsonPath string) {
	fmt.Println("===== Stage 2: JSON → SQLite =====")
	fmt.Printf("Importing items from %s...\n", jsonPath)

	// Check if JSON file exists
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		fmt.Printf("ERROR: JSON file not found: %s\n", jsonPath)
		fmt.Println("\nℹ Run first: go run scripts/db_import/main.go extract-sql")
		os.Exit(1)
	}

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Ensure schema exists
	if err := db.InitSchema(); err != nil {
		fmt.Printf("ERROR: Failed to initialize schema: %v\n", err)
		os.Exit(1)
	}

	// Import base items
	repo := database.NewItemRepository(db)
	count, err := repo.ImportFromJSON(jsonPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to import items: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Imported %d items from base template\n", count)

	// Check for update file
	if _, err := os.Stat(updateJsonPath); err == nil {
		fmt.Printf("\nApplying updates from %s...\n", updateJsonPath)

		updateCount, err := repo.ImportFromJSON(updateJsonPath)
		if err != nil {
			fmt.Printf("WARNING: Failed to apply updates: %v\n", err)
		} else {
			fmt.Printf("✓ Applied %d item updates/additions\n", updateCount)
			count = updateCount // Update total count
		}
	} else {
		fmt.Println("\nℹ No update file found (item_template_update.json)")
		fmt.Println("  Create this file to add custom item modifications")
	}

	fmt.Printf("\n✅ Total items in database: %d\n", count)
	fmt.Println("\n===== Import Complete =====")
	fmt.Println("Next steps:")
	fmt.Println("  1. go run scripts/import_icons/main.go")
	fmt.Println("  2. go run scripts/extract_atlasloot/main.go")
}

func showStats(dbPath string) {
	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	itemRepo := database.NewItemRepository(db)
	catRepo := database.NewCategoryRepository(db)

	itemCount, _ := itemRepo.GetItemCount()
	catCount, _ := catRepo.GetCategoryCount()

	// Count spells
	var spellCount int
	db.DB().QueryRow("SELECT COUNT(*) FROM spells").Scan(&spellCount)

	fmt.Println("=== Database Statistics ===")
	fmt.Printf("Items:      %d\n", itemCount)
	fmt.Printf("Categories: %d\n", catCount)
	fmt.Printf("Spells:     %d\n", spellCount)

	// Show file size
	if info, err := os.Stat(dbPath); err == nil {
		fmt.Printf("DB Size:    %.2f MB\n", float64(info.Size())/(1024*1024))
	}
}

func importSpells(dbPath, sqlPath string) {
	fmt.Println("===== Importing Spells from aowow.sql =====")

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create spells table
	_, err = db.DB().Exec(`
		DROP TABLE IF EXISTS spells;
		CREATE TABLE IF NOT EXISTS spells (
			id INTEGER PRIMARY KEY,
			name TEXT,
			description TEXT,
			icon_id INTEGER DEFAULT 0,
			base_points1 INTEGER DEFAULT 0,
			base_points2 INTEGER DEFAULT 0,
			base_points3 INTEGER DEFAULT 0
		);
		CREATE INDEX IF NOT EXISTS idx_spells_name ON spells(name);
	`)
	if err != nil {
		fmt.Printf("ERROR: Failed to create table: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Open(sqlPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open SQL file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 100*1024*1024)

	insertRegex := regexp.MustCompile(`^INSERT INTO \x60?aowow_spell\x60?`)
	rowRegex := regexp.MustCompile(`\(([^)]*(?:\([^)]*\)[^)]*)*)\)`)

	tx, _ := db.DB().Begin()
	stmt, _ := tx.Prepare("INSERT OR REPLACE INTO spells (id, name, description, base_points1, base_points2, base_points3) VALUES (?, ?, ?, ?, ?, ?)")
	defer stmt.Close()

	lineCount := 0
	spellCount := 0
	inInsert := false
	var valueBuffer strings.Builder

	fmt.Println("Parsing SQL file...")

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if insertRegex.MatchString(line) {
			inInsert = true
			if idx := strings.Index(line, "VALUES"); idx != -1 {
				valueBuffer.WriteString(line[idx+6:])
			}
			continue
		}

		if inInsert {
			valueBuffer.WriteString(" ")
			valueBuffer.WriteString(line)

			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				// Parse rows
				rows := rowRegex.FindAllStringSubmatch(valueBuffer.String(), -1)
				for _, row := range rows {
					if len(row) < 2 {
						continue
					}

					values := parseSQLRow(row[1])
					// We expect at least ~70 columns.
					// Strategy: ID is first. Description is 3rd string from end.
					// Name is 1st string from end of strings (4th from end total?).
					// Let's rely on finding string fields.
					// Based on sample: ..., 'Name', 'Rank', 'Desc', 'Tooltip', 0, 0, 0, 1, 1, 1)
					// So 6 ints at end.
					// Then Tooltip, Desc, Rank, Name.

					if len(values) < 20 {
						continue
					}

					id := parseInt(values[0])

					// Find the string block at the end
					count := len(values)
					// Assuming last 6 are numbers, so index count-7 is ToolTip
					// count-8 is Description
					// count-9 is Rank
					// count-10 is Name

					// Verify they are strings (quoted in SQL, but cleaned by parseSQLRow)
					// We can just grab them by index logic if layout is consistent.
					// aowow_spell has fixed columns.

					nameIdx := count - 10
					descIdx := count - 8

					if nameIdx > 0 && descIdx > 0 && descIdx < count {
						name := cleanStr(values[nameIdx])
						desc := cleanStr(values[descIdx])

						bp1 := 0
						bp2 := 0
						bp3 := 0
						if len(values) > 40 {
							bp1 = parseInt(values[37])
							bp2 = parseInt(values[38])
							bp3 = parseInt(values[39])
						}

						stmt.Exec(id, name, desc, bp1, bp2, bp3)
						spellCount++
					}
				}

				valueBuffer.Reset()
				inInsert = false
			}
		}

		if lineCount%50000 == 0 {
			fmt.Printf("Processed %d lines, imported %d spells...\r", lineCount, spellCount)
		}
	}

	tx.Commit()
	fmt.Printf("\n✓ Imported %d spells\n", spellCount)
}

func importItemSets(dbPath, sqlPath string) {
	fmt.Println("===== Importing Item Sets from aowow.sql =====")

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create item_sets table
	_, err = db.DB().Exec(`
		DROP TABLE IF EXISTS item_sets;
		CREATE TABLE IF NOT EXISTS item_sets (
			id INTEGER PRIMARY KEY,
			name TEXT,
			item_ids TEXT, -- JSON array
			bonuses TEXT   -- JSON array of {threshold, spell_id}
		);
	`)
	if err != nil {
		fmt.Printf("ERROR: Failed to create table: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Open(sqlPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open SQL file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 100*1024*1024)

	insertRegex := regexp.MustCompile(`^INSERT INTO \x60?aowow_itemset\x60?`)
	rowRegex := regexp.MustCompile(`\(([^)]*(?:\([^)]*\)[^)]*)*)\)`)

	tx, _ := db.DB().Begin()
	stmt, _ := tx.Prepare("INSERT OR REPLACE INTO item_sets (id, name, item_ids, bonuses) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	lineCount := 0
	setCount := 0
	inInsert := false
	var valueBuffer strings.Builder

	fmt.Println("Parsing SQL file...")

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if insertRegex.MatchString(line) {
			inInsert = true
			if idx := strings.Index(line, "VALUES"); idx != -1 {
				valueBuffer.WriteString(line[idx+6:])
			}
			continue
		}

		if inInsert {
			valueBuffer.WriteString(" ")
			valueBuffer.WriteString(line)

			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				rows := rowRegex.FindAllStringSubmatch(valueBuffer.String(), -1)
				for _, row := range rows {
					if len(row) < 2 {
						continue
					}

					values := parseSQLRow(row[1])
					// Structure: id, name, item1..10, spell1..8, bonus1..8, skillID, skillLevel
					// Total columns: 1 + 1 + 10 + 8 + 8 + 1 + 1 = 30
					if len(values) < 28 {
						continue
					}

					id := parseInt(values[0])
					name := cleanStr(values[1])

					// Collect Items
					var itemIDs []int
					for i := 2; i <= 11; i++ {
						iid := parseInt(values[i])
						if iid > 0 {
							itemIDs = append(itemIDs, iid)
						}
					}

					// Collect Bonuses
					// Spells are at 12..19 (8 spells)
					// Bonuses are at 20..27 (8 thresholds)
					type SetBonus struct {
						Threshold int `json:"threshold"`
						SpellID   int `json:"spellId"`
					}
					var bonuses []SetBonus
					for i := 0; i < 8; i++ {
						spellID := parseInt(values[12+i])
						threshold := parseInt(values[20+i])
						if spellID > 0 && threshold > 0 {
							bonuses = append(bonuses, SetBonus{Threshold: threshold, SpellID: spellID})
						}
					}

					itemsJson, _ := json.Marshal(itemIDs)
					bonusesJson, _ := json.Marshal(bonuses)

					stmt.Exec(id, name, string(itemsJson), string(bonusesJson))
					setCount++
				}

				valueBuffer.Reset()
				inInsert = false
			}
		}

		if lineCount%50000 == 0 {
			fmt.Printf("Processed %d lines, imported %d sets...\r", lineCount, setCount)
		}
	}

	tx.Commit()
	fmt.Printf("\n✓ Imported %d item sets\n", setCount)
}
