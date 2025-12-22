import { useState, useEffect } from 'react'
import { useItemTooltip } from '../hooks/useItemTooltip'
import ItemTooltip from './ItemTooltip'
import { NPCDetailView, QuestDetailView, ItemDetailView } from './database/detailview'
import { getQualityColor } from '../utils/wow'

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
    const renderTooltip = (item) => {
        if (hoveredItem === item.entry && tooltipCache[item.entry]) {
            return (
                <ItemTooltip
                    item={tooltipCache[item.entry]}
                    style={getTooltipStyle()}
                />
            )
        }
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
            <div className="database-page" style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                <div style={{ padding: '10px', background: '#252525', borderBottom: '1px solid #404040' }}>
                    <button 
                        onClick={goBack}
                        style={{ 
                            background: '#333', 
                            border: '1px solid #555', 
                            color: '#fff', 
                            padding: '5px 15px', 
                            cursor: 'pointer',
                            borderRadius: '3px'
                        }}
                    >
                        ← Back
                    </button>
                    <span style={{ marginLeft: '15px', color: '#888' }}>
                        Viewing: {currentDetail.type.toUpperCase()} #{currentDetail.entry}
                    </span>
                </div>
                <div style={{ flex: 1, overflow: 'auto', padding: '20px' }}>
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
        <div className="database-page" onMouseMove={(e) => handleMouseMove(e, hoveredItem)} style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* Tabs */}
            <div className="tabs" style={{ 
                display: 'flex', 
                gap: 0, 
                padding: '0 10px', 
                background: '#1a1a1a', 
                borderBottom: '1px solid #404040' 
            }}>
                {['Items', 'Sets', 'NPCs', 'Quests', 'Objects', 'Spells', 'Factions'].map(tab => (
                    <button
                        key={tab}
                        onClick={() => setActiveTab(tab.toLowerCase())}
                        style={{
                            padding: '8px 16px',
                            background: activeTab === tab.toLowerCase() ? '#383838' : 'transparent',
                            color: activeTab === tab.toLowerCase() ? '#fff' : '#FFD100',
                            border: activeTab === tab.toLowerCase() ? '1px solid #484848' : '1px solid transparent',
                            cursor: 'pointer',
                            fontWeight: 'bold',
                            borderRadius: '0',
                            fontSize: '13px',
                            textTransform: 'uppercase'
                        }}
                    >
                        {tab}
                    </button>
                ))}
            </div>

            {/* Content Area - 4 columns for three-level classification */}
            <div className="content" style={{ flex: 1, display: 'grid', gridTemplateColumns: '180px 180px 150px 1fr', gap: 0, overflow: 'hidden' }}>
                
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
        </div>
    )
}

export default DatabasePage
