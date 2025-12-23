import React from 'react'
import { getQualityColor } from '../../../utils/wow'

export const LootTile = ({ item, onClick, tooltipHandlers, subtext }) => {
    const name = item.name || item.itemName
    const quality = item.quality
    const icon = item.icon
    
    const renderSubtext = () => {
        if (subtext) return subtext
        
        const parts = []
        if (item.count > 1) parts.push(`x${item.count}`)
        if (item.chance !== undefined) {
             parts.push(`${item.chance.toFixed(1)}%`)
        }
        if ((item.minCount !== undefined && item.maxCount !== undefined) && (item.minCount > 1 || item.maxCount > 1)) {
            parts.push(`(${item.minCount}-${item.maxCount})`)
        }
        
        return parts.join(' ')
    }

    return (
        <div 
            className="relative flex items-center bg-bg-panel p-1.5 rounded cursor-pointer transition-colors hover:bg-bg-hover select-none gap-2"
            onClick={onClick}
            {...tooltipHandlers}
        >
            <div 
                className="w-8 h-8 bg-black/40 flex-shrink-0 rounded overflow-hidden flex items-center justify-center"
                style={{ border: `1px solid ${getQualityColor(quality)}` }}
            >
                {icon ? (
                    <img 
                        src={`/items/icons/${icon}.jpg`} 
                        className="w-full h-full object-cover" 
                        alt="" 
                    />
                ) : '?'}
            </div>
            <div className="flex-1 overflow-hidden min-w-0">
                <div 
                    className="font-bold text-[13px] truncate"
                    style={{ color: getQualityColor(quality) }}
                >
                    {name}
                </div>
                <div className="text-gray-500 text-[11px] mt-0.5">
                    {renderSubtext()}
                </div>
            </div>
        </div>
    )
}

export default LootTile
