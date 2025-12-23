import { useState } from 'react'
import { AdvancedSearch } from '../../wailsjs/go/main/App'
import { getQualityColor } from '../utils/wow'

const QUALITIES = [
    { id: 0, name: 'Poor', color: '#9d9d9d' },
    { id: 1, name: 'Common', color: '#ffffff' },
    { id: 2, name: 'Uncommon', color: '#1eff00' },
    { id: 3, name: 'Rare', color: '#0070dd' },
    { id: 4, name: 'Epic', color: '#a335ee' },
    { id: 5, name: 'Legendary', color: '#ff8000' },
]

function SearchPage({ onItemClick }) {
    const [query, setQuery] = useState('')
    const [minLevel, setMinLevel] = useState('')
    const [maxLevel, setMaxLevel] = useState('')
    const [qualityFilter, setQualityFilter] = useState([])
    const [results, setResults] = useState([])
    const [loading, setLoading] = useState(false)
    const [totalCount, setTotalCount] = useState(0)

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

    return (
        <div className="h-full flex flex-col bg-bg-dark p-4">
            {/* Search Header */}
            <div className="mb-5 p-4 bg-bg-panel rounded-lg border border-border-dark">
                {/* Search Input */}
                <div className="flex gap-3 mb-4">
                    <input 
                        type="text" 
                        value={query} 
                        onChange={e => setQuery(e.target.value)}
                        placeholder="Search items..."
                        className="flex-1 px-4 py-2 bg-bg-main border border-border-dark rounded text-white text-base outline-none focus:border-wow-rare transition-colors"
                        onKeyDown={e => e.key === 'Enter' && handleSearch()}
                    />
                    <button 
                        onClick={handleSearch}
                        className="px-6 py-2 bg-wow-rare text-white rounded font-bold hover:bg-wow-rare/80 transition-colors"
                    >
                        Search
                    </button>
                </div>
                
                {/* Filters */}
                <div className="flex gap-5 items-center flex-wrap">
                    {/* Level Range */}
                    <div className="flex items-center gap-2 text-sm">
                        <span className="text-gray-400">Level:</span>
                        <input 
                            type="number" 
                            placeholder="Min" 
                            value={minLevel} 
                            onChange={e => setMinLevel(e.target.value)}
                            className="w-14 px-2 py-1 bg-bg-main border border-border-dark rounded text-white text-sm outline-none"
                        />
                        <span className="text-gray-600">-</span>
                        <input 
                            type="number" 
                            placeholder="Max" 
                            value={maxLevel} 
                            onChange={e => setMaxLevel(e.target.value)}
                            className="w-14 px-2 py-1 bg-bg-main border border-border-dark rounded text-white text-sm outline-none"
                        />
                    </div>

                    {/* Quality Filter */}
                    <div className="flex gap-1.5">
                        {QUALITIES.map(q => (
                            <button
                                key={q.id}
                                onClick={() => toggleQuality(q.id)}
                                className="px-2 py-1 text-xs font-bold rounded border transition-all"
                                style={{
                                    background: qualityFilter.includes(q.id) ? q.color : 'transparent',
                                    color: qualityFilter.includes(q.id) ? '#000' : q.color,
                                    borderColor: q.color,
                                    opacity: qualityFilter.length === 0 || qualityFilter.includes(q.id) ? 1 : 0.4
                                }}
                            >
                                {q.name}
                            </button>
                        ))}
                    </div>
                </div>
            </div>

            {/* Results Info */}
            <div className="mb-3 text-sm text-gray-400">
                {loading ? (
                    <span className="text-wow-gold animate-pulse">Searching...</span>
                ) : (
                    <span>Found <b className="text-white">{totalCount}</b> items</span>
                )}
            </div>

            {/* Results Grid */}
            <div className="flex-1 overflow-y-auto grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-2 content-start">
                {results.map(item => (
                    <div 
                        key={item.entry} 
                        className="flex items-center bg-bg-panel p-2 border border-border-dark rounded hover:bg-bg-hover transition-colors cursor-pointer"
                        onMouseEnter={() => onItemClick?.(item.entry, true)}
                        onMouseLeave={() => onItemClick?.(item.entry, false)}
                    >
                        {/* Icon */}
                        <div className="w-8 h-8 mr-2 bg-black rounded overflow-hidden flex-shrink-0">
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
                                    className="w-full h-full object-cover"
                                    alt=""
                                />
                            )}
                        </div>
                        
                        {/* Name */}
                        <div className="flex-1 min-w-0">
                            <div 
                                className="font-bold truncate"
                                style={{ color: getQualityColor(item.quality) }}
                            >
                                {item.name}
                            </div>
                            <div className="text-xs text-gray-500">
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
