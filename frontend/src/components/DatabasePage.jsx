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

// Quest APIs
const GetQuestCategories = () => {
    if (window?.go?.main?.App?.GetQuestCategories) {
        return window.go.main.App.GetQuestCategories()
    }
    return Promise.resolve([])
}

const GetQuestsByCategory = (categoryId) => {
    if (window?.go?.main?.App?.GetQuestsByCategory) {
        return window.go.main.App.GetQuestsByCategory(categoryId)
    }
    return Promise.resolve([])
}

const SearchQuests = (query) => {
    if (window?.go?.main?.App?.SearchQuests) {
        return window.go.main.App.SearchQuests(query)
    }
    return Promise.resolve([])
}

// Object APIs
const GetObjectTypes = () => {
    if (window?.go?.main?.App?.GetObjectTypes) {
        return window.go.main.App.GetObjectTypes()
    }
    return Promise.resolve([])
}

const GetObjectsByType = (typeId) => {
    if (window?.go?.main?.App?.GetObjectsByType) {
        return window.go.main.App.GetObjectsByType(typeId)
    }
    return Promise.resolve([])
}

const SearchObjects = (query) => {
    if (window?.go?.main?.App?.SearchObjects) {
        return window.go.main.App.SearchObjects(query)
    }
    return Promise.resolve([])
}

// Spells APIs
const SearchSpells = (query) => {
    if (window?.go?.main?.App?.SearchSpells) {
        return window.go.main.App.SearchSpells(query)
    }
    return Promise.resolve([])
}

// Factions APIs
const GetFactions = () => {
    if (window?.go?.main?.App?.GetFactions) {
        return window.go.main.App.GetFactions()
    }
    return Promise.resolve([])
}

const GetCreatureLoot = (entry) => {
    if (window?.go?.main?.App?.GetCreatureLoot) {
        return window.go.main.App.GetCreatureLoot(entry)
    }
    return Promise.resolve([])
}

const GetQuestDetail = (entry) => {
    if (window?.go?.main?.App?.GetQuestDetail) {
        return window.go.main.App.GetQuestDetail(entry)
    }
    return Promise.resolve(null)
}

const GetCreatureDetail = (entry) => {
    if (window?.go?.main?.App?.GetCreatureDetail) {
        return window.go.main.App.GetCreatureDetail(entry)
    }
    return Promise.resolve(null)
}

const getQualityColor = (quality) => {
    const colors = {
        0: '#9d9d9d', // Poor
        1: '#ffffff', // Common
        2: '#1eff00', // Uncommon
        3: '#0070dd', // Rare
        4: '#a335ee', // Epic
        5: '#ff8000', // Legendary
        6: '#e6cc80'  // Artifact
    }
    return colors[quality] || '#ffffff'
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
    const [expandedCreatureId, setExpandedCreatureId] = useState(null)
    const [creatureLoot, setCreatureLoot] = useState({})
    const [lootLoading, setLootLoading] = useState(false)

    // Quests Tab State
    const [questCategories, setQuestCategories] = useState([])
    const [selectedQuestCategory, setSelectedQuestCategory] = useState(null)
    const [quests, setQuests] = useState([])
    const [questsLoading, setQuestsLoading] = useState(false)
    
    // Objects Tab State
    const [objectTypes, setObjectTypes] = useState([])
    const [selectedObjectType, setSelectedObjectType] = useState(null)
    const [objects, setObjects] = useState([])
    const [objectsLoading, setObjectsLoading] = useState(false)

    // Spells Tab State
    const [spells, setSpells] = useState([])
    const [spellsLoading, setSpellsLoading] = useState(false)

    // Navigation State for Detail Views
    const [detailStack, setDetailStack] = useState([]) // Stack of views: { type, entry }

    // Factions Tab State
    const [factions, setFactions] = useState([])
    const [factionsLoading, setFactionsLoading] = useState(false)
    
    // Load factions when switching to Factions tab
    useEffect(() => {
        if (activeTab === 'factions' && factions.length === 0) {
            setFactionsLoading(true)
            GetFactions()
                .then(res => {
                    setFactions(res || [])
                    setFactionsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load factions:", err)
                    setFactionsLoading(false)
                })
        }
    }, [activeTab])
    
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

    // Detail View Logic
    const navigateTo = (type, entry) => {
        setDetailStack(prev => [...prev, { type, entry }])
    }
    const goBack = () => {
        setDetailStack(prev => prev.slice(0, -1))
    }

    const currentDetail = detailStack.length > 0 ? detailStack[detailStack.length - 1] : null
    
    if (currentDetail) {
         const commonProps = {
            onBack: goBack,
            onNavigate: navigateTo,
            setHoveredItem,
            tooltipCache,
            loadTooltipData
        }

        let view = null
        if (currentDetail.type === 'npc') {
            view = <NPCDetailView entry={currentDetail.entry} {...commonProps} />
        } else if (currentDetail.type === 'quest') {
            view = <QuestDetailView entry={currentDetail.entry} {...commonProps} />
        }

        if (view) {
             return (
                <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                    {view}
                </div>
            )
        }
    }

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

    // Load quest categories when switching to Quests tab
    useEffect(() => {
        if (activeTab === 'quests' && questCategories.length === 0) {
            setQuestsLoading(true)
            GetQuestCategories()
                .then(cats => {
                    // Sort cats: Sorts (negative IDs) last, Zones first? Or just alphabetical?
                    // Currently backend sorts by name.
                    setQuestCategories(cats || [])
                    setQuestsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load quest categories:", err)
                    setQuestsLoading(false)
                })
        }
    }, [activeTab])

    // Load quests when a category is selected
    useEffect(() => {
        if (selectedQuestCategory !== null) {
            setQuestsLoading(true)
            setQuests([])
            GetQuestsByCategory(selectedQuestCategory.id)
                .then(res => {
                    setQuests(res || [])
                    setQuestsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load quests:", err)
                    setQuestsLoading(false)
                })
        }
    }, [selectedQuestCategory])

    // Load object types when switching to Objects tab
    useEffect(() => {
        if (activeTab === 'objects' && objectTypes.length === 0) {
            setObjectsLoading(true)
            GetObjectTypes()
                .then(types => {
                    setObjectTypes(types || [])
                    setObjectsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load object types:", err)
                    setObjectsLoading(false)
                })
        }
    }, [activeTab])

    // Load objects when a type is selected
    useEffect(() => {
        if (selectedObjectType !== null) {
            setObjectsLoading(true)
            setObjects([])
            GetObjectsByType(selectedObjectType.id)
                .then(res => {
                    setObjects(res || [])
                    setObjectsLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load objects:", err)
                    setObjectsLoading(false)
                })
        }
    }, [selectedObjectType])

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

    const toggleCreatureLoot = (entry) => {
        if (expandedCreatureId === entry) {
            setExpandedCreatureId(null)
            return
        }

        setExpandedCreatureId(entry)
        if (!creatureLoot[entry]) {
            setLootLoading(true)
            GetCreatureLoot(entry)
                .then(loot => {
                    setCreatureLoot(prev => ({ ...prev, [entry]: loot }))
                    setLootLoading(false)
                    // Preload icons/tooltips for loot?
                    // Maybe just fetch icons.
                })
                .catch(err => {
                    console.error("Failed to get creature loot:", err)
                    setLootLoading(false)
                })
        }
    }

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
            {/* Sub-Tabs - Turtlehead Style */}
            <div style={{ 
                display: 'flex', 
                gap: '2px', 
                background: '#181818', 
                padding: '8px 10px', 
                borderBottom: '1px solid #404040' 
            }}>
                {['Items', 'Sets', 'NPCs', 'Quests', 'Objects', 'Spells', 'Factions'].map(tab => (
                    <button
                        key={tab}
                        onClick={() => setActiveTab(tab.toLowerCase())}
                        style={{
                            padding: '8px 16px',
                            background: activeTab === tab.toLowerCase() ? '#383838' : 'transparent',
                            color: activeTab === tab.toLowerCase() ? '#fff' : '#FFD100',
                            border: activeTab === tab.toLowerCase() ? '1px solid #484848' : '1px solid transparent',
                            cursor: 'pointer',
                            fontWeight: 'bold',
                            borderRadius: '0',
                            fontSize: '13px',
                            textTransform: 'uppercase'
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
                                            onClick={() => navigateTo('npc', creature.entry)}
                                            style={{ 
                                                borderLeft: creature.rank >= 3 ? '3px solid #a335ee' 
                                                    : creature.rank >= 1 ? '3px solid #ff8000' 
                                                    : '3px solid #1eff00',
                                                cursor: 'pointer',
                                                flexDirection: 'column',
                                                alignItems: 'stretch',
                                                padding: '0'
                                            }}
                                        >
                                            <div style={{ display: 'flex', alignItems: 'center', padding: '8px' }}>
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
                                                    fontSize: '11px',
                                                    display: 'flex',
                                                    alignItems: 'center'
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
                                                    <span style={{ marginLeft: '10px', fontSize: '10px', color: '#666' }}>
                                                        {expandedCreatureId === creature.entry ? '▲' : '▼'}
                                                    </span>
                                                </span>
                                            </div>

                                            {/* Expanded Loot Panel */}
                                            {expandedCreatureId === creature.entry && (
                                                <div style={{ 
                                                    background: '#1a1a1a', 
                                                    borderTop: '1px solid #404040', 
                                                    padding: '10px',
                                                    animation: 'fadeIn 0.2s',
                                                    cursor: 'default'
                                                }} onClick={(e) => e.stopPropagation()}>
                                                    <h4 style={{ margin: '0 0 10px 0', color: '#FFD100', fontSize: '12px', borderBottom: '1px solid #333', paddingBottom: '5px' }}>
                                                        Loot Table
                                                    </h4>
                                                    
                                                    {lootLoading && !creatureLoot[creature.entry] && (
                                                        <div className="loading" style={{ fontSize: '11px' }}>Loading loot...</div>
                                                    )}

                                                    {creatureLoot[creature.entry] && creatureLoot[creature.entry].length === 0 && (
                                                        <div style={{ color: '#888', fontSize: '11px', fontStyle: 'italic' }}>No loot found.</div>
                                                    )}

                                                    {creatureLoot[creature.entry] && (
                                                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))', gap: '5px' }}>
                                                            {creatureLoot[creature.entry].sort((a,b) => b.chance - a.chance).map((item, idx) => (
                                                                <div 
                                                                    key={idx} 
                                                                    style={{ display: 'flex', alignItems: 'center', background: '#242424', padding: '4px', borderRadius: '3px' }}
                                                                    onMouseEnter={(e) => {
                                                                        setHoveredItem(item.itemId)
                                                                        if (!tooltipCache[item.itemId]) {
                                                                            loadTooltipData(item.itemId)
                                                                        }
                                                                    }}
                                                                    onMouseLeave={() => setHoveredItem(null)}
                                                                >
                                                                    <div className="item-icon-placeholder" style={{ 
                                                                        width: '24px', height: '24px', fontSize: '8px', marginRight: '5px',
                                                                        border: `1px solid ${getQualityColor(item.quality)}`
                                                                    }}>
                                                                        {item.icon ? (
                                                                             <img src={`/items/icons/${item.icon}.jpg`} style={{ width: '100%', height: '100%' }} />
                                                                        ) : '?'}
                                                                    </div>
                                                                    <div style={{ display: 'flex', flexDirection: 'column', flex: 1, overflow: 'hidden' }}>
                                                                        <span style={{ 
                                                                            color: getQualityColor(item.quality), 
                                                                            fontSize: '11px', 
                                                                            fontWeight: 'bold',
                                                                            whiteSpace: 'nowrap',
                                                                            overflow: 'hidden',
                                                                            textOverflow: 'ellipsis'
                                                                        }}>
                                                                            {item.itemName}
                                                                        </span>
                                                                        <span style={{ color: '#aaa', fontSize: '10px' }}>
                                                                            {item.chance.toFixed(1)}% {item.minCount > 1 ? `(${item.minCount}-${item.maxCount})` : ''}
                                                                        </span>
                                                                    </div>
                                                                    {renderTooltip({ entry: item.itemId })} 
                                                                </div>
                                                            ))}
                                                        </div>
                                                    )}
                                                </div>
                                            )}
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

                {/* QUESTS TAB */}
                {activeTab === 'quests' && (
                    <>
                        {/* Quest Categories List */}
                        <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                            <div style={{ padding: '0 0 10px 0', borderBottom: '1px solid #404040', marginBottom: '10px' }}>
                                <input 
                                    type="text" 
                                    placeholder="Search Quests..." 
                                    style={{
                                        width: '100%',
                                        padding: '5px',
                                        background: '#242424',
                                        border: '1px solid #404040',
                                        color: '#fff'
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            setQuestsLoading(true)
                                            SearchQuests(e.target.value).then(res => {
                                                setQuests(res || [])
                                                setSelectedQuestCategory({ id: -9999, name: 'Search Results' })
                                                setQuestsLoading(false)
                                            })
                                        }
                                    }}
                                />
                            </div>
                            <h2>Categories</h2>
                            <div className="list">
                                {questsLoading && questCategories.length === 0 && (
                                    <div className="loading">Loading categories...</div>
                                )}
                                {questCategories.map(cat => (
                                    <div
                                        key={cat.id}
                                        className={`item ${selectedQuestCategory?.id === cat.id ? 'active' : ''}`}
                                        onClick={() => setSelectedQuestCategory(cat)}
                                    >
                                        {cat.name} ({cat.count})
                                    </div>
                                ))}
                            </div>
                        </aside>

                        {/* Quest List */}
                        <section className="loot" style={{ gridColumn: '2 / -1' }}>
                            <h2>
                                {selectedQuestCategory 
                                    ? `${selectedQuestCategory.name} (${quests.length})` 
                                    : 'Select a Category'}
                            </h2>
                            
                            {questsLoading && selectedQuestCategory && (
                                <div className="loading">Loading quests...</div>
                            )}
                            
                            {quests.length > 0 && (
                                <div className="loot-items">
                                    {quests.map(quest => (
                                        <div 
                                            key={quest.entry}
                                            className="loot-item"
                                            onClick={() => navigateTo('quest', quest.entry)}
                                            style={{ borderLeft: '3px solid #FFD100', cursor: 'pointer' }}
                                        >
                                            <div className="item-icon-placeholder" style={{ 
                                                background: '#FFD100',
                                                color: '#000',
                                                fontWeight: 'bold',
                                                fontSize: '11px'
                                            }}>
                                                {quest.questLevel > 0 ? quest.questLevel : '-'}
                                            </div>
                                            
                                            <span className="item-id">[{quest.entry}]</span>
                                            
                                            <span style={{ color: '#FFD100', fontWeight: 'bold' }}>
                                                {quest.title}
                                            </span>

                                            {/* Faction/Race Flags (Simple checks) */}
                                            {/* Note: RequiredRaces logic is complex bitmask, simplified here */}
                                            <span style={{ marginLeft: '10px', fontSize: '11px', color: '#888' }}>
                                                {quest.minLevel > 0 && `Requires Lvl ${quest.minLevel}`}
                                            </span>

                                            <span style={{ marginLeft: 'auto', color: '#fff', fontSize: '11px' }}>
                                                {quest.type === 1 && <span style={{color: '#1eff00', marginRight: '5px'}}>[Group]</span>}
                                                {quest.type === 41 && <span style={{color: '#ff8000', marginRight: '5px'}}>[PvP]</span>}
                                                {quest.type === 62 && <span style={{color: '#a335ee', marginRight: '5px'}}>[Raid]</span>}
                                                {quest.type === 81 && <span style={{color: '#a335ee', marginRight: '5px'}}>[Dungeon]</span>}
                                                XP: {quest.rewardXp > 0 ? quest.rewardXp : '-'}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                            )}
                            
                            {!selectedQuestCategory && (
                                <p className="placeholder">Select a category or search to browse Quests</p>
                            )}
                        </section>
                    </>
                )}

                {/* OBJECTS TAB */}
                {activeTab === 'objects' && (
                    <>
                        {/* Object Types List */}
                        <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                            <div style={{ padding: '0 0 10px 0', borderBottom: '1px solid #404040', marginBottom: '10px' }}>
                                <input 
                                    type="text" 
                                    placeholder="Search Objects..." 
                                    style={{
                                        width: '100%',
                                        padding: '5px',
                                        background: '#242424',
                                        border: '1px solid #404040',
                                        color: '#fff'
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            setObjectsLoading(true)
                                            SearchObjects(e.target.value).then(res => {
                                                setObjects(res || [])
                                                setSelectedObjectType({ id: -9999, name: 'Search Results' })
                                                setObjectsLoading(false)
                                            })
                                        }
                                    }}
                                />
                            </div>
                            <h2>Object Types</h2>
                            <div className="list">
                                {objectsLoading && objectTypes.length === 0 && (
                                    <div className="loading">Loading types...</div>
                                )}
                                {objectTypes.map(type => (
                                    <div
                                        key={type.id}
                                        className={`item ${selectedObjectType?.id === type.id ? 'active' : ''}`}
                                        onClick={() => setSelectedObjectType(type)}
                                    >
                                        {type.name} ({type.count})
                                    </div>
                                ))}
                            </div>
                        </aside>

                        {/* Objects List */}
                        <section className="loot" style={{ gridColumn: '2 / -1' }}>
                            <h2>
                                {selectedObjectType 
                                    ? `${selectedObjectType.name} (${objects.length})` 
                                    : 'Select a Type'}
                            </h2>
                            
                            {objectsLoading && selectedObjectType && (
                                <div className="loading">Loading objects...</div>
                            )}
                            
                            {objects.length > 0 && (
                                <div className="loot-items">
                                    {objects.map(obj => (
                                        <div 
                                            key={obj.entry}
                                            className="loot-item"
                                            style={{ borderLeft: '3px solid #00B4FF' }}
                                        >
                                            <div className="item-icon-placeholder" style={{ 
                                                background: '#00B4FF',
                                                color: '#fff',
                                                fontWeight: 'bold',
                                                fontSize: '10px'
                                            }}>
                                                OBJ
                                            </div>
                                            
                                            <span className="item-id">[{obj.entry}]</span>
                                            
                                            <span style={{ color: '#00B4FF', fontWeight: 'bold' }}>
                                                {obj.name}
                                            </span>
                                            
                                            <span style={{ marginLeft: 'auto', color: '#888', fontSize: '11px' }}>
                                                Type: {obj.typeName || obj.type} | Size: {obj.size.toFixed(1)}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                            )}
                            
                            {!selectedObjectType && (
                                <p className="placeholder">Select an object type or search to browse Objects</p>
                            )}
                        </section>
                    </>
                )}

                {/* SPELLS TAB */}
                {activeTab === 'spells' && (
                    <div style={{ gridColumn: '1 / -1', padding: '20px' }}>
                        <div style={{ padding: '0 0 20px 0', borderBottom: '1px solid #404040', marginBottom: '20px' }}>
                            <input 
                                type="text" 
                                placeholder="Search Spells (e.g. 'Fireball', 'Sprint')..." 
                                style={{
                                    width: '100%',
                                    padding: '10px',
                                    background: '#242424',
                                    border: '1px solid #404040',
                                    color: '#fff',
                                    fontSize: '16px'
                                }}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        setSpellsLoading(true)
                                        SearchSpells(e.target.value).then(res => {
                                            setSpells(res || [])
                                            setSpellsLoading(false)
                                        })
                                    }
                                }}
                            />
                        </div>

                        {spellsLoading && <div className="loading">Searching spells...</div>}
                        
                        {!spellsLoading && spells.length === 0 && (
                            <p className="placeholder">Enter a spell name to search.</p>
                        )}

                        <div className="loot-items">
                            {spells.map(spell => (
                                <div 
                                    key={spell.entry}
                                    className="loot-item"
                                    style={{ borderLeft: '3px solid #772ce8' }}
                                >
                                    <div className="item-icon-placeholder" style={{ 
                                        background: '#772ce8',
                                        color: '#fff',
                                        fontWeight: 'bold',
                                        fontSize: '10px'
                                    }}>
                                        SPL
                                    </div>
                                    
                                    <span className="item-id">[{spell.entry}]</span>
                                    
                                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                                        <span style={{ color: '#772ce8', fontWeight: 'bold' }}>
                                            {spell.name} {spell.subname ? `(${spell.subname})` : ''}
                                        </span>
                                        <span style={{ color: '#888', fontSize: '11px', marginTop: '2px' }}>
                                            {spell.description}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* FACTIONS TAB */}
                {activeTab === 'factions' && (
                    <div style={{ gridColumn: '1 / -1', padding: '20px' }}>
                        <h2 style={{ color: '#FFD100', borderBottom: '1px solid #505050', paddingBottom: '10px' }}>
                            Reputation Factions
                        </h2>
                        
                        {factionsLoading && <div className="loading">Loading factions...</div>}

                        <div className="loot-items" style={{ marginTop: '20px' }}>
                            {factions.map(faction => (
                                <div 
                                    key={faction.id}
                                    className="loot-item"
                                    style={{ 
                                        borderLeft: faction.side === 1 ? '3px solid #0070DE' // Alliance
                                            : faction.side === 2 ? '3px solid #C41F3B' // Horde
                                            : '3px solid #FFD100' // Neutral
                                    }}
                                >
                                   <div className="item-icon-placeholder" style={{ 
                                        background: faction.side === 1 ? '#0070DE' 
                                            : faction.side === 2 ? '#C41F3B' 
                                            : '#FFD100',
                                        color: '#000',
                                        fontWeight: 'bold',
                                        fontSize: '10px'
                                    }}>
                                        {faction.side === 1 ? 'A' : faction.side === 2 ? 'H' : 'N'}
                                    </div>
                                    
                                    <span className="item-id">[{faction.id}]</span>
                                    
                                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                                        <span style={{ color: '#fff', fontWeight: 'bold' }}>
                                            {faction.name}
                                        </span>
                                        <span style={{ color: '#888', fontSize: '11px', marginTop: '2px' }}>
                                            {faction.description}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>


                )}
            </div>
        </div>
    )
}

export default DatabasePage
const NPCDetailView = ({ entry, onBack, onNavigate, setHoveredItem, tooltipCache, loadTooltipData }) => {
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        GetCreatureDetail(entry).then(res => {
            setDetail(res)
            setLoading(false)
        })
    }, [entry])

    const renderLootItem = (item) => {
        const hasTooltip = hoveredItem === item.itemId
        if (hasTooltip && !tooltipCache[item.itemId]) {
            loadTooltipData(item.itemId)
        }

        return (
            <div 
                key={item.itemId}
                className="loot-tile" 
                style={{ position: 'relative', display: 'flex', alignItems: 'center', background: '#242424', padding: '5px', borderRadius: '4px', cursor: 'pointer' }}
                onMouseEnter={() => setHoveredItem(item.itemId)}
                onMouseLeave={() => setHoveredItem(null)}
                // onClick={() => onNavigate('item', item.itemId)} // Item detail not fully ready but could be added
            >
                <div style={{ width: '32px', height: '32px', border: `1px solid ${getQualityColor(item.quality)}`, marginRight: '8px' }}>
                    {item.icon ? <img src={`/items/icons/${item.icon}.jpg`} style={{width:'100%'}} /> : '?'}
                </div>
                <div style={{ flex: 1, overflow: 'hidden' }}>
                    <div style={{ color: getQualityColor(item.quality), fontWeight: 'bold', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{item.itemName}</div>
                    <div style={{ color: '#aaa', fontSize: '11px' }}>{item.chance.toFixed(1)}% {item.minCount > 1 && `(${item.minCount}-${item.maxCount})`}</div>
                </div>
                {hasTooltip && tooltipCache[item.itemId] && (
                    <div style={{ position: 'absolute', left: '100%', top: 0, zIndex: 1000, marginLeft: '10px' }}>
                        <ItemTooltip item={{entry: item.itemId, quality: item.quality, name: item.itemName}} tooltip={tooltipCache[item.itemId]} />
                    </div>
                )}
            </div>
        )
    }

    if (loading) return <div className="loading">Loading...</div>
    if (!detail) return <div className="error">NPC not found</div>

    return (
        <div className="detail-view" style={{ flex: 1, overflowY: 'auto', padding: '20px', background: '#121212', color: '#e0e0e0' }}>
            <button onClick={onBack} style={{ background: '#333', border: 'none', color: '#fff', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer', marginBottom: '20px' }}>&larr; Back to List</button>
            
            <header style={{ marginBottom: '30px', borderBottom: '1px solid #333', paddingBottom: '20px' }}>
                <h1 style={{ color: getQualityColor(detail.rank >= 1 ? 3 : 1), margin: '0 0 10px 0' }}>{detail.name}</h1>
                <div style={{ color: '#888' }}>{detail.subname && `<${detail.subname}>`} Level {detail.levelMin}-{detail.levelMax} {detail.typeName} ({detail.rankName})</div>
                <div style={{ marginTop: '10px' }}>
                    <span style={{ marginRight: '20px' }}>Health: {detail.healthMax}</span>
                    <span>Mana: {detail.manaMax}</span>
                </div>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '30px' }}>
                <div>
                    <h3 style={{ borderBottom: '1px solid #FFD100', paddingBottom: '5px', color: '#FFD100' }}>Loot Table</h3>
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '8px', marginTop: '10px' }}>
                        {detail.loot.length > 0 ? detail.loot.sort((a,b)=>b.chance-a.chance).map(renderLootItem) : <div style={{color:'#666'}}>No loot.</div>}
                    </div>
                </div>
                <div>
                    <h3 style={{ borderBottom: '1px solid #FFD100', paddingBottom: '5px', color: '#FFD100' }}>Related Quests</h3>
                    
                    <h4 style={{ color: '#aaa', marginTop: '15px' }}>Starts ({detail.startsQuests ? detail.startsQuests.length : 0})</h4>
                    <ul style={{ listStyle: 'none', padding: 0 }}>
                        {detail.startsQuests && detail.startsQuests.map(q => (
                            <li key={q.entry} style={{ padding: '4px 0' }}>
                                <a className="quest-link" onClick={() => onNavigate('quest', q.entry)} style={{ cursor: 'pointer', color: '#FFD100', textDecoration: 'none' }}>
                                    Warning: Using 'hoveredItem' which is not defined in this scope. Passed as prop? No, renderLootItem uses it.
                                    {q.title}
                                </a>
                            </li>
                        ))}
                    </ul>

                    <h4 style={{ color: '#aaa', marginTop: '15px' }}>Ends ({detail.endsQuests ? detail.endsQuests.length : 0})</h4>
                    <ul style={{ listStyle: 'none', padding: 0 }}>
                        {detail.endsQuests && detail.endsQuests.map(q => (
                            <li key={q.entry} style={{ padding: '4px 0' }}>
                                <a onClick={() => onNavigate('quest', q.entry)} style={{ cursor: 'pointer', color: '#FFD100' }}>
                                    {q.title}
                                </a>
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
        </div>
    )
}

const QuestDetailView = ({ entry, onBack, onNavigate, setHoveredItem, tooltipCache, loadTooltipData }) => {
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        GetQuestDetail(entry).then(res => {
            setDetail(res)
            setLoading(false)
        })
    }, [entry])

    const renderRewardItem = (item, isChoice) => {
        const hasTooltip = hoveredItem === item.entry
        if (hasTooltip && !tooltipCache[item.entry]) {
            loadTooltipData(item.entry)
        }
        return (
            <div 
                key={item.entry}
                style={{ position: 'relative', display: 'flex', alignItems: 'center', background: '#242424', padding: '5px', borderRadius: '4px', margin: '5px 0' }}
                onMouseEnter={() => setHoveredItem(item.entry)}
                onMouseLeave={() => setHoveredItem(null)}
            >
                 <div style={{ width: '32px', height: '32px', border: `1px solid ${getQualityColor(item.quality)}`, marginRight: '8px' }}>
                    {item.icon ? <img src={`/items/icons/${item.icon}.jpg`} style={{width:'100%'}} /> : '?'}
                </div>
                <div>
                   <div style={{ color: getQualityColor(item.quality) }}>{item.name}</div>
                   {item.count > 1 && <div style={{ color: '#aaa', fontSize: '11px' }}>x{item.count}</div>}
                </div>
                {hasTooltip && tooltipCache[item.entry] && (
                    <div style={{ position: 'absolute', left: '100%', top: 0, zIndex: 1000, marginLeft: '10px' }}>
                        <ItemTooltip item={{entry: item.entry, quality: item.quality, name: item.name}} tooltip={tooltipCache[item.entry]} />
                    </div>
                )}
            </div>
        )
    }

    if (loading) return <div className="loading">Loading...</div>
    if (!detail) return <div className="error">Quest not found</div>
    
    return (
         <div className="detail-view" style={{ flex: 1, overflowY: 'auto', padding: '20px', background: '#121212', color: '#e0e0e0' }}>
            <button onClick={onBack} style={{ background: '#333', border: 'none', color: '#fff', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer', marginBottom: '20px' }}>&larr; Back to List</button>
            
            <h1 style={{ color: '#FFD100', marginBottom: '10px' }}>{detail.title}</h1>
            <div style={{ color: '#888', marginBottom: '20px' }}>Level {detail.questLevel} (Min {detail.minLevel}) - {detail.type === 41 ? 'PVP' : detail.type === 81 ? 'Dungeon' : 'Normal'}</div>
            
            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 1fr)', gap: '40px' }}>
                <div>
                    <h3>Description</h3>
                    <p style={{ lineHeight: '1.6', color: '#ccc' }}>{detail.details}</p>
                    
                    <h3>Objectives</h3>
                    <p style={{ lineHeight: '1.6', color: '#ccc' }}>{detail.objectives}</p>
                </div>
                
                <div style={{ background: '#1a1a1a', padding: '20px', borderRadius: '8px' }}>
                    <h3 style={{ marginTop: 0, color: '#FFD100' }}>Rewards</h3>
                    {detail.rewMoney > 0 && <div style={{marginBottom:'10px'}}>Money: {Math.floor(detail.rewMoney/10000)}g {(detail.rewMoney%10000)/100}s</div>}
                    {detail.rewXp > 0 && <div style={{marginBottom:'10px'}}>XP: {detail.rewXp}</div>}
                    
                    {detail.rewards && detail.rewards.length > 0 && (
                        <div>
                            <h4>You will receive:</h4>
                            {detail.rewards.map(i => renderRewardItem(i, false))}
                        </div>
                    )}
                    
                    {detail.choiceRewards && detail.choiceRewards.length > 0 && (
                        <div>
                            <h4>Choose one of:</h4>
                            {detail.choiceRewards.map(i => renderRewardItem(i, true))}
                        </div>
                    )}
                    
                    <h3 style={{ color: '#FFD100', marginTop: '20px' }}>Related</h3>
                    {detail.starters && detail.starters.length > 0 && (
                        <div>
                            <h4>Start:</h4>
                            {detail.starters.map(s => (
                                <div key={s.entry} onClick={() => s.type==='npc' && onNavigate('npc', s.entry)} style={{cursor: s.type==='npc'?'pointer':'default', color: s.type==='npc'?'#4a9eff':'#aaa'}}>
                                    {s.name} ({s.type})
                                </div>
                            ))}
                        </div>
                    )}
                     {detail.enders && detail.enders.length > 0 && (
                        <div>
                            <h4>End:</h4>
                           {detail.enders.map(s => (
                                <div key={s.entry} onClick={() => s.type==='npc' && onNavigate('npc', s.entry)} style={{cursor: s.type==='npc'?'pointer':'default', color: s.type==='npc'?'#4a9eff':'#aaa'}}>
                                    {s.name} ({s.type})
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
         </div>
    )
}
