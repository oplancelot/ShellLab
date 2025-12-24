package models

// Creature represents a WoW NPC
type Creature struct {
	Entry     int    `json:"entry"`
	Name      string `json:"name"`
	Subname   string `json:"subname,omitempty"`
	LevelMin  int    `json:"levelMin"`
	LevelMax  int    `json:"levelMax"`
	HealthMin int    `json:"healthMin"`
	HealthMax int    `json:"healthMax"`
	ManaMin   int    `json:"manaMin"`
	ManaMax   int    `json:"manaMax"`
	Type      int    `json:"type"`
	TypeName  string `json:"typeName"`
	Rank      int    `json:"rank"`
	RankName  string `json:"rankName"`
	Faction   int    `json:"faction"`
	NPCFlags  int    `json:"npcFlags"`
}

// CreatureType represents a creature type category
type CreatureType struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// CreatureDetail includes a creature with its loot and quests
type CreatureDetail struct {
	*Creature
	Loot         []*LootItem      `json:"loot"`
	StartsQuests []*QuestRelation `json:"startsQuests"`
	EndsQuests   []*QuestRelation `json:"endsQuests"`
}

// CreatureTemplateEntry represents a creature for JSON import
type CreatureTemplateEntry struct {
	Entry            int    `json:"entry"`
	Name             string `json:"name"`
	Subname          string `json:"subname"`
	LevelMin         int    `json:"level_min"`
	LevelMax         int    `json:"level_max"`
	HealthMin        int    `json:"health_min"`
	HealthMax        int    `json:"health_max"`
	ManaMin          int    `json:"mana_min"`
	ManaMax          int    `json:"mana_max"`
	CreatureType     int    `json:"creature_type"`
	CreatureRank     int    `json:"creature_rank"`
	Faction          int    `json:"faction"`
	NPCFlags         int    `json:"npc_flags"`
	LootID           int    `json:"loot_id"`
	SkinLootID       int    `json:"skinning_loot_id"`
	PickpocketLootID int    `json:"pickpocket_loot_id"`
}
