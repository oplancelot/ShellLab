import React, { useState, useEffect } from 'react'
import { GetSpellDetail } from '../../../services/api'
import { 
    DetailPageLayout, 
    DetailHeader, 
    DetailSection,
    DetailLoading,
    DetailError
} from '../../ui'
import { getIconPath } from '../../../utils/wow'

const SpellDetailView = ({ entry, onBack, onNavigate, tooltipHook }) => {
    const [detail, setDetail] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        setLoading(true)
        setError(null)
        
        GetSpellDetail(parseInt(entry))
            .then(res => {
                if (!res) {
                    setError("Spell data not found");
                } else {
                    setDetail(res)
                }
                setLoading(false)
            })
            .catch(err => {
                setError(err.toString());
                setLoading(false)
            })
    }, [entry])

    if (loading) return <DetailLoading />
    if (error) return <DetailError message={error} onBack={onBack} />
    if (!detail) return <DetailError message="Spell not found" onBack={onBack} />
    
    // Determine schools
    const schoolMap = {
        0: 'Physical', 1: 'Holy', 2: 'Fire', 3: 'Nature', 4: 'Frost', 5: 'Shadow', 6: 'Arcane'
    }
    const schoolName = schoolMap[detail.school] || 'Unknown'

    // Format power type
    const powerTypes = {
        0: 'Mana', 1: 'Rage', 2: 'Focus', 3: 'Energy', 4: 'Happiness'
    }
    const powerType = powerTypes[detail.powerType] || 'Power'

    return (
        <DetailPageLayout onBack={onBack}>
            <DetailHeader
                title={`${detail.name} [${detail.entry}]`}
                icon={
                    <img 
                        src={getIconPath(detail.icon)} 
                        alt="" 
                        className="w-full h-full object-cover" 
                        onError={(e) => {
                            e.target.style.display = 'none';
                        }}
                    />
                }
                titleColor="#FFD100" 
                subtitle={`Level ${detail.spellLevel} - ${schoolName}`}
            />
            
            <div className="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-10">
                {/* Main Content */}
                <div className="space-y-8">
                     <DetailSection title="Description">
                        <p className="text-gray-300 leading-relaxed whitespace-pre-wrap">
                            {detail.description || 'No description available.'}
                        </p>
                    </DetailSection>

                    {detail.toolTip && (
                        <DetailSection title="Tooltip">
                            <p className="text-gray-300 leading-relaxed whitespace-pre-wrap">
                                {detail.toolTip}
                            </p>
                        </DetailSection>
                    )}
                </div>
                
                {/* Side Panel */}
                <div className="space-y-6">
                    <DetailSection title="Properties">
                        <div className="grid grid-cols-2 gap-y-2 text-sm">
                            <span className="text-gray-500">Duration:</span>
                            <span className="text-gray-300 text-right">{detail.duration}</span>
                            
                            <span className="text-gray-500">Range:</span>
                            <span className="text-gray-300 text-right">{detail.range}</span>
                            
                            <span className="text-gray-500">Cost:</span>
                            <span className="text-gray-300 text-right">
                                {detail.manaCost > 0 ? `${detail.manaCost} ${powerType}` : 'None'}
                            </span>
                            
                            <span className="text-gray-500">Cast Time:</span>
                            <span className="text-gray-300 text-right">{detail.castTime}</span>

                            <span className="text-gray-500">School:</span>
                            <span className="text-gray-300 text-right">{schoolName}</span>
                            
                            <span className="text-gray-500">Level:</span>
                            <span className="text-gray-300 text-right">{detail.spellLevel}</span>
                        </div>
                    </DetailSection>
                </div>
            </div>
        </DetailPageLayout>
    )
}

export default SpellDetailView
