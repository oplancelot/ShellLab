import React, { useState } from 'react'

/**
 * Section header with title and filter input
 */
export const SectionHeader = ({ 
    title, 
    placeholder = 'Filter...', 
    onFilterChange,
    titleColor,
    className = '' 
}) => {
    const [value, setValue] = useState('')
    
    const handleChange = (e) => {
        const newValue = e.target.value
        setValue(newValue)
        onFilterChange?.(newValue)
    }
    
    const handleClear = () => {
        setValue('')
        onFilterChange?.('')
    }
    
    return (
        <div className={`flex flex-col gap-2 p-3 bg-bg-hover border-b border-border-dark min-h-[70px] justify-end ${className}`}>
            <h3 
                className="m-0 text-xs uppercase font-bold tracking-wider"
                style={{ color: titleColor || '#ffd100' }}
            >
                {title}
            </h3>
            <div className="flex items-center bg-bg-main rounded border border-border-dark overflow-hidden">
                <span className="px-2 text-gray-600 flex items-center">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
                    </svg>
                </span>
                <input 
                    type="text"
                    value={value}
                    onChange={handleChange}
                    placeholder={placeholder}
                    className="flex-1 px-2 py-1.5 bg-transparent border-none text-white text-[13px] outline-none min-w-[80px] placeholder:text-gray-600"
                />
                {value && (
                    <button
                        onClick={handleClear}
                        className="px-2 py-1 bg-transparent border-none text-gray-500 cursor-pointer text-sm hover:text-white transition-colors"
                    >
                        ✕
                    </button>
                )}
            </div>
        </div>
    )
}

export default SectionHeader
