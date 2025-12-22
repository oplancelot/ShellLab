// Shared API utilities for database components
// These wrap Wails Go bindings with fallbacks

// Items APIs
export const BrowseItemsByClassAndSlot = (classId, subClass, inventoryType) => {
    if (window?.go?.main?.App?.BrowseItemsByClassAndSlot) {
        return window.go.main.App.BrowseItemsByClassAndSlot(classId, subClass, inventoryType)
    }
    return Promise.resolve([])
}

export const GetItemSets = () => {
    if (window?.go?.main?.App?.GetItemSets) {
        return window.go.main.App.GetItemSets()
    }
    return Promise.resolve([])
}

export const GetItemSetDetail = (itemSetId) => {
    if (window?.go?.main?.App?.GetItemSetDetail) {
        return window.go.main.App.GetItemSetDetail(itemSetId)
    }
    return Promise.resolve(null)
}

// Creature/NPC APIs
export const GetCreatureTypes = () => {
    if (window?.go?.main?.App?.GetCreatureTypes) {
        return window.go.main.App.GetCreatureTypes()
    }
    return Promise.resolve([])
}

export const BrowseCreaturesByType = (creatureType) => {
    if (window?.go?.main?.App?.BrowseCreaturesByType) {
        return window.go.main.App.BrowseCreaturesByType(creatureType)
    }
    return Promise.resolve([])
}

export const GetCreatureLoot = (entry) => {
    if (window?.go?.main?.App?.GetCreatureLoot) {
        return window.go.main.App.GetCreatureLoot(entry)
    }
    return Promise.resolve([])
}

// Quest APIs
export const GetQuestCategories = () => {
    if (window?.go?.main?.App?.GetQuestCategories) {
        return window.go.main.App.GetQuestCategories()
    }
    return Promise.resolve([])
}

export const GetQuestsByCategory = (categoryId) => {
    if (window?.go?.main?.App?.GetQuestsByCategory) {
        return window.go.main.App.GetQuestsByCategory(categoryId)
    }
    return Promise.resolve([])
}

export const SearchQuests = (query) => {
    if (window?.go?.main?.App?.SearchQuests) {
        return window.go.main.App.SearchQuests(query)
    }
    return Promise.resolve([])
}

// Object APIs
export const GetObjectTypes = () => {
    if (window?.go?.main?.App?.GetObjectTypes) {
        return window.go.main.App.GetObjectTypes()
    }
    return Promise.resolve([])
}

export const GetObjectsByType = (typeId) => {
    if (window?.go?.main?.App?.GetObjectsByType) {
        return window.go.main.App.GetObjectsByType(typeId)
    }
    return Promise.resolve([])
}

export const SearchObjects = (query) => {
    if (window?.go?.main?.App?.SearchObjects) {
        return window.go.main.App.SearchObjects(query)
    }
    return Promise.resolve([])
}

// Spells APIs
export const SearchSpells = (query) => {
    if (window?.go?.main?.App?.SearchSpells) {
        return window.go.main.App.SearchSpells(query)
    }
    return Promise.resolve([])
}

// Factions APIs
export const GetFactions = () => {
    if (window?.go?.main?.App?.GetFactions) {
        return window.go.main.App.GetFactions()
    }
    return Promise.resolve([])
}

// Filter helper function
export const filterItems = (items, filter) => {
    if (!filter || !filter.trim()) return items || []
    const searchLower = filter.toLowerCase().trim()
    const searchNum = parseInt(filter)
    const isNumericSearch = !isNaN(searchNum)
    
    return (items || []).filter(item => {
        if (isNumericSearch) {
            if (item.entry === searchNum || item.id === searchNum || item.itemsetId === searchNum) {
                return true
            }
        }
        const name = (item.name || item.title || item.displayName || '').toLowerCase()
        if (name.includes(searchLower)) return true
        if (item.key && item.key.toLowerCase().includes(searchLower)) return true
        return false
    })
}
