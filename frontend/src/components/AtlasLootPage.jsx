import { useState, useEffect } from 'react'
import { GetCategories, GetInstances, GetTables } from '../../wailsjs/go/main/App'
import { useItemTooltip } from '../hooks/useItemTooltip'
import ItemTooltip, { getQualityColor } from './ItemTooltip'

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

    // Preload tooltips for visible items
    const preloadTooltips = (items) => {
        if (!items) return
        items.forEach(item => {
            if (item.itemId && !tooltipCache[item.itemId]) {
                loadTooltipData(item.itemId)
            }
        })
    }

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
            preloadTooltips(loot.items)
        }
    }, [loot])

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

    // Render tooltip content using shared component
    const renderTooltip = (item) => {
        if (hoveredItem !== item.itemId) return null
        
        const tooltip = tooltipCache[item.itemId]
        
        return (
            <ItemTooltip
                item={item}
                tooltip={tooltip}
                style={getTooltipStyle()}
                onMouseEnter={() => setHoveredItem(item.itemId)}
                onMouseLeave={() => setHoveredItem(null)}
            />
        )
    }

    return (
        <div className="app">
            {error && (
                <div className="error-alert">
                    <span>❌</span>
                    <span>{error}</span>
                </div>
            )}

            <div className="content" style={{ display: 'grid', gridTemplateColumns: '150px 200px 200px 1fr', gap: 0 }}>
                {/* Column 1: Categories */}
                <aside className="sidebar">
                    <h2>Categories</h2>
                    {loading && categories.length === 0 && (
                        <div className="loading">Loading...</div>
                    )}
                    <div className="list">
                        {categories.map(cat => (
                            <button
                                key={cat}
                                className={selectedCategory === cat ? 'active' : ''}
                                onClick={() => setSelectedCategory(cat)}
                            >
                                {cat}
                            </button>
                        ))}
                    </div>
                </aside>

                {/* Column 2: Modules/Instances */}
                <section className="instances">
                    <h2>
                        {selectedCategory || 'Select Category'}
                    </h2>
                    {selectedCategory && (
                        <div className="list">
                            {loading && modules.length === 0 ? (
                                <div className="loading">Loading...</div>
                            ) : (
                                <>
                                    {modules.map(mod => (
                                        <div
                                            key={mod}
                                            className={`item ${selectedModule === mod ? 'active' : ''}`}
                                            onClick={() => setSelectedModule(mod)}
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
                    <h2>
                        {selectedModule || 'Select Instance'}
                    </h2>
                    {selectedModule && (
                        <div className="list">
                            {loading && tables.length === 0 ? (
                                <div className="loading">Loading...</div>
                            ) : (
                                <>
                                    {tables.map(tbl => (
                                        <div
                                            key={tbl.key || tbl}
                                            className={`item ${selectedTable === (tbl.key || tbl) ? 'active' : ''}`}
                                            onClick={() => loadLoot(tbl.key || tbl)}
                                        >
                                            {tbl.displayName || tbl}
                                        </div>
                                    ))}
                                </>
                            )}
                        </div>
                    )}
                </section>

                {/* Column 4: Loot Display */}
                <section className="loot">
                    <h2>
                        {loot ? `${loot.bossName} (${loot.items?.length || 0} items)` : 'Loot Table'}
                    </h2>
                    {loading && !loot ? (
                        <div className="loading">Loading loot...</div>
                    ) : loot && loot.items && loot.items.length > 0 ? (
                        <div className="loot-items">
                            {loot.items.map((item, idx) => (
                                <div 
                                    key={idx} 
                                    className="loot-item"
                                    data-quality={item.quality || 0}
                                    onMouseEnter={() => handleItemEnter(item.itemId)}
                                    onMouseMove={(e) => handleMouseMove(e, item.itemId)}
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
                                    
                                    <span className="item-id">[{item.itemId}]</span>
                                    
                                    <span 
                                        className={`item-name ${getQualityClass(item.quality)}`}
                                        style={{color: getQualityColor(item.quality)}}
                                    >
                                        {item.itemName || item.dropChance || 'Unknown Item'}
                                    </span>
                                    
                                    {item.dropChance && (
                                        <span className="item-drop-chance">{item.dropChance}</span>
                                    )}
                                    
                                    {renderTooltip(item)}
                                </div>
                            ))}
                        </div>
                    ) : selectedTable ? (
                        <p className="placeholder">No loot data found for {selectedTable}</p>
                    ) : (
                        <p className="placeholder">Select a boss to view loot</p>
                    )}
                </section>
            </div>
        </div>
    )
}

export default AtlasLootPage
