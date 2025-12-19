import { useState } from 'react'
import { AdvancedSearch } from '../../wailsjs/go/main/App'

function SearchPage({ onItemClick }) {
    const [query, setQuery] = useState('')
    const [minLevel, setMinLevel] = useState('')
    const [maxLevel, setMaxLevel] = useState('')
    const [qualityFilter, setQualityFilter] = useState([])
    const [results, setResults] = useState([])
    const [loading, setLoading] = useState(false)
    const [totalCount, setTotalCount] = useState(0)

    const qualities = [
        { id: 0, name: 'Poor', color: '#9d9d9d' },
        { id: 1, name: 'Common', color: '#ffffff' },
        { id: 2, name: 'Uncommon', color: '#1eff00' },
        { id: 3, name: 'Rare', color: '#0070dd' },
        { id: 4, name: 'Epic', color: '#a335ee' },
        { id: 5, name: 'Legendary', color: '#ff8000' },
    ]

    const handleSearch = () => {
        setLoading(true)
        setResults([])
        
        const filter = {
            query: query,
            minLevel: parseInt(minLevel) || 0,
            maxLevel: parseInt(maxLevel) || 0,
            quality: qualityFilter,
            limit: 50,
            offset: 0
        }

        AdvancedSearch(filter)
            .then(res => {
                setResults(res.items || [])
                setTotalCount(res.totalCount || 0)
                setLoading(false)
            })
            .catch(err => {
                console.error("Search failed:", err)
                setLoading(false)
            })
    }

    const toggleQuality = (id) => {
        if (qualityFilter.includes(id)) {
            setQualityFilter(qualityFilter.filter(q => q !== id))
        } else {
            setQualityFilter([...qualityFilter, id])
        }
    }

    const getQualityColor = (quality) => {
        return qualities.find(q => q.id === quality)?.color || '#ffffff'
    }

    return (
        <div className="search-page" style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <div className="search-header" style={{ marginBottom: '20px', padding: '10px', background: '#222', borderRadius: '4px' }}>
                <div style={{ display: 'flex', gap: '10px', marginBottom: '10px' }}>
                    <input 
                        type="text" 
                        value={query} 
                        onChange={e => setQuery(e.target.value)}
                        placeholder="Search items..."
                        style={{ flex: 1, padding: '8px', fontSize: '16px' }}
                        onKeyDown={e => e.key === 'Enter' && handleSearch()}
                    />
                    <button onClick={handleSearch} style={{ padding: '8px 20px', background: '#0070dd', color: 'white', border: 'none', cursor: 'pointer' }}>
                        Search
                    </button>
                </div>
                
                <div className="filters" style={{ display: 'flex', gap: '20px', alignItems: 'center', flexWrap: 'wrap' }}>
                    {/* Level Range */}
                    <div style={{ display: 'flex', alignItems: 'center', gap: '5px' }}>
                        <span>Lvl:</span>
                        <input 
                            type="number" 
                            placeholder="Min" 
                            value={minLevel} 
                            onChange={e => setMinLevel(e.target.value)}
                            style={{ width: '50px', padding: '4px' }}
                        />
                        <span>-</span>
                        <input 
                            type="number" 
                            placeholder="Max" 
                            value={maxLevel} 
                            onChange={e => setMaxLevel(e.target.value)}
                            style={{ width: '50px', padding: '4px' }}
                        />
                    </div>

                    {/* Quality Filter */}
                    <div style={{ display: 'flex', gap: '5px' }}>
                        {qualities.map(q => (
                            <button
                                key={q.id}
                                onClick={() => toggleQuality(q.id)}
                                style={{
                                    padding: '4px 8px',
                                    background: qualityFilter.includes(q.id) ? q.color : '#333',
                                    color: qualityFilter.includes(q.id) ? '#000' : q.color,
                                    border: `1px solid ${q.color}`,
                                    cursor: 'pointer',
                                    opacity: qualityFilter.length === 0 || qualityFilter.includes(q.id) ? 1 : 0.5
                                }}
                            >
                                {q.name}
                            </button>
                        ))}
                    </div>
                </div>
            </div>

            <div className="results-info" style={{ marginBottom: '10px' }}>
                {loading ? 'Searching...' : `Found ${totalCount} items`}
            </div>

            <div className="results-grid" style={{ flex: 1, overflowY: 'auto', display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '8px', alignContent: 'start' }}>
                {results.map(item => (
                    <div 
                        key={item.entry} 
                        className="loot-item"
                        style={{ 
                            display: 'flex', 
                            alignItems: 'center', 
                            background: '#1a1a1a', 
                            padding: '5px', 
                            border: '1px solid #333' 
                        }}
                        onMouseEnter={() => onItemClick(item.entry, true)}
                        onMouseLeave={() => onItemClick(item.entry, false)}
                    >
                        {/* Icon */}
                        <div style={{ width: '32px', height: '32px', marginRight: '8px', background: '#000' }}>
                            {item.iconPath && (
                                <img 
                                    src={`/items/icons/${item.iconPath}.jpg`}
                                    onError={(e) => {
                                        if (!e.target.src.includes('zamimg.com')) {
                                            e.target.src = `https://wow.zamimg.com/images/wow/icons/medium/${item.iconPath}.jpg`
                                        } else {
                                            e.target.style.display = 'none'
                                        }
                                    }}
                                    style={{ width: '100%', height: '100%' }}
                                />
                            )}
                        </div>
                        
                        {/* Name */}
                        <div style={{ flex: 1 }}>
                            <div style={{ color: getQualityColor(item.quality), fontWeight: 'bold' }}>
                                {item.name}
                            </div>
                            <div style={{ fontSize: '0.8em', color: '#888' }}>
                                Level {item.itemLevel} {item.requiredLevel > 0 ? `(Req ${item.requiredLevel})` : ''}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default SearchPage
