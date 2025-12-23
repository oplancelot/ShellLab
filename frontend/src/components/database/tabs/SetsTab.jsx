import { useState, useEffect, useMemo } from 'react'
import { SidebarPanel, ContentPanel, ScrollList, SectionHeader, ListItem, LootItem } from '../../ui'
import { GetItemSets, GetItemSetDetail, filterItems } from '../../../utils/databaseApi'

function SetsTab({ tooltipHook }) {
    const [itemSets, setItemSets] = useState([])
    const [selectedSet, setSelectedSet] = useState(null)
    const [setDetail, setSetDetail] = useState(null)
    const [loading, setLoading] = useState(false)

    const [setFilter, setSetFilter] = useState('')
    const [itemFilter, setItemFilter] = useState('')

    const { setHoveredItem, loadTooltipData, handleItemEnter, handleMouseMove, tooltipCache } = tooltipHook

    // Load item sets on mount
    useEffect(() => {
        setLoading(true)
        GetItemSets()
            .then(sets => {
                setItemSets(sets || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load item sets:", err)
                setLoading(false)
            })
    }, [])

    // Load set detail when a set is selected
    useEffect(() => {
        if (selectedSet) {
            setLoading(true)
            GetItemSetDetail(selectedSet.itemsetId)
                .then(detail => {
                    setSetDetail(detail)
                    setLoading(false)
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
                    setLoading(false)
                })
        }
    }, [selectedSet])

    const filteredItemSets = useMemo(() => filterItems(itemSets, setFilter), [itemSets, setFilter])
    const filteredSetItems = useMemo(() => {
        if (!setDetail?.items) return []
        return filterItems(setDetail.items, itemFilter)
    }, [setDetail, itemFilter])

    return (
        <>
            {/* Sets List */}
            <SidebarPanel className="col-span-1">
                <SectionHeader 
                    title={`Item Sets (${filteredItemSets.length})`}
                    placeholder="Filter sets..."
                    onFilterChange={setSetFilter}
                />
                <ScrollList>
                    {loading && itemSets.length === 0 && (
                        <div className="p-4 text-center text-wow-gold italic animate-pulse">Loading sets...</div>
                    )}
                    {filteredItemSets.map(set => (
                        <ListItem
                            key={set.itemsetId}
                            active={selectedSet?.itemsetId === set.itemsetId}
                            onClick={() => {
                                setSelectedSet(set)
                                setItemFilter('')
                            }}
                        >
                            <span className="flex justify-between w-full">
                                <span className="truncate">{set.name}</span>
                                <span className="text-gray-600 text-xs">({set.itemCount})</span>
                            </span>
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* Set Details */}
            <ContentPanel className="col-span-3">
                <SectionHeader 
                    title={selectedSet ? `${selectedSet.name} (${filteredSetItems.length})` : 'Select a Set'}
                    placeholder="Filter items..."
                    onFilterChange={setItemFilter}
                />
                
                {loading && selectedSet && (
                    <div className="flex-1 flex items-center justify-center text-wow-gold italic animate-pulse">
                        Loading set details...
                    </div>
                )}
                
                {setDetail && !loading && (
                    <ScrollList className="p-2 space-y-2">
                        {/* Set Items */}
                        <div className="grid grid-cols-1 xl:grid-cols-2 gap-1">
                            {filteredSetItems.map((item, idx) => {
                                const handlers = tooltipHook.getItemHandlers?.(item.entry) || {
                                    onMouseEnter: () => handleItemEnter(item.entry),
                                    onMouseMove: (e) => handleMouseMove(e, item.entry),
                                    onMouseLeave: () => setHoveredItem(null),
                                }
                                
                                return (
                                    <LootItem 
                                        key={item.entry || idx}
                                        item={item}
                                        {...handlers}
                                    />
                                )
                            })}
                        </div>
                        
                        {/* Set Bonuses */}
                        {setDetail.bonuses?.length > 0 && (
                            <div className="mt-4 p-4 bg-bg-main rounded-lg border border-border-dark">
                                <h3 className="text-wow-gold font-bold mb-3 text-sm uppercase tracking-wider">
                                    Set Bonuses
                                </h3>
                                <div className="space-y-2">
                                    {setDetail.bonuses.map((bonus, idx) => (
                                        <div 
                                            key={idx} 
                                            className="text-wow-uncommon text-sm flex items-center gap-2"
                                        >
                                            <span className="bg-wow-uncommon/10 text-wow-uncommon px-2 py-0.5 rounded text-xs font-mono">
                                                {bonus.threshold}pc
                                            </span>
                                            <span>Spell ID: {bonus.spellId}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </ScrollList>
                )}
                
                {!selectedSet && !loading && (
                    <div className="flex-1 flex items-center justify-center text-gray-600 italic">
                        Select an item set to view its items
                    </div>
                )}
            </ContentPanel>
        </>
    )
}

export default SetsTab
