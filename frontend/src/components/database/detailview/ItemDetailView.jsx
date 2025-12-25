import React, { useState, useEffect } from 'react'
import { GetItemDetail } from '../../../services/api'
import { getQualityColor } from '../../../utils/wow'
import { 
    DetailPageLayout, 
    DetailHeader, 
    DetailSection,
    DetailLoading,
    DetailError
} from '../../ui'
import { ItemTooltip, LootItem } from '../../ui'

const ItemDetailView = ({ entry, onBack, onNavigate, tooltipHook }) => {
    const { tooltipCache, loadTooltipData } = tooltipHook
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        GetItemDetail(entry).then(res => {
            setDetail(res)
            setLoading(false)
        })
    }, [entry])

    useEffect(() => {
        if (!tooltipCache[entry]) {
            loadTooltipData(entry)
        }
    }, [entry, tooltipCache, loadTooltipData])

    const renderTooltipBlock = () => {
        if (!detail) return null
        const dummyItem = { 
            entry: detail.entry, 
            quality: detail.quality, 
            name: detail.name 
        }
        
        return (
            <div className="inline-block align-top min-w-[300px]">
                <ItemTooltip 
                    item={dummyItem} 
                    tooltip={tooltipCache[entry]} 
                    style={{ position: 'static' }}
                    interactive={true} 
                />
            </div>
        )
    }

    if (loading) return <DetailLoading />
    if (!detail) return <DetailError message="Item not found" onBack={onBack} />

    const qualityColor = getQualityColor(detail.quality)

    return (
        <DetailPageLayout onBack={onBack}>
            <DetailHeader
                icon={
                    detail.iconPath ? (
                        <img 
                            src={`/items/icons/${detail.iconPath.toLowerCase()}.jpg`} 
                            className="w-full h-full object-cover" 
                            alt="" 
                            onError={(e) => {
                                if (!e.target.src.includes('zamimg.com')) {
                                    e.target.src = `https://wow.zamimg.com/images/wow/icons/medium/${detail.iconPath.toLowerCase()}.jpg`
                                } else {
                                    e.target.style.display = 'none'
                                }
                            }}
                        />
                    ) : '?'
                }
                iconBorderColor={qualityColor}
                title={detail.name}
                titleColor={qualityColor}
                subtitle={`Item Level ${detail.itemLevel}`}
            />

            <div className="flex flex-wrap gap-10">
                {/* Tooltip Block */}
                {renderTooltipBlock()}
                
                {/* Relations */}
                <div className="flex-1 min-w-[300px] space-y-8">
                    {/* Dropped By */}
                    {detail.droppedBy?.length > 0 && (
                        <DetailSection title="Dropped By">
                            <div className="space-y-1">
                                {detail.droppedBy.map(npc => (
                                    <div 
                                        key={npc.entry} 
                                        className="flex items-center justify-between p-2 bg-white/[0.02] hover:bg-white/5 border-b border-white/5 cursor-pointer transition-colors"
                                        onClick={() => onNavigate('npc', npc.entry)}
                                    >
                                        <div>
                                            <div className="text-white font-bold hover:text-wow-gold">
                                                {npc.name}
                                            </div>
                                            <div className="text-xs text-gray-500">
                                                Level {npc.levelMin}{npc.levelMax > npc.levelMin ? `-${npc.levelMax}` : ''}
                                            </div>
                                        </div>
                                        <div className="text-wow-gold font-mono text-sm">
                                            {npc.chance.toFixed(1)}%
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </DetailSection>
                    )}
                    
                    {/* Reward From */}
                    {/* Reward From */}
                    {detail.rewardFrom?.length > 0 && (
                        <DetailSection title="Reward From">
                            <div className="space-y-1">
                                {detail.rewardFrom.map(q => (
                                    <div 
                                        key={q.entry} 
                                        className="flex items-center gap-3 p-2 bg-white/[0.02] hover:bg-white/5 border-b border-white/5 cursor-pointer transition-colors"
                                        onClick={() => onNavigate('quest', q.entry)}
                                    >
                                        <div className="flex-1 min-w-0">
                                            <div className="text-wow-gold font-bold truncate">
                                                {q.title}
                                            </div>
                                            <div className="text-xs text-gray-500">
                                                Level {q.level}
                                                {q.isChoice && (
                                                    <span className="ml-2 text-[10px] border border-white/10 px-1 rounded uppercase">
                                                        Choice
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </DetailSection>
                    )}

                    {/* Contains */}
                    {detail.contains?.length > 0 && (
                        <DetailSection title="Contains">
                            <div className="grid grid-cols-1 gap-1">
                                {detail.contains.map(item => (
                                    <LootItem 
                                        key={item.entry}
                                        item={{ 
                                            ...item, 
                                            dropChance: item.chance ? item.chance.toFixed(1) + '%' : null
                                        }}
                                        showDropChance={true}
                                        onClick={() => onNavigate('item', item.entry)}
                                    />
                                ))}
                            </div>
                        </DetailSection>
                    )}
                </div>
            </div>
        </DetailPageLayout>
    )
}

export default ItemDetailView
