import { useState, useEffect, useMemo } from 'react'
import { SidebarPanel, ContentPanel, ScrollList, SectionHeader, ListItem, EntityIcon } from '../../ui'
import { GetObjectTypes, GetObjectsByType, filterItems } from '../../../utils/databaseApi'

const OBJECT_COLOR = '#00B4FF'

function ObjectsTab() {
    const [objectTypes, setObjectTypes] = useState([])
    const [selectedObjectType, setSelectedObjectType] = useState(null)
    const [objects, setObjects] = useState([])
    const [loading, setLoading] = useState(false)

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
            GetObjectsByType(selectedObjectType.id, '')
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

    const filteredTypes = useMemo(() => filterItems(objectTypes, typeFilter), [objectTypes, typeFilter])
    const filteredObjects = useMemo(() => filterItems(objects, objectFilter), [objects, objectFilter])

    return (
        <>
            {/* Object Types */}
            <SidebarPanel className="col-span-1">
                <SectionHeader 
                    title={`Object Types (${filteredTypes.length})`}
                    placeholder="Filter types..."
                    onFilterChange={setTypeFilter}
                />
                <ScrollList>
                    {loading && objectTypes.length === 0 && (
                        <div className="p-4 text-center text-wow-gold italic animate-pulse">Loading types...</div>
                    )}
                    {filteredTypes.map(type => (
                        <ListItem
                            key={type.id}
                            active={selectedObjectType?.id === type.id}
                            onClick={() => {
                                setSelectedObjectType(type)
                                setObjectFilter('')
                            }}
                        >
                            <span className="flex justify-between w-full">
                                <span>{type.name}</span>
                                <span className="text-gray-600 text-xs">({type.count})</span>
                            </span>
                        </ListItem>
                    ))}
                </ScrollList>
            </SidebarPanel>

            {/* Objects List */}
            <ContentPanel className="col-span-3">
                <SectionHeader 
                    title={selectedObjectType ? `${selectedObjectType.name} (${filteredObjects.length})` : 'Select a Type'}
                    placeholder="Filter objects..."
                    onFilterChange={setObjectFilter}
                />
                
                {loading && selectedObjectType && (
                    <div className="flex-1 flex items-center justify-center text-wow-gold italic animate-pulse">
                        Loading objects...
                    </div>
                )}
                
                {!loading && objects.length > 0 && (
                    <ScrollList className="p-2 space-y-1">
                        {filteredObjects.map(obj => (
                            <div 
                                key={obj.entry}
                                className="flex items-center gap-3 p-2 bg-white/[0.02] hover:bg-white/5 border-l-[3px] cursor-pointer transition-colors rounded-r"
                                style={{ borderLeftColor: OBJECT_COLOR }}
                            >
                                <EntityIcon 
                                    label="OBJ"
                                    color={OBJECT_COLOR}
                                    size="md"
                                />
                                
                                <span className="text-gray-600 text-[11px] font-mono min-w-[50px]">
                                    [{obj.entry}]
                                </span>
                                
                                <span 
                                    className="font-bold flex-1 truncate"
                                    style={{ color: OBJECT_COLOR }}
                                >
                                    {obj.name}
                                </span>
                                
                                <span className="text-gray-500 text-xs ml-auto">
                                    Type: {obj.typeName || obj.type} | Size: {obj.size.toFixed(1)}
                                </span>
                            </div>
                        ))}
                    </ScrollList>
                )}
                
                {!selectedObjectType && !loading && (
                    <div className="flex-1 flex items-center justify-center text-gray-600 italic">
                        Select an object type to browse
                    </div>
                )}
            </ContentPanel>
        </>
    )
}

export default ObjectsTab
