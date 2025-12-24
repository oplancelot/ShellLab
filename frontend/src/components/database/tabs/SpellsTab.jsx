import { useState, useEffect, useMemo } from 'react'
import { SidebarPanel, ContentPanel, ScrollList, SectionHeader, ListItem, EntityIcon } from '../../ui'
import { GetSpellSkillCategories, GetSpellSkillsByCategory, GetSpellsBySkill, filterItems } from '../../../utils/databaseApi'

const SPELL_COLOR = '#772ce8'

function SpellsTab() {
    const [categories, setCategories] = useState([])
    const [skills, setSkills] = useState([])
    const [spells, setSpells] = useState([])
    const [selectedCategory, setSelectedCategory] = useState(null)
    const [selectedSkill, setSelectedSkill] = useState(null)
    const [loading, setLoading] = useState(false)

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
            GetSpellsBySkill(selectedSkill.id, '')
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

    const filteredCategories = useMemo(() => filterItems(categories, categoryFilter), [categories, categoryFilter])
    const filteredSkills = useMemo(() => filterItems(skills, skillFilter), [skills, skillFilter])
    const filteredSpells = useMemo(() => filterItems(spells, spellFilter), [spells, spellFilter])

    return (
        <>
            {/* 1. Categories */}
            <SidebarPanel>
                <SectionHeader 
                    title={`Categories (${filteredCategories.length})`}
                    placeholder="Filter categories..."
                    onFilterChange={setCategoryFilter}
                />
                <ScrollList>
                    {filteredCategories.map(cat => (
                        <ListItem
                            key={cat.id}
                            active={selectedCategory?.id === cat.id}
                            onClick={() => {
                                setSelectedCategory(cat)
                                setSkillFilter('')
                                setSpellFilter('')
                            }}
                        >
                            {cat.name}
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* 2. Skills */}
            <SidebarPanel>
                <SectionHeader 
                    title={selectedCategory ? `${selectedCategory.name} (${filteredSkills.length})` : 'Select Category'}
                    placeholder="Filter skills..."
                    onFilterChange={setSkillFilter}
                />
                <ScrollList>
                    {filteredSkills.map(skill => (
                        <ListItem
                            key={skill.id}
                            active={selectedSkill?.id === skill.id}
                            onClick={() => {
                                setSelectedSkill(skill)
                                setSpellFilter('')
                            }}
                        >
                            <span className="flex justify-between w-full">
                                <span>{skill.name}</span>
                                <span className="text-gray-600 text-xs">({skill.spellCount})</span>
                            </span>
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* 3. Spells List */}
            <ContentPanel className="col-span-2">
                <SectionHeader 
                    title={selectedSkill ? `${selectedSkill.name} (${filteredSpells.length})` : 'Select Skill'}
                    placeholder="Filter spells..."
                    onFilterChange={setSpellFilter}
                    titleColor={SPELL_COLOR}
                />

                {loading && selectedSkill && (
                    <div className="flex-1 flex items-center justify-center italic animate-pulse" style={{ color: SPELL_COLOR }}>
                        Loading spells...
                    </div>
                )}
                
                {!selectedSkill && !loading && (
                    <div className="flex-1 flex items-center justify-center text-gray-600 italic">
                        Select a skill to browse spells.
                    </div>
                )}

                {!loading && spells.length > 0 && (
                    <ScrollList className="p-2 space-y-1">
                        {filteredSpells.map(spell => (
                            <div 
                                key={spell.entry}
                                className="flex items-center gap-3 p-2 bg-white/[0.02] hover:bg-white/5 border-l-[3px] transition-colors rounded-r"
                                style={{ borderLeftColor: SPELL_COLOR }}
                            >
                                {spell.icon ? (
                                    <div 
                                        className="w-8 h-8 rounded border flex-shrink-0 bg-black/40 flex items-center justify-center overflow-hidden"
                                        style={{ borderColor: SPELL_COLOR }}
                                    >
                                        <img 
                                            src={`/items/icons/${spell.icon.toLowerCase()}.jpg`}
                                            alt=""
                                            className="w-full h-full object-cover"
                                            onError={(e) => {
                                                e.target.style.display = 'none'
                                                e.target.parentNode.innerHTML = `<div style="width:100%;height:100%;background-color:${SPELL_COLOR};display:flex;align-items:center;justify-content:center;color:white;font-weight:bold;font-size:11px">SPL</div>`
                                            }}
                                        />
                                    </div>
                                ) : (
                                    <EntityIcon 
                                        label="SPL"
                                        color={SPELL_COLOR}
                                        size="md"
                                    />
                                )}
                                
                                <span className="text-gray-600 text-[11px] font-mono min-w-[50px]">
                                    [{spell.entry}]
                                </span>
                                
                                <div className="flex flex-col flex-1 min-w-0">
                                    <span 
                                        className="font-bold truncate"
                                        style={{ color: SPELL_COLOR }}
                                    >
                                        {spell.name} {spell.subname ? `(${spell.subname})` : ''}
                                    </span>
                                    {spell.description && (
                                        <span className="text-gray-500 text-xs truncate mt-0.5">
                                            {spell.description.length > 100 
                                                ? spell.description.substring(0, 100) + '...' 
                                                : spell.description}
                                        </span>
                                    )}
                                </div>
                            </div>
                        ))}
                    </ScrollList>
                )}
            </ContentPanel>
        </>
    )
}

export default SpellsTab
