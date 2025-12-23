import { useState, useEffect, useMemo } from 'react'
import { SectionHeader } from '../../common/SectionHeader'
import { GetCreatureTypes, BrowseCreaturesByType, GetCreatureLoot, filterItems } from '../../../utils/databaseApi'
import { getQualityColor } from '../../../utils/wow'

function NPCsTab({ onNavigate, tooltipHook }) {
    const [creatureTypes, setCreatureTypes] = useState([])
    const [selectedCreatureType, setSelectedCreatureType] = useState(null)
    const [creatures, setCreatures] = useState([])
    const [loading, setLoading] = useState(false)

    // Independent filter states
    const [typeFilter, setTypeFilter] = useState('')
    const [creatureFilter, setCreatureFilter] = useState('')

    const { setHoveredItem, tooltipCache, loadTooltipData, renderTooltip } = tooltipHook

    // Load creature types on mount
    useEffect(() => {
        setLoading(true)
        GetCreatureTypes()
            .then(types => {
                setCreatureTypes(types || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load creature types:", err)
                setLoading(false)
            })
    }, [])

    // Load creatures when a type is selected
    useEffect(() => {
        if (selectedCreatureType !== null) {
            setLoading(true)
            setCreatures([])
            BrowseCreaturesByType(selectedCreatureType.type, '')
                .then(res => {
                    setCreatures(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load creatures:", err)
                    setLoading(false)
                })
        }
    }, [selectedCreatureType])

    // Filtered lists
    const filteredTypes = useMemo(() => filterItems(creatureTypes, typeFilter), [creatureTypes, typeFilter])
    const filteredCreatures = useMemo(() => filterItems(creatures, creatureFilter), [creatures, creatureFilter])

    return (
        <>
            {/* Creature Types List */}
            <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                <SectionHeader 
                    title={`Creature Types (${filteredTypes.length})`}
                    placeholder="Filter types..."
                    onFilterChange={setTypeFilter}
                />
                <div className="list">
                    {loading && creatureTypes.length === 0 && (
                        <div className="loading">Loading types...</div>
                    )}
                    {filteredTypes.map(type => (
                        <div
                            key={type.type}
                            className={`item ${selectedCreatureType?.type === type.type ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedCreatureType(type)
                                setCreatureFilter('')
                            }}
                        >
                            {type.name} ({type.count})
                        </div>
                    ))}
                </div>
            </aside>

            {/* Creatures List */}
            <section className="loot" style={{ gridColumn: '2 / -1' }}>
                <SectionHeader 
                    title={selectedCreatureType ? `${selectedCreatureType.name} (${filteredCreatures.length})` : 'Select a Type'}
                    placeholder="Filter NPCs..."
                    onFilterChange={setCreatureFilter}
                />
                
                {loading && selectedCreatureType && (
                    <div className="loading">Loading creatures...</div>
                )}
                
                {creatures.length > 0 && (
                    <div className="loot-items">
                        {filteredCreatures.map(creature => (
                            <div 
                                key={creature.entry}
                                className="loot-item"
                                onClick={() => onNavigate('npc', creature.entry)}
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
                                    </span>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
                
                {!selectedCreatureType && (
                    <p className="placeholder">Select a creature type to browse NPCs</p>
                )}
            </section>
        </>
    )
}

export default NPCsTab
