import { useState, useEffect, useMemo } from 'react'
import { FilterInput } from '../../FilterInput'
import { GetQuestCategoryGroups, GetQuestCategoriesByGroup, GetQuestsByEnhancedCategory, filterItems } from '../../../utils/databaseApi'

function QuestsTab({ onNavigate }) {
    const [groups, setGroups] = useState([])
    const [categories, setCategories] = useState([])
    const [quests, setQuests] = useState([])
    const [selectedGroup, setSelectedGroup] = useState(null)
    const [selectedCategory, setSelectedCategory] = useState(null)
    const [loading, setLoading] = useState(false)

    // Independent filter states for each column
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
            GetQuestsByEnhancedCategory(selectedCategory.id)
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

    // Filtered lists
    const filteredGroups = useMemo(() => filterItems(groups, groupFilter), [groups, groupFilter])
    const filteredCategories = useMemo(() => filterItems(categories, categoryFilter), [categories, categoryFilter])
    const filteredQuests = useMemo(() => filterItems(quests, questFilter), [quests, questFilter])

    return (
        <>
            {/* 1. Groups */}
            <aside className="sidebar">
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px', marginBottom: '10px', minHeight: '60px', justifyContent: 'flex-end' }}>
                    <h2 style={{ margin: 0, fontSize: '15px' }}>Quest Types ({filteredGroups.length})</h2>
                    <FilterInput 
                        placeholder="Filter groups..." 
                        onFilterChange={setGroupFilter}
                        style={{ width: '100%' }}
                    />
                </div>
                <div className="list">
                    {filteredGroups.map(group => (
                        <button
                            key={group.id}
                            className={selectedGroup?.id === group.id ? 'active' : ''}
                            onClick={() => {
                                setSelectedGroup(group)
                                setCategoryFilter('')
                                setQuestFilter('')
                            }}
                        >
                            {group.name}
                        </button>
                    ))}
                </div>
            </aside>

            {/* 2. Categories */}
            <section className="instances">
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px', marginBottom: '10px', minHeight: '60px', justifyContent: 'flex-end' }}>
                    <h2 style={{ margin: 0, fontSize: '15px' }}>{selectedGroup ? `${selectedGroup.name} (${filteredCategories.length})` : 'Select Type'}</h2>
                    <FilterInput 
                        placeholder="Filter zones/categories..." 
                        onFilterChange={setCategoryFilter}
                        style={{ width: '100%' }}
                    />
                </div>
                <div className="list">
                    {filteredCategories.map(cat => (
                        <div
                            key={cat.id}
                            className={`item ${selectedCategory?.id === cat.id ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedCategory(cat)
                                setQuestFilter('')
                            }}
                        >
                            {cat.name} ({cat.questCount})
                        </div>
                    ))}
                </div>
            </section>

            {/* 3. Quests List */}
            <section className="loot" style={{ gridColumn: '3 / -1' }}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px', marginBottom: '10px', minHeight: '60px', justifyContent: 'flex-end' }}>
                    <h2 style={{ margin: 0, fontSize: '15px', color: '#FFD100' }}>
                        {selectedCategory ? `${selectedCategory.name} (${filteredQuests.length})` : 'Select Category'}
                    </h2>
                    {quests.length > 0 && (
                        <FilterInput 
                            placeholder="Filter quests..." 
                            onFilterChange={setQuestFilter}
                            style={{ width: '100%' }}
                        />
                    )}
                </div>

                {loading && selectedCategory && <div className="loading">Loading quests...</div>}
                
                {!selectedCategory && (
                    <p className="placeholder">Select a category to browse quests.</p>
                )}

                <div className="loot-items">
                    {filteredQuests.map(quest => (
                        <div 
                            key={quest.entry}
                            className="loot-item"
                            onClick={() => onNavigate('quest', quest.entry)}
                            style={{ borderLeft: '3px solid #FFD100', cursor: 'pointer' }}
                        >
                            <div className="item-icon-placeholder" style={{ 
                                background: '#FFD100',
                                color: '#000',
                                fontWeight: 'bold',
                                fontSize: '11px'
                            }}>
                                {quest.questLevel > 0 ? quest.questLevel : '-'}
                            </div>
                            
                            <span className="item-id">[{quest.entry}]</span>
                            
                            <span style={{ color: '#FFD100', fontWeight: 'bold' }}>
                                {quest.title}
                            </span>

                            <span style={{ marginLeft: '10px', fontSize: '11px', color: '#888' }}>
                                {quest.minLevel > 0 && `Requires Lvl ${quest.minLevel}`}
                            </span>

                            <span style={{ marginLeft: 'auto', color: '#fff', fontSize: '11px' }}>
                                {quest.type === 1 && <span style={{color: '#1eff00', marginRight: '5px'}}>[Group]</span>}
                                {quest.type === 41 && <span style={{color: '#ff8000', marginRight: '5px'}}>[PvP]</span>}
                                {quest.type === 62 && <span style={{color: '#a335ee', marginRight: '5px'}}>[Raid]</span>}
                                {quest.type === 81 && <span style={{color: '#a335ee', marginRight: '5px'}}>[Dungeon]</span>}
                                XP: {quest.rewardXp > 0 ? quest.rewardXp : '-'}
                            </span>
                        </div>
                    ))}
                </div>
            </section>
        </>
    )
}

export default QuestsTab
