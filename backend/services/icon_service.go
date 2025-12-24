package services

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

// IconService handles downloading icons
type IconService struct {
	db        *database.SQLiteDB
	outputDir string
	client    *http.Client
}

// NewIconService creates a new IconService
func NewIconService(db *database.SQLiteDB) *IconService {
	// Default output dir
	outputDir := filepath.Join("frontend", "public", "items", "icons")
	return &IconService{
		db:        db,
		outputDir: outputDir,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// StartDownload initiates the download process in background
func (s *IconService) StartDownload() {
	go func() {
		fmt.Println("[IconService] Starting background icon download...")
		if err := s.downloadProcess(); err != nil {
			fmt.Printf("[IconService] Error: %v\n", err)
		} else {
			fmt.Println("[IconService] Download complete.")
		}
	}()
}

func (s *IconService) downloadProcess() error {
	// Ensure directory exists
	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// 1. Collect all unique icon names from database
	iconNames := make(map[string]bool)

	// Items
	rows, err := s.db.DB().Query(`
		SELECT DISTINCT icon_path 
		FROM items 
		WHERE icon_path IS NOT NULL AND icon_path != ''
	`)
	if err != nil {
		return fmt.Errorf("query items failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil && name != "" {
			iconNames[strings.ToLower(name)] = true
		}
	}

	// Spells
	rows, err = s.db.DB().Query(`
		SELECT DISTINCT icon_name 
		FROM spells 
		WHERE icon_name IS NOT NULL AND icon_name != ''
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil && name != "" {
				iconNames[strings.ToLower(name)] = true
			}
		}
	} else {
		// Log error but continue (column might not exist yet if migration pending)
		fmt.Printf("[IconService] Warning: Failed to query spell icons: %v\n", err)
	}

	fmt.Printf("[IconService] Found %d unique icons to check.\n", len(iconNames))

	// 2. Filter out existing icons
	var toDownload []string
	for name := range iconNames {
		path := filepath.Join(s.outputDir, name+".jpg")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			toDownload = append(toDownload, name)
		}
	}

	if len(toDownload) == 0 {
		fmt.Println("[IconService] All icons exist. Skipping download.")
		return nil
	}

	fmt.Printf("[IconService] Downloading %d missing icons...\n", len(toDownload))

	// 3. Download worker pool
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency limit

	// Sources to try in order
	sources := []string{
		"https://database.turtlecraft.gg/images/icons/large/%s.jpg",
		"https://wow.zamimg.com/images/wow/icons/large/%s.jpg",
		"https://aowow.trinitycore.info/static/images/wow/icons/large/%s.jpg",
	}

	var successCount, failCount int
	var mu sync.Mutex

	for _, name := range toDownload {
		wg.Add(1)
		go func(iconName string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			success := false
			for _, srcFmt := range sources {
				url := fmt.Sprintf(srcFmt, iconName)
				if err := s.downloadFile(url, iconName); err == nil {
					success = true
					break
				}
			}

			mu.Lock()
			if success {
				successCount++
			} else {
				failCount++
				// create a placeholder or just log?
				// fmt.Printf("Failed to download: %s\n", iconName)
			}
			mu.Unlock()

			// Slight delay to be nice
			time.Sleep(50 * time.Millisecond)
		}(name)
	}

	wg.Wait()
	fmt.Printf("[IconService] Downloaded: %d, Failed: %d\n", successCount, failCount)
	return nil
}

func (s *IconService) downloadFile(url, name string) error {
	resp, err := s.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	// Check content type
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "image") {
		return fmt.Errorf("invalid content type: %s", ct)
	}

	filename := filepath.Join(s.outputDir, name+".jpg")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
