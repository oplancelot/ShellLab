import { useState, useEffect, useMemo } from 'react'
import { FilterInput } from '../../FilterInput'
import { GetSpellSkillCategories, GetSpellSkillsByCategory, GetSpellsBySkill, filterItems } from '../../../utils/databaseApi'

function SpellsTab() {
    const [categories, setCategories] = useState([])
    const [skills, setSkills] = useState([])
    const [spells, setSpells] = useState([])
    const [selectedCategory, setSelectedCategory] = useState(null)
    const [selectedSkill, setSelectedSkill] = useState(null)
    const [loading, setLoading] = useState(false)

    // Independent filter states for each column
    const [categoryFilter, setCategoryFilter] = useState('')
    const [skillFilter, setSkillFilter] = useState('')
    const [spellFilter, setSpellFilter] = useState('')

    // Load categories on mount
    useEffect(() => {
        setLoading(true)
        GetSpellSkillCategories()
            .then(cats => {
                setCategories(cats || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load spell categories:", err)
                setLoading(false)
            })
    }, [])

    // Load skills when category is selected
    useEffect(() => {
        if (selectedCategory) {
            setLoading(true)
            setSkills([])
            setSpells([])
            setSelectedSkill(null)
            GetSpellSkillsByCategory(selectedCategory.id)
                .then(res => {
                    setSkills(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load skills:", err)
                    setLoading(false)
                })
        }
    }, [selectedCategory])

    // Load spells when skill is selected
    useEffect(() => {
        if (selectedSkill) {
            setLoading(true)
            setSpells([])
            GetSpellsBySkill(selectedSkill.id)
                .then(res => {
                    setSpells(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load spells:", err)
                    setLoading(false)
                })
        }
    }, [selectedSkill])

    // Filtered lists
    const filteredCategories = useMemo(() => filterItems(categories, categoryFilter), [categories, categoryFilter])
    const filteredSkills = useMemo(() => filterItems(skills, skillFilter), [skills, skillFilter])
    const filteredSpells = useMemo(() => filterItems(spells, spellFilter), [spells, spellFilter])

    return (
        <>
            {/* 1. Categories */}
            <aside className="sidebar">
                <FilterInput 
                    placeholder="Filter categories..." 
                    onFilterChange={setCategoryFilter}
                    style={{ width: '100%', marginBottom: '8px' }}
                />
                <h2>Categories ({filteredCategories.length})</h2>
                <div className="list">
                    {filteredCategories.map(cat => (
                        <button
                            key={cat.id}
                            className={selectedCategory?.id === cat.id ? 'active' : ''}
                            onClick={() => {
                                setSelectedCategory(cat)
                                setSkillFilter('')
                                setSpellFilter('')
                            }}
                        >
                            {cat.name}
                        </button>
                    ))}
                </div>
            </aside>

            {/* 2. Skills */}
            <section className="instances">
                <FilterInput 
                    placeholder="Filter skills..." 
                    onFilterChange={setSkillFilter}
                    style={{ width: '100%', marginBottom: '8px' }}
                />
                <h2>{selectedCategory ? `${selectedCategory.name} (${filteredSkills.length})` : 'Select Category'}</h2>
                <div className="list">
                    {filteredSkills.map(skill => (
                        <div
                            key={skill.id}
                            className={`item ${selectedSkill?.id === skill.id ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedSkill(skill)
                                setSpellFilter('')
                            }}
                        >
                            {skill.name} ({skill.spellCount})
                        </div>
                    ))}
                </div>
            </section>

            {/* 3. Spells List */}
            <section className="loot" style={{ gridColumn: '3 / -1' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '10px' }}>
                    <h2 style={{ margin: 0, color: '#772ce8' }}>
                        {selectedSkill ? `${selectedSkill.name} (${filteredSpells.length})` : 'Select Skill'}
                    </h2>
                    {spells.length > 0 && (
                        <FilterInput 
                            placeholder="Filter spells..." 
                            onFilterChange={setSpellFilter}
                            style={{ maxWidth: '300px' }}
                        />
                    )}
                </div>

                {loading && selectedSkill && <div className="loading">Loading spells...</div>}
                
                {!selectedSkill && (
                    <p className="placeholder">Select a skill to browse spells.</p>
                )}

                <div className="loot-items">
                    {filteredSpells.map(spell => (
                        <div 
                            key={spell.entry}
                            className="loot-item"
                            style={{ borderLeft: '3px solid #772ce8' }}
                        >
                            <div className="item-icon-placeholder" style={{ 
                                background: '#772ce8',
                                color: '#fff',
                                fontWeight: 'bold',
                                fontSize: '10px'
                            }}>
                                SPL
                            </div>
                            
                            <span className="item-id">[{spell.entry}]</span>
                            
                            <div style={{ display: 'flex', flexDirection: 'column', flex: 1 }}>
                                <span style={{ color: '#772ce8', fontWeight: 'bold' }}>
                                    {spell.name} {spell.subname ? `(${spell.subname})` : ''}
                                </span>
                                {spell.description && (
                                    <span style={{ color: '#888', fontSize: '11px', marginTop: '2px' }}>
                                        {spell.description.length > 100 
                                            ? spell.description.substring(0, 100) + '...' 
                                            : spell.description}
                                    </span>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            </section>
        </>
    )
}

export default SpellsTab
