import React from 'react'
import { FilterInput } from './FilterInput/FilterInput'

/**
 * Standardized header for database tab sections.
 * Ensures uniform height (stacks Title above Filter) across all columns.
 */
export const SectionHeader = ({ title, placeholder, onFilterChange, titleColor, style }) => {
    return (
        <div style={{ 
            display: 'flex', 
            flexDirection: 'column', 
            gap: '5px', 
            marginBottom: '10px', 
            minHeight: '60px', 
            justifyContent: 'flex-end',
            padding: '10px 15px',
            ...style 
        }}>
            <h2 style={{ 
                margin: 0, 
                padding: 0,
                fontSize: '15px',
                color: titleColor || 'var(--text-primary)' 
            }}>
                {title}
            </h2>
            <FilterInput 
                placeholder={placeholder || 'Filter...'} 
                onFilterChange={onFilterChange}
                style={{ width: '100%' }}
            />
        </div>
    )
}

export default SectionHeader
