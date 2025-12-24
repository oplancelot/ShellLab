import React from 'react'
import { getQualityColor } from '../../utils/wow'

/**
 * WoW-style item tooltip component
 */
const ItemTooltip = ({ 
    item, 
    tooltip, 
    style,
    onMouseEnter,
    onMouseLeave 
}) => {
    // Loading state
    if (!tooltip) {
        return (
            <div 
                className="flex flex-col gap-1 p-2.5 bg-black/95 border border-border-light rounded pointer-events-none z-[1000] min-w-[200px] shadow-xl"
                style={style}
            >
                <div 
                    className="font-bold text-sm leading-tight"
                    style={{ color: getQualityColor(item?.quality) }}
                >
                    {item?.itemName || item?.name || 'Unknown Item'}
                </div>
                <div className="text-gray-500 italic text-[11px] animate-pulse">Loading...</div>
            </div>
        )
    }

    const qualityColor = getQualityColor(tooltip.quality)

    return (
        <div 
            className="flex flex-col gap-0.5 p-2.5 bg-black/95 border border-border-light rounded pointer-events-none z-[1000] min-w-[240px] max-w-[320px] shadow-2xl font-sans text-xs select-none"
            style={style}
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
        >
            {/* Name */}
            <div className="font-bold text-[14px] leading-tight" style={{ color: qualityColor }}>
                {tooltip.name}
            </div>
            
            {/* Set Name */}
            {tooltip.setName && (
                <div className="text-wow-gold leading-tight">{tooltip.setName}</div>
            )}
            
            {/* Item Level */}
            {tooltip.itemLevel > 0 && (
                <div className="text-wow-gold leading-tight">Item Level {tooltip.itemLevel}</div>
            )}
            
            {/* Binding */}
            {tooltip.binding && (
                <div className="text-white leading-tight">{tooltip.binding}</div>
            )}
            
            {/* Slot / Type */}
            <div className="flex justify-between items-center text-white leading-tight">
                {tooltip.slotName && <span>{tooltip.slotName}</span>}
                {tooltip.typeName && <span>{tooltip.typeName}</span>}
            </div>
            
            {/* Classes / Races */}
            {tooltip.classes && <div className="text-white leading-tight">{tooltip.classes}</div>}
            {tooltip.races && <div className="text-white leading-tight">{tooltip.races}</div>}
            
            {/* Damage */}
            {tooltip.damageText && (
                <div className="flex justify-between items-center text-white leading-tight">
                    <span>{tooltip.damageText}</span>
                    <span className="font-medium">{tooltip.speedText}</span>
                </div>
            )}
            
            {/* DPS */}
            {tooltip.dps && <div className="text-white leading-tight">{tooltip.dps}</div>}
            
            {/* Armor */}
            {tooltip.armor > 0 && (
                <div className="text-white leading-tight">{tooltip.armor} Armor</div>
            )}
            
            {/* Stats */}
            {tooltip.stats?.length > 0 && (
                <div className="flex flex-col">
                    {tooltip.stats.map((stat, i) => (
                        <div key={i} className="text-white leading-tight">{stat}</div>
                    ))}
                </div>
            )}
            
            {/* Resistances */}
            {tooltip.resistances?.length > 0 && (
                <div className="flex flex-col">
                    {tooltip.resistances.map((res, i) => (
                        <div key={i} className="text-white leading-tight">{res}</div>
                    ))}
                </div>
            )}
            
            {/* Durability */}
            {tooltip.durability && (
                <div className="text-white leading-tight text-[11px]">{tooltip.durability}</div>
            )}
            
            {/* Required Level */}
            {tooltip.requiredLevel > 1 && (
                <div className="text-white leading-tight">Requires Level {tooltip.requiredLevel}</div>
            )}
            
            {/* Spell Effects (green) - WoW style: after stats/durability */}
            {tooltip.effects?.length > 0 && (
                <div className="flex flex-col gap-0.5 mt-1">
                    {tooltip.effects.map((effect, i) => (
                        <div key={i} className="text-wow-uncommon leading-tight">{effect}</div>
                    ))}
                </div>
            )}
            
            {/* Set Info */}
            {tooltip.setInfo && (
                <div className="flex flex-col gap-0.5 mt-2 pt-2 border-t border-white/10">
                    <div className="text-wow-gold font-bold">{tooltip.setInfo.name}</div>
                    {tooltip.setInfo.items?.map((setItem, i) => (
                        <div key={i} className="text-gray-500 leading-tight ml-2 text-[11px]">{setItem}</div>
                    ))}
                    <div className="mt-1">
                        {tooltip.setInfo.bonuses?.map((bonus, i) => (
                            <div key={i} className="text-wow-uncommon leading-tight text-[11px]">{bonus}</div>
                        ))}
                    </div>
                </div>
            )}
            
            {/* Description */}
            {tooltip.description && (
                <div className="text-wow-gold italic leading-snug mt-1">"{tooltip.description}"</div>
            )}
            
            {/* Sell Price */}
            {tooltip.sellPrice > 0 && (
                <div className="text-white leading-tight flex items-center gap-1 mt-1 text-[11px]">
                    <span className="text-gray-500">Sell Price:</span> 
                    <span className="flex items-center gap-1">
                        {Math.floor(tooltip.sellPrice / 10000) > 0 && (
                            <span style={{ color: '#FFD700' }} className="drop-shadow-sm">{Math.floor(tooltip.sellPrice / 10000)}g</span>
                        )}
                        {Math.floor((tooltip.sellPrice % 10000) / 100) > 0 && (
                            <span style={{ color: '#C0C0C0' }} className="drop-shadow-sm">{Math.floor((tooltip.sellPrice % 10000) / 100)}s</span>
                        )}
                        {(tooltip.sellPrice % 100) > 0 && (
                            <span style={{ color: '#B87333' }} className="drop-shadow-sm">{tooltip.sellPrice % 100}c</span>
                        )}
                    </span>
                </div>
            )}
        </div>
    )
}

export default ItemTooltip
