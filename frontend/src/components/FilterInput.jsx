import { useState, useMemo } from 'react'

/**
 * A reusable filter input component with client-side filtering logic
 * @param {string} placeholder - Placeholder text for the input
 * @param {function} onFilterChange - Callback when filter value changes
 */
export function FilterInput({ placeholder = 'Filter...', onFilterChange, style = {} }) {
    const [value, setValue] = useState('')

    const handleChange = (e) => {
        const newValue = e.target.value
        setValue(newValue)
        if (onFilterChange) {
            onFilterChange(newValue)
        }
    }

    const handleClear = () => {
        setValue('')
        if (onFilterChange) {
            onFilterChange('')
        }
    }

    return (
        <div style={{ 
            display: 'flex', 
            alignItems: 'center',
            background: '#1a1a1a',
            borderRadius: '4px',
            border: '1px solid #444',
            overflow: 'hidden',
            ...style
        }}>
            <span style={{ padding: '0 8px', color: '#666' }}>🔍</span>
            <input 
                type="text"
                value={value}
                onChange={handleChange}
                placeholder={placeholder}
                style={{
                    flex: 1,
                    padding: '6px 8px',
                    background: 'transparent',
                    border: 'none',
                    color: '#fff',
                    fontSize: '13px',
                    outline: 'none',
                    minWidth: '100px'
                }}
            />
            {value && (
                <button
                    onClick={handleClear}
                    style={{
                        padding: '4px 8px',
                        background: 'transparent',
                        border: 'none',
                        color: '#888',
                        cursor: 'pointer',
                        fontSize: '14px'
                    }}
                    title="Clear filter"
                >
                    ✕
                </button>
            )}
        </div>
    )
}

/**
 * Hook to filter an array of items based on a filter string
 * Matches against name, title, or entry/id fields (case-insensitive)
 */
export function useFilter(items, filterText) {
    return useMemo(() => {
        if (!filterText || !filterText.trim()) {
            return items || []
        }
        
        const searchLower = filterText.toLowerCase().trim()
        const searchNum = parseInt(filterText)
        const isNumericSearch = !isNaN(searchNum)
        
        return (items || []).filter(item => {
            // Match by ID/entry
            if (isNumericSearch) {
                if (item.entry === searchNum || item.id === searchNum || item.itemsetId === searchNum) {
                    return true
                }
            }
            
            // Match by name/title (case-insensitive)
            const name = (item.name || item.title || item.displayName || '').toLowerCase()
            if (name.includes(searchLower)) {
                return true
            }
            
            // Match by key (for categories)
            if (item.key && item.key.toLowerCase().includes(searchLower)) {
                return true
            }
            
            return false
        })
    }, [items, filterText])
}

export default FilterInput
