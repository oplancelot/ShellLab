import { useState, useEffect, useMemo } from 'react'
import { GetCategories, GetInstances, GetTables } from '../../wailsjs/go/main/App'
import { useItemTooltip } from '../hooks/useItemTooltip'
import ItemTooltip, { getQualityColor } from './ItemTooltip'
import { SectionHeader } from './common/SectionHeader'
import { filterItems } from '../utils/databaseApi'
import { GRID_LAYOUT } from './common/layout'
import './common/PageLayout.css'

// Direct call to GetLoot - using window binding
const GetLoot = (category, instance, boss) => {
    if (window?.go?.main?.App?.GetLoot) {
        return window.go.main.App.GetLoot(category, instance, boss)
    }
    return Promise.resolve({ bossName: boss, items: [] })
}

function AtlasLootPage() {
    const [categories, setCategories] = useState([])
    const [modules, setModules] = useState([])
    const [tables, setTables] = useState([])
    const [loot, setLoot] = useState(null)
    
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    const [selectedCategory, setSelectedCategory] = useState('')
    const [selectedModule, setSelectedModule] = useState('')
    const [selectedTable, setSelectedTable] = useState('')

    // Filter states for each column
    const [categoryFilter, setCategoryFilter] = useState('')
    const [moduleFilter, setModuleFilter] = useState('')
    const [tableFilter, setTableFilter] = useState('')
    const [itemFilter, setItemFilter] = useState('')

    // Use shared tooltip hook
    const {
        hoveredItem,
        setHoveredItem,
        tooltipCache,
        loadTooltipData,
        handleMouseMove,
        handleItemEnter,
        getTooltipStyle,
    } = useItemTooltip()

    // Get quality class name
    const getQualityClass = (quality) => `q${quality || 0}`

    // Filtered lists
    const filteredCategories = useMemo(() => filterItems(categories, categoryFilter), [categories, categoryFilter])
    const filteredModules = useMemo(() => filterItems(modules, moduleFilter), [modules, moduleFilter])
    const filteredTables = useMemo(() => {
        // Convert tables to format with name property for filtering
        const tablesWithNames = tables.map(t => {
            if (typeof t === 'string') {
                return { original: t, name: t }
            } else {
                return { original: t, name: t.displayName || t.key || t }
            }
        })
        return filterItems(tablesWithNames, tableFilter)
    }, [tables, tableFilter])
    const filteredItems = useMemo(() => {
        if (!loot?.items) return []
        return filterItems(loot.items.map(item => ({ ...item, name: item.itemName })), itemFilter)
    }, [loot, itemFilter])

    // Load categories on mount
    useEffect(() => {
        setLoading(true)
        GetCategories()
            .then(cats => {
                setCategories(cats || [])
                setLoading(false)
            })
            .catch(err => {
                console.error('Failed to load categories:', err)
                setError('Error loading categories')
                setLoading(false)
            })
    }, [])

    // Load modules when category changes
    useEffect(() => {
        if (selectedCategory) {
            setLoading(true)
            setModules([])
            setTables([])
            setLoot(null)
            setSelectedModule('')
            setSelectedTable('')
            setModuleFilter('')
            setTableFilter('')
            setItemFilter('')

            GetInstances(selectedCategory)
                .then(mods => {
                    setModules(mods || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error('Failed to load modules:', err)
                    setLoading(false)
                })
        }
    }, [selectedCategory])

    // Load tables when module changes
    useEffect(() => {
        if (selectedModule && selectedCategory) {
            setLoading(true)
            setTables([])
            setLoot(null)
            setSelectedTable('')
            setTableFilter('')
            setItemFilter('')

            GetTables(selectedCategory, selectedModule)
                .then(tbls => {
                    setTables(tbls || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error('Failed to load tables:', err)
                    setLoading(false)
                })
        }
    }, [selectedModule])

    // Preload tooltips when loot changes
    useEffect(() => {
        if (loot?.items) {
            loot.items.slice(0, 20).forEach(item => {
                if (item.itemId && !tooltipCache[item.itemId]) {
                    loadTooltipData(item.itemId)
                }
            })
        }
    }, [loot, tooltipCache, loadTooltipData])

    // Load loot when table is clicked
    const loadLoot = (table) => {
        setSelectedTable(table)
        setLoading(true)

        GetLoot(selectedCategory, selectedModule, table)
            .then(result => {
                setLoot(result)
                setLoading(false)
            })
            .catch(err => {
                console.error('Failed to load loot:', err)
                setLoading(false)
            })
    }

    return (
        <div className="database-page">
            {error && (
                <div className="error-alert">
                    <span>❌</span>
                    <span>{error}</span>
                </div>
            )}

            <div className="content" style={{ gridTemplateColumns: GRID_LAYOUT }}>
                {/* Column 1: Categories */}
                <aside className="sidebar">
                    <SectionHeader 
                        title={`Categories (${filteredCategories.length})`}
                        placeholder="Filter categories..."
                        onFilterChange={setCategoryFilter}
                    />
                    <div className="list">
                        {loading && categories.length === 0 && (
                            <div className="loading">Loading...</div>
                        )}
                        {filteredCategories.map(cat => (
                            <button
                                key={cat}
                                className={selectedCategory === cat ? 'active' : ''}
                                onClick={() => {
                                    setSelectedCategory(cat)
                                    setCategoryFilter('')
                                }}
                            >
                                {cat}
                            </button>
                        ))}
                    </div>
                </aside>

                {/* Column 2: Modules/Instances */}
                <section className="instances">
                    <SectionHeader 
                        title={selectedCategory ? `${selectedCategory} (${filteredModules.length})` : 'Select Category'}
                        placeholder="Filter modules..."
                        onFilterChange={setModuleFilter}
                    />
                    {selectedCategory && (
                        <div className="list">
                            {loading && modules.length === 0 ? (
                                <div className="loading">Loading...</div>
                            ) : (
                                <>
                                    {filteredModules.map(mod => (
                                        <div
                                            key={mod}
                                            className={`item ${selectedModule === mod ? 'active' : ''}`}
                                            onClick={() => {
                                                setSelectedModule(mod)
                                                setModuleFilter('')
                                            }}
                                        >
                                            {mod}
                                        </div>
                                    ))}
                                </>
                            )}
                        </div>
                    )}
                </section>

                {/* Column 3: Tables/Bosses */}
                <section className="instances">
                    <SectionHeader 
                        title={selectedModule ? `${selectedModule} (${filteredTables.length})` : 'Select Instance'}
                        placeholder="Filter bosses..."
                        onFilterChange={setTableFilter}
                    />
                    {selectedModule && (
                        <div className="list">
                            {loading && tables.length === 0 ? (
                                <div className="loading">Loading...</div>
                            ) : (
                                <>
                                    {filteredTables.map((tbl, idx) => {
                                        const originalTable = tbl.original
                                        const tableKey = typeof originalTable === 'string' ? originalTable : (originalTable.key || originalTable)
                                        return (
                                        <div
                                            key={tableKey || idx}
                                            className={`item ${selectedTable === tableKey ? 'active' : ''}`}
                                            onClick={() => {
                                                loadLoot(tableKey)
                                                setTableFilter('')
                                            }}
                                        >
                                            {tbl.name}
                                        </div>
                                        )
                                    })}
                                </>
                            )}
                        </div>
                    )}
                </section>

                {/* Column 4: Loot Display */}
                <section className="loot">
                    <SectionHeader 
                        title={loot ? `${loot.bossName} (${filteredItems.length})` : 'Loot Table'}
                        placeholder="Filter items..."
                        onFilterChange={setItemFilter}
                    />
                    
                    {loading && !loot ? (
                        <div className="loading">Loading loot...</div>
                    ) : filteredItems.length > 0 ? (
                        <div className="loot-items">
                            {filteredItems.map((item, idx) => {
                                const itemId = item.itemId || item.entry || item.id
                                return (
                                <div 
                                    key={itemId || idx} 
                                    className="loot-item"
                                    data-quality={item.quality || 0}
                                    onMouseEnter={() => handleItemEnter(itemId)}
                                    onMouseMove={(e) => handleMouseMove(e, itemId)}
                                    onMouseLeave={() => setHoveredItem(null)}
                                >
                                    {item.iconName ? (
                                        <img 
                                            src={`/items/icons/${item.iconName.toLowerCase()}.jpg`}
                                            alt={item.itemName || 'Item'}
                                            className="item-icon"
                                            onError={(e) => {
                                                if (!e.target.src.includes('zamimg.com')) {
                                                    e.target.src = `https://wow.zamimg.com/images/wow/icons/medium/${item.iconName.toLowerCase()}.jpg`
                                                } else {
                                                    e.target.style.display = 'none'
                                                }
                                            }}
                                        />
                                    ) : (
                                        <div className="item-icon-placeholder">?</div>
                                    )}
                                    
                                    <span className="item-id">[{itemId}]</span>
                                    
                                    <span 
                                        className={`item-name ${getQualityClass(item.quality)}`}
                                        style={{color: getQualityColor(item.quality)}}
                                    >
                                        {item.itemName || item.dropChance || 'Unknown Item'}
                                    </span>
                                    
                                    {item.dropChance && (
                                        <span className="item-drop-chance">{item.dropChance}</span>
                                    )}
                                </div>
                                )
                            })}
                        </div>
                    ) : selectedTable ? (
                        <p className="placeholder">No loot data found for {selectedTable}</p>
                    ) : (
                        <p className="placeholder">Select a boss to view loot</p>
                    )}
                </section>
            </div>

            {/* Global Tooltip Layer */}
            {hoveredItem && tooltipCache[hoveredItem] && (
                 <ItemTooltip
                     item={tooltipCache[hoveredItem]}
                     tooltip={tooltipCache[hoveredItem]}
                     style={getTooltipStyle()}
                 />
            )}
        </div>
    )
}

export default AtlasLootPage
