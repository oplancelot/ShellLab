import { useState, useMemo } from 'react'
import { FilterInput } from '../../FilterInput'
import { SearchSpells } from '../../../utils/databaseApi'

// Alphabet for quick navigation
const LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('')

function SpellsTab() {
    const [spells, setSpells] = useState([])
    const [loading, setLoading] = useState(false)
    const [searchQuery, setSearchQuery] = useState('')
    const [selectedLetter, setSelectedLetter] = useState(null)
    const [spellFilter, setSpellFilter] = useState('')

    const handleSearch = () => {
        if (!searchQuery.trim()) return
        setLoading(true)
        setSelectedLetter(null)
        SearchSpells(searchQuery)
            .then(res => {
                setSpells(res || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to search spells:", err)
                setLoading(false)
            })
    }

    const handleLetterClick = (letter) => {
        setSelectedLetter(letter)
        setLoading(true)
        setSearchQuery('')
        // Search for spells starting with this letter
        SearchSpells(letter)
            .then(res => {
                // Filter to only those starting with the letter
                const filtered = (res || []).filter(s => 
                    s.name && s.name.toUpperCase().startsWith(letter)
                )
                setSpells(filtered)
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to browse spells:", err)
                setLoading(false)
            })
    }

    // Filter spells by current filter text
    const filteredSpells = useMemo(() => {
        if (!spellFilter.trim()) return spells
        const filter = spellFilter.toLowerCase()
        return spells.filter(spell => 
            (spell.name && spell.name.toLowerCase().includes(filter)) ||
            (spell.subname && spell.subname.toLowerCase().includes(filter)) ||
            (spell.entry && spell.entry.toString().includes(filter))
        )
    }, [spells, spellFilter])

    return (
        <>
            {/* Left sidebar - alphabet navigation */}
            <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                <h2>Browse A-Z</h2>
                <div className="list" style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                    {LETTERS.map(letter => (
                        <button
                            key={letter}
                            className={selectedLetter === letter ? 'active' : ''}
                            onClick={() => handleLetterClick(letter)}
                            style={{
                                width: '32px',
                                height: '32px',
                                padding: 0,
                                fontSize: '14px',
                                fontWeight: 'bold'
                            }}
                        >
                            {letter}
                        </button>
                    ))}
                </div>
                
                <div style={{ marginTop: '20px', borderTop: '1px solid #404040', paddingTop: '15px' }}>
                    <h2>Search</h2>
                    <input 
                        type="text" 
                        placeholder="Search spells..." 
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                        style={{
                            width: '100%',
                            padding: '8px',
                            background: '#242424',
                            border: '1px solid #404040',
                            color: '#fff',
                            fontSize: '13px',
                            marginBottom: '8px'
                        }}
                    />
                    <button 
                        onClick={handleSearch}
                        style={{
                            width: '100%',
                            padding: '8px',
                            background: '#772ce8',
                            color: '#fff',
                            border: 'none',
                            cursor: 'pointer',
                            fontWeight: 'bold'
                        }}
                    >
                        Search
                    </button>
                </div>
            </aside>

            {/* Right - spell list */}
            <section className="loot" style={{ gridColumn: '2 / -1' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '10px' }}>
                    <h2 style={{ margin: 0, color: '#772ce8' }}>
                        {selectedLetter 
                            ? `Spells: ${selectedLetter} (${filteredSpells.length})`
                            : spells.length > 0 
                                ? `Search Results (${filteredSpells.length})`
                                : 'Spells Database'
                        }
                    </h2>
                    {spells.length > 0 && (
                        <FilterInput 
                            placeholder="Filter results..." 
                            onFilterChange={setSpellFilter}
                            style={{ maxWidth: '300px' }}
                        />
                    )}
                </div>

                {loading && <div className="loading">Searching spells...</div>}
                
                {!loading && spells.length === 0 && (
                    <p className="placeholder">Select a letter or search for spells by name.</p>
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
