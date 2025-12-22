import { useState, useEffect, useMemo } from 'react'
import { SectionHeader } from '../SectionHeader'
import { GetFactions, filterItems } from '../../../utils/databaseApi'

const FACTION_GROUPS = [] // Removed

function FactionsTab() {
    const [factions, setFactions] = useState([])
    const [selectedGroup, setSelectedGroup] = useState(null)
    const [loading, setLoading] = useState(false)

    // Filter state
    const [groupFilter, setGroupFilter] = useState('')
    const [factionFilter, setFactionFilter] = useState('')

    // Load factions on mount
    useEffect(() => {
        setLoading(true)
        GetFactions()
            .then(res => {
                setFactions(res || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load factions:", err)
                setLoading(false)
            })
    }, [])

    // Derive Groups from data
    const groups = useMemo(() => {
        if (factions.length === 0) return []
        
        const map = new Map(factions.map(f => [f.id, f]))
        const parentIds = new Set(factions.map(f => f.categoryId).filter(id => id !== 0))
        
        const g = Array.from(parentIds).map(id => {
             const parent = map.get(id)
             return {
                 id: id,
                 name: parent ? parent.name : `Group ${id}`
             }
        }).sort((a,b) => a.name.localeCompare(b.name))
        
        // Add "Others" group for orphans
        const hasOrphans = factions.some(f => f.categoryId === 0 && !parentIds.has(f.id))
        if (hasOrphans) {
            g.push({ id: 0, name: 'Others' })
        }

        return g
    }, [factions])

    // Filter Groups
    const filteredGroups = useMemo(() => filterItems(groups, groupFilter), [groups, groupFilter])

    // Filter Factions based on Selection
    const filteredFactions = useMemo(() => {
        if (!selectedGroup) return []
        
        let subset = []
        if (selectedGroup.id === 0) {
            // Others: categoryId 0 and NOT a parent (simple check: name is not in groups list? No, id check)
            // Re-calculate isParent set or pass it?
            // Expensive to re-calc. Optimize:
            const parentIds = new Set(factions.map(f => f.categoryId))
            subset = factions.filter(f => f.categoryId === 0 && !parentIds.has(f.id))
        } else {
            subset = factions.filter(f => f.categoryId === selectedGroup.id)
        }
        
        return filterItems(subset, factionFilter)
    }, [factions, selectedGroup, factionFilter])

    return (
        <>
            {/* 1. Groups */}
            <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                <SectionHeader 
                    title={`Faction Groups (${filteredGroups.length})`}
                    placeholder="Filter groups..."
                    onFilterChange={setGroupFilter}
                />
                <div className="list">
                    {filteredGroups.map(group => (
                        <button
                            key={group.id}
                            className={selectedGroup?.id === group.id ? 'active' : ''}
                            onClick={() => {
                                setSelectedGroup(group)
                                setFactionFilter('')
                            }}
                        >
                            {group.name}
                        </button>
                    ))}
                </div>
            </aside>

            {/* 2. Factions List */}
            <section className="loot" style={{ gridColumn: '2 / -1' }}>
                <SectionHeader 
                    title={selectedGroup ? `${selectedGroup.name} (${filteredFactions.length})` : 'Select a Group'}
                    placeholder="Filter factions..."
                    onFilterChange={setFactionFilter}
                    titleColor="#FFD100"
                />
                
                {loading && selectedGroup && <div className="loading">Loading factions...</div>}

                {/* Only show list if group selected */}
                {selectedGroup && (
                    <div className="loot-items">
                        {filteredFactions.map(faction => (
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
                                    color: '#fff',
                                    fontWeight: 'bold',
                                    fontSize: '10px'
                                }}>
                                    {faction.side === 1 ? 'A' : faction.side === 2 ? 'H' : 'N'}
                                </div>
                                
                                <span className="item-id">[{faction.id}]</span>
                                
                                <span style={{ 
                                    color: faction.side === 1 ? '#0070DE' 
                                        : faction.side === 2 ? '#C41F3B'
                                        : '#FFD100',
                                    fontWeight: 'bold' 
                                }}>
                                    {faction.name}
                                </span>
                                
                                <span style={{ marginLeft: 'auto', color: '#888', fontSize: '11px' }}>
                                    {faction.sideName || (faction.side === 1 ? 'Alliance' : faction.side === 2 ? 'Horde' : 'Neutral')}
                                </span>
                            </div>
                        ))}
                    </div>
                )}
                
                {!selectedGroup && (
                    <p className="placeholder">Select a faction group to view reputations</p>
                )}
            </section>
        </>
    )
}

export default FactionsTab
