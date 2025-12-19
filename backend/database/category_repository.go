package database

// CategoryRepository handles category-related database operations
type CategoryRepository struct {
	db *SQLiteDB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *SQLiteDB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Category represents a loot category (instance, boss, set, etc.)
type Category struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	ParentID  *int   `json:"parentId,omitempty"`
	Type      string `json:"type"`
	SortOrder int    `json:"sortOrder"`
}

// CategoryItem represents an item in a category
type CategoryItem struct {
	CategoryID int    `json:"categoryId"`
	ItemID     int    `json:"itemId"`
	DropRate   string `json:"dropRate,omitempty"`
	SortOrder  int    `json:"sortOrder"`
}

// InsertCategory inserts a new category
func (r *CategoryRepository) InsertCategory(cat *Category) (int64, error) {
	result, err := r.db.DB().Exec(`
		INSERT OR REPLACE INTO categories (key, name, parent_id, type, sort_order)
		VALUES (?, ?, ?, ?, ?)
	`, cat.Key, cat.Name, cat.ParentID, cat.Type, cat.SortOrder)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertCategoryItem links an item to a category
func (r *CategoryRepository) InsertCategoryItem(catID, itemID int, dropRate string, sortOrder int) error {
	_, err := r.db.DB().Exec(`
		INSERT INTO category_items (category_id, item_id, drop_rate, sort_order)
		VALUES (?, ?, ?, ?)
	`, catID, itemID, dropRate, sortOrder)
	return err
}

// GetRootCategories returns all top-level categories
func (r *CategoryRepository) GetRootCategories() ([]*Category, error) {
	rows, err := r.db.DB().Query(`
		SELECT id, key, name, parent_id, type, sort_order
		FROM categories
		WHERE parent_id IS NULL
		ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetChildCategories returns child categories of a parent
func (r *CategoryRepository) GetChildCategories(parentID int) ([]*Category, error) {
	rows, err := r.db.DB().Query(`
		SELECT id, key, name, parent_id, type, sort_order
		FROM categories
		WHERE parent_id = ?
		ORDER BY sort_order, name
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetCategoryByKey retrieves a category by its key
func (r *CategoryRepository) GetCategoryByKey(key string) (*Category, error) {
	row := r.db.DB().QueryRow(`
		SELECT id, key, name, parent_id, type, sort_order
		FROM categories
		WHERE key = ?
	`, key)

	cat := &Category{}
	var parentID *int
	err := row.Scan(&cat.ID, &cat.Key, &cat.Name, &parentID, &cat.Type, &cat.SortOrder)
	if err != nil {
		return nil, err
	}
	cat.ParentID = parentID
	return cat, nil
}

// GetCategoryItems returns all items in a category
func (r *CategoryRepository) GetCategoryItems(categoryID int) ([]*Item, error) {
	rows, err := r.db.DB().Query(`
		SELECT i.entry, i.name, i.quality, i.item_level, i.required_level,
			i.class, i.subclass, i.inventory_type, i.icon_path,
			ci.drop_rate
		FROM items i
		JOIN category_items ci ON i.entry = ci.item_id
		WHERE ci.category_id = ?
		ORDER BY ci.sort_order, i.quality DESC, i.item_level DESC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		var dropRate string
		err := rows.Scan(
			&item.Entry, &item.Name, &item.Quality, &item.ItemLevel,
			&item.RequiredLevel, &item.Class, &item.SubClass, &item.InventoryType, &item.IconPath,
			&dropRate,
		)
		if err != nil {
			continue
		}
		item.DropRate = dropRate
		items = append(items, item)
	}

	return items, nil
}

// GetCategoryCount returns the total number of categories
func (r *CategoryRepository) GetCategoryCount() (int, error) {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	return count, err
}

// Helper to scan category rows
func (r *CategoryRepository) scanCategories(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
}) ([]*Category, error) {
	var cats []*Category
	for rows.Next() {
		cat := &Category{}
		var parentID *int
		err := rows.Scan(&cat.ID, &cat.Key, &cat.Name, &parentID, &cat.Type, &cat.SortOrder)
		if err != nil {
			continue
		}
		cat.ParentID = parentID
		cats = append(cats, cat)
	}
	return cats, nil
}
