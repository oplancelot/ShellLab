import { useState, useEffect, useMemo } from 'react'
import { SectionHeader } from '../../common/SectionHeader'
import { GetObjectTypes, GetObjectsByType, filterItems } from '../../../utils/databaseApi'

function ObjectsTab() {
    const [objectTypes, setObjectTypes] = useState([])
    const [selectedObjectType, setSelectedObjectType] = useState(null)
    const [objects, setObjects] = useState([])
    const [loading, setLoading] = useState(false)

    // Independent filter states
    const [typeFilter, setTypeFilter] = useState('')
    const [objectFilter, setObjectFilter] = useState('')

    // Load object types on mount
    useEffect(() => {
        setLoading(true)
        GetObjectTypes()
            .then(types => {
                setObjectTypes(types || [])
                setLoading(false)
            })
            .catch(err => {
                console.error("Failed to load object types:", err)
                setLoading(false)
            })
    }, [])

    // Load objects when a type is selected
    useEffect(() => {
        if (selectedObjectType !== null) {
            setLoading(true)
            setObjects([])
            GetObjectsByType(selectedObjectType.id)
                .then(res => {
                    setObjects(res || [])
                    setLoading(false)
                })
                .catch(err => {
                    console.error("Failed to load objects:", err)
                    setLoading(false)
                })
        }
    }, [selectedObjectType])

    // Filtered lists
    const filteredTypes = useMemo(() => filterItems(objectTypes, typeFilter), [objectTypes, typeFilter])
    const filteredObjects = useMemo(() => filterItems(objects, objectFilter), [objects, objectFilter])

    return (
        <>
            {/* Object Types List */}
            <aside className="sidebar" style={{ gridColumn: '1 / 2' }}>
                <SectionHeader 
                    title={`Object Types (${filteredTypes.length})`}
                    placeholder="Filter types..."
                    onFilterChange={setTypeFilter}
                />
                <div className="list">
                    {loading && objectTypes.length === 0 && (
                        <div className="loading">Loading types...</div>
                    )}
                    {filteredTypes.map(type => (
                        <div
                            key={type.id}
                            className={`item ${selectedObjectType?.id === type.id ? 'active' : ''}`}
                            onClick={() => {
                                setSelectedObjectType(type)
                                setObjectFilter('')
                            }}
                        >
                            {type.name} ({type.count})
                        </div>
                    ))}
                </div>
            </aside>

            {/* Objects List */}
            <section className="loot" style={{ gridColumn: '2 / -1' }}>
                <SectionHeader 
                    title={selectedObjectType ? `${selectedObjectType.name} (${filteredObjects.length})` : 'Select a Type'}
                    placeholder="Filter objects..."
                    onFilterChange={setObjectFilter}
                />
                
                {loading && selectedObjectType && (
                    <div className="loading">Loading objects...</div>
                )}
                
                {objects.length > 0 && (
                    <div className="loot-items">
                        {filteredObjects.map(obj => (
                            <div 
                                key={obj.entry}
                                className="loot-item"
                                style={{ borderLeft: '3px solid #00B4FF' }}
                            >
                                <div className="item-icon-placeholder" style={{ 
                                    background: '#00B4FF',
                                    color: '#fff',
                                    fontWeight: 'bold',
                                    fontSize: '10px'
                                }}>
                                    OBJ
                                </div>
                                
                                <span className="item-id">[{obj.entry}]</span>
                                
                                <span style={{ color: '#00B4FF', fontWeight: 'bold' }}>
                                    {obj.name}
                                </span>
                                
                                <span style={{ marginLeft: 'auto', color: '#888', fontSize: '11px' }}>
                                    Type: {obj.typeName || obj.type} | Size: {obj.size.toFixed(1)}
                                </span>
                            </div>
                        ))}
                    </div>
                )}
                
                {!selectedObjectType && (
                    <p className="placeholder">Select an object type to browse</p>
                )}
            </section>
        </>
    )
}

export default ObjectsTab
