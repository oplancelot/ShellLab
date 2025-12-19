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
			"stackable":         parseInt(values[20]),
			"bonding":           parseInt(values[27]),
			"maxDurability":     parseInt(values[29]),
			"statType1":         parseInt(values[22]),
			"statValue1":        parseInt(values[23]),
			"statType2":         parseInt(values[24]),
			"statValue2":        parseInt(values[25]),
			"statType3":         parseInt(values[26]),
			"statValue3":        parseInt(values[27]),
			"statType4":         parseInt(values[28]),
			"statValue4":        parseInt(values[29]),
			"statType5":         parseInt(values[30]),
			"statValue5":        parseInt(values[31]),
			"statType6":         parseInt(values[32]),
			"statValue6":        parseInt(values[33]),
			"statType7":         parseInt(values[34]),
			"statValue7":        parseInt(values[35]),
			"statType8":         parseInt(values[36]),
			"statValue8":        parseInt(values[37]),
			"statType9":         parseInt(values[38]),
			"statValue9":        parseInt(values[39]),
			"statType10":        parseInt(values[40]),
			"statValue10":       parseInt(values[41]),
			"delay":             parseInt(values[42]),
			"dmgMin1":           parseFloat(values[44]),
			"dmgMax1":           parseFloat(values[45]),
			"dmgType1":          parseInt(values[46]),
			"dmgMin2":           parseFloat(values[47]),
			"dmgMax2":           parseFloat(values[48]),
			"dmgType2":          parseInt(values[49]),
			"armor":             parseInt(values[50]),
			"holyRes":           parseInt(values[51]),
			"fireRes":           parseInt(values[52]),
			"natureRes":         parseInt(values[53]),
			"frostRes":          parseInt(values[54]),
			"shadowRes":         parseInt(values[55]),
			"arcaneRes":         parseInt(values[56]),
			"spellId1":          parseInt(values[57]),
			"spellTrigger1":     parseInt(values[58]),
			"spellId2":          parseInt(values[60]),
			"spellTrigger2":     parseInt(values[61]),
			"spellId3":          parseInt(values[63]),
			"spellTrigger3":     parseInt(values[64]),
			"setId":             parseInt(values[75]),
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

	fmt.Println("=== Database Statistics ===")
	fmt.Printf("Items:      %d\n", itemCount)
	fmt.Printf("Categories: %d\n", catCount)

	// Show file size
	if info, err := os.Stat(dbPath); err == nil {
		fmt.Printf("DB Size:    %.2f MB\n", float64(info.Size())/(1024*1024))
	}
}
