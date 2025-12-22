package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"shelllab/backend/database"
)

// App struct
type App struct {
	ctx           context.Context
	db            *database.SQLiteDB
	itemRepo      *database.ItemRepository
	categoryRepo  *database.CategoryRepository
	atlasLootRepo *database.AtlasLootRepository

	// Cache for category lookups
	categoryCache      map[int]*database.Category
	rootCategoryByName map[string]int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		categoryCache:      make(map[int]*database.Category),
		rootCategoryByName: make(map[string]int),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	fmt.Println("Initializing ShellLab (SQLite Version)...")

	// Initialize SQLite database
	dbPath := filepath.Join("data", "shelllab.db")
	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open database: %v\n", err)
		return
	}

	a.db = db
	a.itemRepo = database.NewItemRepository(db)
	a.categoryRepo = database.NewCategoryRepository(db)
	a.atlasLootRepo = database.NewAtlasLootRepository(db)

	// Print stats
	itemCount, _ := a.itemRepo.GetItemCount()
	catCount, _ := a.categoryRepo.GetCategoryCount()
	fmt.Printf("✓ Database Connected: %s\n", dbPath)
	fmt.Printf("  - Items: %d\n", itemCount)
	fmt.Printf("  - Categories: %d\n", catCount)

	// Build category cache
	a.buildCategoryCache()

	fmt.Println("✓ ShellLab ready!")
}

// buildCategoryCache builds a cache of categories for faster lookups
func (a *App) buildCategoryCache() {
	roots, err := a.categoryRepo.GetRootCategories()
	if err != nil {
		return
	}

	for _, cat := range roots {
		a.categoryCache[cat.ID] = cat
		a.rootCategoryByName[cat.Name] = cat.ID
	}
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// WaitForReady waits for the app to be ready (max 5 seconds)
func (a *App) WaitForReady() bool {
	return a.db != nil
}

// === Frontend API Methods ===

// GetRootCategories returns top-level categories (e.g., "Mage Sets", "Molten Core")
func (a *App) GetRootCategories() []*database.Category {
	fmt.Println("[API] GetRootCategories called")
	cats, err := a.categoryRepo.GetRootCategories()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []*database.Category{}
	}
	return cats
}

// GetChildCategories returns sub-categories (e.g., Bosses in an Instance)
func (a *App) GetChildCategories(parentID int) []*database.Category {
	cats, err := a.categoryRepo.GetChildCategories(parentID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []*database.Category{}
	}
	return cats
}

// GetCategoryItems returns items for a specific category (e.g., drops from Ragnaros)
func (a *App) GetCategoryItems(categoryID int) []*database.Item {
	items, err := a.categoryRepo.GetCategoryItems(categoryID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return []*database.Item{}
	}
	return a.enrichItemsWithIcons(items)
}

// SearchItems searches for items by name (Simple)
func (a *App) SearchItems(query string) []*database.Item {
	items, err := a.itemRepo.SearchItems(query, 50)
	if err != nil {
		fmt.Printf("Error searching items: %v\n", err)
		return []*database.Item{}
	}
	return a.enrichItemsWithIcons(items)
}

// GetItemClasses returns hierarchical item classes
func (a *App) GetItemClasses() []*database.ItemClass {
	classes, err := a.itemRepo.GetItemClasses()
	if err != nil {
		fmt.Printf("Error getting classes: %v\n", err)
		return []*database.ItemClass{}
	}
	return classes
}

// BrowseItemsByClass returns items for a specific class/subclass
func (a *App) BrowseItemsByClass(class, subClass int) []*database.Item {
	// Hardcoded limit for now, maybe pagination later
	items, _, err := a.itemRepo.GetItemsByClass(class, subClass, 200, 0)
	if err != nil {
		fmt.Printf("Error browsing items: %v\n", err)
		return []*database.Item{}
	}
	return a.enrichItemsWithIcons(items)
}

// BrowseItemsByClassAndSlot returns items for a specific class/subclass/inventoryType
func (a *App) BrowseItemsByClassAndSlot(class, subClass, inventoryType int) []*database.Item {
	items, _, err := a.itemRepo.GetItemsByClassAndSlot(class, subClass, inventoryType, 200, 0)
	if err != nil {
		fmt.Printf("Error browsing items by slot: %v\n", err)
		return []*database.Item{}
	}
	return a.enrichItemsWithIcons(items)
}

// AdvancedSearch performs a detailed search
func (a *App) AdvancedSearch(filter database.SearchFilter) *database.SearchResult {
	result, err := a.itemRepo.AdvancedSearch(filter)
	if err != nil {
		fmt.Printf("Error in advanced search: %v\n", err)
		return &database.SearchResult{Items: []*database.Item{}, TotalCount: 0}
	}
	result.Items = a.enrichItemsWithIcons(result.Items)
	return result
}

// GetTooltipData returns detailed item information (no Wails binding generation)
func (a *App) GetTooltipData(itemID int) *database.TooltipData {
	data, err := a.itemRepo.GetTooltipData(itemID)
	if err != nil {
		return nil
	}
	return data
}

// GetItemSets returns all item sets for browsing
func (a *App) GetItemSets() []*database.ItemSetBrowse {
	sets, err := a.itemRepo.GetItemSets()
	if err != nil {
		fmt.Printf("Error getting item sets: %v\n", err)
		return []*database.ItemSetBrowse{}
	}
	return sets
}

// GetItemSetDetail returns detailed information about a specific item set
func (a *App) GetItemSetDetail(itemSetID int) *database.ItemSetDetail {
	detail, err := a.itemRepo.GetItemSetDetail(itemSetID)
	if err != nil {
		fmt.Printf("Error getting item set detail: %v\n", err)
		return nil
	}
	// Enrich items with icons
	detail.Items = a.enrichItemsWithIcons(detail.Items)
	return detail
}

// GetCreatureTypes returns all creature types with counts
func (a *App) GetCreatureTypes() []*database.CreatureType {
	types, err := a.itemRepo.GetCreatureTypes()
	if err != nil {
		fmt.Printf("Error getting creature types: %v\n", err)
		return []*database.CreatureType{}
	}
	return types
}

// BrowseCreaturesByType returns creatures filtered by type
func (a *App) BrowseCreaturesByType(creatureType int) []*database.Creature {
	creatures, _, err := a.itemRepo.GetCreaturesByType(creatureType, 200, 0)
	if err != nil {
		fmt.Printf("Error browsing creatures: %v\n", err)
		return []*database.Creature{}
	}
	return creatures
}

// SearchCreatures searches for creatures by name
func (a *App) SearchCreatures(query string) []*database.Creature {
	creatures, err := a.itemRepo.SearchCreatures(query, 50)
	if err != nil {
		fmt.Printf("Error searching creatures: %v\n", err)
		return []*database.Creature{}
	}
	return creatures
}

// GetQuestCategories returns all quest categories (zones and sorts)
func (a *App) GetQuestCategories() []*database.QuestCategory {
	cats, err := a.itemRepo.GetQuestCategories()
	if err != nil {
		fmt.Printf("Error getting quest categories: %v\n", err)
		return []*database.QuestCategory{}
	}
	return cats
}

// GetQuestsByCategory returns quests filtered by category
func (a *App) GetQuestsByCategory(categoryID int) []*database.Quest {
	quests, err := a.itemRepo.GetQuestsByCategory(categoryID)
	if err != nil {
		fmt.Printf("Error browsing quests: %v\n", err)
		return []*database.Quest{}
	}
	return quests
}

// SearchQuests searches for quests by title
func (a *App) SearchQuests(query string) []*database.Quest {
	quests, err := a.itemRepo.SearchQuests(query)
	if err != nil {
		fmt.Printf("Error searching quests: %v\n", err)
		return []*database.Quest{}
	}
	return quests
}

// GetObjectTypes returns all object types
func (a *App) GetObjectTypes() []*database.ObjectType {
	types, err := a.itemRepo.GetObjectTypes()
	if err != nil {
		fmt.Printf("Error getting object types: %v\n", err)
		return []*database.ObjectType{}
	}
	return types
}

// GetObjectsByType returns objects filtered by type
func (a *App) GetObjectsByType(typeID int) []*database.GameObject {
	objects, err := a.itemRepo.GetObjectsByType(typeID)
	if err != nil {
		fmt.Printf("Error browsing objects: %v\n", err)
		return []*database.GameObject{}
	}
	return objects
}

// SearchObjects searches for objects by name
func (a *App) SearchObjects(query string) []*database.GameObject {
	objects, err := a.itemRepo.SearchObjects(query)
	if err != nil {
		fmt.Printf("Error searching objects: %v\n", err)
		return []*database.GameObject{}
	}
	return objects
}

// === Legacy Compatibility API (for master branch compatibility) ===

// LegacyBossLoot matches the structure from master branch
type LegacyBossLoot struct {
	BossName string           `json:"bossName"`
	Items    []LegacyLootItem `json:"items"`
}

// LegacyLootItem matches the structure from master branch
type LegacyLootItem struct {
	ItemID     int    `json:"itemId"`
	ItemName   string `json:"itemName"`
	IconName   string `json:"iconName"`
	Quality    int    `json:"quality"`
	DropChance string `json:"dropChance,omitempty"`
	SlotType   string `json:"slotType,omitempty"`
}

// GetCategories returns all top-level category names (legacy API)
func (a *App) GetCategories() []string {
	fmt.Println("[API] GetCategories called (AtlasLoot)")
	categories, err := a.atlasLootRepo.GetCategories()
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []string{}
	}
	return categories
}

// GetInstances returns modules for a category (legacy API)
func (a *App) GetInstances(categoryName string) []string {
	fmt.Printf("[API] GetInstances called for: %s\n", categoryName)
	modules, err := a.atlasLootRepo.GetModules(categoryName)
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []string{}
	}
	return modules
}

// GetTables returns tables/bosses for a module (new API for 3-tier structure)
func (a *App) GetTables(categoryName, moduleName string) []database.AtlasTable {
	fmt.Printf("[API] GetTables called for: %s / %s\n", categoryName, moduleName)
	tables, err := a.atlasLootRepo.GetTables(categoryName, moduleName)
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []database.AtlasTable{}
	}
	return tables
}

// GetLoot returns loot for a specific table (legacy API)
func (a *App) GetLoot(categoryName, instanceName, bossKey string) *LegacyBossLoot {
	fmt.Printf("[API] GetLoot called: %s / %s / %s\n", categoryName, instanceName, bossKey)

	// Query atlasloot tables directly
	lootEntries, err := a.atlasLootRepo.GetLootItems(categoryName, instanceName, bossKey)
	if err != nil {
		fmt.Printf("[API] Error getting loot: %v\n", err)
		return &LegacyBossLoot{BossName: bossKey, Items: []LegacyLootItem{}}
	}

	// Figure out boss display name for return (using bossKey as fallback)
	bossName := bossKey

	// Convert to legacy format
	var lootItems []LegacyLootItem
	for _, entry := range lootEntries {
		lootItems = append(lootItems, LegacyLootItem{
			ItemID:     entry.ItemID,
			ItemName:   entry.ItemName,
			IconName:   entry.IconName,
			Quality:    entry.Quality,
			DropChance: entry.DropChance,
		})
	}

	return &LegacyBossLoot{
		BossName: bossName,
		Items:    lootItems,
	}
}

// Helper to add full icon URLs
func (a *App) enrichItemsWithIcons(items []*database.Item) []*database.Item {
	for _, item := range items {
		a.enrichItemIcon(item)
	}
	return items
}

func (a *App) enrichItemIcon(item *database.Item) *database.Item {
	if item == nil {
		return nil
	}
	if item.IconPath != "" && !filepath.IsAbs(item.IconPath) && len(item.IconPath) < 100 {
		item.IconPath = strings.ToLower(item.IconPath)
	}
	return item
}
