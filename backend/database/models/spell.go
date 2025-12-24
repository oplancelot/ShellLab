package models

// Spell represents a WoW spell
type Spell struct {
	Entry       int    `json:"entry"`
	Name        string `json:"name"`
	SubName     string `json:"subname"` // Rank or subtext
	Description string `json:"description"`
	IconID      int    `json:"iconId"`
}

// SpellSkillCategory represents a top-level category for spells
type SpellSkillCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SpellSkill represents a skill that contains spells
type SpellSkill struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"categoryId"`
	Name       string `json:"name"`
	SpellCount int    `json:"spellCount"`
}

// SpellEntry represents a spell for JSON import
type SpellEntry struct {
	Entry             int    `json:"entry"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	EffectBasePoints1 int    `json:"effectBasePoints1"`
	EffectBasePoints2 int    `json:"effectBasePoints2"`
	EffectBasePoints3 int    `json:"effectBasePoints3"`
	EffectDieSides1   int    `json:"effectDieSides1"`
	EffectDieSides2   int    `json:"effectDieSides2"`
	EffectDieSides3   int    `json:"effectDieSides3"`
}
