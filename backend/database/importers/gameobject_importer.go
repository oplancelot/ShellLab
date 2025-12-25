package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"shelllab/backend/database/models"
)

// GameObjectImporter handles game object data imports
type GameObjectImporter struct {
	db *sql.DB
}

// NewGameObjectImporter creates a new game object importer
func NewGameObjectImporter(db *sql.DB) *GameObjectImporter {
	return &GameObjectImporter{db: db}
}

// Import imports locks from JSON
func (g *GameObjectImporter) Import(locksPath string) error {
	// 1. Import Locks
	locksData, err := os.ReadFile(locksPath)
	if err != nil {
		return fmt.Errorf("failed to read locks JSON: %w", err)
	}
	var locks []models.LockEntry
	if err := json.Unmarshal(locksData, &locks); err != nil {
		return fmt.Errorf("failed to parse locks JSON: %w", err)
	}

	tx, err := g.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM locks")

	stmt, err := tx.Prepare(`
		INSERT INTO locks (id, type1, type2, type3, type4, type5, prop1, prop2, prop3, prop4, prop5, req1, req2, req3, req4, req5)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	for _, l := range locks {
		stmt.Exec(l.ID, l.Type1, l.Type2, l.Type3, l.Type4, l.Type5, l.Prop1, l.Prop2, l.Prop3, l.Prop4, l.Prop5, l.Req1, l.Req2, l.Req3, l.Req4, l.Req5)
	}
	stmt.Close()

	return tx.Commit()
}

// CheckAndImport checks if locks table is empty and imports if JSON exists
func (g *GameObjectImporter) CheckAndImport(locksPath string) error {
	var count int
	if err := g.db.QueryRow("SELECT COUNT(*) FROM locks").Scan(&count); err != nil {
		// Table might not exist yet if schema failed, but here we assume it does
		return nil
	}
	if count == 0 {
		return g.Import(locksPath)
	}
	return nil
}
