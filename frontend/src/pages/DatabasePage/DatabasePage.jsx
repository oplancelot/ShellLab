import { useState } from 'react'
import { useItemTooltip } from '../../hooks/useItemTooltip'
import { 
    PageLayout, 
    ContentGrid, 
    TabButton, 
    TabBar,
    ItemTooltip 
} from '../../components/ui'
import { NPCDetailView, QuestDetailView, ItemDetailView, SpellDetailView } from '../../components/database/detailview'
import { GRID_LAYOUT } from '../../components/common/layout'

// Import tab components
import { ItemsTab, SetsTab, NPCsTab, QuestsTab, ObjectsTab, SpellsTab, FactionsTab } from '../../components/database/tabs'

const TABS = ['Items', 'Sets', 'NPCs', 'Quests', 'Objects', 'Spells', 'Factions']

function DatabasePage() {
    const [activeTab, setActiveTab] = useState('items')
    
    // Navigation State for Detail Views
    const [detailStack, setDetailStack] = useState([]) // Stack of views: { type, entry }
    
    // Use shared tooltip hook
    const tooltipHook = useItemTooltip()
    const {
        hoveredItem,
        tooltipCache,
        getTooltipStyle,
    } = tooltipHook

    // Detail View Logic
    const navigateTo = (type, entry) => {
        console.log(`[DatabasePage] Navigating to ${type} with entry: ${entry}`);
        // Clear tooltip before navigation to prevent it from persisting
        tooltipHook.setHoveredItem(null)
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
        renderTooltip: () => null,
    }

    return (
        <PageLayout>
            {/* Tabs View - Hidden when detail active, but kept mounted to preserve state */}
            <div className={`flex flex-col h-full flex-1 overflow-hidden ${currentDetail ? 'hidden' : ''}`}>
                {/* Tab Bar */}
                <TabBar>
                    {TABS.map(tab => (
                        <TabButton
                            key={tab}
                            active={activeTab === tab.toLowerCase()}
                            onClick={() => setActiveTab(tab.toLowerCase())}
                        >
                            {tab}
                        </TabButton>
                    ))}
                </TabBar>

                {/* Content Area - 4 columns for three-level classification */}
                <ContentGrid columns={GRID_LAYOUT}>
                    {activeTab === 'items' && (
                        <ItemsTab 
                            tooltipHook={enhancedTooltipHook} 
                            onNavigate={navigateTo}
                        />
                    )}
                    {activeTab === 'sets' && (
                        <SetsTab tooltipHook={enhancedTooltipHook} />
                    )}
                    {activeTab === 'npcs' && (
                        <NPCsTab 
                            onNavigate={navigateTo}
                            tooltipHook={enhancedTooltipHook}
                        />
                    )}
                    {activeTab === 'quests' && (
                        <QuestsTab onNavigate={navigateTo} />
                    )}
                    {activeTab === 'objects' && (
                        <ObjectsTab />
                    )}
                    {activeTab === 'spells' && (
                        <SpellsTab onNavigate={navigateTo} />
                    )}
                    {activeTab === 'factions' && (
                        <FactionsTab />
                    )}
                </ContentGrid>
            </div>

            {/* Detail View Overlay */}
            {currentDetail && (
                <div className="flex flex-col h-full flex-1 overflow-hidden">
                    {/* Detail Header with breadcrumb */}
                    <div className="bg-bg-hover px-4 py-2 border-b border-border-dark flex items-center gap-4">
                        <button 
                            onClick={goBack}
                            className="bg-bg-panel border border-border-light text-gray-400 px-4 py-1.5 rounded hover:bg-bg-active hover:text-white transition-colors text-sm"
                        >
                            ← Back
                        </button>
                        <span className="text-gray-500 text-sm">
                            Viewing: <b className="text-gray-300 uppercase">{currentDetail.type}</b> 
                            <span className="ml-2 font-mono bg-black/20 px-1.5 py-0.5 rounded">#{currentDetail.entry}</span>
                        </span>
                    </div>
                    
                    {/* Detail Content */}
                    <div className="flex-1 overflow-auto">
                        {currentDetail.type === 'npc' && (
                            <NPCDetailView 
                                entry={currentDetail.entry} 
                                onNavigate={navigateTo}
                                onBack={goBack}
                                tooltipHook={enhancedTooltipHook}
                            />
                        )}
                        {currentDetail.type === 'quest' && (
                            <QuestDetailView 
                                entry={currentDetail.entry} 
                                onNavigate={navigateTo}
                                onBack={goBack}
                                tooltipHook={enhancedTooltipHook}
                            />
                        )}
                        {currentDetail.type === 'item' && (
                            <ItemDetailView 
                                entry={currentDetail.entry} 
                                onNavigate={navigateTo}
                                onBack={goBack}
                                tooltipHook={enhancedTooltipHook}
                            />
                        )}
                        {currentDetail.type === 'spell' && (
                            <SpellDetailView 
                                entry={currentDetail.entry} 
                                onNavigate={navigateTo}
                                onBack={goBack}
                                tooltipHook={enhancedTooltipHook}
                            />
                        )}
                    </div>
                </div>
            )}

            {/* Global Tooltip Layer */}
            {hoveredItem && tooltipCache[hoveredItem] && (
                <ItemTooltip
                    item={tooltipCache[hoveredItem]}
                    tooltip={tooltipCache[hoveredItem]}
                    style={getTooltipStyle()}
                />
            )}
        </PageLayout>
    )
}

export default DatabasePage
