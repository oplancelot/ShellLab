package database

type QuestDetail struct {
	Entry            int                `json:"entry"`
	Title            string             `json:"title"`
	Details          string             `json:"details"`
	Objectives       string             `json:"objectives"`
	OfferRewardText  string             `json:"offerRewardText"`
	EndText          string             `json:"endText"`
	MinLevel         int                `json:"minLevel"`
	QuestLevel       int                `json:"questLevel"`
	Type             int                `json:"type"`
	ZoneOrSort       int                `json:"zoneOrSort"`
	RequiredRaces    int                `json:"requiredRaces"`
	RequiredClasses  int                `json:"requiredClasses"`
	SrcItemID        int                `json:"srcItemId"`
	RewMoney         int                `json:"rewMoney"`
	RewMoneyMaxLevel int                `json:"rewMoneyMaxLevel"`
	RewXP            int                `json:"rewXp"`
	Rewards          []*QuestItem       `json:"rewards"`
	ChoiceRewards    []*QuestItem       `json:"choiceRewards"`
	Reputation       []*QuestReputation `json:"reputation"`
	Starters         []*QuestRelation   `json:"starters"`
	Enders           []*QuestRelation   `json:"enders"`
	Series           []*QuestSeriesItem `json:"series"`
	PrevQuests       []*QuestSeriesItem `json:"prevQuests"`
	ExclusiveQuests  []*QuestSeriesItem `json:"exclusiveQuests"`
}

type QuestSeriesItem struct {
	Entry int    `json:"entry"`
	Title string `json:"title"`
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

	var prevQuestID, nextQuestID, exclusiveGroup, nextQuestInChain int

	err := r.db.DB().QueryRow(`
		SELECT IFNULL(title,''), IFNULL(details,''), IFNULL(objectives,''), IFNULL(offer_reward_text,''), IFNULL(end_text,''),
			IFNULL(min_level,0), IFNULL(quest_level,0), IFNULL(type,0), IFNULL(zone_or_sort,0),
			IFNULL(rew_xp,0), IFNULL(rew_money,0),
			IFNULL(rew_item1,0), IFNULL(rew_item2,0), IFNULL(rew_item3,0), IFNULL(rew_item4,0),
			IFNULL(rew_item_count1,0), IFNULL(rew_item_count2,0), IFNULL(rew_item_count3,0), IFNULL(rew_item_count4,0),
			IFNULL(rew_choice_item1,0), IFNULL(rew_choice_item2,0), IFNULL(rew_choice_item3,0), IFNULL(rew_choice_item4,0), IFNULL(rew_choice_item5,0), IFNULL(rew_choice_item6,0),
			IFNULL(rew_choice_item_count1,0), IFNULL(rew_choice_item_count2,0), IFNULL(rew_choice_item_count3,0), IFNULL(rew_choice_item_count4,0), IFNULL(rew_choice_item_count5,0), IFNULL(rew_choice_item_count6,0),
			IFNULL(rew_rep_faction1,0), IFNULL(rew_rep_faction2,0), IFNULL(rew_rep_faction3,0), IFNULL(rew_rep_faction4,0), IFNULL(rew_rep_faction5,0),
			IFNULL(rew_rep_value1,0), IFNULL(rew_rep_value2,0), IFNULL(rew_rep_value3,0), IFNULL(rew_rep_value4,0), IFNULL(rew_rep_value5,0),
			IFNULL(prev_quest_id,0), IFNULL(next_quest_id,0), IFNULL(exclusive_group,0), IFNULL(next_quest_in_chain,0),
			IFNULL(required_races,0), IFNULL(required_classes,0), IFNULL(src_item_id,0), IFNULL(rew_money_max_level,0)
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
		&prevQuestID, &nextQuestID, &exclusiveGroup, &nextQuestInChain,
		&q.RequiredRaces, &q.RequiredClasses, &q.SrcItemID, &q.RewMoneyMaxLevel,
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

	// 5. Build Quest Chain (Series)
	q.Series = []*QuestSeriesItem{
		{Entry: entry, Title: q.Title},
	}

	// Walk backwards
	currEntry := entry
	for {
		var pID int
		var pTitle string
		err := r.db.DB().QueryRow(`
			SELECT entry, title 
			FROM quests 
			WHERE next_quest_in_chain = ?
			LIMIT 1
		`, currEntry).Scan(&pID, &pTitle)
		if err != nil || pID == 0 {
			break
		}
		// Prepend to series
		q.Series = append([]*QuestSeriesItem{{Entry: pID, Title: pTitle}}, q.Series...)
		currEntry = pID
	}

	// Walk forwards
	currNext := nextQuestInChain
	for currNext != 0 {
		var nID, nextNext int
		var nTitle string
		err := r.db.DB().QueryRow(`
			SELECT entry, title, next_quest_in_chain
			FROM quests 
			WHERE entry = ?
			LIMIT 1
		`, currNext).Scan(&nID, &nTitle, &nextNext)
		if err != nil || nID == 0 {
			break
		}
		q.Series = append(q.Series, &QuestSeriesItem{Entry: nID, Title: nTitle})
		currNext = nextNext
	}

	if len(q.Series) <= 1 {
		q.Series = nil
	}

	// 6. Fetch Prerequisites (prev_quest_id)
	if prevQuestID != 0 {
		var pID int
		var pTitle string
		// prev_quest_id can be single or negative (meaning all in group)
		// For now just handle single.
		if prevQuestID > 0 {
			r.db.DB().QueryRow("SELECT entry, title FROM quests WHERE entry = ?", prevQuestID).Scan(&pID, &pTitle)
			if pID != 0 {
				q.PrevQuests = append(q.PrevQuests, &QuestSeriesItem{Entry: pID, Title: pTitle})
			}
		} else if prevQuestID < 0 {
			// All quests with this same negative ID in exclusive_group must be done?
			// Actually usually means "any of". ExclusiveGroup is different.
		}
	}

	// 7. Fetch Exclusive Group
	if exclusiveGroup != 0 {
		rowsEX, err := r.db.DB().Query(`
			SELECT entry, title FROM quests 
			WHERE exclusive_group = ? AND entry != ?
		`, exclusiveGroup, entry)
		if err == nil {
			defer rowsEX.Close()
			for rowsEX.Next() {
				item := &QuestSeriesItem{}
				rowsEX.Scan(&item.Entry, &item.Title)
				q.ExclusiveQuests = append(q.ExclusiveQuests, item)
			}
		}
	}

	return q, nil
}
