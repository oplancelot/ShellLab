import React from 'react'
import { getQualityColor } from '../../utils/wow'

/**
 * Loot item display with icon, name, and quality color
 */
export const LootItem = ({ 
    item,
    onClick,
    onMouseEnter,
    onMouseMove,
    onMouseLeave,
    showDropChance = false,
    className = '' 
}) => {
    const itemId = item.entry || item.itemId || item.id
    const quality = item.quality || 0
    const qualityColor = getQualityColor(quality)
    const iconName = item.iconPath || item.iconName
    
    return (
        <div 
            className={`
                flex items-center gap-2 p-1.5 
                bg-white/[0.02] hover:bg-white/5 
                border border-white/5 rounded 
                transition-all cursor-pointer group
                ${className}
            `}
            data-quality={quality}
            onClick={onClick}
            onMouseEnter={onMouseEnter}
            onMouseMove={onMouseMove}
            onMouseLeave={onMouseLeave}
        >
            {/* Icon */}
            <div 
                className="w-8 h-8 rounded border flex-shrink-0 bg-black/40 flex items-center justify-center overflow-hidden"
                style={{ borderColor: qualityColor }}
            >
                {iconName ? (
                    <img 
                        src={`/items/icons/${iconName.toLowerCase()}.jpg`}
                        alt=""
                        className="w-full h-full object-cover"
                        onError={(e) => {
                            if (!e.target.src.includes('zamimg.com')) {
                                e.target.src = `https://wow.zamimg.com/images/wow/icons/medium/${iconName.toLowerCase()}.jpg`
                            } else {
                                e.target.style.display = 'none'
                            }
                        }}
                    />
                ) : (
                    <span className="text-gray-600 text-xs">?</span>
                )}
            </div>
            
            {/* ID */}
            <span className="text-gray-600 text-[11px] font-mono min-w-[40px]">
                [{itemId}]
            </span>
            
            {/* Name */}
            <span 
                className="text-[13px] font-bold truncate flex-1 group-hover:drop-shadow-sm"
                style={{ color: qualityColor }}
            >
                {item.name || item.itemName || 'Unknown Item'}
            </span>
            
            {/* Drop Chance (optional) */}
            {showDropChance && item.dropChance && (
                <span className="text-gray-500 text-[10px] uppercase tracking-tight">
                    {item.dropChance}
                </span>
            )}
        </div>
    )
}

/**
 * Icon placeholder for non-item entities (NPC, Object, etc)
 */
export const EntityIcon = ({ 
    label, 
    color = '#555',
    size = 'md' 
}) => {
    const sizes = {
        sm: 'w-6 h-6 text-[10px]',
        md: 'w-8 h-8 text-[11px]',
        lg: 'w-10 h-10 text-xs',
    }
    
    return (
        <div 
            className={`${sizes[size]} rounded flex items-center justify-center font-bold text-white flex-shrink-0`}
            style={{ backgroundColor: color }}
        >
            {label}
        </div>
    )
}

export default { LootItem, EntityIcon }
