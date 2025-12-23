import React from 'react'
import { getQualityColor } from '../../../utils/wow'
import './LootTile.css'

export const LootTile = ({ item, onClick, tooltipHandlers, subtext }) => {
    // Normalize properties
    const name = item.name || item.itemName
    const quality = item.quality
    const icon = item.icon
    
    // Default subtext logic if not provided
    const renderSubtext = () => {
        if (subtext) return subtext
        
        const parts = []
        if (item.count > 1) parts.push(`x${item.count}`)
        if (item.chance !== undefined) {
             // Handle 0 chance or formatting
             parts.push(`${item.chance.toFixed(1)}%`)
        }
        if ((item.minCount !== undefined && item.maxCount !== undefined) && (item.minCount > 1 || item.maxCount > 1)) {
            parts.push(`(${item.minCount}-${item.maxCount})`)
        }
        
        return parts.join(' ')
    }

    return (
        <div 
            className="loot-tile" 
            onClick={onClick}
            {...tooltipHandlers}
        >
            <div className="item-icon-frame" style={{ border: `1px solid ${getQualityColor(quality)}` }}>
                {icon ? <img src={`/items/icons/${icon}.jpg`} className="item-icon-img" alt="" /> : '?'}
            </div>
            <div className="item-info">
                <div className="item-name" style={{ color: getQualityColor(quality) }}>
                    {name}
                </div>
                <div className="item-subtext">
                    {renderSubtext()}
                </div>
            </div>
        </div>
    )
}

export default LootTile
