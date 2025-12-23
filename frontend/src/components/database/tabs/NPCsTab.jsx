import { useState, useEffect, useMemo } from 'react'
import { SidebarPanel, ContentPanel, ScrollList, SectionHeader, ListItem, EntityIcon } from '../../ui'
import { GetCreatureTypes, BrowseCreaturesByType, filterItems } from '../../../utils/databaseApi'

// NPC rank colors
const getRankColor = (rank) => {
    if (rank >= 3) return '#a335ee' // Boss - Epic purple
    if (rank >= 1) return '#ff8000' // Elite - Legendary orange
    return '#1eff00' // Normal - Uncommon green
}

function NPCsTab({ onNavigate, tooltipHook }) {
    const [creatureTypes, setCreatureTypes] = useState([])
    const [selectedCreatureType, setSelectedCreatureType] = useState(null)
    const [creatures, setCreatures] = useState([])
    const [loading, setLoading] = useState(false)

    const [typeFilter, setTypeFilter] = useState('')
    const [creatureFilter, setCreatureFilter] = useState('')

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

    const filteredTypes = useMemo(() => filterItems(creatureTypes, typeFilter), [creatureTypes, typeFilter])
    const filteredCreatures = useMemo(() => filterItems(creatures, creatureFilter), [creatures, creatureFilter])

    return (
        <>
            {/* Creature Types (spans 1 column) */}
            <SidebarPanel className="col-span-1">
                <SectionHeader 
                    title={`Creature Types (${filteredTypes.length})`}
                    placeholder="Filter types..."
                    onFilterChange={setTypeFilter}
                />
                <ScrollList>
                    {loading && creatureTypes.length === 0 && (
                        <div className="p-4 text-center text-wow-gold italic animate-pulse">Loading types...</div>
                    )}
                    {filteredTypes.map(type => (
                        <ListItem
                            key={type.type}
                            active={selectedCreatureType?.type === type.type}
                            onClick={() => {
                                setSelectedCreatureType(type)
                                setCreatureFilter('')
                            }}
                        >
                            <span className="flex justify-between w-full">
                                <span>{type.name}</span>
                                <span className="text-gray-600 text-xs">({type.count})</span>
                            </span>
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* Creatures List (spans remaining columns) */}
            <ContentPanel className="col-span-3">
                <SectionHeader 
                    title={selectedCreatureType ? `${selectedCreatureType.name} (${filteredCreatures.length})` : 'Select a Type'}
                    placeholder="Filter NPCs..."
                    onFilterChange={setCreatureFilter}
                />
                
                {loading && selectedCreatureType && (
                    <div className="flex-1 flex items-center justify-center text-wow-gold italic animate-pulse">
                        Loading creatures...
                    </div>
                )}
                
                {!loading && creatures.length > 0 && (
                    <ScrollList className="p-2 space-y-1">
                        {filteredCreatures.map(creature => {
                            const rankColor = getRankColor(creature.rank)
                            const levelText = creature.levelMin === creature.levelMax 
                                ? `${creature.levelMin}` 
                                : `${creature.levelMin}-${creature.levelMax}`
                            
                            return (
                                <div 
                                    key={creature.entry}
                                    onClick={() => onNavigate('npc', creature.entry)}
                                    className="flex items-center gap-3 p-2 bg-white/[0.02] hover:bg-white/5 border-l-[3px] cursor-pointer transition-colors rounded-r group"
                                    style={{ borderLeftColor: rankColor }}
                                >
                                    {/* Level Badge */}
                                    <EntityIcon 
                                        label={levelText}
                                        color={rankColor}
                                        size="md"
                                    />
                                    
                                    {/* Entry ID */}
                                    <span className="text-gray-600 text-[11px] font-mono min-w-[50px]">
                                        [{creature.entry}]
                                    </span>
                                    
                                    {/* Name & Subname */}
                                    <div className="flex-1 min-w-0">
                                        <span 
                                            className="font-bold group-hover:brightness-110 transition-all"
                                            style={{ color: rankColor }}
                                        >
                                            {creature.name}
                                        </span>
                                        {creature.subname && (
                                            <span className="text-gray-500 ml-2 text-sm">
                                                &lt;{creature.subname}&gt;
                                            </span>
                                        )}
                                    </div>
                                    
                                    {/* Stats */}
                                    <div className="flex items-center gap-3 text-gray-500 text-xs ml-auto">
                                        {creature.rankName !== 'Normal' && (
                                            <span 
                                                className="px-1.5 py-0.5 rounded border"
                                                style={{ color: rankColor, borderColor: `${rankColor}40` }}
                                            >
                                                {creature.rankName}
                                            </span>
                                        )}
                                        <span className="font-mono">
                                            HP: <b className="text-gray-400">{creature.healthMax.toLocaleString()}</b>
                                        </span>
                                    </div>
                                </div>
                            )
                        })}
                    </ScrollList>
                )}
                
                {!selectedCreatureType && !loading && (
                    <div className="flex-1 flex items-center justify-center text-gray-600 italic">
                        Select a creature type to browse NPCs
                    </div>
                )}
            </ContentPanel>
        </>
    )
}

export default NPCsTab
