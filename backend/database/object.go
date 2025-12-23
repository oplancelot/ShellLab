package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// GameObject represents a WoW game object
type GameObject struct {
	Entry     int     `json:"entry"`
	Name      string  `json:"name"`
	Type      int     `json:"type"`
	TypeName  string  `json:"typeName"`
	DisplayID int     `json:"displayId"`
	Size      float64 `json:"size"`
	Data      []int   `json:"data,omitempty"` // For JSON import
}

// ObjectType represents a GO category
type ObjectType struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// LockEntry represents a lock record
type LockEntry struct {
	ID    int `json:"lockID"`
	Type1 int `json:"type1"`
	Type2 int `json:"type2"`
	Type3 int `json:"type3"`
	Type4 int `json:"type4"`
	Type5 int `json:"type5"`
	Prop1 int `json:"lockproperties1"`
	Prop2 int `json:"lockproperties2"`
	Prop3 int `json:"lockproperties3"`
	Prop4 int `json:"lockproperties4"`
	Prop5 int `json:"lockproperties5"`
	Req1  int `json:"requiredskill1"`
	Req2  int `json:"requiredskill2"`
	Req3  int `json:"requiredskill3"`
	Req4  int `json:"requiredskill4"`
	Req5  int `json:"requiredskill5"`
}

// ImportObjects imports game objects and locks from JSON
func (r *ItemRepository) ImportObjects(objectsPath, locksPath string) error {
	// 1. Import Locks
	locksData, err := os.ReadFile(locksPath)
	if err != nil {
		return fmt.Errorf("failed to read locks JSON: %w", err)
	}
	var locks []LockEntry
	if err := json.Unmarshal(locksData, &locks); err != nil {
		return fmt.Errorf("failed to parse locks JSON: %w", err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing
	tx.Exec("DELETE FROM locks")

	stmt, err := tx.Prepare(`
		INSERT INTO locks (id, type1, type2, type3, type4, type5, prop1, prop2, prop3, prop4, prop5, req1, req2, req3, req4, req5)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range locks {
		_, err := stmt.Exec(l.ID, l.Type1, l.Type2, l.Type3, l.Type4, l.Type5, l.Prop1, l.Prop2, l.Prop3, l.Prop4, l.Prop5, l.Req1, l.Req2, l.Req3, l.Req4, l.Req5)
		if err != nil {
			continue // Skip errors
		}
	}

	// 2. Import Objects
	objData, err := os.ReadFile(objectsPath)
	if err != nil {
		return fmt.Errorf("failed to read objects JSON: %w", err)
	}
	var objects []GameObject
	if err := json.Unmarshal(objData, &objects); err != nil {
		return fmt.Errorf("failed to parse objects JSON: %w", err)
	}

	// Clear existing
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
		_, err := objStmt.Exec(o.Entry, o.Name, o.Type, o.DisplayID, o.Size, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7])
		if err != nil {
			continue
		}
	}

	return tx.Commit()
}

// CheckAndImportObjects checks if objects table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportObjects(objectsPath, locksPath string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM objects").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		return r.ImportObjects(objectsPath, locksPath)
	}
	return nil
}

// GetObjectTypes returns derived categories based on Turtlehead logic
func (r *ItemRepository) GetObjectTypes() ([]*ObjectType, error) {
	types := []*ObjectType{}

	// Helper to count derived types (Herbalism, Mining, Lockpicking)
	// These rely on `objects.data0` (lockID) matching `locks.id` and checking `locks.propX`
	countDerived := func(propID int, name string, id int) {
		var count int
		// Turtlehead: type=3 (Chest) AND lockproperties IN (...)
		query := `
			SELECT COUNT(DISTINCT o.entry)
			FROM objects o
			JOIN locks l ON o.data0 = l.id
			WHERE o.type = 3 AND (
				l.prop1 = ? OR l.prop2 = ? OR l.prop3 = ? OR l.prop4 = ? OR l.prop5 = ?
			)
		`
		r.db.DB().QueryRow(query, propID, propID, propID, propID, propID).Scan(&count)
		if count > 0 {
			types = append(types, &ObjectType{ID: id, Name: name, Count: count})
		}
	}

	// Type -3: Herbalism (lockproperties = 2)
	countDerived(2, "Herbalism", -3)
	// Type -4: Mining (lockproperties = 3)
	countDerived(3, "Mining", -4)
	// Type -5: Lockpicking (lockproperties = 1)
	countDerived(1, "Lockpicking", -5)

	// Standard types
	standardTypes := []struct {
		ID   int
		Name string
	}{
		{3, "Chests"},
		{25, "Fishing Pools"},
		{9, "Books & Texts"},
		{2, "Quest Givers"},
		{19, "Mailboxes"},
		{17, "Fishing Nodes"},
		{0, "Doors"},
		{10, "Interactive"},
		{1, "Buttons"},
	}

	for _, st := range standardTypes {
		var count int
		r.db.DB().QueryRow("SELECT COUNT(*) FROM objects WHERE type = ?", st.ID).Scan(&count)
		if count > 0 {
			types = append(types, &ObjectType{ID: st.ID, Name: st.Name, Count: count})
		}
	}

	return types, nil
}

// GetObjectsByType returns objects filtered by type
func (r *ItemRepository) GetObjectsByType(typeID int, nameFilter string) ([]*GameObject, error) {
	var query string
	var args []interface{}

	baseSelect := "SELECT entry, name, type, display_id, size FROM objects o"

	if typeID < 0 {
		var propID int
		switch typeID {
		case -3:
			propID = 2 // Herbalism
		case -4:
			propID = 3 // Mining
		case -5:
			propID = 1 // Lockpicking
		}

		query = baseSelect + `
			JOIN locks l ON o.data0 = l.id
			WHERE o.type = 3 AND (
				l.prop1 = ? OR l.prop2 = ? OR l.prop3 = ? OR l.prop4 = ? OR l.prop5 = ?
			)
		`
		args = append(args, propID, propID, propID, propID, propID)
	} else {
		query = baseSelect + " WHERE o.type = ?"
		args = append(args, typeID)
	}

	if nameFilter != "" {
		query += " AND o.name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	query += " ORDER BY o.name LIMIT 10000"

	rows, err := r.db.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []*GameObject
	for rows.Next() {
		o := &GameObject{}
		if err := rows.Scan(&o.Entry, &o.Name, &o.Type, &o.DisplayID, &o.Size); err != nil {
			continue
		}
		objects = append(objects, o)
	}
	return objects, nil
}

// SearchObjects searches for objects by name
func (r *ItemRepository) SearchObjects(query string) ([]*GameObject, error) {
	rows, err := r.db.DB().Query(`
		SELECT entry, name, type, display_id, size
		FROM objects
		WHERE name LIKE ?
		ORDER BY length(name), name
		LIMIT 50
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []*GameObject
	for rows.Next() {
		o := &GameObject{}
		if err := rows.Scan(&o.Entry, &o.Name, &o.Type, &o.DisplayID, &o.Size); err != nil {
			continue
		}
		objects = append(objects, o)
	}
	return objects, nil
}
