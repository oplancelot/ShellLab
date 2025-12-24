package importers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"shelllab/backend/database/models"
)

// QuestImporter handles quest data imports
type QuestImporter struct {
	db *sql.DB
}

// NewQuestImporter creates a new quest importer
func NewQuestImporter(db *sql.DB) *QuestImporter {
	return &QuestImporter{db: db}
}

// ImportFromJSON imports quests from JSON into SQLite
func (q *QuestImporter) ImportFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var quests []models.QuestTemplateEntry
	if err := json.Unmarshal(data, &quests); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM quests")

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

	for _, qe := range quests {
		_, err := stmt.Exec(
			qe.Entry, qe.Title, qe.MinLevel, qe.QuestLevel, qe.Type, qe.ZoneOrSort, qe.Details, qe.Objectives, qe.OfferRewardText, qe.EndText,
			qe.RewXP, qe.RewOrReqMoney, qe.RewMoneyMaxLevel, qe.RewSpell,
			qe.RewItemId1, qe.RewItemId2, qe.RewItemId3, qe.RewItemId4,
			qe.RewItemCount1, qe.RewItemCount2, qe.RewItemCount3, qe.RewItemCount4,
			qe.RewChoiceItemId1, qe.RewChoiceItemId2, qe.RewChoiceItemId3, qe.RewChoiceItemId4, qe.RewChoiceItemId5, qe.RewChoiceItemId6,
			qe.RewChoiceItemCount1, qe.RewChoiceItemCount2, qe.RewChoiceItemCount3, qe.RewChoiceItemCount4, qe.RewChoiceItemCount5, qe.RewChoiceItemCount6,
			qe.RewRepFaction1, qe.RewRepFaction2, qe.RewRepFaction3, qe.RewRepFaction4, qe.RewRepFaction5,
			qe.RewRepValue1, qe.RewRepValue2, qe.RewRepValue3, qe.RewRepValue4, qe.RewRepValue5,
			qe.PrevQuestId, qe.NextQuestId, qe.ExclusiveGroup, qe.NextQuestInChain,
			qe.RequiredRaces, qe.RequiredClasses, qe.SrcItemId,
		)
		if err != nil {
			continue
		}
	}
	return tx.Commit()
}

// CheckAndImport checks if quests table is empty and imports if JSON exists
func (q *QuestImporter) CheckAndImport(dataDir string) error {
	var count int
	if err := q.db.QueryRow("SELECT COUNT(*) FROM quests").Scan(&count); err != nil {
		return nil
	}
	if count == 0 {
		path := fmt.Sprintf("%s/quests.json", dataDir)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Importing Quests...")
			return q.ImportFromJSON(path)
		}
	}
	return nil
}
