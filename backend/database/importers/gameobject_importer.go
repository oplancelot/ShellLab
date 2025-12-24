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

// Import imports game objects and locks from JSON
func (g *GameObjectImporter) Import(objectsPath, locksPath string) error {
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

	// 2. Import Objects
	objData, err := os.ReadFile(objectsPath)
	if err != nil {
		return fmt.Errorf("failed to read objects JSON: %w", err)
	}
	var objects []models.GameObject
	if err := json.Unmarshal(objData, &objects); err != nil {
		return fmt.Errorf("failed to parse objects JSON: %w", err)
	}

	tx.Exec("DELETE FROM objects")

	objStmt, err := tx.Prepare(`
		INSERT INTO objects (entry, name, type, display_id, size, data0, data1, data2, data3, data4, data5, data6, data7)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer objStmt.Close()

	for _, o := range objects {
		data := make([]int, 8)
		copy(data, o.Data)
		objStmt.Exec(o.Entry, o.Name, o.Type, o.DisplayID, o.Size, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7])
	}

	return tx.Commit()
}

// CheckAndImport checks if objects table is empty and imports if JSON exists
func (g *GameObjectImporter) CheckAndImport(objectsPath, locksPath string) error {
	var count int
	if err := g.db.QueryRow("SELECT COUNT(*) FROM objects").Scan(&count); err != nil {
		return nil
	}
	if count == 0 {
		return g.Import(objectsPath, locksPath)
	}
	return nil
}
