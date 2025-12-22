import { useState, useEffect } from 'react'
import { GetItemClasses, BrowseItemsByClass } from '../../wailsjs/go/main/App'
import { useItemTooltip } from '../hooks/useItemTooltip'
import ItemTooltip, { getQualityColor } from './ItemTooltip'

// Direct call to BrowseItemsByClassAndSlot since binding may not be generated yet
const BrowseItemsByClassAndSlot = (classId, subClass, inventoryType) => {
    if (window?.go?.main?.App?.BrowseItemsByClassAndSlot) {
        return window.go.main.App.BrowseItemsByClassAndSlot(classId, subClass, inventoryType)
    }
    return Promise.resolve([])
}

// Direct call to GetItemSets
const GetItemSets = () => {
    if (window?.go?.main?.App?.GetItemSets) {
        return window.go.main.App.GetItemSets()
    }
    return Promise.resolve([])
}

// Direct call to GetItemSetDetail
const GetItemSetDetail = (itemSetId) => {
    if (window?.go?.main?.App?.GetItemSetDetail) {
        return window.go.main.App.GetItemSetDetail(itemSetId)
    }
    return Promise.resolve(null)
}

// Direct call to GetCreatureTypes
const GetCreatureTypes = () => {
    if (window?.go?.main?.App?.GetCreatureTypes) {
        return window.go.main.App.GetCreatureTypes()
    }
    return Promise.resolve([])
}

// Direct call to BrowseCreaturesByType
const BrowseCreaturesByType = (creatureType) => {
    if (window?.go?.main?.App?.BrowseCreaturesByType) {
        return window.go.main.App.BrowseCreaturesByType(creatureType)
    }
    return Promise.resolve([])
}

function DatabasePage() {
    const [activeTab, setActiveTab] = useState('items')
    
    // Items Tab State - Three-level classification
    const [itemClasses, setItemClasses] = useState([])
    const [selectedClass, setSelectedClass] = useState(null)
    const [selectedSubClass, setSelectedSubClass] = useState(null)
    const [selectedSlot, setSelectedSlot] = useState(null)  // Third level: inventory type
    const [items, setItems] = useState([])
    const [loading, setLoading] = useState(false)
    
    // Sets Tab State
    const [itemSets, setItemSets] = useState([])
    const [selectedSet, setSelectedSet] = useState(null)
    const [setDetail, setSetDetail] = useState(null)
    const [setsLoading, setSetsLoading] = useState(false)
    
    // NPCs Tab State
    const [creatureTypes, setCreatureTypes] = useState([])
    const [selectedCreatureType, setSelectedCreatureType] = useState(null)
    const [creatures, setCreatures] = useState([])
    const [npcsLoading, setNpcsLoading] = useState(false)
    
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

    // Browse items when class/subclass/slot selected
    useEffect(() => {
        if (selectedClass !== null && selectedSubClass !== null) {
            setLoading(true)
            setItems([])
            
            // If slot is selected, use three-level filter
            if (selectedSlot !== null) {
                BrowseItemsByClassAndSlot(selectedClass.class, selectedSubClass.subClass, selectedSlot.inventoryType)
                    .then(res => {
                        setItems(res || [])
                        setLoading(false)
                    })
                    .catch(err => {
                        console.error("Failed to browse items by slot:", err)
                        setLoading(false)
                    })
            } else {
                // Otherwise use two-level filter
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
        }
    }, [selectedSubClass, selectedSlot])

    // Load item sets when switching to Sets tab
    useEffect(() => {
        if (activeTab === 'sets' && itemSets.length === 0) {
            setSetsLoading(true)
            GetItemSets()
                .then(sets => {
                    setItemSets(sets || [])
                    setSetsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load item sets:", err)
                    setSetsLoading(false)
                })
        }
    }, [activeTab])

    // Load creature types when switching to NPCs tab
    useEffect(() => {
        if (activeTab === 'npcs' && creatureTypes.length === 0) {
            setNpcsLoading(true)
            GetCreatureTypes()
                .then(types => {
                    setCreatureTypes(types || [])
                    setNpcsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load creature types:", err)
                    setNpcsLoading(false)
                })
        }
    }, [activeTab])

    // Load creatures when a type is selected
    useEffect(() => {
        if (selectedCreatureType !== null) {
            setNpcsLoading(true)
            setCreatures([])
            BrowseCreaturesByType(selectedCreatureType.type)
                .then(res => {
                    setCreatures(res || [])
                    setNpcsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load creatures:", err)
                    setNpcsLoading(false)
                })
        }
    }, [selectedCreatureType])

    // Load set detail when a set is selected
    useEffect(() => {
        if (selectedSet) {
            setSetsLoading(true)
            GetItemSetDetail(selectedSet.itemsetId)
                .then(detail => {
                    setSetDetail(detail)
                    setSetsLoading(false)
                    // Preload tooltips for set items
                    if (detail?.items) {
                        detail.items.forEach(item => {
                            if (item.entry && !tooltipCache[item.entry]) {
                                loadTooltipData(item.entry)
                            }
                        })
                    }
                })
                .catch(err => {
                    console.error("Failed to load set detail:", err)
                    setSetsLoading(false)
                })
        }
    }, [selectedSet])

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

    const getQualityClass = (quality) => `q${quality || 0}`

    // Render tooltip content using shared component
    const renderTooltip = (item) => {
        if (hoveredItem !== item.entry) return null
        
        const tooltip = tooltipCache[item.entry]
        
        return (
            <ItemTooltip
                item={item}
                tooltip={tooltip}
                style={getTooltipStyle()}
                onMouseEnter={() => setHoveredItem(item.entry)}
                onMouseLeave={() => setHoveredItem(null)}
            />
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

            {/* Content Area - 4 columns for three-level classification */}
            <div className="content" style={{ flex: 1, display: 'grid', gridTemplateColumns: '180px 180px 150px 1fr', gap: 0, overflow: 'hidden' }}>
                
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
                                            setSelectedSlot(null)
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
                                        onClick={() => {
                                            setSelectedSubClass(sc)
                                            setSelectedSlot(null)
                                        }}
                                    >
                                        {sc.name}
                                    </div>
                                ))}
                            </div>
                        </section>

                        {/* 3. Inventory Slots (Third Level) */}
                        <section className="instances">
                            <h2>{selectedSubClass ? 'Slot' : 'Select Type'}</h2>
                            <div className="list">
                                {selectedSubClass?.inventorySlots?.map(slot => (
                                    <div
                                        key={slot.inventoryType}
                                        className={`item ${selectedSlot?.inventoryType === slot.inventoryType ? 'active' : ''}`}
                                        onClick={() => setSelectedSlot(slot)}
                                    >
                                        {slot.name}
                                    </div>
                                ))}
                                {selectedSubClass && selectedSubClass.inventorySlots?.length > 1 && (
                                    <div
                                        className={`item ${selectedSlot === null ? 'active' : ''}`}
                                        onClick={() => setSelectedSlot(null)}
                                        style={{ fontStyle: 'italic', color: '#888' }}
                                    >
                                        All Slots
                                    </div>
                                )}
                            </div>
                        </section>

                        {/* 4. Items List */}
                        <section className="loot">
                            <h2>{selectedSubClass ? `${selectedSlot ? selectedSlot.name : selectedSubClass.name} (${items.length})` : 'Select SubClass'}</h2>
                            {loading && <div className="loading">Loading items...</div>}
                            
                            {items.length > 0 && (
                            <div className="loot-items">
                                    {items.map((item, idx) => (
                                        <div 
                                            key={item.entry} 
                                            className="loot-item"
                                            data-quality={item.quality || 0}
                                            onMouseEnter={() => handleItemEnter(item.entry)}
                                            onMouseMove={(e) => handleMouseMove(e, item.entry)}
                                            onMouseLeave={() => setHoveredItem(null)}
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
                    <>
                        {/* Sets List */}
                        <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                            <h2>Item Sets ({itemSets.length})</h2>
                            <div className="list">
                                {setsLoading && itemSets.length === 0 && (
                                    <div className="loading">Loading sets...</div>
                                )}
                                {itemSets.map(set => (
                                    <div
                                        key={set.itemsetId}
                                        className={`item ${selectedSet?.itemsetId === set.itemsetId ? 'active' : ''}`}
                                        onClick={() => setSelectedSet(set)}
                                    >
                                        {set.name} ({set.itemCount})
                                    </div>
                                ))}
                            </div>
                        </aside>

                        {/* Set Details */}
                        <section className="loot" style={{ gridColumn: '2 / -1' }}>
                            <h2>{selectedSet ? selectedSet.name : 'Select a Set'}</h2>
                            
                            {setsLoading && selectedSet && (
                                <div className="loading">Loading set details...</div>
                            )}
                            
                            {setDetail && !setsLoading && (
                                <div className="loot-items">
                                    {setDetail.items?.map((item, idx) => (
                                        <div 
                                            key={item.entry || idx}
                                            className="loot-item"
                                            data-quality={item.quality || 0}
                                            onMouseEnter={() => handleItemEnter(item.entry)}
                                            onMouseMove={(e) => handleMouseMove(e, item.entry)}
                                            onMouseLeave={() => setHoveredItem(null)}
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
                                    
                                    {/* Set Bonuses */}
                                    {setDetail.bonuses && setDetail.bonuses.length > 0 && (
                                        <div style={{ marginTop: '20px', padding: '10px', background: '#1a1a1a', borderRadius: '4px' }}>
                                            <h3 style={{ color: '#ffd100', marginBottom: '10px' }}>Set Bonuses</h3>
                                            {setDetail.bonuses.map((bonus, idx) => (
                                                <div key={idx} style={{ color: '#1eff00', marginBottom: '5px' }}>
                                                    ({bonus.threshold}) Spell ID: {bonus.spellId}
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            )}
                            
                            {!selectedSet && (
                                <p className="placeholder">Select an item set to view its items</p>
                            )}
                        </section>
                    </>
                )}

                {/* NPCS TAB */}
                {activeTab === 'npcs' && (
                    <>
                        {/* Creature Types List */}
                        <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                            <h2>Creature Types</h2>
                            <div className="list">
                                {npcsLoading && creatureTypes.length === 0 && (
                                    <div className="loading">Loading types...</div>
                                )}
                                {creatureTypes.map(type => (
                                    <div
                                        key={type.type}
                                        className={`item ${selectedCreatureType?.type === type.type ? 'active' : ''}`}
                                        onClick={() => setSelectedCreatureType(type)}
                                    >
                                        {type.name} ({type.count})
                                    </div>
                                ))}
                            </div>
                        </aside>

                        {/* Creatures List */}
                        <section className="loot" style={{ gridColumn: '2 / -1' }}>
                            <h2>
                                {selectedCreatureType 
                                    ? `${selectedCreatureType.name} (${creatures.length})` 
                                    : 'Select a Type'}
                            </h2>
                            
                            {npcsLoading && selectedCreatureType && (
                                <div className="loading">Loading creatures...</div>
                            )}
                            
                            {creatures.length > 0 && (
                                <div className="loot-items">
                                    {creatures.map(creature => (
                                        <div 
                                            key={creature.entry}
                                            className="loot-item"
                                            style={{ 
                                                borderLeft: creature.rank >= 3 ? '3px solid #a335ee' 
                                                    : creature.rank >= 1 ? '3px solid #ff8000' 
                                                    : '3px solid #1eff00'
                                            }}
                                        >
                                            <div className="item-icon-placeholder" style={{ 
                                                background: creature.rank >= 3 ? '#a335ee' 
                                                    : creature.rank >= 1 ? '#ff8000' 
                                                    : '#555',
                                                color: '#fff',
                                                fontWeight: 'bold',
                                                fontSize: '10px'
                                            }}>
                                                {creature.levelMin === creature.levelMax 
                                                    ? creature.levelMin 
                                                    : `${creature.levelMin}-${creature.levelMax}`}
                                            </div>
                                            
                                            <span className="item-id">[{creature.entry}]</span>
                                            
                                            <span style={{ 
                                                color: creature.rank >= 3 ? '#a335ee' 
                                                    : creature.rank >= 1 ? '#ff8000' 
                                                    : '#fff',
                                                fontWeight: creature.rank >= 1 ? 'bold' : 'normal'
                                            }}>
                                                {creature.name}
                                                {creature.subname && (
                                                    <span style={{ color: '#888', fontWeight: 'normal', marginLeft: '5px' }}>
                                                        &lt;{creature.subname}&gt;
                                                    </span>
                                                )}
                                            </span>
                                            
                                            <span style={{ 
                                                marginLeft: 'auto', 
                                                color: '#888',
                                                fontSize: '11px'
                                            }}>
                                                {creature.rankName !== 'Normal' && (
                                                    <span style={{ 
                                                        color: creature.rank >= 3 ? '#a335ee' : '#ff8000',
                                                        marginRight: '8px'
                                                    }}>
                                                        [{creature.rankName}]
                                                    </span>
                                                )}
                                                HP: {creature.healthMax.toLocaleString()}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                            )}
                            
                            {!selectedCreatureType && (
                                <p className="placeholder">Select a creature type to browse NPCs</p>
                            )}
                        </section>
                    </>
                )}

            </div>
        </div>
    )
}

export default DatabasePage
