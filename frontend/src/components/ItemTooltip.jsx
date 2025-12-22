import React from 'react'

/**
 * Get quality color for item quality level
 */
export const getQualityColor = (quality) => {
    const colors = {
        0: '#9d9d9d', // Poor (gray)
        1: '#ffffff', // Common (white)
        2: '#1eff00', // Uncommon (green)
        3: '#0070dd', // Rare (blue)
        4: '#a335ee', // Epic (purple)
        5: '#ff8000', // Legendary (orange)
        6: '#e6cc80', // Artifact (gold)
        7: '#00ccff', // Heirloom (cyan)
    }
    return colors[quality] || '#ffffff'
}

/**
 * Shared Item Tooltip Component
 * @param {Object} props
 * @param {Object} props.item - The item object (must have itemId/entry and quality)
 * @param {Object} props.tooltip - The loaded tooltip data
 * @param {Object} props.style - Position styles for the tooltip
 * @param {Function} props.onMouseEnter - Handler when mouse enters tooltip
 * @param {Function} props.onMouseLeave - Handler when mouse leaves tooltip
 */
function ItemTooltip({ item, tooltip, style, onMouseEnter, onMouseLeave }) {
    // Loading state
    if (!tooltip) {
        return (
            <div className="item-tooltip" style={style}>
                <div className="tooltip-name" style={{ color: getQualityColor(item.quality) }}>
                    {item.itemName || item.name || 'Unknown Item'}
                </div>
                <div className="tooltip-loading">Loading...</div>
            </div>
        )
    }

    return (
        <div 
            className="item-tooltip" 
            style={style}
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
        >
            <div className="tooltip-name" style={{ color: getQualityColor(tooltip.quality) }}>
                {tooltip.name}
            </div>
            
            {tooltip.setName && (
                <div className="tooltip-setname">{tooltip.setName}</div>
            )}
            
            {tooltip.itemLevel > 0 && (
                <div className="tooltip-itemlevel">Item Level {tooltip.itemLevel}</div>
            )}
            
            {tooltip.binding && (
                <div className="tooltip-binding">{tooltip.binding}</div>
            )}
            
            <div className="tooltip-slot-type">
                {tooltip.slotName && <span className="tooltip-slot">{tooltip.slotName}</span>}
                {tooltip.typeName && <span className="tooltip-type">{tooltip.typeName}</span>}
            </div>
            
            {tooltip.classes && (
                <div className="tooltip-classes">{tooltip.classes}</div>
            )}
            
            {tooltip.races && (
                <div className="tooltip-races">{tooltip.races}</div>
            )}
            
            {tooltip.damageText && (
                <div className="tooltip-damage">
                    <span>{tooltip.damageText}</span>
                    <span className="tooltip-speed">{tooltip.speedText}</span>
                </div>
            )}
            
            {tooltip.dps && (
                <div className="tooltip-dps">{tooltip.dps}</div>
            )}
            
            {tooltip.armor > 0 && (
                <div className="tooltip-armor">{tooltip.armor} Armor</div>
            )}
            
            {tooltip.stats && tooltip.stats.length > 0 && (
                <div className="tooltip-stats">
                    {tooltip.stats.map((stat, i) => (
                        <div key={i} className="tooltip-stat">{stat}</div>
                    ))}
                </div>
            )}
            
            {tooltip.resistances && tooltip.resistances.length > 0 && (
                <div className="tooltip-resistances">
                    {tooltip.resistances.map((res, i) => (
                        <div key={i} className="tooltip-resistance">{res}</div>
                    ))}
                </div>
            )}
            
            {tooltip.spellEffects && tooltip.spellEffects.length > 0 && (
                <div className="tooltip-effects">
                    {tooltip.spellEffects.map((effect, i) => (
                        <div key={i} className="tooltip-effect">{effect}</div>
                    ))}
                </div>
            )}
            
            {tooltip.durability && (
                <div className="tooltip-durability">{tooltip.durability}</div>
            )}
            
            {tooltip.requiredLevel > 1 && (
                <div className="tooltip-reqlevel">Requires Level {tooltip.requiredLevel}</div>
            )}
            
            {/* Item Set Info */}
            {tooltip.setInfo && (
                <div className="tooltip-set">
                    <div className="tooltip-set-name">{tooltip.setInfo.name}</div>
                    {tooltip.setInfo.items && tooltip.setInfo.items.map((setItem, i) => (
                        <div key={i} className="tooltip-set-item">{setItem}</div>
                    ))}
                    {tooltip.setInfo.bonuses && tooltip.setInfo.bonuses.map((bonus, i) => (
                        <div key={i} className="tooltip-set-bonus">{bonus}</div>
                    ))}
                </div>
            )}
            
            {tooltip.description && (
                <div className="tooltip-description">"{tooltip.description}"</div>
            )}
            
            {tooltip.sellPrice && (
                <div className="tooltip-sellprice">Sell Price: {tooltip.sellPrice}</div>
            )}
        </div>
    )
}

export default ItemTooltip
