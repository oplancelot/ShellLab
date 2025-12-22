import { useState, useEffect, useMemo } from 'react'
import { SectionHeader } from '../SectionHeader'
import { GetFactions, filterItems } from '../../../utils/databaseApi'

function FactionsTab() {
    const [factions, setFactions] = useState([])
    const [loading, setLoading] = useState(false)

    // Filter state
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

    // Filtered list
    const filteredFactions = useMemo(() => filterItems(factions, factionFilter), [factions, factionFilter])

    return (
        <section className="loot" style={{ gridColumn: '1 / -1' }}>
            <SectionHeader 
                title={`Reputation Factions (${filteredFactions.length})`}
                placeholder="Filter factions..."
                onFilterChange={setFactionFilter}
                titleColor="#FFD100"
                style={{ borderBottom: '1px solid #505050' }} // Optional: keep border if desired, but SectionHeader handles structure
            />
            
            {loading && <div className="loading">Loading factions...</div>}

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
        </section>
    )
}

export default FactionsTab
