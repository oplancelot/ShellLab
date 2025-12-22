import { useState, useEffect, useMemo } from 'react'
import { FilterInput } from '../../FilterInput'
import { GetQuestCategories, GetQuestsByCategory, filterItems } from '../../../utils/databaseApi'

function QuestsTab({ onNavigate }) {
    const [questCategories, setQuestCategories] = useState([])
    const [selectedQuestCategory, setSelectedQuestCategory] = useState(null)
    const [quests, setQuests] = useState([])
    const [loading, setLoading] = useState(false)

    // Independent filter states
    const [categoryFilter, setCategoryFilter] = useState('')
    const [questFilter, setQuestFilter] = useState('')

    // Load quest categories on mount
    useEffect(() => {
        setLoading(true)
        GetQuestCategories()
            .then(cats => {
                setQuestCategories(cats || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load quest categories:", err)
                setLoading(false)
            })
    }, [])

    // Load quests when a category is selected
    useEffect(() => {
        if (selectedQuestCategory !== null) {
            setLoading(true)
            setQuests([])
            GetQuestsByCategory(selectedQuestCategory.id)
                .then(res => {
                    setQuests(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load quests:", err)
                    setLoading(false)
                })
        }
    }, [selectedQuestCategory])

    // Filtered lists
    const filteredQuestCategories = useMemo(() => filterItems(questCategories, categoryFilter), [questCategories, categoryFilter])
    const filteredQuests = useMemo(() => filterItems(quests, questFilter), [quests, questFilter])

    return (
        <>
            {/* Quest Categories List */}
            <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                <FilterInput 
                    placeholder="Filter categories..." 
                    onFilterChange={setCategoryFilter}
                    style={{ width: '100%', marginBottom: '8px' }}
                />
                <h2>Categories ({filteredQuestCategories.length})</h2>
                <div className="list">
                    {loading && questCategories.length === 0 && (
                        <div className="loading">Loading categories...</div>
                    )}
                    {filteredQuestCategories.map(cat => (
                        <div
                            key={cat.id}
                            className={`item ${selectedQuestCategory?.id === cat.id ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedQuestCategory(cat)
                                setQuestFilter('')
                            }}
                        >
                            {cat.name} ({cat.count})
                        </div>
                    ))}
                </div>
            </aside>

            {/* Quest List */}
            <section className="loot" style={{ gridColumn: '2 / -1' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '10px' }}>
                    <h2 style={{ margin: 0 }}>
                        {selectedQuestCategory 
                            ? `${selectedQuestCategory.name} (${filteredQuests.length})` 
                            : 'Select a Category'}
                    </h2>
                    {selectedQuestCategory && (
                        <FilterInput 
                            placeholder="Filter quests..." 
                            onFilterChange={setQuestFilter}
                            style={{ maxWidth: '300px' }}
                        />
                    )}
                </div>
                
                {loading && selectedQuestCategory && (
                    <div className="loading">Loading quests...</div>
                )}
                
                {quests.length > 0 && (
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
                )}
                
                {!selectedQuestCategory && (
                    <p className="placeholder">Select a category to browse Quests</p>
                )}
            </section>
        </>
    )
}

export default QuestsTab
