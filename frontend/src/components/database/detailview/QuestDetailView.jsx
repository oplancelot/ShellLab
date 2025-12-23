import React, { useState, useEffect } from 'react'
import { GetQuestDetail } from '../../../services/api'
import { getQualityColor } from '../../../utils/wow'


const QuestDetailView = ({ entry, onBack, onNavigate, tooltipHook }) => {
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
        return (
            <div 
                key={item.entry}
                style={{ position: 'relative', display: 'flex', alignItems: 'center', background: '#242424', padding: '5px', borderRadius: '4px', margin: '5px 0', cursor: 'pointer' }}
                onClick={() => onNavigate('item', item.entry)}
                {...tooltipHook.getItemHandlers(item.entry)}
            >
                 <div style={{ width: '32px', height: '32px', border: `1px solid ${getQualityColor(item.quality)}`, marginRight: '8px' }}>
                    {item.icon ? <img src={`/items/icons/${item.icon}.jpg`} style={{width:'100%'}} /> : '?'}
                </div>
                <div>
                   <div style={{ color: getQualityColor(item.quality) }}>{item.name}</div>
                   {item.count > 1 && <div style={{ color: '#aaa', fontSize: '11px' }}>x{item.count}</div>}
                </div>
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

            
            <h1 style={{ color: '#FFD100', marginBottom: '10px' }}>{detail.title} [{detail.entry}]</h1>
            <div style={{ color: '#888', marginBottom: '20px' }}>Level {detail.questLevel} (Min {detail.minLevel}) - {detail.type === 41 ? 'PVP' : detail.type === 81 ? 'Dungeon' : 'Normal'}</div>
            
            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 1fr)', gap: '40px' }}>
                <div>
                    <h3>Description</h3>
                    <p style={{ lineHeight: '1.6', color: '#ccc' }}>{detail.details}</p>
                    
                    <h3>Objectives</h3>
                    <p style={{ lineHeight: '1.6', color: '#ccc' }}>{detail.objectives}</p>

                    <div style={{ marginTop: '30px', borderTop: '1px solid #333', paddingTop: '20px' }}>
                        <h3 style={{ color: '#FFD100' }}>Rewards</h3>
                        {detail.rewMoney > 0 && <div style={{marginBottom:'10px'}}>Money: {Math.floor(detail.rewMoney/10000)}g {(detail.rewMoney%10000)/100}s</div>}
                        {detail.rewXp > 0 && <div style={{marginBottom:'10px'}}>XP: {detail.rewXp}</div>}
                        
                        {detail.rewards && detail.rewards.length > 0 && (
                            <div style={{ marginBottom: '20px' }}>
                                <h4 style={{ color: '#aaa' }}>You will receive:</h4>
                                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '10px' }}>
                                    {detail.rewards.map(i => renderRewardItem(i, false))}
                                </div>
                            </div>
                        )}
                        
                        {detail.choiceRewards && detail.choiceRewards.length > 0 && (
                            <div>
                                <h4 style={{ color: '#aaa' }}>Choose one of:</h4>
                                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '10px' }}>
                                    {detail.choiceRewards.map(i => renderRewardItem(i, true))}
                                </div>
                            </div>
                        )}
                    </div>
                </div>
                
                <div style={{ background: '#1a1a1a', padding: '20px', borderRadius: '8px', alignSelf: 'start' }}>
                    <h3 style={{ color: '#FFD100', marginTop: 0 }}>Quest Chain</h3>
                    {detail.series && detail.series.length > 0 ? (
                        <div style={{ marginBottom: '20px' }}>
                            {detail.series.map((s, index) => (
                                <div key={s.entry} style={{ display: 'flex', gap: '8px', marginBottom: '4px', fontSize: '13px' }}>
                                    <span style={{ color: '#666', width: '20px' }}>{index + 1}.</span>
                                    {s.entry === detail.entry ? (
                                        <b style={{ color: '#fff' }}>{s.title} [{s.entry}]</b>
                                    ) : (
                                        <a 
                                            onClick={() => onNavigate('quest', s.entry)}
                                            style={{ color: '#FFD100', cursor: 'pointer', textDecoration: 'none' }}
                                            onMouseEnter={(e) => e.target.style.textDecoration = 'underline'}
                                            onMouseLeave={(e) => e.target.style.textDecoration = 'none'}
                                        >
                                            {s.title} [{s.entry}]
                                        </a>
                                    )}
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div style={{ color: '#666', fontSize: '13px', marginBottom: '20px' }}>This quest is not part of a series.</div>
                    )}

                    {detail.prevQuests && detail.prevQuests.length > 0 && (
                        <div style={{ marginBottom: '20px' }}>
                            <h3 style={{ color: '#FFD100', marginTop: 0 }}>Prerequisites</h3>
                            {detail.prevQuests.map(q => (
                                <div key={q.entry} style={{ fontSize: '13px', marginBottom: '5px' }}>
                                    <a onClick={() => onNavigate('quest', q.entry)} style={{ color: '#FFD100', cursor: 'pointer' }}>
                                        {q.title} [{q.entry}]
                                    </a>
                                </div>
                            ))}
                        </div>
                    )}

                    {detail.exclusiveQuests && detail.exclusiveQuests.length > 0 && (
                        <div style={{ marginBottom: '20px' }}>
                            <h3 style={{ color: '#FFD100', marginTop: 0 }}>Exclusive With</h3>
                            <div style={{ fontSize: '12px', color: '#aaa', marginBottom: '5px' }}>Completion of this quest makes these unavailable:</div>
                            {detail.exclusiveQuests.map(q => (
                                <div key={q.entry} style={{ fontSize: '13px', marginBottom: '5px' }}>
                                    <a onClick={() => onNavigate('quest', q.entry)} style={{ color: '#FFD100', cursor: 'pointer' }}>
                                        {q.title} [{q.entry}]
                                    </a>
                                </div>
                            ))}
                        </div>
                    )}

                    <h3 style={{ color: '#FFD100', marginTop: '20px' }}>Requirements</h3>
                    <div style={{ fontSize: '13px', color: '#ccc' }}>
                        {detail.requiredRaces > 0 && <div style={{ marginBottom: '5px' }}>Races: {detail.requiredRaces} (Mask)</div>}
                        {detail.requiredClasses > 0 && <div style={{ marginBottom: '5px' }}>Classes: {detail.requiredClasses} (Mask)</div>}
                        {detail.srcItemId > 0 && (
                            <div style={{ marginBottom: '5px' }}>
                                Item: <a onClick={() => onNavigate('item', detail.srcItemId)} style={{ color: '#FFD100', cursor: 'pointer' }}>[Item {detail.srcItemId}]</a>
                            </div>
                        )}
                    </div>

                    <h3 style={{ color: '#FFD100', marginTop: '20px' }}>Related</h3>
                    {detail.starters && detail.starters.length > 0 && (
                        <div style={{ marginBottom: '15px' }}>
                            <h4 style={{ color: '#aaa', margin: '0 0 5px 0' }}>Starts:</h4>
                            {detail.starters.map(s => (
                                <div key={s.entry} onClick={() => s.type==='npc' && onNavigate('npc', s.entry)} style={{cursor: s.type==='npc'?'pointer':'default', color: s.type==='npc'?'#4a9eff':'#aaa', padding: '2px 0'}}>
                                    {s.name} ({s.type}) [{s.entry}]
                                </div>
                            ))}
                        </div>
                    )}
                     {detail.enders && detail.enders.length > 0 && (
                        <div>
                            <h4 style={{ color: '#aaa', margin: '0 0 5px 0' }}>Ends:</h4>
                           {detail.enders.map(s => (
                                <div key={s.entry} onClick={() => s.type==='npc' && onNavigate('npc', s.entry)} style={{cursor: s.type==='npc'?'pointer':'default', color: s.type==='npc'?'#4a9eff':'#aaa', padding: '2px 0'}}>
                                    {s.name} ({s.type}) [{s.entry}]
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
