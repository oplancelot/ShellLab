import { useState, useEffect, useMemo } from "react";
import {
  GetCategories,
  GetInstances,
  GetTables,
} from "../../wailsjs/go/main/App";
import { useItemTooltip } from "../hooks/useItemTooltip";
import {
  PageLayout,
  ContentGrid,
  SidebarPanel,
  ContentPanel,
  ScrollList,
  SectionHeader,
  ListItem,
  LootItem,
  ItemTooltip,
} from "./ui";
import { filterItems } from "../utils/databaseApi";
import { GRID_LAYOUT } from "./common/layout";
import { getQualityColor } from "../utils/wow";

// Direct call to GetLoot - using window binding
const GetLoot = (category, instance, boss) => {
  if (window?.go?.main?.App?.GetLoot) {
    return window.go.main.App.GetLoot(category, instance, boss);
  }
  return Promise.resolve({ bossName: boss, items: [] });
};

function AtlasLootPage() {
  const [categories, setCategories] = useState([]);
  const [modules, setModules] = useState([]);
  const [tables, setTables] = useState([]);
  const [loot, setLoot] = useState(null);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [selectedCategory, setSelectedCategory] = useState("");
  const [selectedModule, setSelectedModule] = useState("");
  const [selectedTable, setSelectedTable] = useState("");

  // Filter states for each column
  const [categoryFilter, setCategoryFilter] = useState("");
  const [moduleFilter, setModuleFilter] = useState("");
  const [tableFilter, setTableFilter] = useState("");
  const [itemFilter, setItemFilter] = useState("");

  // Use shared tooltip hook
  const {
    hoveredItem,
    setHoveredItem,
    tooltipCache,
    loadTooltipData,
    handleMouseMove,
    handleItemEnter,
    getTooltipStyle,
  } = useItemTooltip();

  // Filtered lists
  const filteredCategories = useMemo(
    () => filterItems(categories, categoryFilter),
    [categories, categoryFilter]
  );
  const filteredModules = useMemo(
    () => filterItems(modules, moduleFilter),
    [modules, moduleFilter]
  );
  const filteredTables = useMemo(() => {
    const tablesWithNames = tables.map((t) => {
      if (typeof t === "string") {
        return { original: t, name: t };
      } else {
        return { original: t, name: t.displayName || t.key || t };
      }
    });
    return filterItems(tablesWithNames, tableFilter);
  }, [tables, tableFilter]);
  const filteredItems = useMemo(() => {
    if (!loot?.items) return [];
    return filterItems(loot.items, itemFilter);
  }, [loot, itemFilter]);

  // Load categories on mount
  useEffect(() => {
    setLoading(true);
    GetCategories()
      .then((cats) => {
        setCategories(cats || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Failed to load categories:", err);
        setError("Error loading categories");
        setLoading(false);
      });
  }, []);

  // Load modules when category changes
  useEffect(() => {
    if (selectedCategory) {
      setLoading(true);
      setModules([]);
      setTables([]);
      setLoot(null);
      setSelectedModule("");
      setSelectedTable("");
      setModuleFilter("");
      setTableFilter("");
      setItemFilter("");

      GetInstances(selectedCategory)
        .then((mods) => {
          setModules(mods || []);
          setLoading(false);
        })
        .catch((err) => {
          console.error("Failed to load modules:", err);
          setLoading(false);
        });
    }
  }, [selectedCategory]);

  // Load tables when module changes
  useEffect(() => {
    if (selectedModule && selectedCategory) {
      setLoading(true);
      setTables([]);
      setLoot(null);
      setSelectedTable("");
      setTableFilter("");
      setItemFilter("");

      GetTables(selectedCategory, selectedModule)
        .then((tbls) => {
          setTables(tbls || []);
          setLoading(false);
        })
        .catch((err) => {
          console.error("Failed to load tables:", err);
          setLoading(false);
        });
    }
  }, [selectedModule]);

  // Preload tooltips when loot changes
  useEffect(() => {
    if (loot?.items) {
      loot.items.slice(0, 20).forEach((item) => {
        if (item.itemId && !tooltipCache[item.itemId]) {
          loadTooltipData(item.itemId);
        }
      });
    }
  }, [loot, tooltipCache, loadTooltipData]);

  // Load loot when table is clicked
  const loadLoot = (table) => {
    setSelectedTable(table);
    setLoading(true);

    GetLoot(selectedCategory, selectedModule, table)
      .then((result) => {
        setLoot(result);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Failed to load loot:", err);
        setLoading(false);
      });
  };

  return (
    <PageLayout>
      {error && (
        <div className="mx-3 mt-3 p-3 bg-red-900/30 border border-red-500/30 rounded flex items-center gap-3 text-red-400">
          <span>❌</span>
          <span>{error}</span>
        </div>
      )}

      <ContentGrid columns={GRID_LAYOUT}>
        {/* Column 1: Categories */}
        <SidebarPanel>
          <SectionHeader
            title={`Categories (${filteredCategories.length})`}
            placeholder="Filter categories..."
            onFilterChange={setCategoryFilter}
          />
          <ScrollList>
            {loading && categories.length === 0 && (
              <div className="p-4 text-center text-wow-gold italic animate-pulse">
                Loading...
              </div>
            )}
            {filteredCategories.map((cat) => (
              <ListItem
                key={cat}
                active={selectedCategory === cat}
                onClick={() => {
                  setSelectedCategory(cat);
                  setCategoryFilter("");
                }}
              >
                {cat}
              </ListItem>
            ))}
          </ScrollList>
        </SidebarPanel>

        {/* Column 2: Modules/Instances */}
        <SidebarPanel>
          <SectionHeader
            title={
              selectedCategory
                ? `${selectedCategory} (${filteredModules.length})`
                : "Select Category"
            }
            placeholder="Filter modules..."
            onFilterChange={setModuleFilter}
          />
          <ScrollList>
            {loading && modules.length === 0 && selectedCategory && (
              <div className="p-4 text-center text-wow-gold italic animate-pulse">
                Loading...
              </div>
            )}
            {filteredModules.map((mod) => (
              <ListItem
                key={mod}
                active={selectedModule === mod}
                onClick={() => {
                  setSelectedModule(mod);
                  setModuleFilter("");
                }}
              >
                {mod}
              </ListItem>
            ))}
          </ScrollList>
        </SidebarPanel>

        {/* Column 3: Tables/Bosses */}
        <SidebarPanel>
          <SectionHeader
            title={
              selectedModule
                ? `${selectedModule} (${filteredTables.length})`
                : "Select Instance"
            }
            placeholder="Filter bosses..."
            onFilterChange={setTableFilter}
          />
          <ScrollList>
            {loading && tables.length === 0 && selectedModule && (
              <div className="p-4 text-center text-wow-gold italic animate-pulse">
                Loading...
              </div>
            )}
            {filteredTables.map((tbl, idx) => {
              const originalTable = tbl.original;
              const tableKey =
                typeof originalTable === "string"
                  ? originalTable
                  : originalTable.key || originalTable;
              return (
                <ListItem
                  key={tableKey || idx}
                  active={selectedTable === tableKey}
                  onClick={() => {
                    loadLoot(tableKey);
                    setTableFilter("");
                  }}
                >
                  {tbl.name}
                </ListItem>
              );
            })}
          </ScrollList>
        </SidebarPanel>

        {/* Column 4: Loot Display */}
        <ContentPanel>
          <SectionHeader
            title={
              loot ? `${loot.bossName} (${filteredItems.length})` : "Loot Table"
            }
            placeholder="Filter items..."
            onFilterChange={setItemFilter}
          />

          {loading && !loot && selectedTable && (
            <div className="flex-1 flex items-center justify-center text-wow-gold italic animate-pulse">
              Loading loot...
            </div>
          )}

          {filteredItems.length > 0 && (
            <ScrollList className="grid grid-cols-1 xl:grid-cols-2 gap-1 p-2 auto-rows-min">
              {filteredItems.map((item, idx) => {
                const itemId = item.itemId || item.entry || item.id;
                return (
                  <LootItem
                    key={itemId || idx}
                    item={{
                      entry: itemId,
                      name: item.itemName || item.name,
                      quality: item.quality,
                      iconPath: item.iconName || item.iconPath,
                      dropChance: item.dropChance,
                    }}
                    showDropChance
                    onMouseEnter={() => handleItemEnter(itemId)}
                    onMouseMove={(e) => handleMouseMove(e, itemId)}
                    onMouseLeave={() => setHoveredItem(null)}
                  />
                );
              })}
            </ScrollList>
          )}

          {!loading && filteredItems.length === 0 && selectedTable && (
            <div className="flex-1 flex items-center justify-center text-gray-600 italic">
              No loot data found for {selectedTable}
            </div>
          )}

          {!selectedTable && (
            <div className="flex-1 flex items-center justify-center text-gray-600 italic">
              Select a boss to view loot
            </div>
          )}
        </ContentPanel>
      </ContentGrid>

      {/* Global Tooltip Layer */}
      {hoveredItem && tooltipCache[hoveredItem] && (
        <ItemTooltip
          item={tooltipCache[hoveredItem]}
          tooltip={tooltipCache[hoveredItem]}
          style={getTooltipStyle()}
        />
      )}
    </PageLayout>
  );
}

export default AtlasLootPage;
