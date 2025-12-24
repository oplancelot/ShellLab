export namespace main {
	
	export class CreaturePageResult {
	    creatures: models.Creature[];
	    total: number;
	    hasMore: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CreaturePageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.creatures = this.convertValues(source["creatures"], models.Creature);
	        this.total = source["total"];
	        this.hasMore = source["hasMore"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LegacyLootItem {
	    itemId: number;
	    itemName: string;
	    iconName: string;
	    quality: number;
	    dropChance?: string;
	    slotType?: string;
	
	    static createFrom(source: any = {}) {
	        return new LegacyLootItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemId = source["itemId"];
	        this.itemName = source["itemName"];
	        this.iconName = source["iconName"];
	        this.quality = source["quality"];
	        this.dropChance = source["dropChance"];
	        this.slotType = source["slotType"];
	    }
	}
	export class LegacyBossLoot {
	    bossName: string;
	    items: LegacyLootItem[];
	
	    static createFrom(source: any = {}) {
	        return new LegacyBossLoot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bossName = source["bossName"];
	        this.items = this.convertValues(source["items"], LegacyLootItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class AtlasTable {
	    key: string;
	    displayName: string;
	
	    static createFrom(source: any = {}) {
	        return new AtlasTable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.displayName = source["displayName"];
	    }
	}
	export class Category {
	    id: number;
	    key: string;
	    name: string;
	    parentId?: number;
	    type: string;
	    sortOrder: number;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.key = source["key"];
	        this.name = source["name"];
	        this.parentId = source["parentId"];
	        this.type = source["type"];
	        this.sortOrder = source["sortOrder"];
	    }
	}
	export class Creature {
	    entry: number;
	    name: string;
	    subname?: string;
	    levelMin: number;
	    levelMax: number;
	    healthMin: number;
	    healthMax: number;
	    manaMin: number;
	    manaMax: number;
	    type: number;
	    typeName: string;
	    rank: number;
	    rankName: string;
	    faction: number;
	    npcFlags: number;
	
	    static createFrom(source: any = {}) {
	        return new Creature(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.subname = source["subname"];
	        this.levelMin = source["levelMin"];
	        this.levelMax = source["levelMax"];
	        this.healthMin = source["healthMin"];
	        this.healthMax = source["healthMax"];
	        this.manaMin = source["manaMin"];
	        this.manaMax = source["manaMax"];
	        this.type = source["type"];
	        this.typeName = source["typeName"];
	        this.rank = source["rank"];
	        this.rankName = source["rankName"];
	        this.faction = source["faction"];
	        this.npcFlags = source["npcFlags"];
	    }
	}
	export class QuestRelation {
	    entry: number;
	    name: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new QuestRelation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.type = source["type"];
	    }
	}
	export class LootItem {
	    itemId: number;
	    itemName: string;
	    icon: string;
	    quality: number;
	    chance: number;
	    minCount: number;
	    maxCount: number;
	
	    static createFrom(source: any = {}) {
	        return new LootItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemId = source["itemId"];
	        this.itemName = source["itemName"];
	        this.icon = source["icon"];
	        this.quality = source["quality"];
	        this.chance = source["chance"];
	        this.minCount = source["minCount"];
	        this.maxCount = source["maxCount"];
	    }
	}
	export class CreatureDetail {
	    entry: number;
	    name: string;
	    subname?: string;
	    levelMin: number;
	    levelMax: number;
	    healthMin: number;
	    healthMax: number;
	    manaMin: number;
	    manaMax: number;
	    type: number;
	    typeName: string;
	    rank: number;
	    rankName: string;
	    faction: number;
	    npcFlags: number;
	    loot: LootItem[];
	    startsQuests: QuestRelation[];
	    endsQuests: QuestRelation[];
	
	    static createFrom(source: any = {}) {
	        return new CreatureDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.subname = source["subname"];
	        this.levelMin = source["levelMin"];
	        this.levelMax = source["levelMax"];
	        this.healthMin = source["healthMin"];
	        this.healthMax = source["healthMax"];
	        this.manaMin = source["manaMin"];
	        this.manaMax = source["manaMax"];
	        this.type = source["type"];
	        this.typeName = source["typeName"];
	        this.rank = source["rank"];
	        this.rankName = source["rankName"];
	        this.faction = source["faction"];
	        this.npcFlags = source["npcFlags"];
	        this.loot = this.convertValues(source["loot"], LootItem);
	        this.startsQuests = this.convertValues(source["startsQuests"], QuestRelation);
	        this.endsQuests = this.convertValues(source["endsQuests"], QuestRelation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreatureDrop {
	    entry: number;
	    name: string;
	    levelMin: number;
	    levelMax: number;
	    chance: number;
	
	    static createFrom(source: any = {}) {
	        return new CreatureDrop(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.levelMin = source["levelMin"];
	        this.levelMax = source["levelMax"];
	        this.chance = source["chance"];
	    }
	}
	export class CreatureType {
	    type: number;
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new CreatureType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class Faction {
	    id: number;
	    name: string;
	    description: string;
	    side: number;
	    categoryId: number;
	
	    static createFrom(source: any = {}) {
	        return new Faction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.side = source["side"];
	        this.categoryId = source["categoryId"];
	    }
	}
	export class GameObject {
	    entry: number;
	    name: string;
	    type: number;
	    typeName: string;
	    displayId: number;
	    size: number;
	    data?: number[];
	
	    static createFrom(source: any = {}) {
	        return new GameObject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.typeName = source["typeName"];
	        this.displayId = source["displayId"];
	        this.size = source["size"];
	        this.data = source["data"];
	    }
	}
	export class InventorySlot {
	    class: number;
	    subClass: number;
	    inventoryType: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new InventorySlot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.class = source["class"];
	        this.subClass = source["subClass"];
	        this.inventoryType = source["inventoryType"];
	        this.name = source["name"];
	    }
	}
	export class Item {
	    entry: number;
	    name: string;
	    description?: string;
	    quality: number;
	    itemLevel: number;
	    requiredLevel: number;
	    class: number;
	    subClass: number;
	    inventoryType: number;
	    iconPath: string;
	    sellPrice?: number;
	    buyPrice?: number;
	    allowableClass?: number;
	    allowableRace?: number;
	    bonding?: number;
	    maxDurability?: number;
	    armor?: number;
	    statType1?: number;
	    statValue1?: number;
	    statType2?: number;
	    statValue2?: number;
	    statType3?: number;
	    statValue3?: number;
	    statType4?: number;
	    statValue4?: number;
	    statType5?: number;
	    statValue5?: number;
	    statType6?: number;
	    statValue6?: number;
	    statType7?: number;
	    statValue7?: number;
	    statType8?: number;
	    statValue8?: number;
	    statType9?: number;
	    statValue9?: number;
	    statType10?: number;
	    statValue10?: number;
	    delay?: number;
	    dmgMin1?: number;
	    dmgMax1?: number;
	    dmgType1?: number;
	    holyRes?: number;
	    fireRes?: number;
	    natureRes?: number;
	    frostRes?: number;
	    shadowRes?: number;
	    arcaneRes?: number;
	    spellId1?: number;
	    spellTrigger1?: number;
	    spellId2?: number;
	    spellTrigger2?: number;
	    spellId3?: number;
	    spellTrigger3?: number;
	    setId?: number;
	    dropRate?: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.quality = source["quality"];
	        this.itemLevel = source["itemLevel"];
	        this.requiredLevel = source["requiredLevel"];
	        this.class = source["class"];
	        this.subClass = source["subClass"];
	        this.inventoryType = source["inventoryType"];
	        this.iconPath = source["iconPath"];
	        this.sellPrice = source["sellPrice"];
	        this.buyPrice = source["buyPrice"];
	        this.allowableClass = source["allowableClass"];
	        this.allowableRace = source["allowableRace"];
	        this.bonding = source["bonding"];
	        this.maxDurability = source["maxDurability"];
	        this.armor = source["armor"];
	        this.statType1 = source["statType1"];
	        this.statValue1 = source["statValue1"];
	        this.statType2 = source["statType2"];
	        this.statValue2 = source["statValue2"];
	        this.statType3 = source["statType3"];
	        this.statValue3 = source["statValue3"];
	        this.statType4 = source["statType4"];
	        this.statValue4 = source["statValue4"];
	        this.statType5 = source["statType5"];
	        this.statValue5 = source["statValue5"];
	        this.statType6 = source["statType6"];
	        this.statValue6 = source["statValue6"];
	        this.statType7 = source["statType7"];
	        this.statValue7 = source["statValue7"];
	        this.statType8 = source["statType8"];
	        this.statValue8 = source["statValue8"];
	        this.statType9 = source["statType9"];
	        this.statValue9 = source["statValue9"];
	        this.statType10 = source["statType10"];
	        this.statValue10 = source["statValue10"];
	        this.delay = source["delay"];
	        this.dmgMin1 = source["dmgMin1"];
	        this.dmgMax1 = source["dmgMax1"];
	        this.dmgType1 = source["dmgType1"];
	        this.holyRes = source["holyRes"];
	        this.fireRes = source["fireRes"];
	        this.natureRes = source["natureRes"];
	        this.frostRes = source["frostRes"];
	        this.shadowRes = source["shadowRes"];
	        this.arcaneRes = source["arcaneRes"];
	        this.spellId1 = source["spellId1"];
	        this.spellTrigger1 = source["spellTrigger1"];
	        this.spellId2 = source["spellId2"];
	        this.spellTrigger2 = source["spellTrigger2"];
	        this.spellId3 = source["spellId3"];
	        this.spellTrigger3 = source["spellTrigger3"];
	        this.setId = source["setId"];
	        this.dropRate = source["dropRate"];
	    }
	}
	export class ItemSubClass {
	    class: number;
	    subClass: number;
	    name: string;
	    inventorySlots?: InventorySlot[];
	
	    static createFrom(source: any = {}) {
	        return new ItemSubClass(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.class = source["class"];
	        this.subClass = source["subClass"];
	        this.name = source["name"];
	        this.inventorySlots = this.convertValues(source["inventorySlots"], InventorySlot);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ItemClass {
	    class: number;
	    name: string;
	    subClasses?: ItemSubClass[];
	
	    static createFrom(source: any = {}) {
	        return new ItemClass(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.class = source["class"];
	        this.name = source["name"];
	        this.subClasses = this.convertValues(source["subClasses"], ItemSubClass);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QuestReward {
	    entry: number;
	    title: string;
	    level: number;
	    isChoice: boolean;
	
	    static createFrom(source: any = {}) {
	        return new QuestReward(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.title = source["title"];
	        this.level = source["level"];
	        this.isChoice = source["isChoice"];
	    }
	}
	export class ItemDetail {
	    entry: number;
	    name: string;
	    description?: string;
	    quality: number;
	    itemLevel: number;
	    requiredLevel: number;
	    class: number;
	    subClass: number;
	    inventoryType: number;
	    iconPath: string;
	    sellPrice?: number;
	    buyPrice?: number;
	    allowableClass?: number;
	    allowableRace?: number;
	    bonding?: number;
	    maxDurability?: number;
	    armor?: number;
	    statType1?: number;
	    statValue1?: number;
	    statType2?: number;
	    statValue2?: number;
	    statType3?: number;
	    statValue3?: number;
	    statType4?: number;
	    statValue4?: number;
	    statType5?: number;
	    statValue5?: number;
	    statType6?: number;
	    statValue6?: number;
	    statType7?: number;
	    statValue7?: number;
	    statType8?: number;
	    statValue8?: number;
	    statType9?: number;
	    statValue9?: number;
	    statType10?: number;
	    statValue10?: number;
	    delay?: number;
	    dmgMin1?: number;
	    dmgMax1?: number;
	    dmgType1?: number;
	    holyRes?: number;
	    fireRes?: number;
	    natureRes?: number;
	    frostRes?: number;
	    shadowRes?: number;
	    arcaneRes?: number;
	    spellId1?: number;
	    spellTrigger1?: number;
	    spellId2?: number;
	    spellTrigger2?: number;
	    spellId3?: number;
	    spellTrigger3?: number;
	    setId?: number;
	    dropRate?: string;
	    displayId: number;
	    flags: number;
	    buyCount: number;
	    maxCount: number;
	    stackable: number;
	    containerSlots: number;
	    material: number;
	    dmgMin2: number;
	    dmgMax2: number;
	    dmgType2: number;
	    droppedBy: CreatureDrop[];
	    rewardFrom: QuestReward[];
	
	    static createFrom(source: any = {}) {
	        return new ItemDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.quality = source["quality"];
	        this.itemLevel = source["itemLevel"];
	        this.requiredLevel = source["requiredLevel"];
	        this.class = source["class"];
	        this.subClass = source["subClass"];
	        this.inventoryType = source["inventoryType"];
	        this.iconPath = source["iconPath"];
	        this.sellPrice = source["sellPrice"];
	        this.buyPrice = source["buyPrice"];
	        this.allowableClass = source["allowableClass"];
	        this.allowableRace = source["allowableRace"];
	        this.bonding = source["bonding"];
	        this.maxDurability = source["maxDurability"];
	        this.armor = source["armor"];
	        this.statType1 = source["statType1"];
	        this.statValue1 = source["statValue1"];
	        this.statType2 = source["statType2"];
	        this.statValue2 = source["statValue2"];
	        this.statType3 = source["statType3"];
	        this.statValue3 = source["statValue3"];
	        this.statType4 = source["statType4"];
	        this.statValue4 = source["statValue4"];
	        this.statType5 = source["statType5"];
	        this.statValue5 = source["statValue5"];
	        this.statType6 = source["statType6"];
	        this.statValue6 = source["statValue6"];
	        this.statType7 = source["statType7"];
	        this.statValue7 = source["statValue7"];
	        this.statType8 = source["statType8"];
	        this.statValue8 = source["statValue8"];
	        this.statType9 = source["statType9"];
	        this.statValue9 = source["statValue9"];
	        this.statType10 = source["statType10"];
	        this.statValue10 = source["statValue10"];
	        this.delay = source["delay"];
	        this.dmgMin1 = source["dmgMin1"];
	        this.dmgMax1 = source["dmgMax1"];
	        this.dmgType1 = source["dmgType1"];
	        this.holyRes = source["holyRes"];
	        this.fireRes = source["fireRes"];
	        this.natureRes = source["natureRes"];
	        this.frostRes = source["frostRes"];
	        this.shadowRes = source["shadowRes"];
	        this.arcaneRes = source["arcaneRes"];
	        this.spellId1 = source["spellId1"];
	        this.spellTrigger1 = source["spellTrigger1"];
	        this.spellId2 = source["spellId2"];
	        this.spellTrigger2 = source["spellTrigger2"];
	        this.spellId3 = source["spellId3"];
	        this.spellTrigger3 = source["spellTrigger3"];
	        this.setId = source["setId"];
	        this.dropRate = source["dropRate"];
	        this.displayId = source["displayId"];
	        this.flags = source["flags"];
	        this.buyCount = source["buyCount"];
	        this.maxCount = source["maxCount"];
	        this.stackable = source["stackable"];
	        this.containerSlots = source["containerSlots"];
	        this.material = source["material"];
	        this.dmgMin2 = source["dmgMin2"];
	        this.dmgMax2 = source["dmgMax2"];
	        this.dmgType2 = source["dmgType2"];
	        this.droppedBy = this.convertValues(source["droppedBy"], CreatureDrop);
	        this.rewardFrom = this.convertValues(source["rewardFrom"], QuestReward);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ItemSetBrowse {
	    itemsetId: number;
	    name: string;
	    itemIds: number[];
	    itemCount: number;
	    skillId: number;
	    skillLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemSetBrowse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemsetId = source["itemsetId"];
	        this.name = source["name"];
	        this.itemIds = source["itemIds"];
	        this.itemCount = source["itemCount"];
	        this.skillId = source["skillId"];
	        this.skillLevel = source["skillLevel"];
	    }
	}
	export class SetBonus {
	    threshold: number;
	    spellId: number;
	
	    static createFrom(source: any = {}) {
	        return new SetBonus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.threshold = source["threshold"];
	        this.spellId = source["spellId"];
	    }
	}
	export class ItemSetDetail {
	    itemsetId: number;
	    name: string;
	    items: Item[];
	    bonuses: SetBonus[];
	
	    static createFrom(source: any = {}) {
	        return new ItemSetDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemsetId = source["itemsetId"];
	        this.name = source["name"];
	        this.items = this.convertValues(source["items"], Item);
	        this.bonuses = this.convertValues(source["bonuses"], SetBonus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ItemSetInfo {
	    name: string;
	    items: string[];
	    bonuses: string[];
	
	    static createFrom(source: any = {}) {
	        return new ItemSetInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.items = source["items"];
	        this.bonuses = source["bonuses"];
	    }
	}
	
	
	export class ObjectType {
	    id: number;
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new ObjectType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class Quest {
	    entry: number;
	    title: string;
	    questLevel: number;
	    minLevel: number;
	    type: number;
	    zoneOrSort: number;
	    categoryName: string;
	    requiredRaces: number;
	    requiredClasses: number;
	    srcItemId: number;
	    rewardXp: number;
	    rewardMoney: number;
	    prevQuestId: number;
	    nextQuestId: number;
	    exclusiveGroup: number;
	    nextQuestInChain: number;
	
	    static createFrom(source: any = {}) {
	        return new Quest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.title = source["title"];
	        this.questLevel = source["questLevel"];
	        this.minLevel = source["minLevel"];
	        this.type = source["type"];
	        this.zoneOrSort = source["zoneOrSort"];
	        this.categoryName = source["categoryName"];
	        this.requiredRaces = source["requiredRaces"];
	        this.requiredClasses = source["requiredClasses"];
	        this.srcItemId = source["srcItemId"];
	        this.rewardXp = source["rewardXp"];
	        this.rewardMoney = source["rewardMoney"];
	        this.prevQuestId = source["prevQuestId"];
	        this.nextQuestId = source["nextQuestId"];
	        this.exclusiveGroup = source["exclusiveGroup"];
	        this.nextQuestInChain = source["nextQuestInChain"];
	    }
	}
	export class QuestCategory {
	    id: number;
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new QuestCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class QuestCategoryEnhanced {
	    id: number;
	    groupId: number;
	    name: string;
	    questCount: number;
	
	    static createFrom(source: any = {}) {
	        return new QuestCategoryEnhanced(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.groupId = source["groupId"];
	        this.name = source["name"];
	        this.questCount = source["questCount"];
	    }
	}
	export class QuestCategoryGroup {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new QuestCategoryGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class QuestSeriesItem {
	    entry: number;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new QuestSeriesItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.title = source["title"];
	    }
	}
	export class QuestReputation {
	    factionId: number;
	    name: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new QuestReputation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.factionId = source["factionId"];
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class QuestItem {
	    entry: number;
	    name: string;
	    icon: string;
	    count: number;
	    quality: number;
	
	    static createFrom(source: any = {}) {
	        return new QuestItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.count = source["count"];
	        this.quality = source["quality"];
	    }
	}
	export class QuestDetail {
	    entry: number;
	    title: string;
	    details: string;
	    objectives: string;
	    offerRewardText?: string;
	    endText?: string;
	    questLevel: number;
	    minLevel: number;
	    type: number;
	    zoneOrSort: number;
	    categoryName: string;
	    requiredRaces?: number;
	    requiredClasses?: number;
	    rewardXp: number;
	    rewardMoney: number;
	    rewardSpell?: number;
	    rewardItems: QuestItem[];
	    choiceItems: QuestItem[];
	    reputation: QuestReputation[];
	    starters: QuestRelation[];
	    enders: QuestRelation[];
	    series: QuestSeriesItem[];
	    prevQuests: QuestSeriesItem[];
	    exclusiveQuests: QuestSeriesItem[];
	
	    static createFrom(source: any = {}) {
	        return new QuestDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.title = source["title"];
	        this.details = source["details"];
	        this.objectives = source["objectives"];
	        this.offerRewardText = source["offerRewardText"];
	        this.endText = source["endText"];
	        this.questLevel = source["questLevel"];
	        this.minLevel = source["minLevel"];
	        this.type = source["type"];
	        this.zoneOrSort = source["zoneOrSort"];
	        this.categoryName = source["categoryName"];
	        this.requiredRaces = source["requiredRaces"];
	        this.requiredClasses = source["requiredClasses"];
	        this.rewardXp = source["rewardXp"];
	        this.rewardMoney = source["rewardMoney"];
	        this.rewardSpell = source["rewardSpell"];
	        this.rewardItems = this.convertValues(source["rewardItems"], QuestItem);
	        this.choiceItems = this.convertValues(source["choiceItems"], QuestItem);
	        this.reputation = this.convertValues(source["reputation"], QuestReputation);
	        this.starters = this.convertValues(source["starters"], QuestRelation);
	        this.enders = this.convertValues(source["enders"], QuestRelation);
	        this.series = this.convertValues(source["series"], QuestSeriesItem);
	        this.prevQuests = this.convertValues(source["prevQuests"], QuestSeriesItem);
	        this.exclusiveQuests = this.convertValues(source["exclusiveQuests"], QuestSeriesItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class SearchFilter {
	    query: string;
	    quality?: number[];
	    class?: number[];
	    subClass?: number[];
	    inventoryType?: number[];
	    minLevel?: number;
	    maxLevel?: number;
	    minReqLevel?: number;
	    maxReqLevel?: number;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.quality = source["quality"];
	        this.class = source["class"];
	        this.subClass = source["subClass"];
	        this.inventoryType = source["inventoryType"];
	        this.minLevel = source["minLevel"];
	        this.maxLevel = source["maxLevel"];
	        this.minReqLevel = source["minReqLevel"];
	        this.maxReqLevel = source["maxReqLevel"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class SearchResult {
	    items: Item[];
	    totalCount: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], Item);
	        this.totalCount = source["totalCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Spell {
	    entry: number;
	    name: string;
	    subname: string;
	    description: string;
	    iconId: number;
	
	    static createFrom(source: any = {}) {
	        return new Spell(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.subname = source["subname"];
	        this.description = source["description"];
	        this.iconId = source["iconId"];
	    }
	}
	export class SpellSkill {
	    id: number;
	    categoryId: number;
	    name: string;
	    spellCount: number;
	
	    static createFrom(source: any = {}) {
	        return new SpellSkill(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.categoryId = source["categoryId"];
	        this.name = source["name"];
	        this.spellCount = source["spellCount"];
	    }
	}
	export class SpellSkillCategory {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SpellSkillCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class TooltipData {
	    entry: number;
	    name: string;
	    quality: number;
	    itemLevel?: number;
	    binding?: string;
	    unique?: boolean;
	    itemType?: string;
	    slot?: string;
	    armor?: number;
	    damageRange?: string;
	    attackSpeed?: string;
	    dps?: string;
	    stats?: string[];
	    resistances?: string[];
	    effects?: string[];
	    requiredLevel?: number;
	    sellPrice?: number;
	    durability?: string;
	    classes?: string;
	    races?: string;
	    setInfo?: ItemSetInfo;
	
	    static createFrom(source: any = {}) {
	        return new TooltipData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry = source["entry"];
	        this.name = source["name"];
	        this.quality = source["quality"];
	        this.itemLevel = source["itemLevel"];
	        this.binding = source["binding"];
	        this.unique = source["unique"];
	        this.itemType = source["itemType"];
	        this.slot = source["slot"];
	        this.armor = source["armor"];
	        this.damageRange = source["damageRange"];
	        this.attackSpeed = source["attackSpeed"];
	        this.dps = source["dps"];
	        this.stats = source["stats"];
	        this.resistances = source["resistances"];
	        this.effects = source["effects"];
	        this.requiredLevel = source["requiredLevel"];
	        this.sellPrice = source["sellPrice"];
	        this.durability = source["durability"];
	        this.classes = source["classes"];
	        this.races = source["races"];
	        this.setInfo = this.convertValues(source["setInfo"], ItemSetInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

