import { useState, useEffect, useMemo } from 'react'
import { GetItemClasses, BrowseItemsByClass } from '../../../../wailsjs/go/main/App'
import { SidebarPanel, ContentPanel, ScrollList, SectionHeader, ListItem, LootItem } from '../../ui'
import { BrowseItemsByClassAndSlot, filterItems } from '../../../utils/databaseApi'

function ItemsTab({ tooltipHook, onNavigate }) {
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
            
            if (selectedSlot !== null) {
                BrowseItemsByClassAndSlot(selectedClass.class, selectedSubClass.subClass, selectedSlot.inventoryType, '')
                    .then(res => {
                        setItems(res || [])
                        setLoading(false)
                    })
                    .catch(err => {
                        console.error("Failed to browse items by slot:", err)
                        setLoading(false)
                    })
            } else {
                BrowseItemsByClass(selectedClass.class, selectedSubClass.subClass, '')
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

    // Preload tooltips when items change
    useEffect(() => {
        if (items?.length > 0) {
            items.slice(0, 20).forEach(item => {
                if (item.entry && !tooltipCache[item.entry]) {
                    loadTooltipData(item.entry)
                }
            })
        }
    }, [items])

    // Filtered lists
    const filteredClasses = useMemo(() => filterItems(itemClasses, classFilter), [itemClasses, classFilter])
    const filteredSubClasses = useMemo(() => {
        if (!selectedClass?.subClasses) return []
        return filterItems(selectedClass.subClasses, subClassFilter)
    }, [selectedClass, subClassFilter])
    const filteredSlots = useMemo(() => {
        if (!selectedSubClass?.inventorySlots) return []
        return filterItems(selectedSubClass.inventorySlots, slotFilter)
    }, [selectedSubClass, slotFilter])
    const filteredItems = useMemo(() => filterItems(items, itemFilter), [items, itemFilter])

    return (
        <>
            {/* 1. Classes */}
            <SidebarPanel>
                <SectionHeader 
                    title={`Item Class (${filteredClasses.length})`}
                    placeholder="Filter classes..."
                    onFilterChange={setClassFilter}
                />
                <ScrollList>
                    {filteredClasses.map(cls => (
                        <ListItem
                            key={cls.class}
                            active={selectedClass?.class === cls.class}
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
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* 2. SubClasses */}
            <SidebarPanel>
                <SectionHeader 
                    title={selectedClass ? `${selectedClass.name} (${filteredSubClasses.length})` : 'Select Class'}
                    placeholder="Filter types..."
                    onFilterChange={setSubClassFilter}
                />
                <ScrollList>
                    {filteredSubClasses.map(sc => (
                        <ListItem
                            key={sc.subClass}
                            active={selectedSubClass?.subClass === sc.subClass}
                            onClick={() => {
                                setSelectedSubClass(sc)
                                setSelectedSlot(null)
                                setSlotFilter('')
                                setItemFilter('')
                            }}
                        >
                            {sc.name}
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* 3. Inventory Slots */}
            <SidebarPanel>
                <SectionHeader 
                    title={selectedSubClass ? `Slot (${filteredSlots.length})` : 'Select Type'}
                    placeholder="Filter slots..."
                    onFilterChange={setSlotFilter}
                />
                <ScrollList>
                    {filteredSlots.map(slot => (
                        <ListItem
                            key={slot.inventoryType}
                            active={selectedSlot?.inventoryType === slot.inventoryType}
                            onClick={() => {
                                setSelectedSlot(slot)
                                setItemFilter('')
                            }}
                        >
                            {slot.name}
                        </ListItem>
                    ))}
                    {selectedSubClass?.inventorySlots?.length > 1 && (
                        <ListItem
                            active={selectedSlot === null}
                            onClick={() => {
                                setSelectedSlot(null)
                                setItemFilter('')
                            }}
                            className="italic text-gray-500"
                        >
                            All Slots
                        </ListItem>
                    )}
                </ScrollList>
            </SidebarPanel>

            {/* 4. Items List */}
            <ContentPanel>
                <SectionHeader 
                    title={selectedSubClass 
                        ? `${selectedSlot ? selectedSlot.name : selectedSubClass.name} (${filteredItems.length})` 
                        : 'Select SubClass'
                    }
                    placeholder="Filter items..."
                    onFilterChange={setItemFilter}
                />
                
                {loading && (
                    <div className="flex-1 flex items-center justify-center text-wow-gold italic animate-pulse">
                        Loading items...
                    </div>
                )}
                
                {!loading && items.length > 0 && (
                    <ScrollList className="grid grid-cols-1 xl:grid-cols-2 gap-1 p-2 auto-rows-min">
                        {filteredItems.map((item, idx) => {
                            const itemId = item.entry || item.id || item.itemId
                            const handlers = tooltipHook.getItemHandlers?.(itemId) || {
                                onMouseEnter: () => handleItemEnter(itemId),
                                onMouseMove: (e) => handleMouseMove(e, itemId),
                                onMouseLeave: () => setHoveredItem(null),
                            }
                            
                            return (
                                <LootItem 
                                    key={itemId || idx}
                                    item={item}
                                    onClick={() => onNavigate && onNavigate('item', itemId)}
                                    {...handlers}
                                />
                            )
                        })}
                    </ScrollList>
                )}
                
                {!loading && items.length === 0 && selectedSubClass && (
                    <div className="flex-1 flex items-center justify-center text-gray-600 italic">
                        No items found
                    </div>
                )}
            </ContentPanel>
        </>
    )
}

export default ItemsTab
