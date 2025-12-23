package database

import "fmt"

// GameObject represents a WoW game object
type GameObject struct {
	Entry     int     `json:"entry"`
	Name      string  `json:"name"`
	Type      int     `json:"type"`
	TypeName  string  `json:"typeName"`
	DisplayID int     `json:"displayId"`
	Size      float64 `json:"size"`
}

// ObjectType represents a GO category
type ObjectType struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetObjectTypes returns all object types with counts
func (r *ItemRepository) GetObjectTypes() ([]*ObjectType, error) {
	rows, err := r.db.DB().Query(`
		SELECT t.id, t.name, COUNT(o.entry) as count
		FROM object_types t
		LEFT JOIN objects o ON t.id = o.type
		GROUP BY t.id, t.name
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*ObjectType
	for rows.Next() {
		t := &ObjectType{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Count); err != nil {
			continue
		}
		if t.Count > 0 {
			types = append(types, t)
		}
	}
	return types, nil
}

// GetObjectsByType returns objects filtered by type
func (r *ItemRepository) GetObjectsByType(typeID int, nameFilter string) ([]*GameObject, error) {
	whereClause := "WHERE type = ?"
	args := []interface{}{typeID}

	if nameFilter != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}

	query := fmt.Sprintf(`
		SELECT entry, name, type, display_id, size
		FROM objects
		%s
		ORDER BY name
		LIMIT 10000
	`, whereClause)

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
		SELECT o.entry, o.name, o.type, o.display_id, o.size, t.name
		FROM objects o
		LEFT JOIN object_types t ON o.type = t.id
		WHERE o.name LIKE ?
		ORDER BY length(o.name), o.name
		LIMIT 50
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []*GameObject
	for rows.Next() {
		o := &GameObject{}
		var typeName *string
		if err := rows.Scan(&o.Entry, &o.Name, &o.Type, &o.DisplayID, &o.Size, &typeName); err != nil {
			continue
		}
		if typeName != nil {
			o.TypeName = *typeName
		}
		objects = append(objects, o)
	}
	return objects, nil
}
