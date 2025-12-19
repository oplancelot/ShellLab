package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"shelllab/backend/database"
)

func main() {
	fmt.Println("===== Icon Data Importer from aowow.sql =====\n")

	// Connect to database
	dbPath := filepath.Join("data", "shelllab.db")
	fmt.Printf("Connecting to ShellLoot database: %s\n", dbPath)

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Read aowow.sql
	sqlPath := filepath.Join("data", "sql", "aowow.sql")
	fmt.Printf("Reading icon data from: %s\n", sqlPath)

	file, err := os.Open(sqlPath)
	if err != nil {
		fmt.Printf("Failed to open SQL file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Parse INSERT statements from aowow_icons table
	iconMap := make(map[int]string) // displayId -> iconName
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // Increase buffer size

	// Match INSERT INTO `aowow_icons` data rows
	insertRegex := regexp.MustCompile(`INSERT INTO \x60aowow_icons\x60`)
	valueRegex := regexp.MustCompile(`\((\d+),\s*'([^']+)'\)`)

	inIconsInsert := false
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Check if entering aowow_icons INSERT section
		if insertRegex.MatchString(line) {
			inIconsInsert = true
		}

		// If in INSERT section, parse values
		if inIconsInsert {
			// Find all matching (id, 'iconname') pairs
			matches := valueRegex.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 3 {
					var id int
					fmt.Sscanf(match[1], "%d", &id)
					iconName := match[2]
					iconMap[id] = iconName
				}
			}

			// If line ends with semicolon, this INSERT statement has ended
			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				inIconsInsert = false
			}
		}

		if lineCount%100000 == 0 {
			fmt.Printf("Processed %d lines, found %d icons so far...\n", lineCount, len(iconMap))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading SQL file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Found %d icon mappings (displayId -> iconName)\n", len(iconMap))

	if len(iconMap) == 0 {
		fmt.Println("No icon data found! Exiting.")
		return
	}

	// Update icon_path field in items table
	fmt.Println("\nUpdating items table with icon paths...")

	tx, err := db.DB().Begin()
	if err != nil {
		fmt.Printf("Failed to begin transaction: %v\n", err)
		os.Exit(1)
	}
	defer tx.Rollback()

	updateStmt, err := tx.Prepare("UPDATE items SET icon_path = ? WHERE display_id = ?")
	if err != nil {
		fmt.Printf("Failed to prepare update statement: %v\n", err)
		os.Exit(1)
	}
	defer updateStmt.Close()

	updatedCount := 0
	for displayId, iconName := range iconMap {
		result, err := updateStmt.Exec(iconName, displayId)
		if err != nil {
			continue
		}
		affected, _ := result.RowsAffected()
		if affected > 0 {
			updatedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Updated %d items with icon paths\n", updatedCount)

	// Verify results
	var iconCount int
	db.DB().QueryRow("SELECT COUNT(*) FROM items WHERE icon_path IS NOT NULL AND icon_path != ''").Scan(&iconCount)

	fmt.Printf("\n===== Complete =====\n")
	fmt.Printf("Items with icons: %d\n", iconCount)
	fmt.Println("\nYou can now run: go run scripts/download_icons/main.go")
}
