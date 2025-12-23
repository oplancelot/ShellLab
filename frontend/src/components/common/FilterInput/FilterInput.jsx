import { useState, useMemo } from 'react'
import './FilterInput.css'

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
        <div className="filter-input-container" style={style}>
            <div className="filter-icon-wrapper">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
                </svg>
            </div>
            <input 
                type="text"
                className="filter-input-field"
                value={value}
                onChange={handleChange}
                placeholder={placeholder}
            />
            {value && (
                <button
                    className="filter-clear-btn"
                    onClick={handleClear}
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
            
            // Match by subname (for NPCs)
            if (item.subname && item.subname.toLowerCase().includes(searchLower)) {
                return true
            }

            return false
        })
    }, [items, filterText])
}

export default FilterInput
