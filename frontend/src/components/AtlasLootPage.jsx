import { useState, useEffect } from 'react'
import { GetCategories, GetInstances, GetTables } from '../../wailsjs/go/main/App'

// Direct call to GetLoot - using window binding
const GetLoot = (category, instance, boss) => {
    if (window?.go?.main?.App?.GetLoot) {
        return window.go.main.App.GetLoot(category, instance, boss)
    }
    return Promise.resolve({ bossName: boss, items: [] })
}

// Direct call to GetTooltipData
const GetTooltipData = (itemId) => {
    if (window?.go?.main?.App?.GetTooltipData) {
        return window.go.main.App.GetTooltipData(itemId)
    }
    return Promise.resolve(null)
}

function AtlasLootPage() {
    const [categories, setCategories] = useState([])
    const [modules, setModules] = useState([])
    const [tables, setTables] = useState([])
    const [loot, setLoot] = useState(null)
    
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')
    const [tooltipCache, setTooltipCache] = useState({})
    const [hoveredItem, setHoveredItem] = useState(null)
    const [tooltipPos, setTooltipPos] = useState({ top: 0, left: 0 })

    const [selectedCategory, setSelectedCategory] = useState('')
    const [selectedModule, setSelectedModule] = useState('')
    const [selectedTable, setSelectedTable] = useState('')

    // Get quality class name
    const getQualityClass = (quality) => `q${quality || 0}`
    
    // Get quality color for inline style
    const getQualityColor = (quality) => {
        const colors = {
            0: '#9d9d9d', 1: '#ffffff', 2: '#1eff00',
            3: '#0070dd', 4: '#a335ee', 5: '#ff8000', 6: '#e6cc80'
        }
        return colors[quality] || '#ffffff'
    }

    // Handle mouse move - update tooltip position following mouse
    const handleMouseMove = (e, item) => {
        const lootContainer = e.currentTarget.closest('.loot')
        const containerRect = lootContainer ? lootContainer.getBoundingClientRect() : { left: 0, right: window.innerWidth, top: 0, bottom: window.innerHeight }
        const itemRect = e.currentTarget.getBoundingClientRect()
        
        // Tooltip dimensions
        const tooltipWidth = 320
        const tooltipHeight = 400
        
        // Position tooltip to the right and below the cursor
        let left = e.clientX + 15
        let top = e.clientY + 15
        
        // Don't let tooltip cover the item row - keep it below the item
        if (top < itemRect.bottom + 5) {
            top = itemRect.bottom + 5
        }
        
        // Keep within container bounds - horizontal
        if (left + tooltipWidth > containerRect.right - 10) {
            left = e.clientX - tooltipWidth - 15
        }
        if (left < containerRect.left + 10) {
            left = containerRect.left + 10
        }
        
        // Keep within container bounds - vertical
        if (top + tooltipHeight > containerRect.bottom - 10) {
            top = containerRect.bottom - tooltipHeight - 10
        }
        if (top < containerRect.top + 10) {
            top = containerRect.top + 10
        }
        
        setTooltipPos({ top, left })
        setHoveredItem(item.itemId)
    }

    // Handle item enter - load tooltip data
    const handleItemEnter = (item) => {
        loadTooltipData(item.itemId)
    }

    // Load tooltip data for an item
    const loadTooltipData = async (itemId) => {
        if (tooltipCache[itemId]) return tooltipCache[itemId]
        
        try {
            const data = await GetTooltipData(itemId)
            if (data) {
                setTooltipCache(prev => ({ ...prev, [itemId]: data }))
                return data
            }
        } catch (err) {
            console.error('Failed to load tooltip:', err)
        }
        return null
    }

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

    // Render tooltip content
    const renderTooltip = (item) => {
        if (hoveredItem !== item.itemId) return null
        
        const tooltip = tooltipCache[item.itemId]
        
        const tooltipStyle = {
            position: 'fixed',
            left: tooltipPos.left,
            top: tooltipPos.top,
            zIndex: 10000,
        }
        
        if (!tooltip) {
            return (
                <div className="item-tooltip" style={tooltipStyle}>
                    <div className="tooltip-name" style={{color: getQualityColor(item.quality)}}>
                        {item.itemName || 'Unknown Item'}
                    </div>
                    <div className="tooltip-loading">Loading...</div>
                </div>
            )
        }
        
        return (
            <div 
                className="item-tooltip"
                style={tooltipStyle}
                onMouseEnter={() => setHoveredItem(item.itemId)}
                onMouseLeave={() => setHoveredItem(null)}
            >
                <div className="tooltip-name" style={{color: getQualityColor(tooltip.quality)}}>
                    {tooltip.name}
                </div>
                
                {tooltip.setName && (
                    <div className="tooltip-setname">{tooltip.setName}</div>
                )}
                
                {tooltip.itemLevel > 0 && (
                    <div className="tooltip-itemlevel">Item Level {tooltip.itemLevel}</div>
                )}
                
                {tooltip.binding && (
                    <div className="tooltip-binding">{tooltip.binding}</div>
                )}
                
                <div className="tooltip-slot-type">
                    {tooltip.slotName && <span className="tooltip-slot">{tooltip.slotName}</span>}
                    {tooltip.typeName && <span className="tooltip-type">{tooltip.typeName}</span>}
                </div>
                
                {tooltip.classes && (
                    <div className="tooltip-classes">{tooltip.classes}</div>
                )}
                
                {tooltip.races && (
                    <div className="tooltip-races">{tooltip.races}</div>
                )}
                
                {tooltip.damageText && (
                    <div className="tooltip-damage">
                        <span>{tooltip.damageText}</span>
                        <span className="tooltip-speed">{tooltip.speedText}</span>
                    </div>
                )}
                
                {tooltip.dps && (
                    <div className="tooltip-dps">{tooltip.dps}</div>
                )}
                
                {tooltip.armor > 0 && (
                    <div className="tooltip-armor">{tooltip.armor} Armor</div>
                )}
                
                {tooltip.stats && tooltip.stats.length > 0 && (
                    <div className="tooltip-stats">
                        {tooltip.stats.map((stat, i) => (
                            <div key={i} className="tooltip-stat">{stat}</div>
                        ))}
                    </div>
                )}
                
                {tooltip.resistances && tooltip.resistances.length > 0 && (
                    <div className="tooltip-resistances">
                        {tooltip.resistances.map((res, i) => (
                            <div key={i} className="tooltip-resistance">{res}</div>
                        ))}
                    </div>
                )}
                
                {tooltip.spellEffects && tooltip.spellEffects.length > 0 && (
                    <div className="tooltip-effects">
                        {tooltip.spellEffects.map((effect, i) => (
                            <div key={i} className="tooltip-effect">{effect}</div>
                        ))}
                    </div>
                )}
                
                {tooltip.setInfo && (
                    <div className="tooltip-set-info">
                        <div className="tooltip-set-name">{tooltip.setInfo.name}</div>
                        {tooltip.setInfo.items && tooltip.setInfo.items.map((setItem, i) => (
                            <div key={i} className="tooltip-set-item">{setItem}</div>
                        ))}
                        <div className="tooltip-set-spacer"></div>
                        {tooltip.setInfo.bonuses && tooltip.setInfo.bonuses.map((bonus, i) => (
                            <div key={i} className="tooltip-set-bonus">{bonus}</div>
                        ))}
                    </div>
                )}
                
                {tooltip.durability && (
                    <div className="tooltip-durability">{tooltip.durability}</div>
                )}
                
                {tooltip.requiredLevel > 1 && (
                    <div className="tooltip-reqlevel">Requires Level {tooltip.requiredLevel}</div>
                )}
                
                {tooltip.description && (
                    <div className="tooltip-description">"{tooltip.description}"</div>
                )}
                
                {tooltip.sellPrice && (
                    <div className="tooltip-sellprice">Sell Price: {tooltip.sellPrice}</div>
                )}
            </div>
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
                                    onMouseEnter={() => handleItemEnter(item)}
                                    onMouseMove={(e) => handleMouseMove(e, item)}
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
