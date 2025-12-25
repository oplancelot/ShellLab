package repositories

import (
	"database/sql"

	"shelllab/backend/database/models"
)

// GameObjectRepository handles game object-related database operations
type GameObjectRepository struct {
	db *sql.DB
}

// NewGameObjectRepository creates a new game object repository
func NewGameObjectRepository(db *sql.DB) *GameObjectRepository {
	return &GameObjectRepository{db: db}
}

// GetObjectTypes returns derived categories based on Turtlehead logic
func (r *GameObjectRepository) GetObjectTypes() ([]*models.ObjectType, error) {
	types := []*models.ObjectType{}

	// Helper for derived types (Herbalism, Mining, Lockpicking)
	countDerived := func(propID int, name string, id int) {
		var count int
		r.db.QueryRow(`
			SELECT COUNT(DISTINCT o.entry) FROM gameobject_template o
			JOIN locks l ON o.data0 = l.id
			WHERE o.type = 3 AND (l.prop1 = ? OR l.prop2 = ? OR l.prop3 = ? OR l.prop4 = ? OR l.prop5 = ?)
		`, propID, propID, propID, propID, propID).Scan(&count)
		if count > 0 {
			types = append(types, &models.ObjectType{ID: id, Name: name, Count: count})
		}
	}

	countDerived(2, "Herbalism", -3)
	countDerived(3, "Mining", -4)
	countDerived(1, "Lockpicking", -5)

	// Standard types
	standardTypes := []struct {
		ID   int
		Name string
	}{
		{3, "Chests"}, {25, "Fishing Pools"}, {9, "Books & Texts"},
		{2, "Quest Givers"}, {19, "Mailboxes"}, {17, "Fishing Nodes"},
		{0, "Doors"}, {10, "Interactive"}, {1, "Buttons"},
	}

	for _, st := range standardTypes {
		var count int
		r.db.QueryRow("SELECT COUNT(*) FROM gameobject_template WHERE type = ?", st.ID).Scan(&count)
		if count > 0 {
			types = append(types, &models.ObjectType{ID: st.ID, Name: st.Name, Count: count})
		}
	}

	return types, nil
}

// GetObjectsByType returns objects filtered by type
func (r *GameObjectRepository) GetObjectsByType(typeID int, nameFilter string) ([]*models.GameObject, error) {
	var query string
	var args []interface{}

	baseSelect := "SELECT entry, name, type, displayId as display_id, size FROM gameobject_template o"

	if typeID < 0 {
		var propID int
		switch typeID {
		case -3:
			propID = 2
		case -4:
			propID = 3
		case -5:
			propID = 1
		}
		query = baseSelect + `
			JOIN locks l ON o.data0 = l.id
			WHERE o.type = 3 AND (l.prop1 = ? OR l.prop2 = ? OR l.prop3 = ? OR l.prop4 = ? OR l.prop5 = ?)
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

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []*models.GameObject
	for rows.Next() {
		o := &models.GameObject{}
		if err := rows.Scan(&o.Entry, &o.Name, &o.Type, &o.DisplayID, &o.Size); err != nil {
			continue
		}
		objects = append(objects, o)
	}
	return objects, nil
}

// SearchObjects searches for objects by name
func (r *GameObjectRepository) SearchObjects(query string) ([]*models.GameObject, error) {
	rows, err := r.db.Query(`
		SELECT entry, name, type, displayId as display_id, size FROM gameobject_template
		WHERE name LIKE ? ORDER BY length(name), name LIMIT 50
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []*models.GameObject
	for rows.Next() {
		o := &models.GameObject{}
		if err := rows.Scan(&o.Entry, &o.Name, &o.Type, &o.DisplayID, &o.Size); err != nil {
			continue
		}
		objects = append(objects, o)
	}
	return objects, nil
}

// GetObjectCount returns total count
func (r *GameObjectRepository) GetObjectCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM gameobject_template").Scan(&count)
	return count, err
}
