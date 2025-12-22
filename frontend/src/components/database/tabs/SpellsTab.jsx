import { useState } from 'react'
import { SearchSpells } from '../../../utils/databaseApi'

function SpellsTab() {
    const [spells, setSpells] = useState([])
    const [loading, setLoading] = useState(false)
    const [searchQuery, setSearchQuery] = useState('')

    const handleSearch = () => {
        if (!searchQuery.trim()) return
        setLoading(true)
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

    return (
        <div style={{ gridColumn: '1 / -1', padding: '20px' }}>
            <h2 style={{ color: '#772ce8', borderBottom: '1px solid #505050', paddingBottom: '10px' }}>
                Spells Database
            </h2>
            
            <div style={{ display: 'flex', gap: '10px', marginBottom: '20px' }}>
                <input 
                    type="text" 
                    placeholder="Search spells by name..." 
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                    style={{
                        flex: 1,
                        maxWidth: '400px',
                        padding: '8px',
                        background: '#242424',
                        border: '1px solid #404040',
                        color: '#fff',
                        fontSize: '14px'
                    }}
                />
                <button 
                    onClick={handleSearch}
                    style={{
                        padding: '8px 20px',
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

            {loading && <div className="loading">Searching spells...</div>}
            
            {!loading && spells.length === 0 && (
                <p className="placeholder">Enter a spell name to search.</p>
            )}

            <div className="loot-items">
                {spells.map(spell => (
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
                        
                        <div style={{ display: 'flex', flexDirection: 'column' }}>
                            <span style={{ color: '#772ce8', fontWeight: 'bold' }}>
                                {spell.name} {spell.subname ? `(${spell.subname})` : ''}
                            </span>
                            <span style={{ color: '#888', fontSize: '11px', marginTop: '2px' }}>
                                {spell.description}
                            </span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default SpellsTab
