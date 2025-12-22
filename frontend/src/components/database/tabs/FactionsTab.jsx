import { useState, useEffect, useMemo } from 'react'
import { FilterInput } from '../../FilterInput'
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
        <div style={{ gridColumn: '1 / -1', padding: '20px' }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '5px', marginBottom: '10px', minHeight: '60px', justifyContent: 'flex-end', borderBottom: '1px solid #505050', paddingBottom: '10px' }}>
                <h2 style={{ margin: 0, fontSize: '15px', color: '#FFD100' }}>
                    Reputation Factions ({filteredFactions.length})
                </h2>
                <FilterInput 
                    placeholder="Filter factions..." 
                    onFilterChange={setFactionFilter}
                    style={{ width: '100%', maxWidth: '300px' }}
                />
            </div>
            
            {loading && <div className="loading">Loading factions...</div>}

            <div className="loot-items" style={{ marginTop: '20px' }}>
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
        </div>
    )
}

export default FactionsTab
