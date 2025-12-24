import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_quests():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching quests...")
        # Fetching all details needed for the detail page
        query = """
            SELECT 
                entry, Title, MinLevel, QuestLevel, Type, ZoneOrSort, 
                Details, Objectives, OfferRewardText, EndText,
                RewXP, RewOrReqMoney, RewMoneyMaxLevel, RewSpell,
                RewItemId1, RewItemId2, RewItemId3, RewItemId4,
                RewItemCount1, RewItemCount2, RewItemCount3, RewItemCount4,
                RewChoiceItemId1, RewChoiceItemId2, RewChoiceItemId3, RewChoiceItemId4, RewChoiceItemId5, RewChoiceItemId6,
                RewChoiceItemCount1, RewChoiceItemCount2, RewChoiceItemCount3, RewChoiceItemCount4, RewChoiceItemCount5, RewChoiceItemCount6,
                RewRepFaction1, RewRepFaction2, RewRepFaction3, RewRepFaction4, RewRepFaction5,
                RewRepValue1, RewRepValue2, RewRepValue3, RewRepValue4, RewRepValue5,
                PrevQuestId, NextQuestId, ExclusiveGroup, NextQuestInChain,
                RequiredRaces, RequiredClasses, SrcItemId
            FROM quest_template
            ORDER BY entry
        """
        mysql_cursor.execute(query)
        quests = mysql_cursor.fetchall()
        print(f"Found {len(quests)} quests")
        mysql_cursor.close()
        mysql_conn.close()
        
        # Ensure data directory exists
        os.makedirs(DATA_DIR, exist_ok=True)
        
        # Export to JSON
        json_path = os.path.join(DATA_DIR, 'quests.json')
        print(f"Exporting to {json_path}...")
        
        with open(json_path, 'w', encoding='utf-8') as f:
            # We dump the dictionary list directly. 
            # Go struct tags will need to match these keys (e.g. Title, MinLevel...) 
            # or we map them here. For simplicity, let's keep MySQL column names which are mostly PascalCase.
            json.dump(quests, f, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_quests()
