import React, { useState, useEffect } from 'react'
import { GetQuestDetail } from '../../services/api'
import { getQualityColor } from '../../utils/wow'
import ItemTooltip from '../ItemTooltip'

const QuestDetailView = ({ entry, onBack, onNavigate, setHoveredItem, hoveredItem, tooltipCache, loadTooltipData }) => {
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        console.log(`[QuestDetailView] Mounting for entry: ${entry}`);
        setLoading(true)
        setError(null)
        
        GetQuestDetail(entry)
            .then(res => {
                console.log(`[QuestDetailView] Received data:`, res);
                if (!res) {
                    setError("Quest data is empty or invalid.");
                } else {
                    setDetail(res)
                }
                setLoading(false)
            })
            .catch(err => {
                console.error(`[QuestDetailView] Error:`, err);
                setError(err.toString());
                setLoading(false)
            })
    }, [entry])

    const renderRewardItem = (item, isChoice) => {
        const hasTooltip = hoveredItem === item.entry
        if (hasTooltip && !tooltipCache[item.entry]) {
            loadTooltipData(item.entry)
        }
        return (
            <div 
                key={item.entry}
                style={{ position: 'relative', display: 'flex', alignItems: 'center', background: '#242424', padding: '5px', borderRadius: '4px', margin: '5px 0', cursor: 'pointer' }}
                onClick={() => onNavigate('item', item.entry)}
                onMouseEnter={() => setHoveredItem(item.entry)}
                onMouseLeave={() => setHoveredItem(null)}
            >
                 <div style={{ width: '32px', height: '32px', border: `1px solid ${getQualityColor(item.quality)}`, marginRight: '8px' }}>
                    {item.icon ? <img src={`/items/icons/${item.icon}.jpg`} style={{width:'100%'}} /> : '?'}
                </div>
                <div>
                   <div style={{ color: getQualityColor(item.quality) }}>{item.name}</div>
                   {item.count > 1 && <div style={{ color: '#aaa', fontSize: '11px' }}>x{item.count}</div>}
                </div>
                {hasTooltip && tooltipCache[item.entry] && (
                    <div style={{ position: 'absolute', left: '100%', top: 0, zIndex: 1000, marginLeft: '10px' }}>
                        <ItemTooltip item={{entry: item.entry, quality: item.quality, name: item.name}} tooltip={tooltipCache[item.entry]} />
                    </div>
                )}
            </div>
        )
    }

    if (loading) return <div className="loading">Loading quest details (ID: {entry})...</div>
    if (error) return (
        <div className="error-container" style={{ padding: '20px', color: '#ff4444' }}>
            <h3>Error Loading Quest</h3>
            <p>{error}</p>
            <button onClick={onBack} style={{ background: '#333', border: 'none', color: '#fff', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer' }}>&larr; Back</button>
        </div>
    )
    if (!detail) return <div className="error">Quest not found (ID: {entry})</div>
    
    return (
         <div className="detail-view" style={{ flex: 1, overflowY: 'auto', padding: '20px', background: '#121212', color: '#e0e0e0' }}>
            <button onClick={onBack} style={{ background: '#333', border: 'none', color: '#fff', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer', marginBottom: '20px' }}>&larr; Back to List</button>
            
            <h1 style={{ color: '#FFD100', marginBottom: '10px' }}>{detail.title}</h1>
            <div style={{ color: '#888', marginBottom: '20px' }}>Level {detail.questLevel} (Min {detail.minLevel}) - {detail.type === 41 ? 'PVP' : detail.type === 81 ? 'Dungeon' : 'Normal'}</div>
            
            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 1fr)', gap: '40px' }}>
                <div>
                    <h3>Description</h3>
                    <p style={{ lineHeight: '1.6', color: '#ccc' }}>{detail.details}</p>
                    
                    <h3>Objectives</h3>
                    <p style={{ lineHeight: '1.6', color: '#ccc' }}>{detail.objectives}</p>
                </div>
                
                <div style={{ background: '#1a1a1a', padding: '20px', borderRadius: '8px' }}>
                    <h3 style={{ marginTop: 0, color: '#FFD100' }}>Rewards</h3>
                    {detail.rewMoney > 0 && <div style={{marginBottom:'10px'}}>Money: {Math.floor(detail.rewMoney/10000)}g {(detail.rewMoney%10000)/100}s</div>}
                    {detail.rewXp > 0 && <div style={{marginBottom:'10px'}}>XP: {detail.rewXp}</div>}
                    
                    {detail.rewards && detail.rewards.length > 0 && (
                        <div>
                            <h4>You will receive:</h4>
                            {detail.rewards.map(i => renderRewardItem(i, false))}
                        </div>
                    )}
                    
                    {detail.choiceRewards && detail.choiceRewards.length > 0 && (
                        <div>
                            <h4>Choose one of:</h4>
                            {detail.choiceRewards.map(i => renderRewardItem(i, true))}
                        </div>
                    )}
                    
                    <h3 style={{ color: '#FFD100', marginTop: '20px' }}>Related</h3>
                    {detail.starters && detail.starters.length > 0 && (
                        <div>
                            <h4>Start:</h4>
                            {detail.starters.map(s => (
                                <div key={s.entry} onClick={() => s.type==='npc' && onNavigate('npc', s.entry)} style={{cursor: s.type==='npc'?'pointer':'default', color: s.type==='npc'?'#4a9eff':'#aaa'}}>
                                    {s.name} ({s.type})
                                </div>
                            ))}
                        </div>
                    )}
                     {detail.enders && detail.enders.length > 0 && (
                        <div>
                            <h4>End:</h4>
                           {detail.enders.map(s => (
                                <div key={s.entry} onClick={() => s.type==='npc' && onNavigate('npc', s.entry)} style={{cursor: s.type==='npc'?'pointer':'default', color: s.type==='npc'?'#4a9eff':'#aaa'}}>
                                    {s.name} ({s.type})
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
         </div>
    )
}

export default QuestDetailView
