import { useState, useEffect, useMemo } from 'react'
import { FilterInput } from '../../FilterInput'
import { GetItemSets, GetItemSetDetail, filterItems } from '../../../utils/databaseApi'
import { getQualityColor } from '../../../utils/wow'

// Helper for quality class
const getQualityClass = (quality) => `q${quality || 0}`

function SetsTab({ tooltipHook }) {
    const [itemSets, setItemSets] = useState([])
    const [selectedSet, setSelectedSet] = useState(null)
    const [setDetail, setSetDetail] = useState(null)
    const [loading, setLoading] = useState(false)

    // Independent filter states
    const [setFilter, setSetFilter] = useState('')
    const [itemFilter, setItemFilter] = useState('')

    const { setHoveredItem, loadTooltipData, handleItemEnter, handleMouseMove, renderTooltip, tooltipCache } = tooltipHook

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

    // Filtered lists
    const filteredItemSets = useMemo(() => filterItems(itemSets, setFilter), [itemSets, setFilter])
    const filteredSetItems = useMemo(() => {
        if (!setDetail?.items) return []
        return filterItems(setDetail.items, itemFilter)
    }, [setDetail, itemFilter])

    return (
        <>
            {/* Sets List */}
            <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                <FilterInput 
                    placeholder="Filter sets..." 
                    onFilterChange={setSetFilter}
                    style={{ width: '100%', marginBottom: '8px' }}
                />
                <h2>Item Sets ({filteredItemSets.length})</h2>
                <div className="list">
                    {loading && itemSets.length === 0 && (
                        <div className="loading">Loading sets...</div>
                    )}
                    {filteredItemSets.map(set => (
                        <div
                            key={set.itemsetId}
                            className={`item ${selectedSet?.itemsetId === set.itemsetId ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedSet(set)
                                setItemFilter('')
                            }}
                        >
                            {set.name} ({set.itemCount})
                        </div>
                    ))}
                </div>
            </aside>

            {/* Set Details */}
            <section className="loot" style={{ gridColumn: '2 / -1' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '10px' }}>
                    <h2 style={{ margin: 0 }}>{selectedSet ? `${selectedSet.name} (${filteredSetItems.length})` : 'Select a Set'}</h2>
                    {selectedSet && setDetail && (
                        <FilterInput 
                            placeholder="Filter items..." 
                            onFilterChange={setItemFilter}
                            style={{ maxWidth: '300px' }}
                        />
                    )}
                </div>
                
                {loading && selectedSet && (
                    <div className="loading">Loading set details...</div>
                )}
                
                {setDetail && !loading && (
                    <div className="loot-items">
                        {filteredSetItems.map((item, idx) => (
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
    )
}

export default SetsTab
