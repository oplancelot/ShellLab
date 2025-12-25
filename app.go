package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"shelllab/backend/database"
	"shelllab/backend/services"
)

// App struct
type App struct {
	ctx context.Context
	db  *database.SQLiteDB

	// Repositories
	itemRepo      *database.ItemRepository
	creatureRepo  *database.CreatureRepository
	questRepo     *database.QuestRepository
	spellRepo     *database.SpellRepository
	lootRepo      *database.LootRepository
	factionRepo   *database.FactionRepository
	objectRepo    *database.GameObjectRepository
	categoryRepo  *database.CategoryRepository
	atlasLootRepo *database.AtlasLootRepository

	// Cache for category lookups
	categoryCache      map[int]*database.Category
	rootCategoryByName map[string]int

	// Services
	iconService *services.IconService
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

	// Ensure schema exists
	if err := db.InitSchema(); err != nil {
		fmt.Printf("ERROR: Failed to initialize schema: %v\n", err)
		return
	}

	a.db = db

	// Initialize all repositories
	a.itemRepo = database.NewItemRepository(db)
	a.creatureRepo = database.NewCreatureRepository(db)
	a.questRepo = database.NewQuestRepository(db)
	a.spellRepo = database.NewSpellRepository(db)
	a.lootRepo = database.NewLootRepository(db)
	a.factionRepo = database.NewFactionRepository(db)
	a.objectRepo = database.NewGameObjectRepository(db)
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

	// Data import using importers
	dataDir := "data"

	// Import Objects (locks table - still needed)
	fmt.Println("Checking object data...")
	objectImporter := database.NewGameObjectImporter(db)
	if err := objectImporter.CheckAndImport(dataDir + "/locks.json"); err != nil {
		fmt.Printf("ERROR: Failed to import objects: %v\n", err)
	}

	// Import Item Sets
	fmt.Println("Checking item sets...")
	itemSetImporter := database.NewItemSetImporter(db)
	if err := itemSetImporter.CheckAndImport(dataDir); err != nil {
		fmt.Printf("ERROR: Failed to import item sets: %v\n", err)
	}

	// Import Loot
	fmt.Println("Checking loot data...")
	lootImporter := database.NewLootImporter(db)
	lootImporter.CheckAndImport(dataDir)

	// Import Factions
	fmt.Println("Checking faction data...")
	factionImporter := database.NewFactionImporter(db)
	factionImporter.CheckAndImport(dataDir)

	// Import Metadata (Zones, Skills)
	fmt.Println("Checking metadata...")
	metadataImporter := database.NewMetadataImporter(db)
	metadataImporter.ImportAll(dataDir)

	// Import AtlasLoot Data
	fmt.Println("Checking AtlasLoot data...")
	atlasImporter := database.NewAtlasLootImporter(db)
	if err := atlasImporter.CheckAndImport(dataDir); err != nil {
		fmt.Printf("ERROR: Failed to import AtlasLoot: %v\n", err)
	}

	// Import Full 1:1 MySQL tables (if not already imported)
	fmt.Println("Checking full MySQL table data...")
	genImporter := database.NewGeneratedImporter(db)
	a.importFullTables(genImporter, dataDir)

	// Initialize services
	a.iconService = services.NewIconService(db)
	a.iconService.StartDownload()

	fmt.Println("✓ ShellLab ready!")
}

// importFullTables imports the 1:1 MySQL tables if they are empty
func (a *App) importFullTables(importer *database.GeneratedImporter, dataDir string) {
	tables := []struct {
		name     string
		jsonFile string
		importFn func(string) error
	}{
		{"item_template", "item_template.json", importer.ImportItemTemplate},
		{"creature_template", "creature_template.json", importer.ImportCreatureTemplate},
		{"quest_template", "quest_template.json", importer.ImportQuestTemplate},
		{"spell_template", "spell_template.json", importer.ImportSpellTemplate},
		{"gameobject_template", "gameobject_template.json", importer.ImportGameobjectTemplate},
	}

	for _, t := range tables {
		var count int
		a.db.DB().QueryRow("SELECT COUNT(*) FROM " + t.name).Scan(&count)
		if count == 0 {
			jsonPath := filepath.Join(dataDir, t.jsonFile)
			fmt.Printf("  Importing %s from %s...\n", t.name, t.jsonFile)
			if err := t.importFn(jsonPath); err != nil {
				fmt.Printf("  ERROR importing %s: %v\n", t.name, err)
			} else {
				var newCount int
				a.db.DB().QueryRow("SELECT COUNT(*) FROM " + t.name).Scan(&newCount)
				fmt.Printf("  ✓ Imported %d rows into %s\n", newCount, t.name)
			}
		}

		// ALWAYS check/refresh icons if the files exist
		if t.name == "item_template" {
			fmt.Println("  Refreshing item icons...")
			if err := importer.ImportItemIcons(filepath.Join(dataDir, "item_icons.json")); err != nil {
				fmt.Printf("  ERROR updating item icons: %v\n", err)
			}
		}
		if t.name == "spell_template" {
			fmt.Println("  Refreshing spell icons...")
			if err := importer.ImportSpellIcons(filepath.Join(dataDir, "spells_enhanced.json")); err != nil {
				fmt.Printf("  ERROR updating spell icons: %v\n", err)
			}
		}
	}
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
	fmt.Println("[API] GetItemClasses called")
	classes, err := a.itemRepo.GetItemClasses()
	if err != nil {
		fmt.Printf("[API] Error getting classes: %v\n", err)
		return []*database.ItemClass{}
	}
	fmt.Printf("[API] GetItemClasses returning %d classes\n", len(classes))
	return classes
}

// BrowseItemsByClass returns items for a specific class/subclass
func (a *App) BrowseItemsByClass(class, subClass int, nameFilter string) []*database.Item {
	fmt.Printf("[API] BrowseItemsByClass called: class=%d, subClass=%d, filter='%s'\n", class, subClass, nameFilter)
	// No limit - return all matching items
	items, _, err := a.itemRepo.GetItemsByClass(class, subClass, nameFilter, 999999, 0)
	if err != nil {
		fmt.Printf("[API] Error browsing items: %v\n", err)
		return []*database.Item{}
	}
	fmt.Printf("[API] BrowseItemsByClass returning %d items\n", len(items))
	return a.enrichItemsWithIcons(items)
}

// BrowseItemsByClassAndSlot returns items for a specific class/subclass/inventoryType
func (a *App) BrowseItemsByClassAndSlot(class, subClass, inventoryType int, nameFilter string) []*database.Item {
	// No limit - return all matching items
	items, _, err := a.itemRepo.GetItemsByClassAndSlot(class, subClass, inventoryType, nameFilter, 999999, 0)
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
	fmt.Println("[API] GetItemSets called")
	sets, err := a.itemRepo.GetItemSets()
	if err != nil {
		fmt.Printf("[API] Error getting item sets: %v\n", err)
		return []*database.ItemSetBrowse{}
	}
	fmt.Printf("[API] GetItemSets returning %d sets\n", len(sets))
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
	fmt.Println("[API] GetCreatureTypes called")
	types, err := a.creatureRepo.GetCreatureTypes()
	if err != nil {
		fmt.Printf("[API] Error getting creature types: %v\n", err)
		return []*database.CreatureType{}
	}
	fmt.Printf("[API] GetCreatureTypes returning %d types\n", len(types))
	return types
}

// BrowseCreaturesByType returns creatures filtered by type
func (a *App) BrowseCreaturesByType(creatureType int, nameFilter string) []*database.Creature {
	fmt.Printf("[API] BrowseCreaturesByType called: type=%d, filter='%s'\n", creatureType, nameFilter)
	// No limit - return all matching creatures
	creatures, _, err := a.creatureRepo.GetCreaturesByType(creatureType, nameFilter, 999999, 0)
	if err != nil {
		fmt.Printf("[API] Error browsing creatures: %v\n", err)
		return []*database.Creature{}
	}
	return creatures
}

// CreaturePageResult is the result of paginated creature query
type CreaturePageResult struct {
	Creatures []*database.Creature `json:"creatures"`
	Total     int                  `json:"total"`
	HasMore   bool                 `json:"hasMore"`
}

// BrowseCreaturesByTypePaged returns creatures with pagination support
func (a *App) BrowseCreaturesByTypePaged(creatureType int, nameFilter string, limit, offset int) *CreaturePageResult {
	fmt.Printf("[API] BrowseCreaturesByTypePaged called: type=%d, limit=%d, offset=%d\n", creatureType, limit, offset)
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	creatures, total, err := a.creatureRepo.GetCreaturesByType(creatureType, nameFilter, limit, offset)
	if err != nil {
		fmt.Printf("[API] Error browsing creatures: %v\n", err)
		return &CreaturePageResult{Creatures: []*database.Creature{}, Total: 0, HasMore: false}
	}
	return &CreaturePageResult{
		Creatures: creatures,
		Total:     total,
		HasMore:   offset+len(creatures) < total,
	}
}

// SearchCreatures searches for creatures by name
func (a *App) SearchCreatures(query string) []*database.Creature {
	creatures, err := a.creatureRepo.SearchCreatures(query, 50)
	if err != nil {
		fmt.Printf("Error searching creatures: %v\n", err)
		return []*database.Creature{}
	}
	return creatures
}

// GetQuestCategories returns all quest categories (zones and sorts)
func (a *App) GetQuestCategories() []*database.QuestCategory {
	cats, err := a.questRepo.GetQuestCategories()
	if err != nil {
		fmt.Printf("Error getting quest categories: %v\n", err)
		return []*database.QuestCategory{}
	}
	return cats
}

// GetQuestsByCategory returns quests filtered by category
func (a *App) GetQuestsByCategory(categoryID int) ([]*database.Quest, error) {
	quests, err := a.questRepo.GetQuestsByCategory(categoryID)
	if err != nil {
		fmt.Printf("Error browsing quests: %v\n", err)
		return nil, err
	}
	return quests, nil
}

// SearchQuests searches for quests by title
func (a *App) SearchQuests(query string) ([]*database.Quest, error) {
	quests, err := a.questRepo.SearchQuests(query)
	if err != nil {
		fmt.Printf("Error searching quests: %v\n", err)
		return nil, err
	}
	return quests, nil
}

// GetObjectTypes returns all object types
func (a *App) GetObjectTypes() []*database.ObjectType {
	fmt.Println("[API] GetObjectTypes called")
	types, err := a.objectRepo.GetObjectTypes()
	if err != nil {
		fmt.Printf("[API] Error getting object types: %v\n", err)
		return []*database.ObjectType{}
	}
	fmt.Printf("[API] GetObjectTypes returning %d types\n", len(types))
	return types
}

// GetObjectsByType returns objects filtered by type
func (a *App) GetObjectsByType(typeID int, nameFilter string) []*database.GameObject {
	fmt.Printf("[API] GetObjectsByType called: type=%d, filter='%s'\n", typeID, nameFilter)
	objects, err := a.objectRepo.GetObjectsByType(typeID, nameFilter)
	if err != nil {
		fmt.Printf("[API] Error browsing objects: %v\n", err)
		return []*database.GameObject{}
	}
	fmt.Printf("[API] GetObjectsByType returning %d objects\n", len(objects))
	return objects
}

// SearchObjects searches for objects by name
func (a *App) SearchObjects(query string) []*database.GameObject {
	objects, err := a.objectRepo.SearchObjects(query)
	if err != nil {
		fmt.Printf("Error searching objects: %v\n", err)
		return []*database.GameObject{}
	}
	return objects
}

// SearchSpells searches for spells by name
func (a *App) SearchSpells(query string) []*database.Spell {
	spells, err := a.spellRepo.SearchSpells(query)
	if err != nil {
		fmt.Printf("Error searching spells: %v\n", err)
		return []*database.Spell{}
	}
	return spells
}

// GetFactions returns all factions
func (a *App) GetFactions() []*database.Faction {
	fmt.Println("[API] GetFactions called")
	factions, err := a.factionRepo.GetFactions()
	if err != nil {
		fmt.Printf("[API] Error getting factions: %v\n", err)
		return []*database.Faction{}
	}
	fmt.Printf("[API] GetFactions returning %d factions\n", len(factions))
	return factions
}

// GetCreatureLoot returns the loot for a creature
func (a *App) GetCreatureLoot(entry int) []*database.LootItem {
	loot, err := a.lootRepo.GetCreatureLoot(entry)
	if err != nil {
		fmt.Printf("Error getting creature loot: %v\n", err)
		return []*database.LootItem{}
	}
	return loot
}

// === Legacy Compatibility API (for master branch compatibility) ===

// LegacyBossLoot matches the structure from master branch
type LegacyBossLoot struct {
	BossName string           `json:"bossName"`
	Items    []LegacyLootItem `json:"items"`
}

// GetCreatureDetail returns full details for a creature
func (a *App) GetCreatureDetail(entry int) (*database.CreatureDetail, error) {
	c, err := a.creatureRepo.GetCreatureDetail(entry)
	if err != nil {
		fmt.Printf("Error getting creature detail [%d]: %v\n", entry, err)
		return nil, err
	}
	return c, nil
}

// GetQuestDetail returns full details for a quest
func (a *App) GetQuestDetail(entry int) (*database.QuestDetail, error) {
	q, err := a.questRepo.GetQuestDetail(entry)
	if err != nil {
		fmt.Printf("Error getting quest detail [%d]: %v\n", entry, err)
		return nil, err
	}
	return q, nil
}

// GetItemDetail returns full details for an item
func (a *App) GetItemDetail(entry int) (*database.ItemDetail, error) {
	i, err := a.itemRepo.GetItemDetail(entry)
	if err != nil {
		fmt.Printf("Error getting item detail [%d]: %v\n", entry, err)
		return nil, err
	}
	if i != nil && i.Item != nil {
		a.enrichItemIcon(i.Item)
	}
	return i, nil
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
			ItemName:   entry.Name,
			IconName:   entry.IconPath,
			Quality:    entry.Quality,
			DropChance: entry.DropChance,
		})
	}

	result := &LegacyBossLoot{
		BossName: bossName,
		Items:    lootItems,
	}
	fmt.Printf("[API] GetLoot returning %d items for %s\n", len(lootItems), bossKey)
	return result
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

// === Spell Skills API (3-level navigation) ===

// GetSpellSkillCategories returns spell skill categories (Class Skills, Professions, etc.)
func (a *App) GetSpellSkillCategories() []*database.SpellSkillCategory {
	fmt.Println("[API] GetSpellSkillCategories called")
	cats, err := a.spellRepo.GetSpellSkillCategories()
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []*database.SpellSkillCategory{}
	}
	return cats
}

// GetSpellSkillsByCategory returns skills for a category
func (a *App) GetSpellSkillsByCategory(categoryID int) []*database.SpellSkill {
	fmt.Printf("[API] GetSpellSkillsByCategory called: %d\n", categoryID)
	skills, err := a.spellRepo.GetSpellSkillsByCategory(categoryID)
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []*database.SpellSkill{}
	}
	return skills
}

// GetSpellsBySkill returns spells for a skill
func (a *App) GetSpellsBySkill(skillID int, nameFilter string) []*database.Spell {
	fmt.Printf("[API] GetSpellsBySkill called: %d\n", skillID)
	spells, err := a.spellRepo.GetSpellsBySkill(skillID, nameFilter)
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []*database.Spell{}
	}
	return spells
}

// === Enhanced Quest Categories API (3-level navigation) ===

// GetQuestCategoryGroups returns quest category groups (Zones, Class Quests, etc.)
func (a *App) GetQuestCategoryGroups() []*database.QuestCategoryGroup {
	fmt.Println("[API] GetQuestCategoryGroups called")
	groups, err := a.questRepo.GetQuestCategoryGroups()
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []*database.QuestCategoryGroup{}
	}
	return groups
}

// GetQuestCategoriesByGroup returns categories for a group
func (a *App) GetQuestCategoriesByGroup(groupID int) []*database.QuestCategoryEnhanced {
	fmt.Printf("[API] GetQuestCategoriesByGroup called: %d\n", groupID)
	cats, err := a.questRepo.GetQuestCategoriesByGroup(groupID)
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []*database.QuestCategoryEnhanced{}
	}
	return cats
}

// GetQuestsByEnhancedCategory returns quests for a category (ZoneOrSort value)
func (a *App) GetQuestsByEnhancedCategory(categoryID int, nameFilter string) []*database.Quest {
	fmt.Printf("[API] GetQuestsByEnhancedCategory called: %d\n", categoryID)
	quests, err := a.questRepo.GetQuestsByEnhancedCategory(categoryID, nameFilter)
	if err != nil {
		fmt.Printf("[API] Error: %v\n", err)
		return []*database.Quest{}
	}
	return quests
}
