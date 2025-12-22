import React, { useState, useEffect } from 'react'
import { GetCreatureDetail } from '../../services/api'
import { getQualityColor } from '../../utils/wow'
import ItemTooltip from '../ItemTooltip'

const NPCDetailView = ({ entry, onBack, onNavigate, setHoveredItem, hoveredItem, tooltipCache, loadTooltipData }) => {
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
        const hasTooltip = hoveredItem === item.itemId
        if (hasTooltip && !tooltipCache[item.itemId]) {
            loadTooltipData(item.itemId)
        }

        return (
            <div 
                key={item.itemId}
                className="loot-tile" 
                style={{ position: 'relative', display: 'flex', alignItems: 'center', background: '#242424', padding: '5px', borderRadius: '4px', cursor: 'pointer' }}
                onMouseEnter={() => setHoveredItem(item.itemId)}
                onMouseLeave={() => setHoveredItem(null)}
                onClick={() => onNavigate('item', item.itemId)}
            >
                <div style={{ width: '32px', height: '32px', border: `1px solid ${getQualityColor(item.quality)}`, marginRight: '8px' }}>
                    {item.icon ? <img src={`/items/icons/${item.icon}.jpg`} style={{width:'100%'}} /> : '?'}
                </div>
                <div style={{ flex: 1, overflow: 'hidden' }}>
                    <div style={{ color: getQualityColor(item.quality), fontWeight: 'bold', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{item.itemName}</div>
                    <div style={{ color: '#aaa', fontSize: '11px' }}>{item.chance.toFixed(1)}% {item.minCount > 1 && `(${item.minCount}-${item.maxCount})`}</div>
                </div>
                {hasTooltip && tooltipCache[item.itemId] && (
                    <div style={{ position: 'absolute', left: '100%', top: 0, zIndex: 1000, marginLeft: '10px' }}>
                        <ItemTooltip item={{entry: item.itemId, quality: item.quality, name: item.itemName}} tooltip={tooltipCache[item.itemId]} />
                    </div>
                )}
            </div>
        )
    }

    if (loading) return <div className="loading">Loading...</div>
    if (!detail) return <div className="error">NPC not found</div>

    return (
        <div className="detail-view" style={{ flex: 1, overflowY: 'auto', padding: '20px', background: '#121212', color: '#e0e0e0' }}>
            <button onClick={onBack} style={{ background: '#333', border: 'none', color: '#fff', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer', marginBottom: '20px' }}>&larr; Back to List</button>
            
            <header style={{ marginBottom: '30px', borderBottom: '1px solid #333', paddingBottom: '20px' }}>
                <h1 style={{ color: getQualityColor(detail.rank >= 1 ? 3 : 1), margin: '0 0 10px 0' }}>{detail.name}</h1>
                <div style={{ color: '#888' }}>{detail.subname && `<${detail.subname}>`} Level {detail.levelMin}-{detail.levelMax} {detail.typeName} ({detail.rankName})</div>
                <div style={{ marginTop: '10px' }}>
                    <span style={{ marginRight: '20px' }}>Health: {detail.healthMax}</span>
                    <span>Mana: {detail.manaMax}</span>
                </div>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '30px' }}>
                <div>
                    <h3 style={{ borderBottom: '1px solid #FFD100', paddingBottom: '5px', color: '#FFD100' }}>Loot Table</h3>
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '8px', marginTop: '10px' }}>
                        {detail.loot.length > 0 ? detail.loot.sort((a,b)=>b.chance-a.chance).map(renderLootItem) : <div style={{color:'#666'}}>No loot.</div>}
                    </div>
                </div>
                <div>
                    <h3 style={{ borderBottom: '1px solid #FFD100', paddingBottom: '5px', color: '#FFD100' }}>Related Quests</h3>
                    
                    <h4 style={{ color: '#aaa', marginTop: '15px' }}>Starts ({detail.startsQuests ? detail.startsQuests.length : 0})</h4>
                    <ul style={{ listStyle: 'none', padding: 0 }}>
                        {detail.startsQuests && detail.startsQuests.map(q => (
                            <li key={q.entry} style={{ padding: '4px 0' }}>
                                <a className="quest-link" onClick={() => onNavigate('quest', q.entry)} style={{ cursor: 'pointer', color: '#FFD100', textDecoration: 'none' }}>
                                    {q.title}
                                </a>
                            </li>
                        ))}
                    </ul>

                    <h4 style={{ color: '#aaa', marginTop: '15px' }}>Ends ({detail.endsQuests ? detail.endsQuests.length : 0})</h4>
                    <ul style={{ listStyle: 'none', padding: 0 }}>
                        {detail.endsQuests && detail.endsQuests.map(q => (
                            <li key={q.entry} style={{ padding: '4px 0' }}>
                                <a onClick={() => onNavigate('quest', q.entry)} style={{ cursor: 'pointer', color: '#FFD100' }}>
                                    {q.title}
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
