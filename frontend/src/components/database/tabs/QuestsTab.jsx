import { useState, useEffect, useMemo } from 'react'
import { SidebarPanel, ContentPanel, ScrollList, SectionHeader, ListItem, EntityIcon } from '../../ui'
import { GetQuestCategoryGroups, GetQuestCategoriesByGroup, GetQuestsByEnhancedCategory, filterItems } from '../../../utils/databaseApi'

// Quest type badge colors
const getQuestTypeInfo = (type) => {
    const types = {
        1: { label: 'Group', color: '#1eff00' },
        41: { label: 'PvP', color: '#ff8000' },
        62: { label: 'Raid', color: '#a335ee' },
        81: { label: 'Dungeon', color: '#a335ee' },
    }
    return types[type] || null
}

function QuestsTab({ onNavigate }) {
    const [groups, setGroups] = useState([])
    const [categories, setCategories] = useState([])
    const [quests, setQuests] = useState([])
    const [selectedGroup, setSelectedGroup] = useState(null)
    const [selectedCategory, setSelectedCategory] = useState(null)
    const [loading, setLoading] = useState(false)

    const [groupFilter, setGroupFilter] = useState('')
    const [categoryFilter, setCategoryFilter] = useState('')
    const [questFilter, setQuestFilter] = useState('')

    // Load groups on mount
    useEffect(() => {
        setLoading(true)
        GetQuestCategoryGroups()
            .then(res => {
                setGroups(res || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load quest groups:", err)
                setLoading(false)
            })
    }, [])

    // Load categories when group is selected
    useEffect(() => {
        if (selectedGroup) {
            setLoading(true)
            setCategories([])
            setQuests([])
            setSelectedCategory(null)
            GetQuestCategoriesByGroup(selectedGroup.id)
                .then(res => {
                    setCategories(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load categories:", err)
                    setLoading(false)
                })
        }
    }, [selectedGroup])

    // Load quests when category is selected
    useEffect(() => {
        if (selectedCategory) {
            setLoading(true)
            setQuests([])
            GetQuestsByEnhancedCategory(selectedCategory.id, '')
                .then(res => {
                    setQuests(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load quests:", err)
                    setLoading(false)
                })
        }
    }, [selectedCategory])

    const filteredGroups = useMemo(() => filterItems(groups, groupFilter), [groups, groupFilter])
    const filteredCategories = useMemo(() => filterItems(categories, categoryFilter), [categories, categoryFilter])
    const filteredQuests = useMemo(() => filterItems(quests, questFilter), [quests, questFilter])

    return (
        <>
            {/* 1. Groups */}
            <SidebarPanel>
                <SectionHeader 
                    title={`Quest Types (${filteredGroups.length})`}
                    placeholder="Filter groups..."
                    onFilterChange={setGroupFilter}
                />
                <ScrollList>
                    {filteredGroups.map(group => (
                        <ListItem
                            key={group.id}
                            active={selectedGroup?.id === group.id}
                            onClick={() => {
                                setSelectedGroup(group)
                                setCategoryFilter('')
                                setQuestFilter('')
                            }}
                        >
                            {group.name}
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* 2. Categories */}
            <SidebarPanel>
                <SectionHeader 
                    title={selectedGroup ? `${selectedGroup.name} (${filteredCategories.length})` : 'Select Type'}
                    placeholder="Filter zones..."
                    onFilterChange={setCategoryFilter}
                />
                <ScrollList>
                    {filteredCategories.map(cat => (
                        <ListItem
                            key={cat.id}
                            active={selectedCategory?.id === cat.id}
                            onClick={() => {
                                setSelectedCategory(cat)
                                setQuestFilter('')
                            }}
                        >
                            <span className="flex justify-between w-full">
                                <span>{cat.name}</span>
                                <span className="text-gray-600 text-xs">({cat.questCount})</span>
                            </span>
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* 3. Quests List (spans 2 columns) */}
            <ContentPanel className="col-span-2">
                <SectionHeader 
                    title={selectedCategory ? `${selectedCategory.name} (${filteredQuests.length})` : 'Select Category'}
                    placeholder="Filter quests..."
                    onFilterChange={setQuestFilter}
                    titleColor="#FFD100"
                />

                {loading && selectedCategory && (
                    <div className="flex-1 flex items-center justify-center text-wow-gold italic animate-pulse">
                        Loading quests...
                    </div>
                )}
                
                {!selectedCategory && !loading && (
                    <div className="flex-1 flex items-center justify-center text-gray-600 italic">
                        Select a category to browse quests.
                    </div>
                )}

                {!loading && quests.length > 0 && (
                    <ScrollList className="p-2 space-y-1">
                        {filteredQuests.map(quest => {
                            const typeInfo = getQuestTypeInfo(quest.type)
                            
                            return (
                                <div 
                                    key={quest.entry}
                                    onClick={() => onNavigate('quest', quest.entry)}
                                    className="flex items-center gap-3 p-2 bg-white/[0.02] hover:bg-white/5 border-l-[3px] border-wow-gold cursor-pointer transition-colors rounded-r group"
                                >
                                    {/* Level Badge */}
                                    <EntityIcon 
                                        label={quest.questLevel > 0 ? quest.questLevel : '-'}
                                        color="#FFD100"
                                        size="md"
                                    />
                                    
                                    {/* Entry ID */}
                                    <span className="text-gray-600 text-[11px] font-mono min-w-[50px]">
                                        [{quest.entry}]
                                    </span>
                                    
                                    {/* Title */}
                                    <span className="text-wow-gold font-bold flex-1 group-hover:brightness-110 transition-all truncate">
                                        {quest.title}
                                    </span>
                                    
                                    {/* Min Level */}
                                    {quest.minLevel > 0 && (
                                        <span className="text-gray-500 text-xs">
                                            Req Lvl {quest.minLevel}
                                        </span>
                                    )}
                                    
                                    {/* Type Badge */}
                                    {typeInfo && (
                                        <span 
                                            className="px-1.5 py-0.5 rounded text-[10px] uppercase border"
                                            style={{ 
                                                color: typeInfo.color, 
                                                borderColor: `${typeInfo.color}40` 
                                            }}
                                        >
                                            {typeInfo.label}
                                        </span>
                                    )}
                                    
                                    {/* XP */}
                                    <span className="text-gray-500 text-xs font-mono">
                                        XP: <b className="text-gray-400">{quest.rewardXp > 0 ? quest.rewardXp : '-'}</b>
                                    </span>
                                </div>
                            )
                        })}
                    </ScrollList>
                )}
            </ContentPanel>
        </>
    )
}

export default QuestsTab
