import { useState, useEffect, useMemo } from 'react'
import { GetItemClasses, BrowseItemsByClass } from '../../../../wailsjs/go/main/App'
import { SectionHeader } from '../../common/SectionHeader'
import { BrowseItemsByClassAndSlot, filterItems } from '../../../utils/databaseApi'
import { getQualityColor } from '../../../utils/wow'

// Helper for quality class
const getQualityClass = (quality) => `q${quality || 0}`

function ItemsTab({ tooltipHook }) {
    const [itemClasses, setItemClasses] = useState([])
    const [selectedClass, setSelectedClass] = useState(null)
    const [selectedSubClass, setSelectedSubClass] = useState(null)
    const [selectedSlot, setSelectedSlot] = useState(null)
    const [items, setItems] = useState([])
    const [loading, setLoading] = useState(false)

    // Independent filter states for each column
    const [classFilter, setClassFilter] = useState('')
    const [subClassFilter, setSubClassFilter] = useState('')
    const [slotFilter, setSlotFilter] = useState('')
    const [itemFilter, setItemFilter] = useState('')

    const { setHoveredItem, loadTooltipData, handleItemEnter, handleMouseMove, tooltipCache } = tooltipHook

    // Load item classes on mount
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
                BrowseItemsByClassAndSlot(selectedClass.class, selectedSubClass.subClass, selectedSlot.inventoryType, itemFilter)
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
    }, [selectedSubClass, selectedSlot, itemFilter])

    // Preload tooltips when items change
    useEffect(() => {
        if (items && items.length > 0) {
            items.slice(0, 20).forEach(item => {
                if (item.entry && !tooltipCache[item.entry]) {
                    loadTooltipData(item.entry)
                }
            })
        }
    }, [items])

    // Filtered lists - each column has its own filter
    const filteredClasses = useMemo(() => filterItems(itemClasses, classFilter), [itemClasses, classFilter])
    const filteredSubClasses = useMemo(() => {
        if (!selectedClass?.subClasses) return []
        return filterItems(selectedClass.subClasses, subClassFilter)
    }, [selectedClass, subClassFilter])
    const filteredSlots = useMemo(() => {
        if (!selectedSubClass?.inventorySlots) return []
        return filterItems(selectedSubClass.inventorySlots, slotFilter)
    }, [selectedSubClass, slotFilter])
    // No more frontend filtering for items - backend does it
    const filteredItems = items

    return (
        <>
            {/* 1. Classes */}
            <aside className="sidebar">
                <SectionHeader 
                    title={`Item Class (${filteredClasses.length})`}
                    placeholder="Filter classes..."
                    onFilterChange={setClassFilter}
                />
                <div className="list">
                    {filteredClasses.map(cls => (
                        <button
                            key={cls.class}
                            className={selectedClass?.class === cls.class ? 'active' : ''}
                            onClick={() => {
                                setSelectedClass(cls)
                                setSelectedSubClass(null)
                                setSelectedSlot(null)
                                setItems([])
                                setSubClassFilter('')
                                setSlotFilter('')
                                setItemFilter('')
                            }}
                        >
                            {cls.name}
                        </button>
                    ))}
                </div>
            </aside>

            {/* 2. SubClasses */}
            <section className="instances">
                <SectionHeader 
                    title={selectedClass ? `${selectedClass.name} (${filteredSubClasses.length})` : 'Select Class'}
                    placeholder="Filter types..."
                    onFilterChange={setSubClassFilter}
                />
                <div className="list">
                    {filteredSubClasses.map(sc => (
                        <div
                            key={sc.subClass}
                            className={`item ${selectedSubClass?.subClass === sc.subClass ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedSubClass(sc)
                                setSelectedSlot(null)
                                setSlotFilter('')
                                setItemFilter('')
                            }}
                        >
                            {sc.name}
                        </div>
                    ))}
                </div>
            </section>

            {/* 3. Inventory Slots (Third Level) */}
            <section className="instances">
                <SectionHeader 
                    title={selectedSubClass ? `Slot (${filteredSlots.length})` : 'Select Type'}
                    placeholder="Filter slots..."
                    onFilterChange={setSlotFilter}
                />
                <div className="list">
                    {filteredSlots.map(slot => (
                        <div
                            key={slot.inventoryType}
                            className={`item ${selectedSlot?.inventoryType === slot.inventoryType ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedSlot(slot)
                                setItemFilter('')
                            }}
                        >
                            {slot.name}
                        </div>
                    ))}
                    {selectedSubClass && selectedSubClass.inventorySlots?.length > 1 && (
                        <div
                            className={`item ${selectedSlot === null ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedSlot(null)
                                setItemFilter('')
                            }}
                            style={{ fontStyle: 'italic', color: '#888' }}
                        >
                            All Slots
                        </div>
                    )}
                </div>
            </section>

            {/* 4. Items List */}
            <section className="loot">
                <SectionHeader 
                    title={selectedSubClass ? `${selectedSlot ? selectedSlot.name : selectedSubClass.name} (${filteredItems.length})` : 'Select SubClass'}
                    placeholder="Filter items..."
                    onFilterChange={setItemFilter}
                />
                {loading && <div className="loading">Loading items...</div>}
                
                {items.length > 0 && (
                <div className="loot-items">
                        {filteredItems.map((item, idx) => {
                            const itemId = item.entry || item.id || item.itemId
                            
                            return (
                            <div 
                                key={itemId || idx} 
                                className="loot-item"
                                data-quality={item.quality || 0}
                                onMouseEnter={() => handleItemEnter(itemId)}
                                onMouseMove={(e) => handleMouseMove(e, itemId)}
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
                            </div>
                            )
                        })}
                    </div>
                )}
            </section>
        </>
    )
}

export default ItemsTab
