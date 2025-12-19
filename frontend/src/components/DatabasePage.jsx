import { useState, useEffect } from 'react'
import { GetItemClasses, BrowseItemsByClass } from '../../wailsjs/go/main/App'

// Direct call to GetTooltipData since Wails doesn't auto-generate binding for it
const GetTooltipData = (itemId) => {
    if (window?.go?.main?.App?.GetTooltipData) {
        return window.go.main.App.GetTooltipData(itemId)
    }
    return Promise.resolve(null)
}

function DatabasePage() {
    const [activeTab, setActiveTab] = useState('items')
    
    // Items Tab State
    const [itemClasses, setItemClasses] = useState([])
    const [selectedClass, setSelectedClass] = useState(null)
    const [selectedSubClass, setSelectedSubClass] = useState(null)
    const [items, setItems] = useState([])
    const [loading, setLoading] = useState(false)
    const [tooltipCache, setTooltipCache] = useState({})
    const [tooltipBelow, setTooltipBelow] = useState({})

    // Load Item Classes on mount
    useEffect(() => {
        setLoading(true)
        GetItemClasses()
            .then(classes => {
                setItemClasses(classes || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load item classes:", err)
                setLoading(false)
            })
    }, [])

    // Browse items when class/subclass selected
    useEffect(() => {
        if (selectedClass !== null && selectedSubClass !== null) {
            setLoading(true)
            setItems([])
            BrowseItemsByClass(selectedClass.class, selectedSubClass.subClass)
                .then(res => {
                    setItems(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to browse items:", err)
                    setLoading(false)
                })
        }
    }, [selectedSubClass])

    // Preload tooltips when items change
    useEffect(() => {
        if (items && items.length > 0) {
            items.forEach(item => {
                if (item.entry && !tooltipCache[item.entry]) {
                    loadTooltipData(item.entry)
                }
            })
        }
    }, [items])

    const getQualityColor = (quality) => {
        const colors = {
            0: '#9d9d9d', 1: '#ffffff', 2: '#1eff00', 
            3: '#0070dd', 4: '#a335ee', 5: '#ff8000', 6: '#e6cc80'
        }
        return colors[quality] || '#ffffff'
    }

    const getQualityClass = (quality) => `q${quality || 0}`

    // Load tooltip data
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

    // Check tooltip position
    const checkTooltipPosition = (event, idx) => {
        const element = event.currentTarget
        const container = element.closest('.loot-items')
        if (!element || !container) return
        
        const elementRect = element.getBoundingClientRect()
        const containerRect = container.getBoundingClientRect()
        
        const spaceAbove = elementRect.top - containerRect.top
        const needsBelow = spaceAbove < 320
        
        if (needsBelow !== tooltipBelow[idx]) {
            setTooltipBelow(prev => ({ ...prev, [idx]: needsBelow }))
        }
    }

    // Render tooltip content
    const renderTooltip = (item) => {
        const tooltip = tooltipCache[item.entry]
        
        if (!tooltip) {
            return (
                <div className="item-tooltip">
                    <div className="tooltip-name" style={{color: getQualityColor(item.quality)}}>
                        {item.name || 'Unknown Item'}
                    </div>
                    <div className="tooltip-loading">Loading...</div>
                </div>
            )
        }
        
        return (
            <div className="item-tooltip">
                <div className="tooltip-name" style={{color: getQualityColor(tooltip.quality)}}>
                    {tooltip.name}
                </div>
                
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
                
                {tooltip.durability && (
                    <div className="tooltip-durability">{tooltip.durability}</div>
                )}
                
                {tooltip.requiredLevel > 1 && (
                    <div className="tooltip-reqlevel">Requires Level {tooltip.requiredLevel}</div>
                )}
                
                {tooltip.description && (
                    <div className="tooltip-description">"{tooltip. description}"</div>
                )}
                
                {tooltip.sellPrice && (
                    <div className="tooltip-sellprice">Sell Price: {tooltip.sellPrice}</div>
                )}
            </div>
        )
    }

    return (
        <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* Sub-Tabs */}
            <div style={{ display: 'flex', gap: '2px', background: '#1a1a1a', padding: '5px', borderBottom: '1px solid #333' }}>
                {['Items', 'Sets', 'NPCs'].map(tab => (
                    <button
                        key={tab}
                        onClick={() => setActiveTab(tab.toLowerCase())}
                        style={{
                            padding: '8px 20px',
                            background: activeTab === tab.toLowerCase() ? '#404040' : '#2a2a2a',
                            color: activeTab === tab.toLowerCase() ? '#fff' : '#aaa',
                            border: 'none',
                            cursor: 'pointer',
                            fontWeight: 'bold',
                            borderRadius: '3px'
                        }}
                    >
                        {tab}
                    </button>
                ))}
            </div>

            {/* Content Area */}
            <div className="content" style={{ flex: 1, display: 'grid', gridTemplateColumns: '200px 200px 1fr', gap: 0, overflow: 'hidden' }}>
                
                {/* ITEMS TAB */}
                {activeTab === 'items' && (
                    <>
                        {/* 1. Classes */}
                        <aside className="sidebar">
                            <h2>Item Class</h2>
                            <div className="list">
                                {itemClasses.map(cls => (
                                    <button
                                        key={cls.class}
                                        className={selectedClass?.class === cls.class ? 'active' : ''}
                                        onClick={() => {
                                            setSelectedClass(cls)
                                            setSelectedSubClass(null)
                                            setItems([])
                                        }}
                                    >
                                        {cls.name}
                                    </button>
                                ))}
                            </div>
                        </aside>

                        {/* 2. SubClasses */}
                        <section className="instances">
                            <h2>{selectedClass ? selectedClass.name : 'Select Class'}</h2>
                            <div className="list">
                                {selectedClass?.subClasses?.map(sc => (
                                    <div
                                        key={sc.subClass}
                                        className={`item ${selectedSubClass?.subClass === sc.subClass ? 'active' : ''}`}
                                        onClick={() => setSelectedSubClass(sc)}
                                    >
                                        {sc.name}
                                    </div>
                                ))}
                            </div>
                        </section>

                        {/* 3. Items List */}
                        <section className="loot">
                            <h2>{selectedSubClass ? `${selectedSubClass.name} (${items.length})` : 'Select SubClass'}</h2>
                            {loading && <div className="loading">Loading items...</div>}
                            
                            {items.length > 0 && (
                                <div className="loot-items">
                                    {items.map((item, idx) => (
                                        <div 
                                            key={item.entry} 
                                            className={`loot-item ${tooltipBelow[idx] ? 'tooltip-below' : ''}`}
                                            data-quality={item.quality || 0}
                                            onMouseEnter={(e) => {
                                                checkTooltipPosition(e, idx)
                                                loadTooltipData(item.entry)
                                            }}
                                        >
                                            {item.iconPath ? (
                                                <img 
                                                    className="item-icon"
                                                    src={`/items/icons/${item.iconPath}.jpg`}
                                                    alt={item.name}
                                                    onError={(e) => {
                                                        if (!e.target.src.includes('zamimg.com')) {
                                                            e.target.src = `https://wow.zamimg.com/images/wow/icons/medium/${item.iconPath}.jpg`
                                                        } else {
                                                            e.target.style.display = 'none'
                                                        }
                                                    }}
                                                />
                                            ) : (
                                                <div className="item-icon-placeholder">?</div>
                                            )}
                                            
                                            <span className="item-id">[{item.entry}]</span>
                                            
                                            <span 
                                                className={`item-name ${getQualityClass(item.quality)}`}
                                                style={{color: getQualityColor(item.quality)}}
                                            >
                                                {item.name}
                                            </span>
                                            
                                            {renderTooltip(item)}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </section>
                    </>
                )}

                {/* SETS TAB */}
                {activeTab === 'sets' && (
                    <div style={{ gridColumn: '1 / -1', padding: '20px', color: '#aaa' }}>
                        <h3>Item Sets</h3>
                        <p>Browsing by Item Sets will be implemented soon. Please use "AtlasLoot" tab to find sets for now.</p>
                    </div>
                )}

                {/* NPCS TAB */}
                {activeTab === 'npcs' && (
                    <div style={{ gridColumn: '1 / -1', padding: '20px', color: '#aaa' }}>
                        <h3>NPCs</h3>
                        <p>NPC Database not yet imported.</p>
                    </div>
                )}

            </div>
        </div>
    )
}

export default DatabasePage
