package database

type QuestDetail struct {
	Entry           int                `json:"entry"`
	Title           string             `json:"title"`
	Details         string             `json:"details"`
	Objectives      string             `json:"objectives"`
	OfferRewardText string             `json:"offerRewardText"`
	EndText         string             `json:"endText"`
	MinLevel        int                `json:"minLevel"`
	QuestLevel      int                `json:"questLevel"`
	Type            int                `json:"type"`
	ZoneOrSort      int                `json:"zoneOrSort"`
	RewMoney        int                `json:"rewMoney"`
	RewXP           int                `json:"rewXp"`
	Rewards         []*QuestItem       `json:"rewards"`
	ChoiceRewards   []*QuestItem       `json:"choiceRewards"`
	Reputation      []*QuestReputation `json:"reputation"`
	Starters        []*QuestRelation   `json:"starters"`
	Enders          []*QuestRelation   `json:"enders"`
}

type QuestItem struct {
	Entry   int    `json:"entry"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Count   int    `json:"count"`
	Quality int    `json:"quality"`
}

type QuestReputation struct {
	FactionID int    `json:"factionId"`
	Name      string `json:"name"`
	Value     int    `json:"value"`
}

type QuestRelation struct {
	Entry int    `json:"entry"`
	Name  string `json:"name"`
	Type  string `json:"type"` // "npc" or "object"
}

// GetQuestDetail returns full details for a quest
func (r *ItemRepository) GetQuestDetail(entry int) (*QuestDetail, error) {
	q := &QuestDetail{
		Entry:         entry,
		Rewards:       []*QuestItem{},
		ChoiceRewards: []*QuestItem{},
		Reputation:    []*QuestReputation{},
		Starters:      []*QuestRelation{},
		Enders:        []*QuestRelation{},
	}

	// 1. Fetch Quest Data
	var rewIt1, rewIt2, rewIt3, rewIt4 int
	var rewCt1, rewCt2, rewCt3, rewCt4 int
	var chIt1, chIt2, chIt3, chIt4, chIt5, chIt6 int
	var chCt1, chCt2, chCt3, chCt4, chCt5, chCt6 int
	var repF1, repF2, repF3, repF4, repF5 int
	var repV1, repV2, repV3, repV4, repV5 int

	err := r.db.DB().QueryRow(`
		SELECT title, details, objectives, offer_reward_text, end_text,
			min_level, quest_level, type, zone_or_sort,
			rew_xp, rew_money,
			rew_item1, rew_item2, rew_item3, rew_item4,
			rew_item_count1, rew_item_count2, rew_item_count3, rew_item_count4,
			rew_choice_item1, rew_choice_item2, rew_choice_item3, rew_choice_item4, rew_choice_item5, rew_choice_item6,
			rew_choice_item_count1, rew_choice_item_count2, rew_choice_item_count3, rew_choice_item_count4, rew_choice_item_count5, rew_choice_item_count6,
			rew_rep_faction1, rew_rep_faction2, rew_rep_faction3, rew_rep_faction4, rew_rep_faction5,
			rew_rep_value1, rew_rep_value2, rew_rep_value3, rew_rep_value4, rew_rep_value5
		FROM quests WHERE entry = ?
	`, entry).Scan(
		&q.Title, &q.Details, &q.Objectives, &q.OfferRewardText, &q.EndText,
		&q.MinLevel, &q.QuestLevel, &q.Type, &q.ZoneOrSort,
		&q.RewXP, &q.RewMoney,
		&rewIt1, &rewIt2, &rewIt3, &rewIt4,
		&rewCt1, &rewCt2, &rewCt3, &rewCt4,
		&chIt1, &chIt2, &chIt3, &chIt4, &chIt5, &chIt6,
		&chCt1, &chCt2, &chCt3, &chCt4, &chCt5, &chCt6,
		&repF1, &repF2, &repF3, &repF4, &repF5,
		&repV1, &repV2, &repV3, &repV4, &repV5,
	)
	if err != nil {
		return nil, err
	}

	// Helper to fetch item info
	fetchItem := func(id, count int) *QuestItem {
		if id == 0 {
			return nil
		}
		var name, icon string
		var quality int
		r.db.DB().QueryRow("SELECT name, icon, quality FROM items WHERE entry = ?", id).Scan(&name, &icon, &quality)
		return &QuestItem{Entry: id, Name: name, Icon: icon, Count: count, Quality: quality}
	}

	// 2. Process Rewards
	for i, id := range []int{rewIt1, rewIt2, rewIt3, rewIt4} {
		if item := fetchItem(id, []int{rewCt1, rewCt2, rewCt3, rewCt4}[i]); item != nil {
			q.Rewards = append(q.Rewards, item)
		}
	}
	for i, id := range []int{chIt1, chIt2, chIt3, chIt4, chIt5, chIt6} {
		if item := fetchItem(id, []int{chCt1, chCt2, chCt3, chCt4, chCt5, chCt6}[i]); item != nil {
			q.ChoiceRewards = append(q.ChoiceRewards, item)
		}
	}

	// 3. Process Reputation
	fetchRep := func(id, val int) {
		if id == 0 {
			return
		}
		var name string
		r.db.DB().QueryRow("SELECT name FROM factions WHERE id = ?", id).Scan(&name)
		q.Reputation = append(q.Reputation, &QuestReputation{FactionID: id, Name: name, Value: val})
	}
	fetchRep(repF1, repV1)
	fetchRep(repF2, repV2)
	fetchRep(repF3, repV3)
	fetchRep(repF4, repV4)
	fetchRep(repF5, repV5)

	// 4. Fetch Relations (Starters/Enders)
	// NPC Starters
	rows, err := r.db.DB().Query(`
		SELECT c.entry, c.name FROM npc_quest_start n
		JOIN creatures c ON n.entry = c.entry
		WHERE n.quest = ?
	`, entry)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			rel := &QuestRelation{Type: "npc"}
			rows.Scan(&rel.Entry, &rel.Name)
			q.Starters = append(q.Starters, rel)
		}
	}

	// GO Starters
	rowsG, err := r.db.DB().Query(`
		SELECT g.entry, g.name FROM go_quest_start n
		JOIN game_objects g ON n.entry = g.entry
		WHERE n.quest = ?
	`, entry)
	if err == nil {
		defer rowsG.Close()
		for rowsG.Next() {
			rel := &QuestRelation{Type: "object"}
			rowsG.Scan(&rel.Entry, &rel.Name)
			q.Starters = append(q.Starters, rel)
		}
	}

	// NPC Enders
	rowsE, err := r.db.DB().Query(`
		SELECT c.entry, c.name FROM npc_quest_end n
		JOIN creatures c ON n.entry = c.entry
		WHERE n.quest = ?
	`, entry)
	if err == nil {
		defer rowsE.Close()
		for rowsE.Next() {
			rel := &QuestRelation{Type: "npc"}
			rowsE.Scan(&rel.Entry, &rel.Name)
			q.Enders = append(q.Enders, rel)
		}
	}

	// GO Enders
	rowsGE, err := r.db.DB().Query(`
		SELECT g.entry, g.name FROM go_quest_end n
		JOIN game_objects g ON n.entry = g.entry
		WHERE n.quest = ?
	`, entry)
	if err == nil {
		defer rowsGE.Close()
		for rowsGE.Next() {
			rel := &QuestRelation{Type: "object"}
			rowsGE.Scan(&rel.Entry, &rel.Name)
			q.Enders = append(q.Enders, rel)
		}
	}

	return q, nil
}
