import { useState, useEffect } from 'react'
import { useItemTooltip } from '../hooks/useItemTooltip'
import './DatabasePage.css'
import ItemTooltip from './ItemTooltip'
import { NPCDetailView, QuestDetailView, ItemDetailView } from './database/detailview'
import { getQualityColor } from '../utils/wow'
import { GRID_LAYOUT } from './common/layout'

// Import tab components
import { ItemsTab, SetsTab, NPCsTab, QuestsTab, ObjectsTab, SpellsTab, FactionsTab } from './database/tabs'

function DatabasePage() {
    const [activeTab, setActiveTab] = useState('items')
    
    // Navigation State for Detail Views
    const [detailStack, setDetailStack] = useState([]) // Stack of views: { type, entry }
    
    // Use shared tooltip hook
    const tooltipHook = useItemTooltip()
    const {
        hoveredItem,
        setHoveredItem,
        tooltipCache,
        loadTooltipData,
        handleMouseMove,
        handleItemEnter,
        getTooltipStyle,
    } = tooltipHook

    // Helper for quality class
    const getQualityClass = (quality) => `q${quality || 0}`

    // Tooltip renderer helper to pass to tabs
    // Note: Tooltip is now rendered globally in DatabasePage to avoid overflow clipping
    const renderTooltip = (item) => {
        return null
    }

    // Detail View Logic
    const navigateTo = (type, entry) => {
        console.log(`[DatabasePage] Navigating to ${type} with entry: ${entry}`);
        setDetailStack(prev => [...prev, { type, entry }])
    }
    const goBack = () => {
        console.log(`[DatabasePage] Going back. Previous stack size: ${detailStack.length}`);
        setDetailStack(prev => prev.slice(0, -1))
    }

    const currentDetail = detailStack.length > 0 ? detailStack[detailStack.length - 1] : null
    
    // Enhanced tooltip hook to pass to tabs
    const enhancedTooltipHook = {
        ...tooltipHook,
        renderTooltip,
    }

    // Render Detail View if active
    if (currentDetail) {
        return (
            <div className="database-page">
                <div className="detail-header">
                    <button className="back-button" onClick={goBack}>
                        ← Back
                    </button>
                    <span className="detail-info">
                        Viewing: {currentDetail.type.toUpperCase()} #{currentDetail.entry}
                    </span>
                </div>
                <div className="detail-content">
                    {currentDetail.type === 'npc' && (
                        <NPCDetailView 
                            entry={currentDetail.entry} 
                            onNavigate={navigateTo}
                            onBack={goBack}
                        />
                    )}
                    {currentDetail.type === 'quest' && (
                        <QuestDetailView 
                            entry={currentDetail.entry} 
                            onNavigate={navigateTo}
                            onBack={goBack}
                        />
                    )}
                    {currentDetail.type === 'item' && (
                        <ItemDetailView 
                            entry={currentDetail.entry} 
                            onNavigate={navigateTo}
                            onBack={goBack}
                        />
                    )}
                </div>
            </div>
        )
    }

    return (
        <div className="database-page">
            {/* Tabs */}
            <div className="tabs">
                {['Items', 'Sets', 'NPCs', 'Quests', 'Objects', 'Spells', 'Factions'].map(tab => (
                    <button
                        key={tab}
                        onClick={() => setActiveTab(tab.toLowerCase())}
                        className={`tab-button ${activeTab === tab.toLowerCase() ? 'active' : ''}`}
                    >
                        {tab}
                    </button>
                ))}
            </div>

            {/* Content Area - 4 columns for three-level classification */}
            <div className="content" style={{ gridTemplateColumns: GRID_LAYOUT }}>
                
                {/* ITEMS TAB */}
                {activeTab === 'items' && (
                    <ItemsTab tooltipHook={enhancedTooltipHook} />
                )}

                {/* SETS TAB */}
                {activeTab === 'sets' && (
                    <SetsTab tooltipHook={enhancedTooltipHook} />
                )}

                {/* NPCS TAB */}
                {activeTab === 'npcs' && (
                    <NPCsTab 
                        onNavigate={navigateTo}
                        tooltipHook={enhancedTooltipHook}
                    />
                )}

                {/* QUESTS TAB */}
                {activeTab === 'quests' && (
                    <QuestsTab onNavigate={navigateTo} />
                )}

                {/* OBJECTS TAB */}
                {activeTab === 'objects' && (
                    <ObjectsTab />
                )}

                {/* SPELLS TAB */}
                {activeTab === 'spells' && (
                    <SpellsTab />
                )}

                {/* FACTIONS TAB */}
                {activeTab === 'factions' && (
                    <FactionsTab />
                )}
            </div>

            {/* Global Tooltip Layer */}
            {hoveredItem && tooltipCache[hoveredItem] && (
                 <ItemTooltip
                     item={tooltipCache[hoveredItem]}
                     tooltip={tooltipCache[hoveredItem]}
                     style={getTooltipStyle()}
                 />
            )}
        </div>
    )
}

export default DatabasePage
