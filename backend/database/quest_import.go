package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// QuestTemplateEntry represents a quest record from JSON (MySQL export)
type QuestTemplateEntry struct {
	Entry               int    `json:"entry"`
	Title               string `json:"Title"`
	MinLevel            int    `json:"MinLevel"`
	QuestLevel          int    `json:"QuestLevel"`
	Type                int    `json:"Type"`
	ZoneOrSort          int    `json:"ZoneOrSort"`
	Details             string `json:"Details"`
	Objectives          string `json:"Objectives"`
	OfferRewardText     string `json:"OfferRewardText"`
	EndText             string `json:"EndText"`
	RewXP               int    `json:"RewXP"`
	RewOrReqMoney       int    `json:"RewOrReqMoney"`
	RewMoneyMaxLevel    int    `json:"RewMoneyMaxLevel"`
	RewSpell            int    `json:"RewSpell"`
	RewItemId1          int    `json:"RewItemId1"`
	RewItemId2          int    `json:"RewItemId2"`
	RewItemId3          int    `json:"RewItemId3"`
	RewItemId4          int    `json:"RewItemId4"`
	RewItemCount1       int    `json:"RewItemCount1"`
	RewItemCount2       int    `json:"RewItemCount2"`
	RewItemCount3       int    `json:"RewItemCount3"`
	RewItemCount4       int    `json:"RewItemCount4"`
	RewChoiceItemId1    int    `json:"RewChoiceItemId1"`
	RewChoiceItemId2    int    `json:"RewChoiceItemId2"`
	RewChoiceItemId3    int    `json:"RewChoiceItemId3"`
	RewChoiceItemId4    int    `json:"RewChoiceItemId4"`
	RewChoiceItemId5    int    `json:"RewChoiceItemId5"`
	RewChoiceItemId6    int    `json:"RewChoiceItemId6"`
	RewChoiceItemCount1 int    `json:"RewChoiceItemCount1"`
	RewChoiceItemCount2 int    `json:"RewChoiceItemCount2"`
	RewChoiceItemCount3 int    `json:"RewChoiceItemCount3"`
	RewChoiceItemCount4 int    `json:"RewChoiceItemCount4"`
	RewChoiceItemCount5 int    `json:"RewChoiceItemCount5"`
	RewChoiceItemCount6 int    `json:"RewChoiceItemCount6"`
	RewRepFaction1      int    `json:"RewRepFaction1"`
	RewRepFaction2      int    `json:"RewRepFaction2"`
	RewRepFaction3      int    `json:"RewRepFaction3"`
	RewRepFaction4      int    `json:"RewRepFaction4"`
	RewRepFaction5      int    `json:"RewRepFaction5"`
	RewRepValue1        int    `json:"RewRepValue1"`
	RewRepValue2        int    `json:"RewRepValue2"`
	RewRepValue3        int    `json:"RewRepValue3"`
	RewRepValue4        int    `json:"RewRepValue4"`
	RewRepValue5        int    `json:"RewRepValue5"`
	PrevQuestId         int    `json:"PrevQuestId"`
	NextQuestId         int    `json:"NextQuestId"`
	ExclusiveGroup      int    `json:"ExclusiveGroup"`
	NextQuestInChain    int    `json:"NextQuestInChain"`
	RequiredRaces       int    `json:"RequiredRaces"`
	RequiredClasses     int    `json:"RequiredClasses"`
	SrcItemId           int    `json:"SrcItemId"`
}

// ImportQuestsFromJSON imports quests from JSON into SQLite
func (r *ItemRepository) ImportQuestsFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", jsonPath, err)
	}

	var quests []QuestTemplateEntry
	if err := json.Unmarshal(data, &quests); err != nil {
		return fmt.Errorf("failed to parse JSON %s: %w", jsonPath, err)
	}

	tx, err := r.db.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM quests")

	// Note: We use named matching for simplicity or verbose positional args
	stmt, err := tx.Prepare(`
		INSERT INTO quests (
			entry, title, min_level, quest_level, type, zone_or_sort, details, objectives, offer_reward_text, end_text,
			rew_xp, rew_money, rew_money_max_level, rew_spell,
			rew_item1, rew_item2, rew_item3, rew_item4,
			rew_item_count1, rew_item_count2, rew_item_count3, rew_item_count4,
			rew_choice_item1, rew_choice_item2, rew_choice_item3, rew_choice_item4, rew_choice_item5, rew_choice_item6,
			rew_choice_item_count1, rew_choice_item_count2, rew_choice_item_count3, rew_choice_item_count4, rew_choice_item_count5, rew_choice_item_count6,
			rew_rep_faction1, rew_rep_faction2, rew_rep_faction3, rew_rep_faction4, rew_rep_faction5,
			rew_rep_value1, rew_rep_value2, rew_rep_value3, rew_rep_value4, rew_rep_value5,
			prev_quest_id, next_quest_id, exclusive_group, next_quest_in_chain,
			required_races, required_classes, src_item_id
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range quests {
		_, err := stmt.Exec(
			q.Entry, q.Title, q.MinLevel, q.QuestLevel, q.Type, q.ZoneOrSort, q.Details, q.Objectives, q.OfferRewardText, q.EndText,
			q.RewXP, q.RewOrReqMoney, q.RewMoneyMaxLevel, q.RewSpell,
			q.RewItemId1, q.RewItemId2, q.RewItemId3, q.RewItemId4,
			q.RewItemCount1, q.RewItemCount2, q.RewItemCount3, q.RewItemCount4,
			q.RewChoiceItemId1, q.RewChoiceItemId2, q.RewChoiceItemId3, q.RewChoiceItemId4, q.RewChoiceItemId5, q.RewChoiceItemId6,
			q.RewChoiceItemCount1, q.RewChoiceItemCount2, q.RewChoiceItemCount3, q.RewChoiceItemCount4, q.RewChoiceItemCount5, q.RewChoiceItemCount6,
			q.RewRepFaction1, q.RewRepFaction2, q.RewRepFaction3, q.RewRepFaction4, q.RewRepFaction5,
			q.RewRepValue1, q.RewRepValue2, q.RewRepValue3, q.RewRepValue4, q.RewRepValue5,
			q.PrevQuestId, q.NextQuestId, q.ExclusiveGroup, q.NextQuestInChain,
			q.RequiredRaces, q.RequiredClasses, q.SrcItemId,
		)
		if err != nil {
			fmt.Printf("Warning: Failed to import quest %d: %v\n", q.Entry, err)
			continue
		}
	}

	return tx.Commit()
}

// ImportAllQuests checks and imports quests
func (r *ItemRepository) ImportAllQuests(dataDir string) error {
	path := fmt.Sprintf("%s/quests.json", dataDir)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Importing Quests from %s...\n", path)
		if err := r.ImportQuestsFromJSON(path); err != nil {
			fmt.Printf("Error importing quests: %v\n", err)
			return err
		}
		fmt.Println("✓ Quests imported successfully!")
	}
	return nil
}

// CheckAndImportQuests checks if quests table is empty and imports if JSON exists
func (r *ItemRepository) CheckAndImportQuests(dataDir string) error {
	var count int
	err := r.db.DB().QueryRow("SELECT COUNT(*) FROM quests").Scan(&count)
	if err != nil {
		return nil
	}

	if count == 0 {
		return r.ImportAllQuests(dataDir)
	}
	return nil
}
