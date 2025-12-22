import React, { useState, useEffect } from 'react'
import { GetItemDetail } from '../../services/api'
import { getQualityColor } from '../../utils/wow'
import ItemTooltip from '../ItemTooltip'

const ItemDetailView = ({ entry, onBack, onNavigate, setHoveredItem, tooltipCache, loadTooltipData }) => {
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        GetItemDetail(entry).then(res => {
            setDetail(res)
            setLoading(false)
        })
    }, [entry])

    const renderTooltipBlock = () => {
        if (!detail) return null
        const dummyItem = { 
            entry: detail.entry, 
            quality: detail.quality, 
            name: detail.name 
        }
        if (!tooltipCache[entry]) {
            loadTooltipData(entry)
        }
        
        return (
             <div style={{ display: 'inline-block', verticalAlign: 'top', minWidth: '300px' }}>
                <ItemTooltip 
                    item={dummyItem} 
                    tooltip={tooltipCache[entry]} 
                    style={{ position: 'static', border: '1px solid #444', background: '#000' }} 
                />
             </div>
        )
    }

    if (loading) return <div className="loading">Loading Item...</div>
    if (!detail) return <div className="error">Item not found</div>

    return (
        <div className="detail-view" style={{ flex: 1, overflowY: 'auto', padding: '20px', background: '#121212', color: '#e0e0e0' }}>
            <button onClick={onBack} style={{ background: '#333', border: 'none', color: '#fff', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer', marginBottom: '20px' }}>&larr; Back to List</button>
            
            <header style={{ marginBottom: '30px', display: 'flex', gap: '20px' }}>
                <div style={{ width: '56px', height: '56px', border: `2px solid ${getQualityColor(detail.quality)}`, borderRadius: '4px' }}>
                     {detail.iconPath ? <img src={`/items/icons/${detail.iconPath}.jpg`} style={{width:'100%', height:'100%'}} /> : '?'}
                </div>
                <div>
                     <h1 style={{ color: getQualityColor(detail.quality), margin: '0 0 5px 0' }}>{detail.name}</h1>
                     <div style={{ color: '#888' }}>Item Level {detail.itemLevel}</div>
                </div>
            </header>

            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '40px' }}>
                {renderTooltipBlock()}
                
                <div style={{ flex: 1, minWidth: '300px' }}>
                     {detail.droppedBy && detail.droppedBy.length > 0 && (
                        <div style={{ marginBottom: '30px' }}>
                            <h3 style={{ borderBottom: '1px solid #FFD100', paddingBottom: '5px', color: '#FFD100' }}>Dropped By</h3>
                            <ul style={{ listStyle: 'none', padding: 0 }}>
                                {detail.droppedBy.map(npc => (
                                    <li key={npc.entry} style={{ padding: '8px', borderBottom: '1px solid #333', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                        <div>
                                            <a onClick={() => onNavigate('npc', npc.entry)} style={{ cursor: 'pointer', color: '#fff', fontWeight: 'bold' }}>
                                                {npc.name}
                                            </a>
                                            <div style={{ fontSize: '12px', color: '#888' }}>Level {npc.levelMin}-{npc.levelMax}</div>
                                        </div>
                                        <div style={{ color: '#FFD100' }}>{npc.chance.toFixed(1)}%</div>
                                    </li>
                                ))}
                            </ul>
                        </div>
                     )}
                     
                     {detail.rewardFrom && detail.rewardFrom.length > 0 && (
                        <div>
                            <h3 style={{ borderBottom: '1px solid #FFD100', paddingBottom: '5px', color: '#FFD100' }}>Reward From</h3>
                            <ul style={{ listStyle: 'none', padding: 0 }}>
                                {detail.rewardFrom.map(q => (
                                    <li key={q.entry} style={{ padding: '8px', borderBottom: '1px solid #333' }}>
                                        <a onClick={() => onNavigate('quest', q.entry)} style={{ cursor: 'pointer', color: '#FFD100' }}>
                                            {q.title}
                                        </a>
                                        <span style={{ marginLeft: '10px', fontSize: '12px', color: '#888' }}>Level {q.level}</span>
                                        {q.isChoice && <span style={{ marginLeft: '10px', fontSize: '10px', color: '#aaa', border: '1px solid #555', padding: '0 3px' }}>Choice</span>}
                                    </li>
                                ))}
                            </ul>
                        </div>
                     )}
                </div>
            </div>
        </div>
    )
}

export default ItemDetailView
