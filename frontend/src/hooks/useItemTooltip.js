import { useState, useCallback } from 'react'

/**
 * Custom hook for item tooltip with mouse-following behavior
 * @returns {Object} Tooltip state and handlers
 */
export function useItemTooltip() {
    const [hoveredItem, setHoveredItem] = useState(null)
    const [tooltipPos, setTooltipPos] = useState({ top: 0, left: 0 })
    const [tooltipCache, setTooltipCache] = useState({})

    // Get tooltip data from backend
    const getTooltipData = useCallback(async (itemId) => {
        if (window?.go?.main?.App?.GetTooltipData) {
            return window.go.main.App.GetTooltipData(itemId)
        }
        return null
    }, [])

    // Load tooltip data for an item
    const loadTooltipData = useCallback(async (itemId) => {
        if (tooltipCache[itemId]) return tooltipCache[itemId]
        
        try {
            const data = await getTooltipData(itemId)
            if (data) {
                setTooltipCache(prev => ({ ...prev, [itemId]: data }))
                return data
            }
        } catch (err) {
            console.error('Failed to load tooltip:', err)
        }
        return null
    }, [tooltipCache, getTooltipData])

    // Handle mouse move - update tooltip position following mouse
    const handleMouseMove = useCallback((e, itemId) => {
        const lootContainer = e.currentTarget.closest('.loot')
        const containerRect = lootContainer 
            ? lootContainer.getBoundingClientRect() 
            : { left: 0, right: window.innerWidth, top: 0, bottom: window.innerHeight }
        const itemRect = e.currentTarget.getBoundingClientRect()
        
        // Tooltip dimensions
        const tooltipWidth = 320
        const tooltipHeight = 400
        
        // Position tooltip to the right and below the cursor
        let left = e.clientX + 15
        let top = e.clientY + 15
        
        // Don't let tooltip cover the item row - keep it below the item
        if (top < itemRect.bottom + 5) {
            top = itemRect.bottom + 5
        }
        
        // Keep within container bounds - horizontal
        if (left + tooltipWidth > containerRect.right - 10) {
            left = e.clientX - tooltipWidth - 15
        }
        if (left < containerRect.left + 10) {
            left = containerRect.left + 10
        }
        
        // Keep within container bounds - vertical
        if (top + tooltipHeight > containerRect.bottom - 10) {
            top = containerRect.bottom - tooltipHeight - 10
        }
        if (top < containerRect.top + 10) {
            top = containerRect.top + 10
        }
        
        setTooltipPos({ top, left })
        setHoveredItem(itemId)
    }, [])

    // Handle item enter - load tooltip data
    const handleItemEnter = useCallback((itemId) => {
        loadTooltipData(itemId)
    }, [loadTooltipData])

    // Handle item leave - hide tooltip
    const handleItemLeave = useCallback(() => {
        setHoveredItem(null)
    }, [])

    // Get event handlers for an item element
    const getItemHandlers = useCallback((itemId) => ({
        onMouseEnter: () => handleItemEnter(itemId),
        onMouseMove: (e) => handleMouseMove(e, itemId),
        onMouseLeave: handleItemLeave,
    }), [handleItemEnter, handleMouseMove, handleItemLeave])

    // Get styles for the tooltip container
    const getTooltipStyle = useCallback(() => ({
        position: 'fixed',
        left: tooltipPos.left,
        top: tooltipPos.top,
        zIndex: 10000,
    }), [tooltipPos])

    return {
        hoveredItem,
        setHoveredItem,
        tooltipPos,
        tooltipCache,
        loadTooltipData,
        handleMouseMove,
        handleItemEnter,
        handleItemLeave,
        getItemHandlers,
        getTooltipStyle,
    }
}

export default useItemTooltip
