package importers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "os"
    "shelllab/backend/database/models"
)

type GeneratedImporter struct {
    db *sql.DB
}

func NewGeneratedImporter(db *sql.DB) *GeneratedImporter {
    return &GeneratedImporter{db: db}
}


func (i *GeneratedImporter) ImportItemTemplate(jsonPath string) error {
    file, err := os.Open(jsonPath)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    // Expect array
    if _, err := decoder.Token(); err != nil {
        return err
    }

    tx, err := i.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `INSERT OR REPLACE INTO item_template (entry,class,subclass,name,description,display_id,quality,flags,buy_count,buy_price,sell_price,inventory_type,allowable_class,allowable_race,item_level,required_level,required_skill,required_skill_rank,required_spell,required_honor_rank,required_city_rank,required_reputation_faction,required_reputation_rank,max_count,stackable,container_slots,stat_type1,stat_value1,stat_type2,stat_value2,stat_type3,stat_value3,stat_type4,stat_value4,stat_type5,stat_value5,stat_type6,stat_value6,stat_type7,stat_value7,stat_type8,stat_value8,stat_type9,stat_value9,stat_type10,stat_value10,delay,range_mod,ammo_type,dmg_min1,dmg_max1,dmg_type1,dmg_min2,dmg_max2,dmg_type2,dmg_min3,dmg_max3,dmg_type3,dmg_min4,dmg_max4,dmg_type4,dmg_min5,dmg_max5,dmg_type5,block,armor,holy_res,fire_res,nature_res,frost_res,shadow_res,arcane_res,spellid_1,spelltrigger_1,spellcharges_1,spellppmrate_1,spellcooldown_1,spellcategory_1,spellcategorycooldown_1,spellid_2,spelltrigger_2,spellcharges_2,spellppmrate_2,spellcooldown_2,spellcategory_2,spellcategorycooldown_2,spellid_3,spelltrigger_3,spellcharges_3,spellppmrate_3,spellcooldown_3,spellcategory_3,spellcategorycooldown_3,spellid_4,spelltrigger_4,spellcharges_4,spellppmrate_4,spellcooldown_4,spellcategory_4,spellcategorycooldown_4,spellid_5,spelltrigger_5,spellcharges_5,spellppmrate_5,spellcooldown_5,spellcategory_5,spellcategorycooldown_5,bonding,page_text,page_language,page_material,start_quest,lock_id,material,sheath,random_property,set_id,max_durability,area_bound,map_bound,duration,bag_family,disenchant_id,food_type,min_money_loot,max_money_loot,wrapped_gift,extra_flags,other_team_entry,script_name) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    count := 0
    for decoder.More() {
        var item models.ItemTemplateFull
        if err := decoder.Decode(&item); err != nil {
            continue
        }

        _, err = stmt.Exec(item.Entry, item.Class, item.Subclass, item.Name, item.Description, item.DisplayId, item.Quality, item.Flags, item.BuyCount, item.BuyPrice, item.SellPrice, item.InventoryType, item.AllowableClass, item.AllowableRace, item.ItemLevel, item.RequiredLevel, item.RequiredSkill, item.RequiredSkillRank, item.RequiredSpell, item.RequiredHonorRank, item.RequiredCityRank, item.RequiredReputationFaction, item.RequiredReputationRank, item.MaxCount, item.Stackable, item.ContainerSlots, item.StatType1, item.StatValue1, item.StatType2, item.StatValue2, item.StatType3, item.StatValue3, item.StatType4, item.StatValue4, item.StatType5, item.StatValue5, item.StatType6, item.StatValue6, item.StatType7, item.StatValue7, item.StatType8, item.StatValue8, item.StatType9, item.StatValue9, item.StatType10, item.StatValue10, item.Delay, item.RangeMod, item.AmmoType, item.DmgMin1, item.DmgMax1, item.DmgType1, item.DmgMin2, item.DmgMax2, item.DmgType2, item.DmgMin3, item.DmgMax3, item.DmgType3, item.DmgMin4, item.DmgMax4, item.DmgType4, item.DmgMin5, item.DmgMax5, item.DmgType5, item.Block, item.Armor, item.HolyRes, item.FireRes, item.NatureRes, item.FrostRes, item.ShadowRes, item.ArcaneRes, item.Spellid1, item.Spelltrigger1, item.Spellcharges1, item.Spellppmrate1, item.Spellcooldown1, item.Spellcategory1, item.Spellcategorycooldown1, item.Spellid2, item.Spelltrigger2, item.Spellcharges2, item.Spellppmrate2, item.Spellcooldown2, item.Spellcategory2, item.Spellcategorycooldown2, item.Spellid3, item.Spelltrigger3, item.Spellcharges3, item.Spellppmrate3, item.Spellcooldown3, item.Spellcategory3, item.Spellcategorycooldown3, item.Spellid4, item.Spelltrigger4, item.Spellcharges4, item.Spellppmrate4, item.Spellcooldown4, item.Spellcategory4, item.Spellcategorycooldown4, item.Spellid5, item.Spelltrigger5, item.Spellcharges5, item.Spellppmrate5, item.Spellcooldown5, item.Spellcategory5, item.Spellcategorycooldown5, item.Bonding, item.PageText, item.PageLanguage, item.PageMaterial, item.StartQuest, item.LockId, item.Material, item.Sheath, item.RandomProperty, item.SetId, item.MaxDurability, item.AreaBound, item.MapBound, item.Duration, item.BagFamily, item.DisenchantId, item.FoodType, item.MinMoneyLoot, item.MaxMoneyLoot, item.WrappedGift, item.ExtraFlags, item.OtherTeamEntry, item.ScriptName)
        if err != nil {
            // fmt.Printf("Error importing %s: %v\n", "item_template", err)
            continue
        }
        count++
    }

    return tx.Commit()
}

func (i *GeneratedImporter) ImportCreatureTemplate(jsonPath string) error {
    file, err := os.Open(jsonPath)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    // Expect array
    if _, err := decoder.Token(); err != nil {
        return err
    }

    tx, err := i.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `INSERT OR REPLACE INTO creature_template (entry,display_id1,display_id2,display_id3,display_id4,mount_display_id,name,subname,gossip_menu_id,level_min,level_max,health_min,health_max,mana_min,mana_max,armor,faction,npc_flags,speed_walk,speed_run,scale,detection_range,call_for_help_range,leash_range,rank,xp_multiplier,dmg_min,dmg_max,dmg_school,attack_power,dmg_multiplier,base_attack_time,ranged_attack_time,unit_class,unit_flags,dynamic_flags,beast_family,trainer_type,trainer_spell,trainer_class,trainer_race,ranged_dmg_min,ranged_dmg_max,ranged_attack_power,type,type_flags,loot_id,pickpocket_loot_id,skinning_loot_id,holy_res,fire_res,nature_res,frost_res,shadow_res,arcane_res,spell_id1,spell_id2,spell_id3,spell_id4,spell_list_id,pet_spell_list_id,spawn_spell_id,auras,gold_min,gold_max,ai_name,movement_type,inhabit_type,civilian,racial_leader,regeneration,equipment_id,trainer_id,vendor_id,mechanic_immune_mask,school_immune_mask,immunity_flags,flags_extra,phase_quest_id,script_name) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    count := 0
    for decoder.More() {
        var item models.CreatureTemplateFull
        if err := decoder.Decode(&item); err != nil {
            continue
        }

        _, err = stmt.Exec(item.Entry, item.DisplayId1, item.DisplayId2, item.DisplayId3, item.DisplayId4, item.MountDisplayId, item.Name, item.Subname, item.GossipMenuId, item.LevelMin, item.LevelMax, item.HealthMin, item.HealthMax, item.ManaMin, item.ManaMax, item.Armor, item.Faction, item.NpcFlags, item.SpeedWalk, item.SpeedRun, item.Scale, item.DetectionRange, item.CallForHelpRange, item.LeashRange, item.Rank, item.XpMultiplier, item.DmgMin, item.DmgMax, item.DmgSchool, item.AttackPower, item.DmgMultiplier, item.BaseAttackTime, item.RangedAttackTime, item.UnitClass, item.UnitFlags, item.DynamicFlags, item.BeastFamily, item.TrainerType, item.TrainerSpell, item.TrainerClass, item.TrainerRace, item.RangedDmgMin, item.RangedDmgMax, item.RangedAttackPower, item.Type, item.TypeFlags, item.LootId, item.PickpocketLootId, item.SkinningLootId, item.HolyRes, item.FireRes, item.NatureRes, item.FrostRes, item.ShadowRes, item.ArcaneRes, item.SpellId1, item.SpellId2, item.SpellId3, item.SpellId4, item.SpellListId, item.PetSpellListId, item.SpawnSpellId, item.Auras, item.GoldMin, item.GoldMax, item.AiName, item.MovementType, item.InhabitType, item.Civilian, item.RacialLeader, item.Regeneration, item.EquipmentId, item.TrainerId, item.VendorId, item.MechanicImmuneMask, item.SchoolImmuneMask, item.ImmunityFlags, item.FlagsExtra, item.PhaseQuestId, item.ScriptName)
        if err != nil {
            // fmt.Printf("Error importing %s: %v\n", "creature_template", err)
            continue
        }
        count++
    }

    return tx.Commit()
}

func (i *GeneratedImporter) ImportQuestTemplate(jsonPath string) error {
    file, err := os.Open(jsonPath)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    // Expect array
    if _, err := decoder.Token(); err != nil {
        return err
    }

    tx, err := i.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `INSERT OR REPLACE INTO quest_template (entry,Method,ZoneOrSort,MinLevel,MaxLevel,QuestLevel,Type,RequiredClasses,RequiredRaces,RequiredSkill,RequiredSkillValue,RequiredCondition,RepObjectiveFaction,RepObjectiveValue,RequiredMinRepFaction,RequiredMinRepValue,RequiredMaxRepFaction,RequiredMaxRepValue,SuggestedPlayers,LimitTime,QuestFlags,SpecialFlags,PrevQuestId,NextQuestId,ExclusiveGroup,NextQuestInChain,SrcItemId,SrcItemCount,SrcSpell,Title,Details,Objectives,OfferRewardText,RequestItemsText,EndText,ObjectiveText1,ObjectiveText2,ObjectiveText3,ObjectiveText4,ReqItemId1,ReqItemId2,ReqItemId3,ReqItemId4,ReqItemCount1,ReqItemCount2,ReqItemCount3,ReqItemCount4,ReqSourceId1,ReqSourceId2,ReqSourceId3,ReqSourceId4,ReqSourceCount1,ReqSourceCount2,ReqSourceCount3,ReqSourceCount4,ReqCreatureOrGOId1,ReqCreatureOrGOId2,ReqCreatureOrGOId3,ReqCreatureOrGOId4,ReqCreatureOrGOCount1,ReqCreatureOrGOCount2,ReqCreatureOrGOCount3,ReqCreatureOrGOCount4,ReqSpellCast1,ReqSpellCast2,ReqSpellCast3,ReqSpellCast4,RewChoiceItemId1,RewChoiceItemId2,RewChoiceItemId3,RewChoiceItemId4,RewChoiceItemId5,RewChoiceItemId6,RewChoiceItemCount1,RewChoiceItemCount2,RewChoiceItemCount3,RewChoiceItemCount4,RewChoiceItemCount5,RewChoiceItemCount6,RewItemId1,RewItemId2,RewItemId3,RewItemId4,RewItemCount1,RewItemCount2,RewItemCount3,RewItemCount4,RewRepFaction1,RewRepFaction2,RewRepFaction3,RewRepFaction4,RewRepFaction5,RewRepValue1,RewRepValue2,RewRepValue3,RewRepValue4,RewRepValue5,RewXP,RewOrReqMoney,RewMoneyMaxLevel,RewSpell,RewSpellCast,RewMailTemplateId,RewMailDelaySecs,RewMailMoney,PointMapId,PointX,PointY,PointOpt,DetailsEmote1,DetailsEmote2,DetailsEmote3,DetailsEmote4,DetailsEmoteDelay1,DetailsEmoteDelay2,DetailsEmoteDelay3,DetailsEmoteDelay4,IncompleteEmote,CompleteEmote,OfferRewardEmote1,OfferRewardEmote2,OfferRewardEmote3,OfferRewardEmote4,OfferRewardEmoteDelay1,OfferRewardEmoteDelay2,OfferRewardEmoteDelay3,OfferRewardEmoteDelay4,StartScript,CompleteScript) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    count := 0
    for decoder.More() {
        var item models.QuestTemplateFull
        if err := decoder.Decode(&item); err != nil {
            continue
        }

        _, err = stmt.Exec(item.Entry, item.Method, item.Zoneorsort, item.Minlevel, item.Maxlevel, item.Questlevel, item.Type, item.Requiredclasses, item.Requiredraces, item.Requiredskill, item.Requiredskillvalue, item.Requiredcondition, item.Repobjectivefaction, item.Repobjectivevalue, item.Requiredminrepfaction, item.Requiredminrepvalue, item.Requiredmaxrepfaction, item.Requiredmaxrepvalue, item.Suggestedplayers, item.Limittime, item.Questflags, item.Specialflags, item.Prevquestid, item.Nextquestid, item.Exclusivegroup, item.Nextquestinchain, item.Srcitemid, item.Srcitemcount, item.Srcspell, item.Title, item.Details, item.Objectives, item.Offerrewardtext, item.Requestitemstext, item.Endtext, item.Objectivetext1, item.Objectivetext2, item.Objectivetext3, item.Objectivetext4, item.Reqitemid1, item.Reqitemid2, item.Reqitemid3, item.Reqitemid4, item.Reqitemcount1, item.Reqitemcount2, item.Reqitemcount3, item.Reqitemcount4, item.Reqsourceid1, item.Reqsourceid2, item.Reqsourceid3, item.Reqsourceid4, item.Reqsourcecount1, item.Reqsourcecount2, item.Reqsourcecount3, item.Reqsourcecount4, item.Reqcreatureorgoid1, item.Reqcreatureorgoid2, item.Reqcreatureorgoid3, item.Reqcreatureorgoid4, item.Reqcreatureorgocount1, item.Reqcreatureorgocount2, item.Reqcreatureorgocount3, item.Reqcreatureorgocount4, item.Reqspellcast1, item.Reqspellcast2, item.Reqspellcast3, item.Reqspellcast4, item.Rewchoiceitemid1, item.Rewchoiceitemid2, item.Rewchoiceitemid3, item.Rewchoiceitemid4, item.Rewchoiceitemid5, item.Rewchoiceitemid6, item.Rewchoiceitemcount1, item.Rewchoiceitemcount2, item.Rewchoiceitemcount3, item.Rewchoiceitemcount4, item.Rewchoiceitemcount5, item.Rewchoiceitemcount6, item.Rewitemid1, item.Rewitemid2, item.Rewitemid3, item.Rewitemid4, item.Rewitemcount1, item.Rewitemcount2, item.Rewitemcount3, item.Rewitemcount4, item.Rewrepfaction1, item.Rewrepfaction2, item.Rewrepfaction3, item.Rewrepfaction4, item.Rewrepfaction5, item.Rewrepvalue1, item.Rewrepvalue2, item.Rewrepvalue3, item.Rewrepvalue4, item.Rewrepvalue5, item.Rewxp, item.Reworreqmoney, item.Rewmoneymaxlevel, item.Rewspell, item.Rewspellcast, item.Rewmailtemplateid, item.Rewmaildelaysecs, item.Rewmailmoney, item.Pointmapid, item.Pointx, item.Pointy, item.Pointopt, item.Detailsemote1, item.Detailsemote2, item.Detailsemote3, item.Detailsemote4, item.Detailsemotedelay1, item.Detailsemotedelay2, item.Detailsemotedelay3, item.Detailsemotedelay4, item.Incompleteemote, item.Completeemote, item.Offerrewardemote1, item.Offerrewardemote2, item.Offerrewardemote3, item.Offerrewardemote4, item.Offerrewardemotedelay1, item.Offerrewardemotedelay2, item.Offerrewardemotedelay3, item.Offerrewardemotedelay4, item.Startscript, item.Completescript)
        if err != nil {
            // fmt.Printf("Error importing %s: %v\n", "quest_template", err)
            continue
        }
        count++
    }

    return tx.Commit()
}

func (i *GeneratedImporter) ImportSpellTemplate(jsonPath string) error {
    file, err := os.Open(jsonPath)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    // Expect array
    if _, err := decoder.Token(); err != nil {
        return err
    }

    tx, err := i.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `INSERT OR REPLACE INTO spell_template (entry,school,category,castUI,dispel,mechanic,attributes,attributesEx,attributesEx2,attributesEx3,attributesEx4,stances,stancesNot,targets,targetCreatureType,requiresSpellFocus,casterAuraState,targetAuraState,castingTimeIndex,recoveryTime,categoryRecoveryTime,interruptFlags,auraInterruptFlags,channelInterruptFlags,procFlags,procChance,procCharges,maxLevel,baseLevel,spellLevel,durationIndex,powerType,manaCost,manCostPerLevel,manaPerSecond,manaPerSecondPerLevel,rangeIndex,speed,modelNextSpell,stackAmount,totem1,totem2,reagent1,reagent2,reagent3,reagent4,reagent5,reagent6,reagent7,reagent8,reagentCount1,reagentCount2,reagentCount3,reagentCount4,reagentCount5,reagentCount6,reagentCount7,reagentCount8,equippedItemClass,equippedItemSubClassMask,equippedItemInventoryTypeMask,effect1,effect2,effect3,effectDieSides1,effectDieSides2,effectDieSides3,effectBaseDice1,effectBaseDice2,effectBaseDice3,effectDicePerLevel1,effectDicePerLevel2,effectDicePerLevel3,effectRealPointsPerLevel1,effectRealPointsPerLevel2,effectRealPointsPerLevel3,effectBasePoints1,effectBasePoints2,effectBasePoints3,effectBonusCoefficient1,effectBonusCoefficient2,effectBonusCoefficient3,effectMechanic1,effectMechanic2,effectMechanic3,effectImplicitTargetA1,effectImplicitTargetA2,effectImplicitTargetA3,effectImplicitTargetB1,effectImplicitTargetB2,effectImplicitTargetB3,effectRadiusIndex1,effectRadiusIndex2,effectRadiusIndex3,effectApplyAuraName1,effectApplyAuraName2,effectApplyAuraName3,effectAmplitude1,effectAmplitude2,effectAmplitude3,effectMultipleValue1,effectMultipleValue2,effectMultipleValue3,effectChainTarget1,effectChainTarget2,effectChainTarget3,effectItemType1,effectItemType2,effectItemType3,effectMiscValue1,effectMiscValue2,effectMiscValue3,effectTriggerSpell1,effectTriggerSpell2,effectTriggerSpell3,effectPointsPerComboPoint1,effectPointsPerComboPoint2,effectPointsPerComboPoint3,spellVisual1,spellVisual2,spellIconId,activeIconId,spellPriority,name,nameFlags,nameSubtext,nameSubtextFlags,description,descriptionFlags,auraDescription,auraDescriptionFlags,manaCostPercentage,startRecoveryCategory,startRecoveryTime,minTargetLevel,maxTargetLevel,spellFamilyName,spellFamilyFlags,maxAffectedTargets,dmgClass,preventionType,stanceBarOrder,dmgMultiplier1,dmgMultiplier2,dmgMultiplier3,minFactionId,minReputation,requiredAuraVision,customFlags) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    count := 0
    for decoder.More() {
        var item models.SpellTemplateFull
        if err := decoder.Decode(&item); err != nil {
            continue
        }

        _, err = stmt.Exec(item.Entry, item.School, item.Category, item.Castui, item.Dispel, item.Mechanic, item.Attributes, item.Attributesex, item.Attributesex2, item.Attributesex3, item.Attributesex4, item.Stances, item.Stancesnot, item.Targets, item.Targetcreaturetype, item.Requiresspellfocus, item.Casteraurastate, item.Targetaurastate, item.Castingtimeindex, item.Recoverytime, item.Categoryrecoverytime, item.Interruptflags, item.Aurainterruptflags, item.Channelinterruptflags, item.Procflags, item.Procchance, item.Proccharges, item.Maxlevel, item.Baselevel, item.Spelllevel, item.Durationindex, item.Powertype, item.Manacost, item.Mancostperlevel, item.Manapersecond, item.Manapersecondperlevel, item.Rangeindex, item.Speed, item.Modelnextspell, item.Stackamount, item.Totem1, item.Totem2, item.Reagent1, item.Reagent2, item.Reagent3, item.Reagent4, item.Reagent5, item.Reagent6, item.Reagent7, item.Reagent8, item.Reagentcount1, item.Reagentcount2, item.Reagentcount3, item.Reagentcount4, item.Reagentcount5, item.Reagentcount6, item.Reagentcount7, item.Reagentcount8, item.Equippeditemclass, item.Equippeditemsubclassmask, item.Equippediteminventorytypemask, item.Effect1, item.Effect2, item.Effect3, item.Effectdiesides1, item.Effectdiesides2, item.Effectdiesides3, item.Effectbasedice1, item.Effectbasedice2, item.Effectbasedice3, item.Effectdiceperlevel1, item.Effectdiceperlevel2, item.Effectdiceperlevel3, item.Effectrealpointsperlevel1, item.Effectrealpointsperlevel2, item.Effectrealpointsperlevel3, item.Effectbasepoints1, item.Effectbasepoints2, item.Effectbasepoints3, item.Effectbonuscoefficient1, item.Effectbonuscoefficient2, item.Effectbonuscoefficient3, item.Effectmechanic1, item.Effectmechanic2, item.Effectmechanic3, item.Effectimplicittargeta1, item.Effectimplicittargeta2, item.Effectimplicittargeta3, item.Effectimplicittargetb1, item.Effectimplicittargetb2, item.Effectimplicittargetb3, item.Effectradiusindex1, item.Effectradiusindex2, item.Effectradiusindex3, item.Effectapplyauraname1, item.Effectapplyauraname2, item.Effectapplyauraname3, item.Effectamplitude1, item.Effectamplitude2, item.Effectamplitude3, item.Effectmultiplevalue1, item.Effectmultiplevalue2, item.Effectmultiplevalue3, item.Effectchaintarget1, item.Effectchaintarget2, item.Effectchaintarget3, item.Effectitemtype1, item.Effectitemtype2, item.Effectitemtype3, item.Effectmiscvalue1, item.Effectmiscvalue2, item.Effectmiscvalue3, item.Effecttriggerspell1, item.Effecttriggerspell2, item.Effecttriggerspell3, item.Effectpointspercombopoint1, item.Effectpointspercombopoint2, item.Effectpointspercombopoint3, item.Spellvisual1, item.Spellvisual2, item.Spelliconid, item.Activeiconid, item.Spellpriority, item.Name, item.Nameflags, item.Namesubtext, item.Namesubtextflags, item.Description, item.Descriptionflags, item.Auradescription, item.Auradescriptionflags, item.Manacostpercentage, item.Startrecoverycategory, item.Startrecoverytime, item.Mintargetlevel, item.Maxtargetlevel, item.Spellfamilyname, item.Spellfamilyflags, item.Maxaffectedtargets, item.Dmgclass, item.Preventiontype, item.Stancebarorder, item.Dmgmultiplier1, item.Dmgmultiplier2, item.Dmgmultiplier3, item.Minfactionid, item.Minreputation, item.Requiredauravision, item.Customflags)
        if err != nil {
            // fmt.Printf("Error importing %s: %v\n", "spell_template", err)
            continue
        }
        count++
    }

    return tx.Commit()
}

func (i *GeneratedImporter) ImportGameobjectTemplate(jsonPath string) error {
    file, err := os.Open(jsonPath)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    // Expect array
    if _, err := decoder.Token(); err != nil {
        return err
    }

    tx, err := i.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `INSERT OR REPLACE INTO gameobject_template (entry,type,displayId,name,faction,flags,size,data0,data1,data2,data3,data4,data5,data6,data7,data8,data9,data10,data11,data12,data13,data14,data15,data16,data17,data18,data19,data20,data21,data22,data23,mingold,maxgold,phase_quest_id,script_name) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    count := 0
    for decoder.More() {
        var item models.GameobjectTemplateFull
        if err := decoder.Decode(&item); err != nil {
            continue
        }

        _, err = stmt.Exec(item.Entry, item.Type, item.Displayid, item.Name, item.Faction, item.Flags, item.Size, item.Data0, item.Data1, item.Data2, item.Data3, item.Data4, item.Data5, item.Data6, item.Data7, item.Data8, item.Data9, item.Data10, item.Data11, item.Data12, item.Data13, item.Data14, item.Data15, item.Data16, item.Data17, item.Data18, item.Data19, item.Data20, item.Data21, item.Data22, item.Data23, item.Mingold, item.Maxgold, item.PhaseQuestId, item.ScriptName)
        if err != nil {
            // fmt.Printf("Error importing %s: %v\n", "gameobject_template", err)
            continue
        }
        count++
    }

    return tx.Commit()
}

// ImportItemIcons loads icon paths from item_icons.json and updates item_template
func (i *GeneratedImporter) ImportItemIcons(jsonPath string) error {
	fmt.Printf("  -> Reading item icons from %s...\n", jsonPath)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil // Icons are optional
	}

	var iconMap map[string]string
	if err := json.Unmarshal(data, &iconMap); err != nil {
		fmt.Printf("  ERROR parsing item_icons.json: %v\n", err)
		return nil
	}

	fmt.Printf("  -> Updating database with %d icon mappings...\n", len(iconMap))
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE item_template SET icon_path = ? WHERE display_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for displayIDStr, iconName := range iconMap {
		var displayID int
		fmt.Sscanf(displayIDStr, "%d", &displayID)
		if displayID > 0 {
			res, err := stmt.Exec(iconName, displayID)
			if err != nil {
				continue
			}
			if rows, _ := res.RowsAffected(); rows > 0 {
				count++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	fmt.Printf("  ✓ Successfully updated %d items with icons\n", count)
	return nil
}

// SpellEnhanced represents a spell record from spells_enhanced.json
type SpellEnhanced struct {
	SpellIconId int    `json:"spellIconId"`
	IconName    string `json:"iconName"`
}

// ImportSpellIcons loads spell icons from spells_enhanced.json and updates spell_template
func (i *GeneratedImporter) ImportSpellIcons(jsonPath string) error {
	fmt.Printf("  -> Reading spell icons from %s...\n", jsonPath)
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil // Optional
	}
	defer file.Close()

	var spells []SpellEnhanced
	if err := json.NewDecoder(file).Decode(&spells); err != nil {
		fmt.Printf("  ERROR parsing spells_enhanced.json: %v\n", err)
		return nil
	}

	// Build unique icon map
	iconMap := make(map[int]string)
	for _, s := range spells {
		if s.SpellIconId > 0 && s.IconName != "" && s.IconName != "temp" {
			iconMap[s.SpellIconId] = s.IconName
		}
	}

	fmt.Printf("  -> Updating database with %d spell icon mappings...\n", len(iconMap))
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE spell_template SET iconName = ? WHERE spellIconId = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for iconId, iconName := range iconMap {
		res, err := stmt.Exec(iconName, iconId)
		if err != nil {
			continue
		}
		if rows, _ := res.RowsAffected(); rows > 0 {
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	fmt.Printf("  ✓ Successfully updated %d spells with icons\n", count)
	return nil
}
