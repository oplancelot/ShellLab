package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"shelllab/backend/database"
)

func main() {
	fmt.Println("===== Icon Downloader =====\n")

	// Connect to database
	dbPath := filepath.Join("data", "shelllab.db")
	fmt.Printf("Connecting to database: %s\n", dbPath)

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Query all item icon names
	fmt.Println("Loading items from database...")

	rows, err := db.DB().Query(`
		SELECT DISTINCT icon_path 
		FROM items 
		WHERE icon_path IS NOT NULL AND icon_path != ''
		ORDER BY icon_path
	`)
	if err != nil {
		fmt.Printf("Failed to query items: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	// Collect all unique icon names
	iconNames := make(map[string]bool)
	for rows.Next() {
		var iconPath string
		if err := rows.Scan(&iconPath); err != nil {
			continue
		}
		if iconPath != "" {
			iconNames[iconPath] = true
		}
	}

	fmt.Printf("Found %d unique icons\n\n", len(iconNames))

	if len(iconNames) == 0 {
		fmt.Println("No icons to download!")
		return
	}

	// Create output directory
	outputDir := "frontend/public/items/icons"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Download icons
	fmt.Println("Downloading icons...")

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Limit concurrency to 10

	downloaded := 0
	skipped := 0
	failed := 0

	var mu sync.Mutex

	for iconName := range iconNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			filename := filepath.Join(outputDir, strings.ToLower(name)+".jpg")

			// Check if file already exists
			if _, err := os.Stat(filename); err == nil {
				mu.Lock()
				skipped++
				mu.Unlock()
				return
			}

			// Download (Try Turtle WoW first)
			url := fmt.Sprintf("https://database.turtle-wow.org/images/icons/medium/%s.jpg", strings.ToLower(name))
			resp, err := http.Get(url)

			// If Turtle WoW fails or returns 404, try Zamimg
			if err != nil || resp.StatusCode != 200 {
				if resp != nil {
					resp.Body.Close()
				}
				// Fallback to Zamimg
				url = fmt.Sprintf("https://wow.zamimg.com/images/wow/icons/medium/%s.jpg", strings.ToLower(name))
				resp, err = http.Get(url)
			}

			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				mu.Lock()
				failed++
				mu.Unlock()
				return
			}

			// Save file
			file, err := os.Create(filename)
			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				return
			}

			mu.Lock()
			downloaded++
			if downloaded%100 == 0 {
				fmt.Printf("Progress: %d/%d downloaded\n", downloaded, len(iconNames))
			}
			mu.Unlock()

			time.Sleep(100 * time.Millisecond) // Avoid requests too fast
		}(iconName)
	}

	wg.Wait()

	fmt.Printf("\n===== Complete =====\n")
	fmt.Printf("Downloaded: %d\n", downloaded)
	fmt.Printf("Skipped (existing): %d\n", skipped)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("\nIcons saved to: %s\n", outputDir)
}
