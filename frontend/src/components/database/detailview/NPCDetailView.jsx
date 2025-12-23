import React, { useState, useEffect } from 'react'
import { GetCreatureDetail } from '../../../services/api'
import { getQualityColor } from '../../../utils/wow'
import { 
    DetailPageLayout, 
    DetailHeader, 
    DetailSection,
    DetailGrid,
    LootGrid,
    StatBadge,
    DetailLoading,
    DetailError
} from '../../ui'
import { LootItem } from '../../ui'

const NPCDetailView = ({ entry, onBack, onNavigate, tooltipHook }) => {
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        GetCreatureDetail(entry).then(res => {
            setDetail(res)
            setLoading(false)
        })
    }, [entry])

    const renderLootItem = (item) => {
        const handlers = tooltipHook?.getItemHandlers?.(item.itemId) || {}
        return (
            <LootItem
                key={item.itemId}
                item={{
                    entry: item.itemId,
                    name: item.name,
                    quality: item.quality,
                    iconPath: item.iconPath,
                    dropChance: `${item.chance.toFixed(1)}%`
                }}
                onClick={() => onNavigate('item', item.itemId)}
                showDropChance
                {...handlers}
            />
        )
    }

    if (loading) return <DetailLoading />
    if (!detail) return <DetailError message="NPC not found" onBack={onBack} />

    const rankColor = detail.rank >= 3 ? getQualityColor(4) : detail.rank >= 1 ? getQualityColor(5) : getQualityColor(1)

    return (
        <DetailPageLayout onBack={onBack}>
            <DetailHeader
                title={detail.name}
                titleColor={rankColor}
                subtitle={
                    <>
                        {detail.subname && <span className="italic">&lt;{detail.subname}&gt;</span>}
                        {' '}Level {detail.levelMin}-{detail.levelMax} {detail.typeName} ({detail.rankName})
                    </>
                }
                stats={
                    <>
                        <StatBadge label="Health" value={detail.healthMax.toLocaleString()} />
                        <StatBadge label="Mana" value={detail.manaMax.toLocaleString()} />
                    </>
                }
            />

            <DetailGrid>
                {/* Loot Table */}
                <DetailSection title="Loot Table">
                    {detail.loot?.length > 0 ? (
                        <LootGrid>
                            {detail.loot
                                .sort((a, b) => b.chance - a.chance)
                                .map(renderLootItem)
                            }
                        </LootGrid>
                    ) : (
                        <div className="text-gray-600 italic">No loot.</div>
                    )}
                </DetailSection>

                {/* Related Quests */}
                <div className="space-y-6">
                    <DetailSection title="Related Quests">
                        {/* Starts */}
                        <h4 className="text-gray-400 text-sm font-semibold mb-2">
                            Starts ({detail.startsQuests?.length || 0})
                        </h4>
                        {detail.startsQuests?.length > 0 ? (
                            <ul className="space-y-1 mb-6">
                                {detail.startsQuests.map(q => (
                                    <li key={q.entry}>
                                        <a 
                                            className="text-wow-gold hover:underline cursor-pointer"
                                            onClick={() => onNavigate('quest', q.entry)}
                                        >
                                            {q.name} <span className="text-gray-600 text-xs">[{q.entry}]</span>
                                        </a>
                                    </li>
                                ))}
                            </ul>
                        ) : (
                            <div className="text-gray-600 text-xs italic mb-6">None</div>
                        )}

                        {/* Ends */}
                        <h4 className="text-gray-400 text-sm font-semibold mb-2">
                            Ends ({detail.endsQuests?.length || 0})
                        </h4>
                        {detail.endsQuests?.length > 0 ? (
                            <ul className="space-y-1">
                                {detail.endsQuests.map(q => (
                                    <li key={q.entry}>
                                        <a 
                                            className="text-wow-gold hover:underline cursor-pointer"
                                            onClick={() => onNavigate('quest', q.entry)}
                                        >
                                            {q.name} <span className="text-gray-600 text-xs">[{q.entry}]</span>
                                        </a>
                                    </li>
                                ))}
                            </ul>
                        ) : (
                            <div className="text-gray-600 text-xs italic">None</div>
                        )}
                    </DetailSection>
                </div>
            </DetailGrid>
        </DetailPageLayout>
    )
}

export default NPCDetailView
