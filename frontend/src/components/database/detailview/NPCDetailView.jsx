import React, { useState, useEffect } from 'react'
import { GetCreatureDetail } from '../../../services/api'
import { getQualityColor } from '../../../utils/wow'


import './DetailView.css'

import LootTile from '../../common/LootTile/LootTile'

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
        return (
            <LootTile
                key={item.itemId}
                item={item}
                onClick={() => onNavigate('item', item.itemId)}
                tooltipHandlers={tooltipHook.getItemHandlers(item.itemId)}
            />
        )
    }

    if (loading) return <div className="loading-view">Loading...</div>
    if (!detail) return <div className="error-view">NPC not found</div>

    return (
        <div className="detail-view">
            
            <header className="detail-page-header">
                <h1 className="detail-title" style={{ color: getQualityColor(detail.rank >= 1 ? 3 : 1) }}>{detail.name}</h1>
                <div className="detail-subtitle">{detail.subname && `<${detail.subname}>`} Level {detail.levelMin}-{detail.levelMax} {detail.typeName} ({detail.rankName})</div>
                <div className="detail-stats">
                    <span>Health: {detail.healthMax}</span>
                    <span>Mana: {detail.manaMax}</span>
                </div>
            </header>

            <div className="detail-grid">
                <div>
                    <h3 className="section-header">Loot Table</h3>
                    <div className="loot-grid">
                        {detail.loot && detail.loot.length > 0 ? detail.loot.sort((a,b)=>b.chance-a.chance).map(renderLootItem) : <div style={{color:'#666'}}>No loot.</div>}
                    </div>
                </div>
                <div>
                    <h3 className="section-header">Related Quests</h3>
                    
                    <h4 style={{ color: '#aaa', marginTop: '15px' }}>Starts ({detail.startsQuests ? detail.startsQuests.length : 0})</h4>
                    <ul className="quest-list">
                        {detail.startsQuests && detail.startsQuests.map(q => (
                            <li key={q.entry} className="quest-list-item">
                                <a className="quest-link" onClick={() => onNavigate('quest', q.entry)}>
                                    {q.name} [{q.entry}]
                                </a>
                            </li>
                        ))}
                    </ul>

                    <h4 style={{ color: '#aaa', marginTop: '15px' }}>Ends ({detail.endsQuests ? detail.endsQuests.length : 0})</h4>
                    <ul className="quest-list">
                        {detail.endsQuests && detail.endsQuests.map(q => (
                            <li key={q.entry} className="quest-list-item">
                                <a className="quest-link" onClick={() => onNavigate('quest', q.entry)}>
                                    {q.name} [{q.entry}]
                                </a>
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
        </div>
    )
}

export default NPCDetailView
