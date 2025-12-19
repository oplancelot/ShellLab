package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"shelllab/backend/database"
	"strings"
)

func main() {
	fmt.Println("=== AtlasLoot Full Data Extractor (All Categories) ===")
	fmt.Println()

	// Open database
	db, err := database.NewSQLiteDB("data/shelllab.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Ensure schema exists
	if err := db.InitSchema(); err != nil {
		log.Fatal("Failed to initialize schema:", err)
	}

	repo := database.NewAtlasLootRepository(db)

	// Clear existing data
	fmt.Println("Clearing existing AtlasLoot data...")
	if err := repo.ClearAllData(); err != nil {
		log.Fatal("Failed to clear data:", err)
	}

	// Parse TableRegister.lua for display names
	registerPath := "addons/AtlasLoot/Database/TableRegister.lua"
	fmt.Printf("\nParsing %s for display names...\n", registerPath)
	tableDisplayNames, err := parseTableRegister(registerPath)
	if err != nil {
		log.Fatal("Failed to parse TableRegister.lua:", err)
	}
	fmt.Printf("Found %d table name mappings\n\n", len(tableDisplayNames))

	// Define categories and their corresponding Lua files
	categories := []struct {
		Key         string
		DisplayName string
		FilePath    string
		SortOrder   int
	}{
		{"AtlasLootInstances", "Instances", "addons/AtlasLoot/Database/Instances.lua", 1},
		{"AtlasLootSets", "Sets", "addons/AtlasLoot/Database/Sets.lua", 2},
		{"AtlasLootFactions", "Factions", "addons/AtlasLoot/Database/Factions.lua", 3},
		{"AtlasLootPvP", "PvP", "addons/AtlasLoot/Database/PvP.lua", 4},
		{"AtlasLootWorldBosses", "World Bosses", "addons/AtlasLoot/Database/WorldBosses.lua", 5},
		{"AtlasLootWorldEvents", "World Events", "addons/AtlasLoot/Database/WorldEvents.lua", 6},
		{"AtlasLootCrafting", "Crafting", "addons/AtlasLoot/Database/Crafting.lua", 7},
	}

	totalItems := 0

	// Process each category
	for _, cat := range categories {
		fmt.Printf("=== Processing Category: %s ===\n", cat.DisplayName)

		// Check if file exists
		if _, err := os.Stat(cat.FilePath); os.IsNotExist(err) {
			fmt.Printf("⚠ File not found: %s, skipping...\n\n", cat.FilePath)
			continue
		}

		// Create category
		catID, err := repo.InsertCategory(cat.Key, cat.DisplayName, cat.SortOrder)
		if err != nil {
			log.Printf("Warning: Failed to create category %s: %v\n", cat.DisplayName, err)
			continue
		}
		fmt.Printf("✓ Created category: %s (ID: %d)\n", cat.DisplayName, catID)

		// Parse the Lua file
		fmt.Printf("Parsing %s...\n", cat.FilePath)
		tables, err := parseLuaFile(cat.FilePath)
		if err != nil {
			log.Printf("Warning: Failed to parse %s: %v\n", cat.FilePath, err)
			continue
		}
		fmt.Printf("Found %d loot tables\n", len(tables))

		// Group tables by module
		var moduleGroups map[string][]string
		if cat.Key == "AtlasLootInstances" {
			moduleGroups = groupTablesByInstance(tables)
		} else {
			// For non-instance categories, create a single module
			moduleGroups = map[string][]string{cat.DisplayName: {}}
			for tableName := range tables {
				moduleGroups[cat.DisplayName] = append(moduleGroups[cat.DisplayName], tableName)
			}
		}

		// Process each module
		for moduleName, tableNames := range moduleGroups {
			modID, err := repo.InsertModule(int(catID), moduleName, moduleName, 0)
			if err != nil {
				log.Printf("  Warning: Failed to create module %s: %v\n", moduleName, err)
				continue
			}
			fmt.Printf("  ✓ Module: %s\n", moduleName)

			// Process each table
			for _, tableName := range tableNames {
				tableData := tables[tableName]
				displayName := getDisplayName(tableName, tableDisplayNames)

				tblID, err := repo.InsertTable(int(modID), tableName, displayName, 0)
				if err != nil {
					log.Printf("    Warning: Failed to create table %s: %v\n", tableName, err)
					continue
				}

				// Insert items
				itemCount := 0
				for i, item := range tableData.Items {
					if item.ItemID > 0 {
						err := repo.InsertItem(int(tblID), item.ItemID, item.DropRate, i)
						if err != nil {
							continue
						}
						itemCount++
					}
				}
				totalItems += itemCount
				if itemCount > 0 {
					fmt.Printf("    → %s: %d items\n", displayName, itemCount)
				}
			}
		}
		fmt.Println()
	}

	// Print final stats
	fmt.Println("=== Import Complete ===")
	stats, _ := repo.GetStats()
	for key, count := range stats {
		fmt.Printf("  %s: %d\n", key, count)
	}
	fmt.Printf("\n✅ Successfully imported %d items total!\n", totalItems)
}

func getDisplayName(tableName string, displayNames map[string]string) string {
	displayName := displayNames[tableName]
	if displayName == "" {
		displayName = cleanTableName(tableName)
	}

	// Add wing suffix for Dire Maul to distinguish East/North/West
	if strings.HasPrefix(tableName, "DME") {
		displayName = displayName + " (East)"
	} else if strings.HasPrefix(tableName, "DMN") {
		displayName = displayName + " (North)"
	} else if strings.HasPrefix(tableName, "DMW") {
		displayName = displayName + " (West)"
	}

	return displayName
}

type LootTable struct {
	Name  string
	Items []LootItem
}

type LootItem struct {
	ItemID   int
	DropRate string
}

func parseLuaFile(path string) (map[string]*LootTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	tables := make(map[string]*LootTable)

	// Regex to match table start: TableName = {
	tableStartRegex := regexp.MustCompile(`^\s*(\w+)\s*=\s*\{`)
	// Regex to match item line: { 12345, "icon", "name", "desc", "10%" }
	itemRegex := regexp.MustCompile(`\{\s*(\d+)\s*,.*?"([^"]*%)`)
	// Simple item without drop rate
	simpleItemRegex := regexp.MustCompile(`\{\s*(\d+)\s*,`)

	lines := strings.Split(content, "\n")
	var currentTable string
	var currentItems []LootItem

	for _, line := range lines {
		// Check for table start
		if matches := tableStartRegex.FindStringSubmatch(line); len(matches) > 1 {
			// Save previous table
			if currentTable != "" && len(currentItems) > 0 {
				tables[currentTable] = &LootTable{
					Name:  currentTable,
					Items: currentItems,
				}
			}

			currentTable = matches[1]
			currentItems = []LootItem{}
			continue
		}

		// Check for item entry
		if strings.Contains(line, "{") && currentTable != "" {
			// Try to match item with drop rate
			if matches := itemRegex.FindStringSubmatch(line); len(matches) > 2 {
				itemID := 0
				fmt.Sscanf(matches[1], "%d", &itemID)
				dropRate := matches[2]
				currentItems = append(currentItems, LootItem{
					ItemID:   itemID,
					DropRate: dropRate,
				})
			} else if matches := simpleItemRegex.FindStringSubmatch(line); len(matches) > 1 {
				// Item without explicit drop rate
				itemID := 0
				fmt.Sscanf(matches[1], "%d", &itemID)
				currentItems = append(currentItems, LootItem{
					ItemID:   itemID,
					DropRate: "",
				})
			}
		}

		// Check for table end
		if strings.Contains(line, "};") && currentTable != "" {
			if len(currentItems) > 0 {
				tables[currentTable] = &LootTable{
					Name:  currentTable,
					Items: currentItems,
				}
			}
			currentTable = ""
			currentItems = []LootItem{}
		}
	}

	return tables, nil
}

func groupTablesByInstance(tables map[string]*LootTable) map[string][]string {
	groups := make(map[string][]string)

	instancePrefixes := map[string]string{
		"MC":     "Molten Core",
		"Ony":    "Onyxia's Lair",
		"BWL":    "Blackwing Lair",
		"ZG":     "Zul'Gurub",
		"AQ20":   "Ruins of Ahn'Qiraj",
		"AQ40":   "Temple of Ahn'Qiraj",
		"NAX":    "Naxxramas",
		"BRD":    "Blackrock Depths",
		"LBRS":   "Lower Blackrock Spire",
		"UBRS":   "Upper Blackrock Spire",
		"Strat":  "Stratholme",
		"Scholo": "Scholomance",
		"DM":     "Dire Maul",
		"ST":     "Sunken Temple",
		"Mara":   "Maraudon",
		"Uld":    "Uldaman",
		"RFK":    "Razorfen Kraul",
		"RFD":    "Razorfen Downs",
		"SM":     "Scarlet Monastery",
		"WC":     "Wailing Caverns",
		"SFK":    "Shadowfang Keep",
		"RFC":    "Ragefire Chasm",
		"DM2":    "Deadmines",
		"VC":     "Deadmines",
	}

	for tableName := range tables {
		matched := false
		for prefix, instanceName := range instancePrefixes {
			if strings.HasPrefix(tableName, prefix) {
				groups[instanceName] = append(groups[instanceName], tableName)
				matched = true
				break
			}
		}

		if !matched {
			// Put in "Other" group
			groups["Other"] = append(groups["Other"], tableName)
		}
	}

	return groups
}

// parseTableRegister extracts display names from TableRegister.lua using line-by-line parsing
func parseTableRegister(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	displayNames := make(map[string]string)
	scanner := bufio.NewScanner(file)

	// Regex patterns
	keyRegex := regexp.MustCompile(`\["(\w+)"\]\s*=`)
	atlasRegex := regexp.MustCompile(`"AtlasLoot\w*Items"`)
	alRegex := regexp.MustCompile(`AL\["([^"]+)"\]`)

	var currentKey string
	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line starts a new table entry
		if matches := keyRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentKey = matches[1]
			buffer.Reset()
			buffer.WriteString(line)
		} else if currentKey != "" {
			// Continue accumulating lines for current key
			buffer.WriteString(" ")
			buffer.WriteString(strings.TrimSpace(line))
		}

		// Check if we found the AtlasLoot marker
		if currentKey != "" && atlasRegex.MatchString(buffer.String()) {
			// Extract display name from accumulated buffer
			content := buffer.String()
			displayName := extractDisplayNameFromEntry(content, alRegex)
			if displayName != "" {
				displayNames[currentKey] = displayName
			}
			currentKey = ""
			buffer.Reset()
		}
	}

	return displayNames, scanner.Err()
}

// extractDisplayNameFromEntry parses an entry to get the display name
func extractDisplayNameFromEntry(content string, alRegex *regexp.Regexp) string {
	// Extract all AL["..."] references
	alMatches := alRegex.FindAllStringSubmatch(content, -1)

	if len(alMatches) == 0 {
		// Try to find a simple quoted string
		quoteRegex := regexp.MustCompile(`\{\s*"([^"]+)"`)
		if qm := quoteRegex.FindStringSubmatch(content); len(qm) > 1 {
			return qm[1]
		}
		return ""
	}

	if len(alMatches) == 1 {
		return alMatches[0][1]
	}

	// For multiple AL[] references, determine which one is the display name
	// by checking if it's followed by suffix indicators
	for i, match := range alMatches {
		text := match[1]
		// Skip common suffixes
		if text == "Rare" || text == "Summon" || text == "Quest" {
			continue
		}
		// Skip if this looks like a category/instance name and there's a next one
		if i < len(alMatches)-1 {
			nextText := alMatches[i+1][1]
			if nextText != "Rare" && nextText != "Summon" && nextText != "Quest" {
				// Both are content, take the second (likely boss name)
				return nextText
			}
		}
		// This is the display name
		return text
	}

	return ""
}

func cleanTableName(name string) string {
	// Remove common prefixes
	for _, prefix := range []string{"MC", "BWL", "Ony", "ZG", "AQ20", "AQ40", "NAX", "BRD", "LBRS", "UBRS", "Strat", "Scholo", "DM", "ST", "Mara", "Uld"} {
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
			break
		}
	}

	// Insert spaces before capital letters
	result := ""
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += " "
		}
		result += string(r)
	}

	if result == "" {
		return name
	}
	return result
}
