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
        <div 
            className="flex items-center bg-bg-main rounded border border-border-dark overflow-hidden transition-colors focus-within:border-border-light"
            style={style}
        >
            <div className="px-2 text-gray-600 flex items-center">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
                </svg>
            </div>
            <input 
                type="text"
                className="flex-1 px-2 py-1.5 bg-transparent border-none text-white text-[13px] outline-none min-w-[80px] placeholder:text-gray-600"
                value={value}
                onChange={handleChange}
                placeholder={placeholder}
            />
            {value && (
                <button
                    className="px-2 py-1 bg-transparent border-none text-gray-500 cursor-pointer text-sm hover:text-white transition-colors"
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
